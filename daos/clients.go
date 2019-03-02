package daos

import (
	"sync"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"websocket-server/util/logger"
)

var log = logger.CreateLog("daos")

//
//type Clients struct {
//	mu         	sync.RWMutex
//	Clients           	map[string] map[string] map[string] sockjs.Session //topic -> username -> sessionID : Session
//	ClientSessions    	map[string] map[string] string                     //topic -> sessionId : username
//	Intervals 			map[string] map[string] map[string] sockjs.Session 	//topic -> interval -> sessionID : Session
//	IntervalSessions 	map[string] map[string] string    					//topic -> sessionId : interval
//}
//
//
//func CreateClients() Clients{
//
//	var MyClients Clients;
//
//	MyClients.mu.Lock()
//	MyClients.Clients = make(map[string] map[string] map[string] sockjs.Session)
//	MyClients.ClientSessions = make(map[string] map[string] string)
//	MyClients.Intervals = make(map[string] map[string] map[string] sockjs.Session)
//	MyClients.IntervalSessions = make(map[string] map[string] string)
//	MyClients.mu.Unlock()
//
//	return MyClients
//
//}
//
//func (MyClients Clients) SetTopic(topic string) {
//	MyClients.mu.Lock()
//	MyClients.Clients[topic] = make(map[string] map[string] sockjs.Session)
//	MyClients.ClientSessions[topic] = make(map[string] string)
//	MyClients.Intervals[topic] = make(map[string] map[string] sockjs.Session)
//	MyClients.IntervalSessions[topic] = make(map[string] string)
//	MyClients.mu.Unlock()
//}
//
//func (MyClients Clients) AddSubscriber(subscriber WebsocketRequest,session sockjs.Session){
//	//CLients
//	MyClients.mu.Lock()
//	clientSession := MyClients.GetListSessionByTopicAndUsername(subscriber.Topic, subscriber.Username)
//
//	if clientSession == nil{
//		MyClients.Clients[subscriber.Topic][subscriber.Username] = make(map[string] sockjs.Session)
//	}
//	MyClients.Clients[subscriber.Topic][subscriber.Username][session.ID()] = session
//	MyClients.ClientSessions[subscriber.Topic][session.ID()] = subscriber.Username
//
//	//Interval
//	intervalSession := MyClients.GetListSessionByTopicAndInterval(subscriber.Topic, subscriber.Interval)
//	if intervalSession == nil{
//		MyClients.Intervals[subscriber.Topic][subscriber.Interval] = make(map[string] sockjs.Session)
//	}
//	MyClients.Intervals[subscriber.Topic][subscriber.Interval][session.ID()] = session
//	MyClients.IntervalSessions[subscriber.Topic][session.ID()] = subscriber.Interval
//	MyClients.mu.Unlock()
//}
//
//func (MyClients Clients) RemoveSubscriber(topic string, sessionID string){
//	username := MyClients.GetUsernameByTopicAndSessionId(topic, sessionID)
//
//	MyClients.mu.Lock()
//	delete(MyClients.Clients[topic][username], sessionID)
//
//	clientSession := MyClients.GetListSessionByTopicAndUsername(topic, username)
//	if len(clientSession) == 0 {
//		delete(MyClients.Clients[topic], username)
//	}
//	delete(MyClients.ClientSessions[topic], sessionID)
//
//	interval := MyClients.GetIntervalByTopicAndSessionId(topic, sessionID)
//
//	delete(MyClients.Intervals[topic][interval], sessionID)
//
//	intervalSession := MyClients.GetListSessionByTopicAndInterval(topic, interval)
//
//	if len(intervalSession) == 0 {
//		delete(MyClients.Intervals[topic], interval)
//	}
//
//	delete(MyClients.IntervalSessions[topic], sessionID)
//	MyClients.mu.Unlock()
//
//}
//
//func (MyClients Clients) GetUsernameByTopicAndSessionId(topic string, sessionId string) string{
//	var username string
//	MyClients.mu.RLock()
//	username = MyClients.ClientSessions[topic][sessionId]
//	MyClients.mu.RUnlock()
//	return username
//}
//
//func (MyClients Clients) GetIntervalByTopicAndSessionId(topic string, sessionId string) string{
//	var interval string
//	MyClients.mu.RLock()
//	interval = MyClients.IntervalSessions[topic][sessionId]
//	MyClients.mu.RUnlock()
//	return interval
//}
//
//func (MyClients Clients) GetAllClientsByTopic(topic string) map[string] map[string] sockjs.Session{
//	var c map[string] map[string] sockjs.Session
//	MyClients.mu.RLock()
//	c = MyClients.Clients[topic]
//	MyClients.mu.RUnlock()
//	return c
//}
//
//func (MyClients Clients) GetListSessionByTopicAndUsername(topic string, username string) map[string] sockjs.Session{
//	sessions := make(map[string] sockjs.Session)
//	MyClients.mu.RLock()
//	sessions = MyClients.Clients[topic][username]
//	MyClients.mu.RUnlock()
//	return sessions
//}
//
//func (MyClients Clients) GetListSessionByTopicAndInterval(topic string, interval string) map[string] sockjs.Session{
//	sessions := make(map[string] sockjs.Session)
//	MyClients.mu.RLock()
//	sessions = MyClients.Intervals[topic][interval]
//	MyClients.mu.RUnlock()
//	return sessions
//}


