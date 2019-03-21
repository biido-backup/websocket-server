package daos

import (
	"encoding/json"
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
	trollBoxJson := redis.GetValueByKey("TROLL_BOX")
	json.Unmarshal(trollBoxJson, &TrollBox)

	log.Debug("CREATE TROLL BOX", TrollBox)
}

func InsertToTrollBox(payload reqpayload.TrollBox) {
	chat := chat{Username: payload.Username, Message: payload.Message, Time: time.Now().UTC().Unix() * 1000}
	if len(TrollBox) >= 10 {
		TrollBox = TrollBox[1:]
	}
	TrollBox = append(TrollBox, chat)

	trollBoxJson, _ := json.Marshal(TrollBox)
	redis.SetValueByKey("TROLL_BOX", trollBoxJson)

	log.Debug("INSERT TO TROLL BOX", TrollBox)
}

func GetTrollBoxBroadcastMessage() string {
	broadcastMessage := BroadcastMessage{Type: "TROLL_BOX", Payload: TrollBox}
	broadcastMessageJson, _ := json.Marshal(broadcastMessage)
	log.Debug("TROLL BOX BROADCAST MESSAGE", string(broadcastMessageJson))
	return string(broadcastMessageJson)
}