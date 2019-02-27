package trading

import (
	"encoding/json"
	"github.com/lib/pq"
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
	CreatedAt				pq.NullTime				`json:"created_at"`
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

type OrderHistories struct {
	Type 					string					`json:"type"`
	Size					int64					`json:"size"`
	Payload					[]OrderHistory			`json:"payload"`
}