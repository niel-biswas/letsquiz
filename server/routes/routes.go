package routes

import (
	"letsquiz/server/controllers"
)

func RegisterRoutes(router *Router) {
	router.Handle("GET", "/users", controllers.GetUsers)
	router.Handle("POST", "/users", controllers.CreateUser)
	router.Handle("GET", "/users/{id}", controllers.GetUserByID)
	router.Handle("GET", "/users/byname/{name}", controllers.GetUserIDByName)
	router.Handle("PUT", "/users/{id}", controllers.UpdateUser)

	router.Handle("GET", "/categories", controllers.GetCategories)
	router.Handle("POST", "/categories", controllers.CreateCategory)
	router.Handle("GET", "/categories/{id}", controllers.GetCategoryByID)
	router.Handle("GET", "/categories/byname/{name}", controllers.GetCategoryIDByName)
	router.Handle("PUT", "/categories/{id}", controllers.UpdateCategory)

	router.Handle("GET", "/quizzes", controllers.GetQuizzes)
	router.Handle("POST", "/quizzes", controllers.CreateQuiz)
	router.Handle("GET", "/quizzes/{id}", controllers.GetQuizByID)
	router.Handle("PUT", "/quizzes/{id}", controllers.UpdateQuiz)
	router.Handle("GET", "/quizzes/{id}/questions", controllers.GetQuestionsByQuizID) // Added route to fetch questions by quiz ID

	router.Handle("GET", "/questions", controllers.GetQuestions)
	router.Handle("POST", "/questions", controllers.CreateQuestion)
	router.Handle("GET", "/questions/{id}", controllers.GetQuestionByID)
	router.Handle("PUT", "/questions/{id}", controllers.UpdateQuestion)
	router.Handle("GET", "/questions/{id}/answers", controllers.GetAnswersByQuestionID) // Added route to fetch answers by question ID

	router.Handle("GET", "/answers", controllers.GetAnswers)
	router.Handle("POST", "/answers", controllers.CreateAnswer)
	router.Handle("GET", "/answers/{id}", controllers.GetAnswerByID)
	router.Handle("PUT", "/answers/{id}", controllers.UpdateAnswer)

	router.Handle("GET", "/attempts", controllers.GetUserQuizAttempts)
	router.Handle("POST", "/attempts", controllers.CreateUserQuizAttempt)
	router.Handle("GET", "/attempts/{id}", controllers.GetUserQuizAttemptByID)
	router.Handle("PUT", "/attempts/{id}", controllers.UpdateUserQuizAttempt)

	router.Handle("GET", "/user-answers", controllers.GetUserAnswers)
	router.Handle("POST", "/user-answers", controllers.CreateUserAnswer)
	router.Handle("GET", "/user-answers/{id}", controllers.GetUserAnswerByID)
	router.Handle("PUT", "/user-answers/{id}", controllers.UpdateUserAnswer)

	router.Handle("GET", "/leaderboards", controllers.GetLeaderboards)
	router.Handle("POST", "/leaderboards", controllers.CreateLeaderboard)
	router.Handle("GET", "/leaderboards/{id}", controllers.GetLeaderboardByID)
	router.Handle("PUT", "/leaderboards/{id}", controllers.UpdateLeaderboard)

	router.Handle("GET", "/feedbacks", controllers.GetFeedbacks)
	router.Handle("POST", "/feedbacks", controllers.CreateFeedback)
	router.Handle("GET", "/feedbacks/{id}", controllers.GetFeedbackByID)
	router.Handle("PUT", "/feedbacks/{id}", controllers.UpdateFeedback)
}
