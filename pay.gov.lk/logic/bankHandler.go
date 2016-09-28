package logic

import (
	"pay.gov.lk/model"
)

type BankHandler struct{

}


func (p BankHandler) ConfirmAcc(Id string) model.ConfirmedDetails {
	return model.ConfirmedDetails {}
}


func NewBankHandler() BankHandler{
	return BankHandler{}
}