package models

type Quiz struct {
	Title string `json:"title"`
	Code string `json:"code"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Type string `json:"type"`
	Description string `json:"description"`
	Answers []Answer `json:"answers"`
	CorrectAnswer string `json:"correct_answer"`
}

type Answer struct {
	Label string `json:"label"`
	Value string `json:"value"`
}