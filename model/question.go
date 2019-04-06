package model

type Option string

type Question struct {
	QuestionText string                 `json:"question_text"`
	Options      map[Option]interface{} `json:"options"`
	Answer       Option                 `json:"-"`
}
