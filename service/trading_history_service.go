package service

import (
	"database/sql"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
	"websocket-server/const/formtter"
	"websocket-server/daos/trading"
	"websocket-server/util/database"
)

type TradingHistory struct {
	Id				int64 					`json:"id"`
	AskId			sql.NullInt64 			`json:"ask_id"`
	BidId			sql.NullInt64 			`json:"bid_id"`
	Price			decimal.NullDecimal 	`json:"price"`
	Amount			decimal.NullDecimal 	`json:"amount"`
	Rate			sql.NullString 			`json:"rate"`
	CreatedAt		pq.NullTime 			`json:"created_at"`
	CreatedBy		sql.NullString  		`json:"created_by"`
	AppliedOnWallet	sql.NullBool			`json:"applied_on_wallet"`
	TradingType		sql.NullString			`json:"trading_type"`
}

func GetLastTradingHistoryBySize(listHistory *trading.ListHistory, rate string, size int64) error {
	var histories []trading.History

	q := database.Db.Select("id AS id", "amount AS amount", "price AS price", "trading_type AS type", "created_at AS time").
			From("trading_history").
			Where(dbx.Like("rate", rate)).
			OrderBy("id DESC").
			Limit(size)

	err := q.All(&histories)
	if err != nil {
		log.Error(err)
		return err
	}

	listHistory.Histories = histories

	return nil
}


func GetLast24HTransactionByRate(last24h *trading.Last24h, rate string) error {

	//var last24h trading.Last24h

	tLast24h := time.Now().AddDate(0,0,-1)
	tLast24hStr := tLast24h.Format(formtter.SQLTIMEFORMAT)

	q := database.Db.Select("MAX(price) AS high", "MIN(price) AS low", "SUM(amount) as volume").
		From("trading_history").
		Where(dbx.And(dbx.Like("rate", rate), dbx.NewExp("created_at > '"+tLast24hStr+"'")))

	//log.Println(q.Build().SQL())

	err := q.One(last24h)
	if err != nil {
		log.Error(err)
		return err
	}


	var tradingHistory []TradingHistory

	q2 := database.Db.NewQuery("SELECT * FROM trading_history " +
								"WHERE id=(SELECT min(id) FROM trading_history WHERE created_at > '"+tLast24hStr+"') " +
								"OR id=(SELECT max(id) FROM trading_history WHERE created_at > '"+tLast24hStr+"')" +
								"AND rate = '"+rate+"'")


	err = q2.All(&tradingHistory)

	if err != nil {
		log.Error(err)
		return err
	}


	var first, last decimal.Decimal

	if len(tradingHistory) > 1 {
		if (tradingHistory[0].Id < tradingHistory[1].Id){
			first = tradingHistory[0].Price.Decimal
			last = tradingHistory[1].Price.Decimal
		} else {
			first = tradingHistory[1].Price.Decimal
			last = tradingHistory[0].Price.Decimal
		}
	} else if len(tradingHistory) == 1 {
		first = tradingHistory[0].Price.Decimal
		last = tradingHistory[0].Price.Decimal
	} else {
		first = decimal.Zero
		last = decimal.Zero
	}


	(*last24h).LastPrice = last
	(*last24h).Change = (last.Sub(first)).Div(last)
	if (*last24h).Change.GreaterThanOrEqual(decimal.Zero){
		(*last24h).State = "increasing"
	} else {
		(*last24h).State = "decreasing"
		(*last24h).Change = (*last24h).Change.Neg()
	}

	//log.Println((*last24h).High)
	//log.Println((*last24h).Low)
	//log.Println((*last24h).Change)
	//log.Println((*last24h).LastPrice)
	//log.Println((*last24h).State)

	return nil

}

