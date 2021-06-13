package entity

type Email struct {
	Id            int    `json:"id"`
	Receiver      string `json:"receiver"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	SendingUser   string `json:"sendingUser"`
	SendingSystem string `json:"sendingSystem"`
}
