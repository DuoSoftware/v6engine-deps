package lib

import (
	"duov6.com/gorest"
	"pay.gov.lk/model"
	"pay.gov.lk/logic"
)

type PayService struct {
	gorest.RestService

	pay gorest.EndPoint `method:"POST" path:"/accounts/pay/" postdata:"PaymentInfo"`
}

func (p PayService) Pay(u model.PaymentInfo) {
	h := logic.NewPaymentHandler();
	h.Pay(u)
}