func GetOpenOrdersByUsernameAndRate(openOrders *trading.OpenOrders, username string, rate string) error{
	iq := "" +
		"(SELECT ta.id AS id, " +
		"'ASK' AS trading_type, " +
		"ta.rate AS rate, " +
		"ta.price AS price, " +
		"ta.amount AS amount, " +
		"ta.filled_amount AS filled_amount, " +
		"ta.created_at AS created_at, " +
		"ta.admin_fee AS admin_fee " +
		"FROM trading_ask ta " +
		"WHERE ta.member_id = (SELECT id from member where username = '"+username+"') " +
		"AND UPPER(ta.rate) = UPPER('"+rate+"') " +
		"AND ta.open = true " +
		"AND ta.type_id = 1) " +
		"UNION " +
		"(SELECT " +
		"tb.id AS id, " +
		"'BID' AS trading_type, " +
		"tb.rate AS rate, " +
		"tb.price AS price, " +
		"tb.amount AS amount, " +
		"tb.filled_amount AS filled_amount, " +
		"tb.created_at AS created_at, " +
		"tb.admin_fee AS admin_fee " +
		"FROM trading_bid tb " +
		"WHERE tb.member_id = (SELECT id from member where username = '"+username+"') " +
		"AND UPPER(tb.rate) = UPPER('"+rate+"') " +
		"AND tb.open = true " +
		"AND tb.type_id = 1) " +
		"ORDER BY created_at DESC "

	q := database.Db.NewQuery(iq)

	//log.Println(q.SQL())

	var listOpenOrder []trading.OpenOrder

	err := q.All(&listOpenOrder)
	if err != nil {
		log.Error(err)
		return err
	}

	openOrders.Payload = listOpenOrder

	return nil
}

func GetListOpenOrderByMemberIdAndRate(listOpenOrder *trading.ListOpenOrder, memberId uint64, rate string) error{
	iq := "" +
		"(SELECT ta.id AS id, " +
				"'ASK' AS trading_type, " +
				"ta.rate AS rate, " +
				"ta.price AS price, " +
				"ta.amount AS amount, " +
				"ta.filled_amount AS filled_amount, " +
				"ta.created_at AS created_at, " +
				"ta.admin_fee AS admin_fee " +
		"FROM trading_ask ta " +
		"WHERE ta.member_id = "+strconv.FormatUint(memberId, 10)+" " +
			"AND UPPER(ta.rate) = UPPER('"+rate+"') " +
			"AND ta.open = true " +
			"AND ta.type_id = 1) " +
		"UNION " +
		"(SELECT " +
				"tb.id AS id, " +
				"'BID' AS trading_type, " +
				"tb.rate AS rate, " +
				"tb.price AS price, " +
				"tb.amount AS amount, " +
				"tb.filled_amount AS filled_amount, " +
				"tb.created_at AS created_at, " +
				"tb.admin_fee AS admin_fee " +
		"FROM trading_bid tb " +
		"WHERE tb.member_id = "+strconv.FormatUint(memberId, 10)+" " +
			"AND UPPER(tb.rate) = UPPER('"+rate+"') " +
			"AND tb.open = true " +
			"AND tb.type_id = 1) " +
		"ORDER BY created_at DESC "

	q := database.Db.NewQuery(iq)

	var openOrders []trading.OpenOrder

	err := q.All(&openOrders)
	if err != nil {
		log.Error(err)
		return err
	}

	listOpenOrder.OpenOrders = openOrders

	return nil
}

