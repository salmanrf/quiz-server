package dto

type JoinRoomDto struct {
	RoomCode string `json:"room_code"`
	Username string `json:"username"`
	Password string `json:"password"`
}