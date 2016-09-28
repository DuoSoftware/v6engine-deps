package logic

import (
	"pay.gov.lk/model"
)

type PaymentHandler struct{

}

func (p PaymentHandler) Pay(u model.PaymentInfo) {
	
}

func NewPaymentHandler() PaymentHandler{
	return PaymentHandler{}
}