func GetOrderHistoriesByUsernameAndRateAndOffsetAndLimit(orderHistories *trading.OrderHistories, username string, rate string, offset uint64, limit int) error {
	iq := "" +
		"(SELECT h.id AS id, " +
		"h.rate AS rate, " +
		"'ASK' AS trading_type, " +
		"h.price AS price, " +
		"h.amount AS amount, " +
		"a.admin_fee AS admin_fee, " +
		"h.created_at AS created_at " +
		"FROM trading_history  h " +
		"LEFT JOIN trading_ask a ON a.id = h.ask_id " +
		"WHERE a.member_id = (SELECT id from member where username = '"+username+"') AND UPPER(a.rate) = UPPER('"+rate+"'))" +
		"UNION" +
		"(SELECT h.id AS id, " +
		"h.rate AS rate, " +
		"'BID' AS trading_type, " +
		"h.price AS price, " +
		"h.amount AS amount, " +
		"b.admin_fee AS admin_fee, " +
		"h.created_at AS created_at " +
		"FROM trading_history h " +
		"LEFT JOIN trading_bid b ON b.id = h.bid_id " +
		"WHERE b.member_id = (SELECT id from member where username = '"+username+"') AND UPPER(b.rate) = UPPER('"+rate+"'))"


	var listOrderHistory []trading.OrderHistory

	q := database.Db.NewQuery("" +
		"SELECT outer_query.* " +
		"FROM (" +
		"SELECT " +
		"row_number() over() rn, " +
		"inner_query.* " +
		"FROM ("+iq+" ORDER BY created_at DESC) inner_query) " +
		"outer_query " +
		"WHERE outer_query.rn >= "+strconv.FormatUint(offset, 10)+" + 1 " +
		"AND outer_query.rn <= "+strconv.FormatUint(offset, 10)+" + "+strconv.Itoa(limit))

	err := q.All(&listOrderHistory)
	if err != nil {
		log.Error(err)
		return err
	}

	qCount := database.Db.NewQuery("SELECT COUNT(id) as size FROM  ("+iq+") as order_history")

	qCount.One(orderHistories)
	if err != nil {
		log.Error(err)
		return err
	}

	orderHistories.Payload = listOrderHistory

	return nil

}


func GetListOrderHistoryByMemberIdAndRateAndOffsetAndLimit(listOrderHistory *trading.ListOrderHistory, memberId uint64, rate string, offset uint64, limit int) error {
	iq := "" +
		"(SELECT h.id AS id, " +
		"h.rate AS rate, " +
		"'ASK' AS trading_type, " +
		"h.price AS price, " +
		"h.amount AS amount, " +
		"a.admin_fee AS admin_fee, " +
		"h.created_at AS created_at " +
		"FROM trading_history  h " +
		"LEFT JOIN trading_ask a ON a.id = h.ask_id " +
		"WHERE a.member_id = "+strconv.FormatUint(memberId, 10)+" AND UPPER(a.rate) = UPPER('"+rate+"'))" +
		"UNION" +
		"(SELECT h.id AS id, " +
		"h.rate AS rate, " +
		"'BID' AS trading_type, " +
		"h.price AS price, " +
		"h.amount AS amount, " +
		"b.admin_fee AS admin_fee, " +
		"h.created_at AS created_at " +
		"FROM trading_history h " +
		"LEFT JOIN trading_bid b ON b.id = h.bid_id " +
		"WHERE b.member_id = "+strconv.FormatUint(memberId, 10)+" AND UPPER(b.rate) = UPPER('"+rate+"'))"


	var orderHistories []trading.OrderHistory

	q := database.Db.NewQuery("" +
		"SELECT outer_query.* " +
		"FROM (" +
			"SELECT " +
			"row_number() over() rn, " +
			"inner_query.* " +
			"FROM ("+iq+" ORDER BY created_at DESC) inner_query) " +
		"outer_query " +
		"WHERE outer_query.rn >= "+strconv.FormatUint(offset, 10)+" + 1 " +
		"AND outer_query.rn <= "+strconv.FormatUint(offset, 10)+" + "+strconv.Itoa(limit))

	err := q.All(&orderHistories)
	if err != nil {
		log.Error(err)
		return err
	}

	qCount := database.Db.NewQuery("SELECT COUNT(id) as size FROM  ("+iq+") as order_history")

	qCount.One(listOrderHistory)
	if err != nil {
		log.Error(err)
		return err
	}

	listOrderHistory.OrderHistories = orderHistories

	return nil

}