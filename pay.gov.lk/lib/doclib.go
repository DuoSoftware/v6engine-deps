package lib

import (
	"duov6.com/gorest"
	"pay.gov.lk/model"
	"pay.gov.lk/logic"
)

type DocService struct {
	gorest.RestService

	documentAccConfirm gorest.EndPoint `method:"GET" path:"/documents/confirmacc/{Id:string}/" output:"PrintDocument"`
	documentTranReciept gorest.EndPoint `method:"GET" path:"/documents/tranreciept/{Id:string}/" output:"PrintDocument"`
}

func (p DocService) DocumentAccConfirm(Id string) model.PrintDocument {
	h := logic.NewDocumentHandler();
	return h.DocumentAccConfirm(Id)
}

func (p DocService) DocumentTranReciept(Id string) model.PrintDocument {
	h := logic.NewDocumentHandler();
	return h.DocumentTranReciept(Id)
}