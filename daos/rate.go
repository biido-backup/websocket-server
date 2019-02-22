package daos

type Rate struct{
	MainCurrency		string 		`json:"mainCurrency"`
	PivotCurrency		string 		`json:"pivotCurrency"`
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