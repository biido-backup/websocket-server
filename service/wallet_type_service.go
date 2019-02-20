package service

import (
	"database/sql"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"websocket-server/util/database"
)

var log = logrus.New()

func GetAllWalletType() ([]sql.NullString, error) {
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

