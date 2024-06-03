package models

// Определяем структуру для фактов
type Fact struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMoId     int    `json:"indicator_to_mo_id"`
	IndicatorToMoFactId int    `json:"indicator_to_mo_fact_id"`
	Value               int    `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              int    `json:"is_plan"`
	AuthUserId          int    `json:"auth_user_id"`
	Comment             string `json:"comment"`
}
