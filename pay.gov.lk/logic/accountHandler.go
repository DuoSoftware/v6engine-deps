package logic

import (
	"pay.gov.lk/model"
)

type AccountHandler struct{

}

func (p AccountHandler) GetAllAccounts() model.Account {
	return model.Account {}
}

func (p AccountHandler) GetAccount(Id string) model.Account {
	return model.Account {}
}

func (p AccountHandler) AddAccount(u model.Account) {
	//return Account{"", "", "", "", "", false}
}

func (p AccountHandler) DeactivateAccount(Id string) {
	
}

func NewAccountHandler() AccountHandler{
	return AccountHandler{}
}