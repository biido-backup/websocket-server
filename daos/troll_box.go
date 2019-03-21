package daos

import (
	"encoding/json"
	"github.com/spf13/viper"
	"time"
	"websocket-server/daos/reqpayload"
	"websocket-server/util/redis"
)

var TrollBox []chat

type chat struct {
	Username	string 		`json:"username"`
	Message		string		`json:"message"`
	Time		int64		`json:"time"`
}

func CreateTrollBox() {
	TrollBox = make([]chat, 0, 1)
	trollBoxJson := redis.GetValueByKey(viper.GetString("redis.trollbox.key"))
	json.Unmarshal(trollBoxJson, &TrollBox)

	log.Debug("CREATE TROLL BOX", TrollBox)
}

func InsertToTrollBox(payload reqpayload.TrollBox) {
	chat := chat{Username: payload.Username, Message: payload.Message, Time: time.Now().UTC().Unix() * 1000}
	if len(TrollBox) >= viper.GetInt("trollbox.size") {
		TrollBox = TrollBox[1:]
	}
	TrollBox = append(TrollBox, chat)

	trollBoxJson, _ := json.Marshal(TrollBox)
	redis.SetValueByKey(viper.GetString("redis.trollbox.key"), trollBoxJson)

	//log.Debug("INSERT TO TROLL BOX", TrollBox)
}

func GetTrollBoxBroadcastMessage() string {
	broadcastMessage := BroadcastMessage{Type: "TROLL_BOX", Payload: TrollBox}
	broadcastMessageJson, _ := json.Marshal(broadcastMessage)
	//log.Debug("TROLL BOX BROADCAST MESSAGE", string(broadcastMessageJson))
	return string(broadcastMessageJson)
}