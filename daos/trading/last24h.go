package trading

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Last24h struct {
	Low 			decimal.NullDecimal			`json:"low"`
	High 			decimal.NullDecimal			`json:"high"`
	Volume			decimal.NullDecimal			`json:"volume"`
	LastPrice		decimal.Decimal				`json:"lastPrice"`
	Change			decimal.Decimal				`json:"change"`
	State			string						`json:"state"`
}

func (last24h Last24h) JsonLast24h() (string, error) {
	last24hJson, err := json.Marshal(last24h)
	if err != nil {
		return "", err
	}

	return string(last24hJson), nil
}

type TradingLast24h struct {
	Type 		string 						`json:"type"`
	Payload		Last24h						`json:"payload"`
}
