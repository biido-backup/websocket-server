package daos

type WebsocketRequest struct {
	Method 		string			`json:"method"`
	Username	string 			`json:"username"`
	Topic		string			`json:"topic"`
	Interval	string			`json:"interval"`  //1m, 1h, 1d, 1w, 1M
}
