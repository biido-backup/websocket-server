package daos

var TRDChart *TradingChart

type TradingChart struct {
	chartListMap map[string]map[string][]Chart
}

func CreateTradingChart() {
	chartListMap := make(map[string]map[string][]Chart)
	TRDChart = &TradingChart{chartListMap}
}

func InitRateTradingChart(rate string){
	TRDChart.chartListMap[rate] = make(map[string][]Chart)
}

func SetChartList(rate, unitOfTime string, chartList []Chart) {
	TRDChart.chartListMap[rate][unitOfTime] = chartList
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