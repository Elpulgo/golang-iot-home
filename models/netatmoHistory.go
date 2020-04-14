package models

type NetatmoHistory struct {
	Steps []Step `json:"body"`
	Name  string
}

type Step struct {
	Start    int         `json:"beg_time"`
	Duration int         `json:"step_time"`
	Values   [][]float32 `json:"value"`
}
