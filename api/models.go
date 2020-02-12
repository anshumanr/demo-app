package api

type cbResponse struct {
	UUID       string `json:"uuid"`
	Accountid  string `json:"accountid"`
	From       string `json:"from"`
	To         string `json:"to"`
	CallStatus string `json:"callstatus"`
}
