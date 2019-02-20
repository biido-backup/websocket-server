package daos

import (
	"github.com/shopspring/decimal"
)

type PriceAmount struct {
	Price		decimal.Decimal
	Amount		decimal.Decimal
}
