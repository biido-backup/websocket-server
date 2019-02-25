package daos

import "strings"

type Rate struct{
	MainCurrency		string 		`json:"mainCurrency"`
	PivotCurrency		string 		`json:"pivotCurrency"`
}

func GetRateFromStringDash(rateStr string) Rate{
	split := strings.Split(rateStr,"-")
	rate := Rate{MainCurrency:split[0], PivotCurrency:split[1]}
	return rate
}

func GetRateFromStringSlash(rateStr string) Rate{
	split := strings.Split(rateStr,"/")
	rate := Rate{MainCurrency:split[0], PivotCurrency:split[1]}
	return rate
}

func (rate Rate) String() string {
	return rate.MainCurrency+"/"+rate.PivotCurrency
}

func (rate Rate) StringDash() string {
	return rate.MainCurrency+"-"+rate.PivotCurrency
}

func (rate Rate) StringSlah() string {
	return rate.MainCurrency+"/"+rate.PivotCurrency
}