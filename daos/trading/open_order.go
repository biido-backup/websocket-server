package trading

import (
	"encoding/json"
	"github.com/lib/pq"
)

type OpenOrder struct {
	Id						int64 					`json:"id"`
	Rate					string 					`json:"rate"`
	Type					string 					`json:"type"`
	Price					string 					`json:"price"`
	Amount					string 					`json:"amount"`
	FilledAmount			string 					`json:"filled_amount"`
	RemainderAmount			string					`json:"remainder_amount"`
	TotalRemainderAmount	string					`json:"total_remainder_amount"`
	AdminFee				string					`json:"admin_fee"`
	MainPrecision			string					`json:"main_precision"`
	PivotPrecision			string					`json:"pivot_precision"`
	CreatedAt				pq.NullTime				`json:"created_at"`
	CreatedAtUnix			int64					`json:"created_at_unix"`
}

type ListOpenOrder struct {
	Username 				string					`json:"username"`
	OpenOrders				[]OpenOrder				`json:"openOrders"`
}

func (listOpenOrder ListOpenOrder) JsonListOpenOrder() (string, error) {
	listOpenOrderJson, err := json.Marshal(listOpenOrder)
	if err != nil {
		return "", err
	}

	return string(listOpenOrderJson), nil
}

type OpenOrders struct {
	Type 				string					`json:"type"`
	Payload				[]OpenOrder				`json:"payload"`
}
