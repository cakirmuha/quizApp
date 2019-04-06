package database

import (
	"quizApp/model"
)

var questions = []model.Question{
	{
		QuestionText: "Which is greater than 5 ?",
		Options: map[model.Option]interface{}{
			"a": 1,
			"b": 7,
			"c": 3,
			"d": 4,
		},
		Answer: "b",
	},
	{
		QuestionText: "What is the capital of Malta ?",
		Options: map[model.Option]interface{}{
			"a": "Birgu",
			"b": "Sliema",
			"c": "Valletta",
			"d": "Cospicua",
		},
		Answer: "c",
	},
	{
		QuestionText: "What is the result of 5 / 2 ?",
		Options: map[model.Option]interface{}{
			"a": 2,
			"b": 3.5,
			"c": 0.4,
			"d": 2.5,
		},
		Answer: "d",
	},
}

func (db *DB) GetQuestions() []model.Question {
	return questions
}

func (db *DB) CalculateScoreAndDegree(username string, answers []model.Option) (int, float64) {
	userScore := 0
	for i, a := range answers {
		if a == questions[i].Answer {
			userScore++
		}
	}

	db.cache.userMu.RLock()
	lessScoreCount, userCount := 0, len(db.cache.userScore)
	for _, us := range db.cache.userScore {
		if userScore > us {
			lessScoreCount++
		}
	}
	db.cache.userMu.RUnlock()

	db.cache.userMu.Lock()
	db.cache.userScore[username] = userScore
	db.cache.userMu.Unlock()

	successRate := 100.0 // first user success rate: 100 (can be 0)
	if userCount > 0 {
		successRate = float64(lessScoreCount) / float64(userCount) * 100
	}

	return userScore, successRate
}
