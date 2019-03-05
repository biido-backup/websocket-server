package cache

import (
	"sync"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/service"
	"websocket-server/util/logger"
)

var log = logger.CreateLog("cache")

var muCacheOrderBook sync.RWMutex
var CacheOrderBook map[string] trading.Orderbook

var muCacheListHistory sync.RWMutex
var CacheListHistory map[string] trading.TradingListHistory

var muCacheLast24h sync.RWMutex
var CacheLast24h map[string] trading.TradingLast24h

func InitCache(){
	CacheOrderBook = make(map[string] trading.Orderbook)
	CacheListHistory = make(map[string] trading.TradingListHistory)
	CacheLast24h = make(map[string] trading.TradingLast24h)

}

func FillCacheByRate(rate daos.Rate){
	var listHistory trading.ListHistory
	err := service.GetLastTradingHistoryBySize(&listHistory, rate.StringSlah(), trdconst.SIZE_TRDHIST)
	if err != nil {
		log.Error("Error query to trading_history", err)
	}
	trdListHistory := trading.TradingListHistory{Type: trdconst.TRADINGHISTORY, Payload: listHistory}

	muCacheListHistory.Lock()
	CacheListHistory[rate.StringDash()] = trdListHistory
	muCacheListHistory.Unlock()

	log.Println(rate.StringSlah())

	var last24h trading.Last24h

	err = service.GetLast24HTransactionByRate(&last24h, rate.StringSlah())

	if err != nil{
		log.Fatal("Error query to last24h", err)
	}

	trdLast24h := trading.TradingLast24h{Type: trdconst.LAST24H, Payload: last24h}
	muCacheLast24h.Lock()
	CacheLast24h[rate.StringDash()] = trdLast24h
	muCacheLast24h.Unlock()

	//log.Debug(CacheOrderBook)
	//log.Debug(CacheListHistory)
	//log.Debug(CacheLast24h)

}

func GetCacheByTopicAndType(topic string, cacheType string) interface{}{
	if cacheType == trdconst.ORDERBOOK{
		var cacheOrderBook trading.Orderbook
		muCacheOrderBook.RLock()
		cacheOrderBook = CacheOrderBook[topic]
		muCacheOrderBook.RUnlock()
		return cacheOrderBook
	} else if cacheType == trdconst.TRADINGHISTORY {
		var cacheListHistory trading.TradingListHistory
		muCacheListHistory.RLock()
		cacheListHistory = CacheListHistory[topic]
		muCacheListHistory.RUnlock()
		return cacheListHistory
	} else if cacheType == trdconst.LAST24H {
		var cacheLast24h trading.TradingLast24h
		muCacheLast24h.RLock()
		cacheLast24h = CacheLast24h[topic]
		muCacheLast24h.RUnlock()
		return cacheLast24h
	} else {
		return nil
	}
}

func SetCacheByTopicAndType(topic string, cacheType string, cache interface{}){
	if cacheType == trdconst.ORDERBOOK{
		muCacheOrderBook.Lock()
		CacheOrderBook[topic] = cache.(trading.Orderbook)
		muCacheOrderBook.Unlock()
	} else if cacheType == trdconst.TRADINGHISTORY {
		muCacheListHistory.Lock()
		CacheListHistory[topic] = cache.(trading.TradingListHistory)
		muCacheListHistory.Unlock()
	} else if cacheType == trdconst.LAST24H {
		muCacheLast24h.Lock()
		CacheLast24h[topic] = cache.(trading.TradingLast24h)
		muCacheLast24h.Unlock()
	}
}

