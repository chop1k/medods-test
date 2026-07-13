package models

type ConnectivityCheck struct {
	IP     string `json:"ip"`
	Scheme string `json:"scheme"`
}

type Calendar struct {
	From string                   `json:"from"`
	To   string                   `json:"to"`
	Days map[string][]CalendarDay `json:"days"`
}

type CalendarDay struct {
	ID *int `json:"id" binding:"omitempty"`
	TaskBody
}
