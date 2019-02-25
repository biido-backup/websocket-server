package daos

import "time"

var TRDChart *TradingChart

type TradingChart struct {
	chartListMap map[string]map[string][]Chart
	nextChartMap map[string]map[string][]Chart
}

func CreateTradingChart() {
	chartListMap := make(map[string]map[string][]Chart)
	lastChartMap := make(map[string]map[string][]Chart)
	TRDChart = &TradingChart{chartListMap, lastChartMap}
}

func InitRateTradingChart(rate string){
	TRDChart.chartListMap[rate] = make(map[string][]Chart)
	TRDChart.nextChartMap[rate] = make(map[string][]Chart)
}

func InitUnitOfTimeTradingChart(rate, unitOfTime string) {
	TRDChart.chartListMap[rate][unitOfTime] = make([]Chart, 0, 1)
	TRDChart.nextChartMap[rate][unitOfTime] = make([]Chart, 0, 1)
}

func SetChartList(rate, unitOfTime string, chartList []Chart) {
	TRDChart.chartListMap[rate][unitOfTime] = chartList
}

func SetNextChart(rate, unitOfTime string, lastTime int64) {
	lastTimeUnix := time.Unix(lastTime/1000, 0)

	var nextTime int64 = 0
	if unitOfTime == "5M" {
		nextTime = lastTimeUnix.Add(5 * time.Minute).Unix()
	} else if unitOfTime == "15M" {
		nextTime = lastTimeUnix.Add(15 * time.Minute).Unix()
	} else if unitOfTime == "30M" {
		nextTime = lastTimeUnix.Add(30 * time.Minute).Unix()
	} else if unitOfTime == "1H" {
		nextTime = lastTimeUnix.Add(1 * time.Hour).Unix()
	} else if unitOfTime == "6H" {
		nextTime = lastTimeUnix.Add(6 * time.Hour).Unix()
	} else if unitOfTime == "12H" {
		nextTime = lastTimeUnix.Add(12 * time.Hour).Unix()
	} else if unitOfTime == "1D" {
		nextTime = lastTimeUnix.AddDate(0,0,1).Unix()
	} else if unitOfTime == "1W" {
		nextTime = lastTimeUnix.AddDate(0,0,7).Unix()
	} else if unitOfTime == "1MO" {
		nextTime = lastTimeUnix.AddDate(0,1,0).Unix()
	}
	nextTime = nextTime * 1000

	chart := Chart{0,0,0,0,0,0,nextTime}
	TRDChart.nextChartMap[rate][unitOfTime] = []Chart{chart}
}

func CalculateNextChart(rate string, chart Chart) {
	unitOfTimeList := []string{"5M", "15M", "30M", "1H", "6H", "12H", "1D", "1W", "1MO"}
	for _, unitOfTime := range unitOfTimeList {
		nextTime := TRDChart.nextChartMap[rate][unitOfTime][0].Time

		var nextEndTime int64 = 0
		if unitOfTime == "5M" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(5 * time.Minute).Unix()
		} else if unitOfTime == "15M" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(15 * time.Minute).Unix()
		} else if unitOfTime == "30M" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(30 * time.Minute).Unix()
		} else if unitOfTime == "1H" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(1 * time.Hour).Unix()
		} else if unitOfTime == "6H" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(6 * time.Hour).Unix()
		} else if unitOfTime == "12H" {
			nextEndTime = time.Unix(nextTime/1000,0).Add(12 * time.Hour).Unix()
		} else if unitOfTime == "1D" {
			nextEndTime = time.Unix(nextTime/1000,0).AddDate(0,0,1).Unix()
		} else if unitOfTime == "1W" {
			nextEndTime = time.Unix(nextTime/1000,0).AddDate(0,0,7).Unix()
		} else if unitOfTime == "1MO" {
			nextEndTime = time.Unix(nextTime/1000,0).AddDate(0,1,0).Unix()
		}
		nextEndTime = nextEndTime * 1000

		if chart.Time >= nextTime {
			if chart.Time < nextEndTime {
				if TRDChart.nextChartMap[rate][unitOfTime][0].High == 0 || chart.High > TRDChart.nextChartMap[rate][unitOfTime][0].High {

				}
				if TRDChart.nextChartMap[rate][unitOfTime][0].Low == 0 || chart.Low < TRDChart.nextChartMap[rate][unitOfTime][0].Low {

				}
				if TRDChart.nextChartMap[rate][unitOfTime][0].Open == 0 {

				}
				TRDChart.nextChartMap[rate][unitOfTime][0].Close = chart.Close
				if TRDChart.nextChartMap[rate][unitOfTime][0].Volume > 0 || chart.Volume > 0{
					TRDChart.nextChartMap[rate][unitOfTime][0].Mean =
						((TRDChart.nextChartMap[rate][unitOfTime][0].Mean * TRDChart.nextChartMap[rate][unitOfTime][0].Volume) +
							(chart.Mean * chart.Volume)) / (TRDChart.nextChartMap[rate][unitOfTime][0].Volume + chart.Volume)
				}
				TRDChart.nextChartMap[rate][unitOfTime][0].Volume = TRDChart.nextChartMap[rate][unitOfTime][0].Volume + chart.Volume
			} else {
				chart.Time = nextEndTime
				TRDChart.nextChartMap[rate][unitOfTime][0] = chart
			}
		}
	}
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

func GetNextChart(rate, unitOfTime string) (chartList []Chart) {
	return TRDChart.nextChartMap[rate][unitOfTime]
}