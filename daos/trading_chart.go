package daos

var TRDChart *TradingChart

type TradingChart struct {
	chartListMap map[string]map[string][]Chart
}

func CreateTradingChart() {
	chartListMap := make(map[string]map[string][]Chart)
	TRDChart = &TradingChart{chartListMap}
}

func SetChartList(rate, unitOfTime string, chartList []Chart) {
	unitTimeChatListMap := make(map[string][]Chart)
	unitTimeChatListMap[unitOfTime] = chartList
	TRDChart.chartListMap[rate] = unitTimeChatListMap
}

func InsertChart(rate, unitOfTime string, chart Chart, quantity int64) {
	if int64(len(TRDChart.chartListMap[rate][unitOfTime])) >= quantity {
		TRDChart.chartListMap[rate][unitOfTime] = TRDChart.chartListMap[rate][unitOfTime][1:len(TRDChart.chartListMap[rate][unitOfTime])]
	}
	TRDChart.chartListMap[rate][unitOfTime] = append(TRDChart.chartListMap[rate][unitOfTime], chart)
}

func GetChartList(rate, unitOfTime string) (chartList []Chart) {
	return TRDChart.chartListMap[rate][unitOfTime]
}