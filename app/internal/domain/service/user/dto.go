package user

type UserChargeMoneyDTO struct {
	Id     int `json:"id"`
	Amount int `json:"amount"`
}

type UserRemoveMoneyDTO struct {
	Id     int `json:"id"`
	Amount int `json:"amount"`
}

type TransferMoneyDTO struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Amount int `json:"amount"`
}
