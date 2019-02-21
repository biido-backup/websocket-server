package daos

var Charts *TradingCharts

type TradingCharts struct {
	ListMap map[string]map[string][]string
	OffsetMap map[string]map[string]int64
}

func CreateTradingCharts() {
	listMap := make(map[string]map[string][]string)
	offsetMap := make(map[string]map[string]int64)
	Charts = &TradingCharts{ListMap: listMap, OffsetMap: offsetMap}
}

func GetTradingChartListMap() map[string]map[string][]string {
	return Charts.ListMap
}

func GetTradingChartOffsetMap() map[string]map[string]int64 {
	return Charts.OffsetMap
}