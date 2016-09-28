package model

type Account struct {
	Number 			string
	Type 			string
	Name 			string
	
	TransactionID 	string
	CardNumber 		string
	CardType 		string
	NameonCard 		string
	Expiry 			string
	DisplayName 	string
}



type AccountStatus struct {
	Number 			string
	Status 			string
}

type PaymentInfo struct {
	ToInstituteID	string
	FromInstituteID string
	AccountID 		string
	BankId			string
	CUSDECNumber	string
	CUSDECInfo		string
	AmountPayable	float64
	AmountToPay		float64
	IsAccepted		bool
}

type PrintDocument struct {
	Title 	string
	Header 	map[string]interface{}
	Body 	map[string] interface{}
}

type ConfirmedDetails struct {
	AccountID		string
	IsVerified		bool
}

type Trnasction struct {
	TransactionID	string
	IsVerified		bool
}

type Institute struct {
	InstituteID 		string
	InstituteName		string
}