package service

import (
	"database/sql"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"websocket-server/daos"
	"websocket-server/util/database"
)

var log = logrus.New()

func GetAllWalletTypeCode() ([]sql.NullString, error) {
	var walletType []sql.NullString

	err := database.Db.Select("code").
					From("wallet_type").
					Where(dbx.HashExp{"flag_active":1}).
					All(&walletType)

	if (err != nil){
		log.Error(err)
		return nil, err
	}

	return walletType, nil
}

func GetAllWalletType() []daos.WalletType{
	var walletType []daos.WalletType

	err := database.Db.Select("id", "code", "name", "flag_active", "pivot_priority", "is_fiat").
		From("wallet_type").
		Where(dbx.HashExp{"flag_active":1}).
		All(&walletType)

	if (err != nil){
		panic(err)
	}

	return walletType
}

