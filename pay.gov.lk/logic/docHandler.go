package logic

import (
	"pay.gov.lk/model"
)

type DocumentHandler struct{

}

func (p DocumentHandler) DocumentAccConfirm(Id string) model.PrintDocument {
	return model.PrintDocument {}
}

func (p DocumentHandler) DocumentTranReciept(Id string) model.PrintDocument {
	return model.PrintDocument {}
}

func NewDocumentHandler() DocumentHandler {
	return DocumentHandler{}
}