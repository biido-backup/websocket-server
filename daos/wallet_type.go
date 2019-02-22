package daos

import (
	"database/sql"
)

type WalletType struct {
	Id				uint64				`json:"id"`
	Code 			string				`json:"code"`
	Name 			string				`json:"name"`
	FlagActive		uint8				`json:"flag_active"`
	PivotPriority	sql.NullInt64		`json:"pivot_priority"`
	IsFiat			sql.NullBool		`json:"is_fiat"`
}
