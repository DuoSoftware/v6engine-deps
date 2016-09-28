package lib

import (
	"duov6.com/gorest"
	"pay.gov.lk/model"
	"pay.gov.lk/logic"
)

type AccountService struct {
	gorest.RestService

	getStatus gorest.EndPoint `method:"GET" path:"/account/status/{Id:string}/" output:"AccountStatus"`
	setStatus gorest.EndPoint `method:"POST" path:"/account/status/" postdata:"AccountStatus"`

	getAll gorest.EndPoint `method:"GET" path:"/account/" output:"[]Account"`
	getAccount gorest.EndPoint `method:"GET" path:"/account/{Id:string}/" output:"Account"`

	addAccount gorest.EndPoint `method:"POST" path:"/account/" postdata:"Account"`
	deactivateAccount gorest.EndPoint `method:"DELETE" path:"/account/{Id:string}/"`
}

func (p AccountService) GetStatus(Id string) model.AccountStatus {
	return  model.AccountStatus {}
}

func (p AccountService) SetStatus(u model.AccountStatus) {

}

func (p AccountService) GetAll() []model.Account {
	//h := logic.NewAccountHandler();
	a :=make([]model.Account,0)
	return a
}

func (p AccountService) GetAccount(Id string) model.Account {
	h := logic.NewAccountHandler();
	return  h.GetAccount(Id)
}

func (p AccountService) AddAccount(u model.Account) {
	h := logic.NewAccountHandler();
	h.AddAccount(u)
}

func (p AccountService) DeactivateAccount(Id string) {
	h := logic.NewAccountHandler();
	h.DeactivateAccount(Id)
}
