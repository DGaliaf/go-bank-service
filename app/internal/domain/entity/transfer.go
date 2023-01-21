package entity

type TransferMoney struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Amount int `json:"amount"`
}
