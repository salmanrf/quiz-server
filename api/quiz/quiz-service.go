package quiz

import (
	"github.com/aidarkhanov/nanoid"
	"github.com/salmanrf/svelte-go-quiz-server/api/quiz/dto"
	"github.com/salmanrf/svelte-go-quiz-server/api/quiz/models"
)

var quizzes = make(map[string]*models.Quiz)

func Create(input dto.CreateQuizDto) (models.Quiz, error) {
	alphabet := nanoid.DefaultAlphabet

	code, _ := nanoid.Generate(alphabet, 6)
	
	for _, exists := quizzes[code]; exists; {
		code, _ = nanoid.Generate(alphabet, 6)
	}

	newQuiz := &models.Quiz{
		Title: input.Title, 
		Code: code, 
		Questions: input.Questions,
	}

	quizzes[newQuiz.Code] = newQuiz

	return *newQuiz, nil
}

func Find() []models.Quiz {
	quizzesLength := len(quizzes)
	result := make([]models.Quiz, quizzesLength)

	index := 0;
	for _, q := range quizzes {
		result[index] = *q
		index++
	}
	
	return result;
}

func Get(code string) (*models.Quiz, bool) {
	q, e := quizzes[code]
	
	return q, e;
}
