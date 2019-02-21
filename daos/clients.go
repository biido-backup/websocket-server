package daos

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Clients struct {
	Clients map[string] map[string] map[string] sockjs.Session  	//topic -> username -> sessionID : Session
	Sessions map[string] map[string] string 						//topic -> sessionId : username
	Interval map[string] map[string] string							//topic -> interval : sessionId
}


func CreateClients() *Clients{
	var clients Clients

	clients.Clients = make(map[string] map[string] map[string] sockjs.Session)
	clients.Sessions = make(map[string] map[string] string)

	return &clients
}

func (clients Clients) SetTopic(topic string) {
	clients.Clients[topic] = make(map[string] map[string] sockjs.Session)
	clients.Sessions[topic] = make(map[string] string)
}

func (clients Clients) AddClient(topic string, username string,session sockjs.Session){
	if clients.Clients[topic][username] == nil{
		clients.Clients[topic][username] = make(map[string] sockjs.Session)
	}
	clients.Clients[topic][username][session.ID()] = session
	clients.Sessions[topic][session.ID()] = username
}

func (clients Clients) DeleteClient(topic string, sessionID string){
	username := clients.Sessions[topic][sessionID]

	delete(clients.Clients[topic][username], sessionID)
	if len(clients.Clients[topic][username]) == 0 {
		delete(clients.Clients[topic], username)
	}
	delete(clients.Sessions[topic], sessionID)
}