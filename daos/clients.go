package daos

import (
	"sync"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Clients struct {
	mu         	sync.RWMutex
	Clients           	map[string] map[string] map[string] sockjs.Session //topic -> username -> sessionID : Session
	ClientSessions    	map[string] map[string] string                     //topic -> sessionId : username
	Intervals 			map[string] map[string] map[string] sockjs.Session 	//topic -> interval -> sessionID : Session
	IntervalSessions 	map[string] map[string] string    					//topic -> sessionId : interval
}


func CreateClients() Clients{

	var clients Clients;

	clients.mu.Lock()
	clients.Clients = make(map[string] map[string] map[string] sockjs.Session)
	clients.ClientSessions = make(map[string] map[string] string)
	clients.Intervals = make(map[string] map[string] map[string] sockjs.Session)
	clients.IntervalSessions = make(map[string] map[string] string)
	clients.mu.Unlock()

	return clients

}

func (clients Clients) SetTopic(topic string) {
	clients.mu.Lock()
	clients.Clients[topic] = make(map[string] map[string] sockjs.Session)
	clients.ClientSessions[topic] = make(map[string] string)
	clients.Intervals[topic] = make(map[string] map[string] sockjs.Session)
	clients.IntervalSessions[topic] = make(map[string] string)
	clients.mu.Unlock()
}

func (clients Clients) AddSubscriber(subscriber Subscriber,session sockjs.Session){
	//CLients
	clients.mu.Lock()
	clientSession := clients.GetListSessionByTopicAndUsername(subscriber.Topic, subscriber.Username)

	if clientSession == nil{
		clients.Clients[subscriber.Topic][subscriber.Username] = make(map[string] sockjs.Session)
	}
	clients.Clients[subscriber.Topic][subscriber.Username][session.ID()] = session
	clients.ClientSessions[subscriber.Topic][session.ID()] = subscriber.Username

	//Interval
	intervalSession := clients.GetListSessionByTopicAndInterval(subscriber.Topic, subscriber.Interval)
	if intervalSession == nil{
		clients.Intervals[subscriber.Topic][subscriber.Interval] = make(map[string] sockjs.Session)
	}
	clients.Intervals[subscriber.Topic][subscriber.Interval][session.ID()] = session
	clients.IntervalSessions[subscriber.Topic][session.ID()] = subscriber.Interval
	clients.mu.Unlock()
}

func (clients Clients) RemoveSubscriber(topic string, sessionID string){
	username := clients.GetUsernameByTopicAndSessionId(topic, sessionID)

	clients.mu.Lock()
	delete(clients.Clients[topic][username], sessionID)

	clientSession := clients.GetListSessionByTopicAndUsername(topic, username)
	if len(clientSession) == 0 {
		delete(clients.Clients[topic], username)
	}
	delete(clients.ClientSessions[topic], sessionID)

	interval := clients.GetIntervalByTopicAndSessionId(topic, sessionID)

	delete(clients.Intervals[topic][interval], sessionID)

	intervalSession := clients.GetListSessionByTopicAndInterval(topic, interval)

	if len(intervalSession) == 0 {
		delete(clients.Intervals[topic], interval)
	}

	delete(clients.IntervalSessions[topic], sessionID)
	clients.mu.Unlock()

}

func (clients Clients) GetUsernameByTopicAndSessionId(topic string, sessionId string) string{
	var username string
	clients.mu.RLock()
	username = clients.ClientSessions[topic][sessionId]
	clients.mu.RUnlock()
	return username
}

func (clients Clients) GetIntervalByTopicAndSessionId(topic string, sessionId string) string{
	var interval string
	clients.mu.RLock()
	interval = clients.IntervalSessions[topic][sessionId]
	clients.mu.RUnlock()
	return interval
}

func (clients Clients) GetAllClientsByTopic(topic string) map[string] map[string] sockjs.Session{
	var c map[string] map[string] sockjs.Session
	clients.mu.RLock()
	c = clients.Clients[topic]
	clients.mu.RUnlock()
	return c
}

func (clients Clients) GetListSessionByTopicAndUsername(topic string, username string) map[string] sockjs.Session{
	sessions := make(map[string] sockjs.Session)
	clients.mu.RLock()
	sessions = clients.Clients[topic][username]
	clients.mu.RUnlock()
	return sessions
}

func (clients Clients) GetListSessionByTopicAndInterval(topic string, interval string) map[string] sockjs.Session{
	sessions := make(map[string] sockjs.Session)
	clients.mu.RLock()
	sessions = clients.Intervals[topic][interval]
	clients.mu.RUnlock()
	return sessions
}


//type Clients struct {
//	muClients         	sync.RWMutex
//	Clients           	map[string] map[string] map[string] sockjs.Session //topic -> username -> sessionID : Session
//
//	muClientsSessions 	sync.RWMutex
//	ClientSessions    	map[string] map[string] string                     //topic -> sessionId : username
//
//	muIntervals 		sync.RWMutex
//	Intervals 			map[string] map[string] map[string] sockjs.Session 	//topic -> interval -> sessionID : Session
//
//	muIntervalSessions 	sync.RWMutex
//	IntervalSessions 	map[string] map[string] string    					//topic -> sessionId : interval
//}
//
//
//func CreateClients() Clients{
//
//	var clients Clients;
//
//	clients.muClients.Lock()
//	clients.Clients = make(map[string] map[string] map[string] sockjs.Session)
//	clients.muClients.Unlock()
//
//	clients.muClientsSessions.Lock()
//	clients.ClientSessions = make(map[string] map[string] string)
//	clients.muClientsSessions.Unlock()
//
//	clients.muIntervals.Lock()
//	clients.Intervals = make(map[string] map[string] map[string] sockjs.Session)
//	clients.muIntervals.Unlock()
//
//	clients.muIntervalSessions.Lock()
//	clients.IntervalSessions = make(map[string] map[string] string)
//	clients.muIntervalSessions.Unlock()
//
//	return clients
//
//}
//
//func (clients Clients) SetTopic(topic string) {
//	clients.muClients.Lock()
//	clients.Clients[topic] = make(map[string] map[string] sockjs.Session)
//	clients.muClients.Unlock()
//
//	clients.muClientsSessions.Lock()
//	clients.ClientSessions[topic] = make(map[string] string)
//	clients.muClientsSessions.Unlock()
//
//	clients.muIntervals.Lock()
//	clients.Intervals[topic] = make(map[string] map[string] sockjs.Session)
//	clients.muIntervals.Unlock()
//
//	clients.muIntervalSessions.Lock()
//	clients.IntervalSessions[topic] = make(map[string] string)
//	clients.muIntervalSessions.Unlock()
//}
//
//func (clients Clients) AddSubscriber(subscriber Subscriber,session sockjs.Session){
//	//CLients
//	clientSession := clients.GetListSessionByTopicAndUsername(subscriber.Topic, subscriber.Username)
//	if clientSession == nil{
//		clients.muClients.Lock()
//		clients.Clients[subscriber.Topic][subscriber.Username] = make(map[string] sockjs.Session)
//		clients.muClients.Unlock()
//	}
//	clients.muClients.Lock()
//	clients.Clients[subscriber.Topic][subscriber.Username][session.ID()] = session
//	clients.muClients.Unlock()
//
//	//ClientSession
//	clients.muClientsSessions.Lock()
//	clients.ClientSessions[subscriber.Topic][session.ID()] = subscriber.Username
//	clients.muClientsSessions.Unlock()
//
//	//Interval
//	intervalSession := clients.GetListSessionByTopicAndInterval(subscriber.Topic, subscriber.Interval)
//	if intervalSession == nil{
//		clients.muIntervals.Lock()
//		clients.Intervals[subscriber.Topic][subscriber.Interval] = make(map[string] sockjs.Session)
//		clients.muIntervals.Unlock()
//	}
//	clients.muIntervals.Lock()
//	clients.Intervals[subscriber.Topic][subscriber.Interval][session.ID()] = session
//	clients.muIntervals.Unlock()
//
//	//Interval Session
//	clients.muIntervalSessions.Lock()
//	clients.IntervalSessions[subscriber.Topic][session.ID()] = subscriber.Interval
//	clients.muIntervalSessions.Unlock()
//}
//
//func (clients Clients) RemoveSubscriber(topic string, sessionID string){
//	username := clients.GetUsernameByTopicAndSessionId(topic, sessionID)
//
//	clients.muClients.Lock()
//	delete(clients.Clients[topic][username], sessionID)
//	clients.muClients.Unlock()
//
//	clientSession := clients.GetListSessionByTopicAndUsername(topic, username)
//	if len(clientSession) == 0 {
//		clients.muClients.Lock()
//		delete(clients.Clients[topic], username)
//		clients.muClients.Unlock()
//	}
//	clients.muClientsSessions.Lock()
//	delete(clients.ClientSessions[topic], sessionID)
//	clients.muClientsSessions.Unlock()
//
//	interval := clients.GetIntervalByTopicAndSessionId(topic, sessionID)
//
//	clients.muIntervals.Lock()
//	delete(clients.Intervals[topic][interval], sessionID)
//	clients.muIntervals.Unlock()
//
//	intervalSession := clients.GetListSessionByTopicAndInterval(topic, interval)
//
//	if len(intervalSession) == 0 {
//		clients.muIntervals.Lock()
//		delete(clients.Intervals[topic], interval)
//		clients.muIntervals.Unlock()
//	}
//
//	clients.muIntervalSessions.Lock()
//	delete(clients.IntervalSessions[topic], sessionID)
//	clients.muIntervalSessions.Unlock()
//
//}
//
//func (clients Clients) GetUsernameByTopicAndSessionId(topic string, sessionId string) string{
//	var username string
//	clients.muClientsSessions.RLock()
//	username = clients.ClientSessions[topic][sessionId]
//	clients.muClientsSessions.RUnlock()
//	return username
//}
//
//func (clients Clients) GetIntervalByTopicAndSessionId(topic string, sessionId string) string{
//	var interval string
//	clients.muIntervalSessions.RLock()
//	interval = clients.IntervalSessions[topic][sessionId]
//	clients.muIntervalSessions.RUnlock()
//	return interval
//}
//
//func (clients Clients) GetAllClientsByTopic(topic string) map[string] map[string] sockjs.Session{
//	var c map[string] map[string] sockjs.Session
//	clients.muClients.RLock()
//	c = clients.Clients[topic]
//	clients.muClients.RUnlock()
//	return c
//}
//
//func (clients Clients) GetListSessionByTopicAndUsername(topic string, username string) map[string] sockjs.Session{
//	sessions := make(map[string] sockjs.Session)
//	clients.muClients.RLock()
//	sessions = clients.Clients[topic][username]
//	clients.muClients.RUnlock()
//	return sessions
//}
//
//func (clients Clients) GetListSessionByTopicAndInterval(topic string, interval string) map[string] sockjs.Session{
//	sessions := make(map[string] sockjs.Session)
//	clients.muIntervals.RLock()
//	sessions = clients.Intervals[topic][interval]
//	clients.muIntervals.RUnlock()
//	return sessions
//}

