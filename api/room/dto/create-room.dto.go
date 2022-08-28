package dto

type CreateRoomDto struct {
	QuizCode string `json:"quiz_code"`
	AdminName string `json:"admin_name"`
	Password string `json:"password"`
	NextQuestionInterval int `json:"next_question_interval"`
}