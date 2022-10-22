package wshandler

type JoinMessage struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Role int `json:"role"`
	Room Room `json:"room"`
}