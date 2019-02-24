package trading

import (
	"encoding/json"
	"time"
)

type OrderHistory struct {
	Id						int64 					`json:"id"`
	Rate					string 					`json:"rate"`
	Type					string 					`json:"type"`
	Price					string 					`json:"price"`
	Amount					string 					`json:"amount"`
	TotalAmount				string 					`json:"total_amount"`
	AdminFee				string					`json:"admin_fee"`
	MainPrecision			string					`json:"main_precision"`
	PivotPrecision			string					`json:"pivot_precision"`
	CreatedAt				time.Time				`json:"created_at"`
}

type ListOrderHistory struct {
	Username 				string					`json:"username"`
	Size					int64					`json:"size"`
	OrderHistories			[]OrderHistory			`json:"orderHistories"`
}

func (listOrderHistory ListOrderHistory) JsonListOrderHistory() (string, error) {
	listOrderHistoryJson, err := json.Marshal(listOrderHistory)
	if err != nil {
		return "", err
	}

	return string(listOrderHistoryJson), nil
}