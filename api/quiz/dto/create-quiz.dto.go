package dto

import "github.com/salmanrf/svelte-go-quiz-server/api/quiz/models"


type CreateQuizDto struct {
	Title string `json:"title"`
	Questions []models.Question `json:"questions"`
}