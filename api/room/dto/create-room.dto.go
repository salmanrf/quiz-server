package dto

type CreateRoomDto struct {
	Title string `json:"title"`
	QuizCode string `json:"quiz_code"`
	Password string `json:"password"`
	Quota int `json:"quota"`
	QuestionTimeout int `json:"question_timeout"`
	NextQuestionInterval int `json:"next_question_interval"`
}