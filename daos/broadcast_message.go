package daos

type BroadcastMessage struct {
	Type 			string 			`json:"type"`
	Payload			interface{}		`json:"payload"`
}