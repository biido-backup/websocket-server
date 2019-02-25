package trading

import (
	"database/sql"
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

type History struct {
	Id 					uint64					`json:"id"`
	Amount       		decimal.Decimal			`json:"amount"`
	Price        		decimal.Decimal			`json:"price"`
	Type				sql.NullString			`json:"type"`
	Time				time.Time				`json:"time"`
}

type ListHistory struct {
	Histories 			[]History				`json:"histories"`
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
	Payload			ListHistory						`json:"payload"`
}
