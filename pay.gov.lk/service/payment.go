package main

import (
	"pay.gov.lk/lib"
	"duov6.com/cebadapter"
	"duov6.com/gorest"
	"duov6.com/term"
	"net/http"
)


func main() {

	cebadapter.Attach("DuoAuth", func(s bool){
		cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
			term.Write("Store Configuration Successfully Loaded...", term.Information)

			agent := cebadapter.GetAgent();
			
			agent.Client.OnEvent("globalConfigChanged.StoreConfig", func(from string, name string, data map[string]interface{}, resources map[string]interface{}){
				cebadapter.GetLatestGlobalConfig("StoreConfig", func(data []interface{}) {
					term.Write("Store Configuration Successfully Updated...", term.Information)
				});
			});
		})
		term.Write ("Successfully registered in CEB", term.Information)
	});

	//paylib.SetupConfig()
	term.GetConfig()
	go runRestFul()


	term.SplashScreen("splash.art")
	term.Write("================================================================", term.Splash)
	term.Write("|     Admintration Console running on  :9000                   |", term.Splash)
	term.Write("|     https RestFul Service running on :3048                   |", term.Splash)
	term.Write("|     Duo v6 Auth Service 6.0                                  |", term.Splash)
	term.Write("================================================================", term.Splash)
	term.StartCommandLine()

}

func runRestFul() {
	gorest.RegisterService(new(lib.PayService))
	gorest.RegisterService(new(lib.BankService))
	gorest.RegisterService(new(lib.DocService))
	gorest.RegisterService(new(lib.AccountService))

/*
	c := authlib.GetConfig()
	if c.Https_Enabled {
		err := http.ListenAndServeTLS(":4048", c.Cirtifcate, c.PrivateKey, gorest.Handle())
		if err != nil {
			term.Write(err.Error(), term.Error)
			return
		}
	} else {
*/
		err := http.ListenAndServe(":4048", gorest.Handle())
		if err != nil {
			term.Write(err.Error(), term.Error)
			return
		}
//	}

}