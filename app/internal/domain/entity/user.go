package entity

import "encoding/json"

type User struct {
	Id      int `json:"id"`
	Balance int `json:"balance"`
}

func (u User) Marshal() ([]byte, error) {
	userMap := map[string]interface{}{
		"balance": u.Balance,
	}

	marshaledUser, err := json.Marshal(userMap)
	if err != nil {
		return []byte{}, err
	}

	return marshaledUser, nil
}
