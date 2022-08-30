package wshandler

type JoinMessage struct {
	Username string `json:"username"`
	RoomCode string `json:"room_code"`	
	Role int `json:"role"`
	Room Room `json:"room"`
}