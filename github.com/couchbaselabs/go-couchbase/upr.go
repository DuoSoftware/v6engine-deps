package couchbase

import (
	"log"
	"time"

	"fmt"
	"github.com/couchbase/gomemcached"
	"github.com/couchbase/gomemcached/client"
	"sync"
)

// A UprFeed streams mutation events from a bucket.
//
// Events from the bucket can be read from the channel 'C'.  Remember
// to call Close() on it when you're done, unless its channel has
// closed itself already.
type UprFeed struct {
	C <-chan *memcached.UprEvent

	bucket     *Bucket
	nodeFeeds  map[string]*FeedInfo     // The UPR feeds of the individual nodes
	output     chan *memcached.UprEvent // Same as C but writeably-typed
	quit       chan bool
	name       string // name of this UPR feed
	sequence   uint32 // sequence number for this feed
	connected  bool
	killSwitch chan bool
	lock       sync.Mutex // synchronize access to feed.output CBIDXT-237
}

// UprFeed from a single connection
type FeedInfo struct {
	uprFeed   *memcached.UprFeed // UPR feed handle
	host      string             // hostname
	connected bool               // connected
}

type FailoverLog map[uint16]memcached.FailoverLog

// GetFailoverLogs, get the failover logs for a set of vbucket ids
func (b *Bucket) GetFailoverLogs(vBuckets []uint16) (FailoverLog, error) {

	// map vbids to their corresponding hosts
	vbHostList := make(map[string][]uint16)
	vbm := b.VBServerMap()
	if len(vbm.VBucketMap) < len(vBuckets) {
		return nil, fmt.Errorf("vbmap smaller than vbucket list: %v vs. %v",
			vbm.VBucketMap, vBuckets)
	}

	for _, vb := range vBuckets {
		masterID := vbm.VBucketMap[vb][0]
		master := b.getMasterNode(masterID)
		if master == "" {
			return nil, fmt.Errorf("No master found for vb %d", vb)
		}

		vbList := vbHostList[master]
		if vbList == nil {
			vbList = make([]uint16, 0)
		}
		vbList = append(vbList, vb)
		vbHostList[master] = vbList
	}

	failoverLogMap := make(FailoverLog)
	for _, serverConn := range b.getConnPools() {

		vbList := vbHostList[serverConn.host]
		if vbList == nil {
			continue
		}

		mc, err := serverConn.Get()
		if err != nil {
			log.Printf("No Free connections for vblist %v", vbList)
			return nil, fmt.Errorf("No Free connections for host %s",
				serverConn.host)

		}
		// close the connection so that it doesn't get reused for upr data
		// connection
		defer mc.Close()
		failoverlogs, err := mc.UprGetFailoverLog(vbList)
		if err != nil {
			return nil, fmt.Errorf("Error getting failover log %s host %s",
				err.Error(), serverConn.host)

		}

		for vb, log := range failoverlogs {
			failoverLogMap[vb] = *log
		}
	}

	return failoverLogMap, nil
}

// StartUprFeed creates and starts a new Upr feed
// No data will be sent on the channel unless vbuckets streams are requested
func (b *Bucket) StartUprFeed(name string, sequence uint32) (*UprFeed, error) {

	feed := &UprFeed{
		bucket:     b,
		output:     make(chan *memcached.UprEvent, 10),
		quit:       make(chan bool),
		nodeFeeds:  make(map[string]*FeedInfo, 0),
		name:       name,
		sequence:   sequence,
		killSwitch: make(chan bool),
	}

	err := feed.connectToNodes()
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to bucket %s", err.Error())
	}
	feed.connected = true
	go feed.run()

	feed.C = feed.output
	return feed, nil
}

// UprRequestStream starts a stream for a vb on a feed
func (feed *UprFeed) UprRequestStream(vb uint16, opaque uint16, flags uint32,
	vuuid, startSequence, endSequence, snapStart, snapEnd uint64) error {

	vbm := feed.bucket.VBServerMap()
	if len(vbm.VBucketMap) < int(vb) {
		return fmt.Errorf("vbmap smaller than vbucket list: %v vs. %v",
			vb, vbm.VBucketMap)
	}

	if int(vb) >= len(vbm.VBucketMap) {
		return fmt.Errorf("Invalid vbucket id %d", vb)
	}

	masterID := vbm.VBucketMap[vb][0]
	master := feed.bucket.getMasterNode(masterID)
	if master == "" {
		return fmt.Errorf("Master node not found for vbucket %d", vb)
	}
	singleFeed := feed.nodeFeeds[master]
	if singleFeed == nil {
		return fmt.Errorf("UprFeed for this host not found")
	}

	if err := singleFeed.uprFeed.UprRequestStream(vb, opaque, flags,
		vuuid, startSequence, endSequence, snapStart, snapEnd); err != nil {
		return err
	}

	return nil
}

