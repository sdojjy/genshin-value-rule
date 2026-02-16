package eval

type AccountEvaluator interface {
	CalculateValuation(details Assets) ValuationResult
}

type Assets struct {
	Characters     map[string]int
	Weapons        map[string]int
	YuanShi        int
	JiuChanZhiYuan int
	YellowCount    int
}

type ValuationResult struct {
	FinalTotal float64 `json:"finalTotal"`
	Breakdown  string  `json:"breakdown"`
}