type Clients struct {
	muClients         	*sync.RWMutex
	Clients           	map[string] map[string] map[string] sockjs.Session //topic -> username -> sessionID : Session

	muClientsSessions 	*sync.RWMutex
	ClientSessions    	map[string] map[string] string                     //topic -> sessionId : username

	muIntervals 		*sync.RWMutex
	Intervals 			map[string] map[string] map[string] sockjs.Session 	//topic -> interval -> sessionID : Session

	muIntervalSessions 	*sync.RWMutex
	IntervalSessions 	map[string] map[string] string    					//topic -> sessionId : interval
}

var MyClients Clients


func CreateClients(){

	//var MyClients Clients;

	MyClients.muClients = &sync.RWMutex{}
	MyClients.muClientsSessions = &sync.RWMutex{}
	MyClients.muIntervals = &sync.RWMutex{}
	MyClients.muIntervalSessions = &sync.RWMutex{}

	MyClients.muClients.Lock()
	MyClients.Clients = make(map[string] map[string] map[string] sockjs.Session)
	MyClients.muClients.Unlock()

	MyClients.muClientsSessions.Lock()
	MyClients.ClientSessions = make(map[string] map[string] string)
	MyClients.muClientsSessions.Unlock()

	MyClients.muIntervals.Lock()
	MyClients.Intervals = make(map[string] map[string] map[string] sockjs.Session)
	MyClients.muIntervals.Unlock()

	MyClients.muIntervalSessions.Lock()
	MyClients.IntervalSessions = make(map[string] map[string] string)
	MyClients.muIntervalSessions.Unlock()

	//return MyClients

}

func (clients Clients) SetTopic(topic string) {
	clients.muClients.Lock()
	clients.Clients[topic] = make(map[string] map[string] sockjs.Session)
	clients.muClients.Unlock()

	clients.muClientsSessions.Lock()
	clients.ClientSessions[topic] = make(map[string] string)
	clients.muClientsSessions.Unlock()

	clients.muIntervals.Lock()
	clients.Intervals[topic] = make(map[string] map[string] sockjs.Session)
	clients.muIntervals.Unlock()

	clients.muIntervalSessions.Lock()
	clients.IntervalSessions[topic] = make(map[string] string)
	clients.muIntervalSessions.Unlock()
}