// UprCloseStream ends a vbucket stream.
func (feed *UprFeed) UprCloseStream(vb, opaqueMSB uint16) error {
	vbm := feed.bucket.VBServerMap()
	if len(vbm.VBucketMap) < int(vb) {
		return fmt.Errorf("vbmap smaller than vbucket list: %v vs. %v",
			vb, vbm.VBucketMap)
	}

	if int(vb) >= len(vbm.VBucketMap) {
		return fmt.Errorf("Invalid vbucket id %d", vb)
	}

	masterID := vbm.VBucketMap[vb][0]
	master := feed.bucket.getMasterNode(masterID)
	if master == "" {
		return fmt.Errorf("Master node not found for vbucket %d", vb)
	}
	singleFeed := feed.nodeFeeds[master]
	if singleFeed == nil {
		return fmt.Errorf("UprFeed for this host not found")
	}

	if err := singleFeed.uprFeed.CloseStream(vb, opaqueMSB); err != nil {
		return err
	}
	return nil
}

// Goroutine that runs the feed
func (feed *UprFeed) run() {
	retryInterval := initialRetryInterval
	bucketOK := true
	for {
		// Connect to the UPR feed of each server node:
		if bucketOK {
			// Run until one of the sub-feeds fails:
			select {
			case <-feed.killSwitch:
			case <-feed.quit:
				return
			}
			//feed.closeNodeFeeds()
			retryInterval = initialRetryInterval
		}

		// On error, try to refresh the bucket in case the list of nodes changed:
		log.Printf("go-couchbase: UPR connection lost; reconnecting to bucket %q in %v",
			feed.bucket.Name, retryInterval)

		if err := feed.bucket.Refresh(); err != nil {
			log.Printf("Unable to refresh bucket %s ", err.Error())
			feed.closeNodeFeeds()
		}
		// this will only connect to nodes that are not connected or changed
		// user will have to reconnect the stream
		err := feed.connectToNodes()
		bucketOK = err == nil

		select {
		case <-time.After(retryInterval):
		case <-feed.quit:
			return
		}
		if retryInterval *= 2; retryInterval > maximumRetryInterval {
			retryInterval = maximumRetryInterval
		}
	}
}

func (feed *UprFeed) connectToNodes() (err error) {
	for _, serverConn := range feed.bucket.getConnPools() {

		// this maybe a reconnection, so check if the connection to the node
		// already exists. Connect only if the node is not found in the list
		// or connected == false
		nodeFeed := feed.nodeFeeds[serverConn.host]

		if nodeFeed != nil && nodeFeed.connected == true {
			continue
		}

		var singleFeed *memcached.UprFeed
		var name string
		if feed.name == "" {
			name = "DefaultUprClient"
		} else {
			name = feed.name
		}
		singleFeed, err = serverConn.StartUprFeed(name, feed.sequence)
		if err != nil {
			log.Printf("go-couchbase: Error connecting to upr feed of %s: %v", serverConn.host, err)
			feed.closeNodeFeeds()
			return
		}
		// add the node to the connection map
		feedInfo := &FeedInfo{
			uprFeed:   singleFeed,
			connected: true,
			host:      serverConn.host,
		}
		feed.nodeFeeds[serverConn.host] = feedInfo
		go feed.forwardUprEvents(feedInfo, feed.killSwitch, serverConn.host)
	}
	return
}

// Goroutine that forwards Upr events from a single node's feed to the aggregate feed.
func (feed *UprFeed) forwardUprEvents(nodeFeed *FeedInfo, killSwitch chan bool, host string) {
	singleFeed := nodeFeed.uprFeed

	for {
		select {
		case event, ok := <-singleFeed.C:
			if !ok {
				if singleFeed.Error != nil {
					log.Printf("go-couchbase: Upr feed from %s failed: %v", host, singleFeed.Error)
				}
				killSwitch <- true
				return
			}
			feed.lock.Lock()
			feed.output <- event
			feed.lock.Unlock()
			if event.Status == gomemcached.NOT_MY_VBUCKET {
				log.Printf(" Got a not my vbucket error !! ")
				if err := feed.bucket.Refresh(); err != nil {
					log.Printf("Unable to refresh bucket %s ", err.Error())
					feed.closeNodeFeeds()
					return
				}
				// this will only connect to nodes that are not connected or changed
				// user will have to reconnect the stream
				if err := feed.connectToNodes(); err != nil {
					log.Printf("Unable to connect to nodes %s", err.Error())
					return
				}

			}
		case <-feed.quit:
			nodeFeed.connected = false
			return
		}
	}
}

func (feed *UprFeed) closeNodeFeeds() {
	for _, f := range feed.nodeFeeds {
		f.uprFeed.Close()
	}
	feed.nodeFeeds = nil
}

// Close a Upr feed.
func (feed *UprFeed) Close() error {
	select {
	case <-feed.quit:
		return nil
	default:
	}

	feed.closeNodeFeeds()
	close(feed.quit)

	feed.lock.Lock()
	defer feed.lock.Unlock()
	close(feed.output)

	return nil
}
