package quiz

var quizzes = make(map[string]*Quiz)

type Quiz struct {
	Title string `json:"title"`
	Code string `json:"code"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Type string `json:"type"`
	Description string `json:"description"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	Value string `json:"value"`
	Correct bool `json:"type"`
}

func Create(input Quiz) bool {
	if _, exists := quizzes[input.Code]; exists {
		return !exists
	}

	newQuiz := &Quiz{Title: input.Title, Code: input.Code, Questions: input.Questions}

	quizzes[newQuiz.Code] = newQuiz

	return true
}

func Find() []Quiz {
	quizzesLength := len(quizzes)
	result := make([]Quiz, quizzesLength)

	index := 0;
	for _, q := range quizzes {
		result[index] = *q
		index++
	}
	
	return result;
}

func Get(code string) (*Quiz, bool) {
	q, e := quizzes[code]
	
	return q, e;
}
