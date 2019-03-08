package trading

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type History struct {
	Id 					uint64					`json:"id"`
	AskId				uint64		 			`json:"ask_id"`
	BidId				uint64		 			`json:"bid_id"`
	Amount       		decimal.Decimal			`json:"amount"`
	Price        		decimal.Decimal			`json:"price"`
	Type				sql.NullString			`json:"type"`
	Time				pq.NullTime				`json:"time"`
}

type ListHistory struct {
	Rate      	string    `json:"rate"`
	Histories 	[]History `json:"histories"`
}

func (listHistory ListHistory) JsonListHistory() (string, error) {
	listHistoryJson, err := json.Marshal(listHistory)
	if err != nil {
		return "", err
	}

	return string(listHistoryJson), nil
}

type TradingListHistory struct {
	Type 			string 						`json:"type"`
	Payload			[]History					`json:"payload"`
}
