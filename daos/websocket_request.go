package daos

type WebsocketRequest struct {
	Type 		string			`json:"type"`
	Payload		interface{}		`json:"payload"`
}
