package lib

import (
	"duov6.com/gorest"
	"pay.gov.lk/model"
	"pay.gov.lk/logic"
)

type BankService struct {
	gorest.RestService

	confirmAcc gorest.EndPoint `method:"GET" path:"/bank/confirmacc/{Id:string}/" output:"ConfirmedDetails"`
	rejectAcc gorest.EndPoint `method:"GET" path:"/bank/rejectacc/{Id:string}/" output:"ConfirmedDetails"`

	getAll gorest.EndPoint `method:"GET" path:"/bank/" output:"[]Institute"`
	getOne gorest.EndPoint `method:"GET" path:"/bank/{Id:string}/" output:"Institute"`
}

func (p BankService) ConfirmAcc(Id string) model.ConfirmedDetails {
	h := logic.NewBankHandler()
	return h.ConfirmAcc(Id)
}

func (p BankService) RejectAcc(Id string) model.ConfirmedDetails {
	h := logic.NewBankHandler()
	return h.ConfirmAcc(Id)
}

func (p BankService) GetAll() []model.Institute {
	//h := logic.NewBankHandler()
	a :=make([]model.Institute,0)
	return a
}

func (p BankService) GetOne(Id string) model.Institute {
	//h := logic.NewBankHandler()
	return model.Institute{}
}