func (clients Clients) AddSubscriber(request WebsocketRequest,session sockjs.Session){
	//CLients
	clientSession := clients.GetListSessionByTopicAndUsername(request.Topic, request.Username)
	if clientSession == nil{
		clients.muClients.Lock()
		clients.Clients[request.Topic][request.Username] = make(map[string] sockjs.Session)
		clients.muClients.Unlock()
	}
	clients.muClients.Lock()
	clients.Clients[request.Topic][request.Username][session.ID()] = session
	clients.muClients.Unlock()

	//ClientSession
	clients.muClientsSessions.Lock()
	clients.ClientSessions[request.Topic][session.ID()] = request.Username
	clients.muClientsSessions.Unlock()

	//Interval
	intervalSession := clients.GetListSessionByTopicAndInterval(request.Topic, request.Interval)
	if intervalSession == nil{
		clients.muIntervals.Lock()
		clients.Intervals[request.Topic][request.Interval] = make(map[string] sockjs.Session)
		clients.muIntervals.Unlock()
	}
	clients.muIntervals.Lock()
	clients.Intervals[request.Topic][request.Interval][session.ID()] = session
	clients.muIntervals.Unlock()

	//Interval Session
	clients.muIntervalSessions.Lock()
	clients.IntervalSessions[request.Topic][session.ID()] = request.Interval
	clients.muIntervalSessions.Unlock()
}

func (clients Clients) RemoveSubscriber(topic string, sessionID string){
	username := clients.GetUsernameByTopicAndSessionId(topic, sessionID)

	clients.muClients.Lock()
	delete(clients.Clients[topic][username], sessionID)
	clients.muClients.Unlock()

	clientSession := clients.GetListSessionByTopicAndUsername(topic, username)
	if len(clientSession) == 0 {
		clients.muClients.Lock()
		delete(clients.Clients[topic], username)
		clients.muClients.Unlock()
	}
	clients.muClientsSessions.Lock()
	delete(clients.ClientSessions[topic], sessionID)
	clients.muClientsSessions.Unlock()

	interval := clients.GetIntervalByTopicAndSessionId(topic, sessionID)

	clients.muIntervals.Lock()
	delete(clients.Intervals[topic][interval], sessionID)
	clients.muIntervals.Unlock()

	intervalSession := clients.GetListSessionByTopicAndInterval(topic, interval)

	if len(intervalSession) == 0 {
		clients.muIntervals.Lock()
		delete(clients.Intervals[topic], interval)
		clients.muIntervals.Unlock()
	}

	clients.muIntervalSessions.Lock()
	delete(clients.IntervalSessions[topic], sessionID)
	clients.muIntervalSessions.Unlock()

}

func (clients Clients) GetUsernameByTopicAndSessionId(topic string, sessionId string) string{
	var username string
	clients.muClientsSessions.RLock()
	username = clients.ClientSessions[topic][sessionId]
	clients.muClientsSessions.RUnlock()
	return username
}

func (clients Clients) GetIntervalByTopicAndSessionId(topic string, sessionId string) string{
	var interval string
	clients.muIntervalSessions.RLock()
	interval = clients.IntervalSessions[topic][sessionId]
	clients.muIntervalSessions.RUnlock()
	return interval
}

func (clients Clients) GetAllClientsByTopic(topic string) map[string] map[string] sockjs.Session{
	var c map[string] map[string] sockjs.Session
	clients.muClients.RLock()
	c = clients.Clients[topic]
	clients.muClients.RUnlock()
	return c
}

func (clients Clients) GetListSessionByTopicAndUsername(topic string, username string) map[string] sockjs.Session{
	sessions := make(map[string] sockjs.Session)
	clients.muClients.RLock()
	sessions = clients.Clients[topic][username]
	clients.muClients.RUnlock()
	return sessions
}

func (clients Clients) GetListSessionByTopicAndInterval(topic string, interval string) map[string] sockjs.Session{
	sessions := make(map[string] sockjs.Session)
	clients.muIntervals.RLock()
	sessions = clients.Intervals[topic][interval]
	clients.muIntervals.RUnlock()
	return sessions
}

func (clients Clients) GetSessionByTopicAndUsernameAndSessionId(topic string, username string, sessionId string) (sockjs.Session, bool){
	clients.muClients.RLock()
	session, exist := clients.Clients[topic][username][sessionId]
	clients.muClients.RUnlock()
	return session, exist
}

