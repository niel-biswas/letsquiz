package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"letsquiz/common"
	"letsquiz/config"
	"letsquiz/logger"
	"net/http"
)

type Answer struct {
	ID               int       `json:"id"`
	QuestionId       int       `json:"question_id"`
	Text             string    `json:"text"`
	IsCorrect        bool      `json:"is_correct"`
	CreationDate     time.Time `json:"creation_date"`
	LastModifiedDate time.Time `json:"last_modified_date"`
}

type Question struct {
	ID                  int       `json:"id"`
	QuizId              int       `json:"quiz_id"`
	Text                string    `json:"text"`
	Type                string    `json:"type"` // "single" or "multiple"
	Points              int       `json:"points"`
	MultiChoiceAnsLimit int       `json:"multi_choice_ans_limit"`
	HintExplanation     string    `json:"hint_explanation"`
	DifficultyLevel     string    `json:"difficulty_level"`
	CreationDate        time.Time `json:"creation_date"`
	LastModifiedDate    time.Time `json:"last_modified_date"`
}

type questionForms struct {
	*huh.Form
	Question      Question `json:"question"`
	Answers       []Answer `json:"answers"`
	CorrectAnswer string   `json:"correct_answer"`
}

type DynamicQuizModel struct {
	common.Model
	QuizID           int
	CurrentFormGroup int
	QuestionForms    []questionForms
	TotalFormGroups  int
}

var answersCommaSeparated, questionScorePoint, multiChoiceAnsLimit []string

// Initialize forms for each question
func InitialDynamicQuizModel(quizID, questionCount int, focused string, buttons []common.Button) DynamicQuizModel {
	logger.Info("InitialDynamicQuizModel called", "quizID", quizID, "questionCount", questionCount)
	m := DynamicQuizModel{
		QuizID:           quizID,
		CurrentFormGroup: 0,
		QuestionForms:    make([]questionForms, questionCount),
		TotalFormGroups:  questionCount,
		Model:            common.Model{CurrentScreen: "QuizMetadata", Focused: focused, Buttons: buttons},
	}
	answersCommaSeparated, questionScorePoint, multiChoiceAnsLimit = make([]string, questionCount), make([]string, questionCount), make([]string, questionCount)
	m.initForm()
	logger.Info("Initialized DynamicQuizModel", "quizID", quizID, "totalSteps", m.TotalFormGroups)
	return m
}

// Initialize forms and fetch existing data if available
func (m *DynamicQuizModel) initForm() {
	logger.Info("initForm called")

	// Fetch existing questions and answers from the backend
	existingQuestionForms, err := m.fetchQuestionsAndAnswers()
	if err != nil {
		logger.Error("Failed to fetch existing questions and answers", "error", err)
	}

	for i := 0; i < m.TotalFormGroups; i++ {
		if m.QuestionForms[i].Form == nil {
			m.QuestionForms[i] = questionForms{
				Form: huh.NewForm(),
			}
		}

		q := &m.QuestionForms[i]

		if i < len(existingQuestionForms) {
			*q = existingQuestionForms[i]

			var answers []string
			for _, ans := range q.Answers {
				answers = append(answers, ans.Text)
			}
			answersCommaSeparated[i] = strings.Join(answers, ",")
			questionScorePoint[i] = strconv.Itoa(q.Question.Points)
			multiChoiceAnsLimit[i] = strconv.Itoa(q.Question.MultiChoiceAnsLimit)
		}

		group := huh.NewGroup(
			huh.NewInput().Key("text").
				Title(fmt.Sprintf("Enter text for Question %d", i+1)).
				Placeholder("Question text").
				Value(&q.Question.Text),
			huh.NewInput().Key("hint_explanation").
				Title("Enter question hints for help").
				Placeholder("Hints here").
				Value(&q.Question.HintExplanation),
			huh.NewSelect[string]().Key("type").
				Title("Type of Question").
				Options(
					huh.NewOption("Single Choice", "single"),
					huh.NewOption("Multiple Choice", "multiple"),
				).
				Value(&q.Question.Type),
			huh.NewInput().Key("multi_choice_ans_limit").
				Title("Enter multiple choice answer limit for the question").
				Placeholder("Enter multi-choice answer limit here").
				Validate(func(val string) error {
					if q.Question.Type == "single" && val != "" {
						return errors.New("This is not applicable to single choice question, please leave it blank.")
					} else if q.Question.Type == "multiple" && val == "" {
						return errors.New("This is a required field for multiple choice question.")
					}
					return nil
				}).
				Value(&multiChoiceAnsLimit[i]),
			huh.NewInput().Key("points").
				Title("Enter the score for correct answer").
				Placeholder("Score points here").
				Value(&questionScorePoint[i]),
			huh.NewSelect[string]().Key("difficulty_level").
				Title("Difficulty Level").
				Options(
					huh.NewOption("Easy", "easy"),
					huh.NewOption("Medium", "medium"),
					huh.NewOption("Hard", "hard"),
				).
				Value(&q.Question.DifficultyLevel),
			huh.NewInput().Key("answers").
				Title("Enter your answers").
				Placeholder("Comma separated values within \"\" here").
				Value(&answersCommaSeparated[i]),
			huh.NewInput().Key("correct_answer").
				Title("Enter the correct answer").
				Placeholder("Correct answer (match case) here").
				Value(&q.CorrectAnswer),
		)

		form := huh.NewForm(group)
		m.QuestionForms[i].Form = form
	}

	logger.Info("Initialized form with groups for all questions")
}

// Check if answers exist for the current question
func (m *DynamicQuizModel) checkIfAnswersExist() (bool, []Answer, error) {
	url := fmt.Sprintf("%s/questions/%d/answers", config.AppConfig.BackendURL, m.QuestionForms[m.CurrentFormGroup].Question.ID)
	logger.Info("Checking if answers exist", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to check if answers exist", "url", url, "error", err)
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var answers []Answer
		err = json.NewDecoder(resp.Body).Decode(&answers)
		if err != nil {
			logger.Error("Failed to decode answers response", "error", err)
			return false, nil, err
		}
		logger.Info("Answers exist", "answers", answers)
		return true, answers, nil
	}
	logger.Info("Answers do not exist")
	return false, nil, nil
}

// Check if quiz exists and fetch its questions
func (m *DynamicQuizModel) checkIfQuizExists() (bool, []Question, error) {
	url := fmt.Sprintf("%s/quizzes/%d/questions", config.AppConfig.BackendURL, m.QuizID)
	logger.Info("Checking if quiz exists", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Failed to check if quiz exists", "url", url, "error", err)
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var questions []Question
		err = json.NewDecoder(resp.Body).Decode(&questions)
		if err != nil {
			logger.Error("Failed to decode questions response", "error", err)
			return false, nil, err
		}
		logger.Info("Quiz exists", "questions", questions)
		return true, questions, nil
	}
	logger.Info("Quiz does not exist")
	return false, nil, nil
}

// Fetch existing questions and answers from the backend
func (m *DynamicQuizModel) fetchQuestionsAndAnswers() ([]questionForms, error) {
	quizExists, questions, err := m.checkIfQuizExists()
	if err != nil {
		return nil, err
	}
	if !quizExists {
		return nil, nil
	}

	var questionFormsList []questionForms
	for _, question := range questions {
		answersExist, answers, err := m.checkIfAnswersExist()
		if err != nil {
			return nil, err
		}

		qForm := questionForms{
			Question: question,
		}
		if answersExist {
			qForm.Answers = answers
			// Find the correct answer
			for _, ans := range answers {
				if ans.IsCorrect {
					qForm.CorrectAnswer = ans.Text
					break
				}
			}
		}

		questionFormsList = append(questionFormsList, qForm)
	}

	logger.Info("Fetched existing questions and answers", "questionFormsList", questionFormsList)
	return questionFormsList, nil
}

// Save responses to the backend
func (m *DynamicQuizModel) saveResponsesToBackend() {
	logger.Info("Saving responses to backend", "responsesCount", len(m.QuestionForms))
	quizExists, _, err := m.checkIfQuizExists()
	if err != nil {
		logger.Error("Error checking if quiz exists", "error", err)
	}

	if !quizExists {
		m.postResponseToBackend()
		return
	}

	for i, q := range m.QuestionForms {
		logger.Info("Processing form data", "formIndex", i, "formData", q)
		currentDate := time.Now().UTC()
		q.Question.CreationDate = currentDate
		q.Question.LastModifiedDate = currentDate

		// Create question response data
		questionResponse := map[string]interface{}{
			"quiz_id":                m.QuizID,
			"id":                     q.Question.ID,
			"text":                   q.Question.Text,
			"type":                   q.Question.Type,
			"points":                 q.Question.Points,
			"hint_explanation":       q.Question.HintExplanation,
			"difficulty_level":       q.Question.DifficultyLevel,
			"multi_choice_ans_limit": q.Question.MultiChoiceAnsLimit,
			"creation_date":          q.Question.CreationDate.Format(time.RFC3339),
			"last_modified_date":     q.Question.LastModifiedDate.Format(time.RFC3339),
		}
		logger.Info("Sending question to backend (PUT)", "questionID", q.Question.ID, "questionResponse", questionResponse)
		questionID, err := m.postQuestionToBackend(questionResponse)
		if err != nil {
			logger.Error("Failed to send question to backend", "error", err)
			continue
		}

		// Update the question ID in the model
		m.QuestionForms[i].Question.ID = questionID

		// Clear the Answers list to avoid duplication
		q.Answers = nil

		// Process and append answers
		answerTexts := strings.Split(answersCommaSeparated[i], ",")
		for _, answerText := range answerTexts {
			if strings.TrimSpace(answerText) == q.CorrectAnswer {
				q.Answers = append(q.Answers, Answer{QuestionId: questionID, Text: strings.TrimSpace(answerText), IsCorrect: true})
			} else {
				q.Answers = append(q.Answers, Answer{QuestionId: questionID, Text: strings.TrimSpace(answerText), IsCorrect: false})
			}
		}

		logger.Info("Processed answers for question", "questionID", questionID, "answers", q.Answers)

		// Send answers to the backend
		for _, ans := range q.Answers {
			ans.CreationDate = currentDate
			ans.LastModifiedDate = currentDate
			ans.QuestionId = questionID

			answerResponse := map[string]interface{}{
				"question_id":        questionID,
				"text":               ans.Text,
				"is_correct":         ans.IsCorrect,
				"creation_date":      ans.CreationDate.Format(time.RFC3339),
				"last_modified_date": ans.LastModifiedDate.Format(time.RFC3339),
			}
			logger.Info("Sending answer to backend", "answer", ans, "answerResponse", answerResponse)
			err := m.postAnswerToBackend(answerResponse)
			if err != nil {
				logger.Error("Failed to send answer to backend", "error", err)
			}
		}
	}
	logger.Info("All questions and answers saved to backend")
}

// Post a question to the backend and return its ID
func (m *DynamicQuizModel) postQuestionToBackend(questionData map[string]interface{}) (int, error) {
	url := config.AppConfig.BackendURL + "/questions"
	jsonData, err := json.Marshal(questionData)
	if err != nil {
		logger.Error("Failed to marshal question data", "questionData", questionData, "error", err)
		return 0, err
	}
	logger.Info("POSTing question data", "url", url, "data", string(jsonData))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to send question to backend", "url", url, "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Handle non-successful status codes
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		logger.Error("Failed to save question, invalid status code", "statusCode", resp.StatusCode, "body", bodyString)
		return 0, fmt.Errorf("failed to save question, status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	// Parse the response to extract the question ID
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err == io.EOF {
		// Handle the case where the response body is empty
		logger.Error("Failed to decode response, response body is empty", "error", err)
		return 0, fmt.Errorf("failed to decode response, response body is empty")
	} else if err != nil {
		logger.Error("Failed to decode response", "error", err)
		return 0, err
	}

	// Extract and return the question ID
	id, ok := response["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response format: %v", response)
	}

	logger.Info("Successfully posted question", "questionID", int(id))
	return int(id), nil
}

// Post an answer to the backend
func (m *DynamicQuizModel) postAnswerToBackend(answerData map[string]interface{}) error {
	url := config.AppConfig.BackendURL + "/answers"
	jsonData, err := json.Marshal(answerData)
	if err != nil {
		logger.Error("Failed to marshal answer data", "answerData", answerData, "error", err)
		return err
	}

	logger.Info("POSTing answer data", "url", url, "data", string(jsonData))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to send answer to backend", "url", url, "error", err)
		return err
	}
	defer resp.Body.Close()

	// Handle non-successful status codes
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		logger.Error("Failed to save answer, invalid status code", "statusCode", resp.StatusCode, "body", bodyString)
		return fmt.Errorf("failed to save answer, status code: %d, body: %s", resp.StatusCode, bodyString)
	}
	logger.Info("Successfully posted answer")
	return nil
}

// UpdateDynamicQuizModel handles the updates and state transitions
func UpdateDynamicQuizModel(m DynamicQuizModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		logger.Info("UpdateDynamicQuizModel called", "CurrentFormGroup", m.CurrentFormGroup, "msgType", fmt.Sprintf("%T", msg), "Key pressed", msg.String())
		switch msg.String() {
		case "ctrl+right":
			logger.Info("Handling next form step", "CurrentFormGroup", m.CurrentFormGroup)
			m.saveCurrentFormData()
			m.NextForm()
			logger.Info("Model updated after NextForm", "CurrentFormGroup", m.CurrentFormGroup)
			// Check if we have reached the end of forms and the form state is completed
			if m.CurrentFormGroup == m.TotalFormGroups-1 && m.QuestionForms[m.CurrentFormGroup].Form.State == huh.StateCompleted {
				logger.Info("End of forms reached and form state is completed, quitting application")
				return m, tea.Quit
			}
		case "enter":
			// Check if we have reached the end of forms and the form state is completed and enter is hit
			if m.CurrentFormGroup == m.TotalFormGroups-1 && m.QuestionForms[m.CurrentFormGroup].Form.State == huh.StateCompleted {
				logger.Info("End of forms reached and form state is completed, quitting application")
				return m, tea.Quit
			}
		case "ctrl+left":
			logger.Info("Handling previous form step", "CurrentFormGroup", m.CurrentFormGroup)
			m.saveCurrentFormData()
			m.PreviousForm()
			logger.Info("Model updated after PreviousForm", "CurrentFormGroup", m.CurrentFormGroup)
		case "esc":
			logger.Info("Handling esc key press", "CurrentFormGroup", m.CurrentFormGroup)
			return m, tea.Quit
		}
	}

	// Update the current form based on the msg
	if m.CurrentFormGroup < len(m.QuestionForms) {
		form := m.QuestionForms[m.CurrentFormGroup].Form
		newFormModel, cmd := form.Update(msg)
		if newForm, ok := newFormModel.(*huh.Form); ok {
			m.QuestionForms[m.CurrentFormGroup].Form = newForm
			logger.Info("Form updated", "currentForm", m.CurrentFormGroup)
		}
		return m, cmd
	}

	logger.Info("Model update completed", "CurrentFormGroup", m.CurrentFormGroup)
	return m, nil
}

// Save the current form data to the model
func (m *DynamicQuizModel) saveCurrentFormData() {
	logger.Info("Saving current form data", "CurrentFormGroup", m.CurrentFormGroup)
	q := &m.QuestionForms[m.CurrentFormGroup].Question
	multiChoiceAnsLimit[m.CurrentFormGroup] = m.QuestionForms[m.CurrentFormGroup].Form.GetString("multi_choice_ans_limit")
	q.MultiChoiceAnsLimit, _ = strconv.Atoi(multiChoiceAnsLimit[m.CurrentFormGroup])
	questionScorePoint[m.CurrentFormGroup] = m.QuestionForms[m.CurrentFormGroup].Form.GetString("points")
	q.Points, _ = strconv.Atoi(questionScorePoint[m.CurrentFormGroup])
	answersCommaSeparated[m.CurrentFormGroup] = m.QuestionForms[m.CurrentFormGroup].Form.GetString("answers")
	answerTexts := strings.Split(answersCommaSeparated[m.CurrentFormGroup], ",")
	m.QuestionForms[m.CurrentFormGroup].Answers = nil // Clear the Answers list before appending new ones
	for _, answerText := range answerTexts {
		if strings.TrimSpace(answerText) == m.QuestionForms[m.CurrentFormGroup].CorrectAnswer {
			m.QuestionForms[m.CurrentFormGroup].Answers = append(m.QuestionForms[m.CurrentFormGroup].Answers, Answer{QuestionId: m.QuestionForms[m.CurrentFormGroup].Question.ID, Text: strings.TrimSpace(answerText), IsCorrect: true})
		} else {
			m.QuestionForms[m.CurrentFormGroup].Answers = append(m.QuestionForms[m.CurrentFormGroup].Answers, Answer{QuestionId: m.QuestionForms[m.CurrentFormGroup].Question.ID, Text: strings.TrimSpace(answerText), IsCorrect: false})
		}
	}
	// Log all fields of the current question form data
	logger.Info("Saved form data", "formIndex", m.CurrentFormGroup, "formData", map[string]interface{}{
		"QuestionID":            q.ID,
		"QuizID":                q.QuizId,
		"Text":                  q.Text,
		"Type":                  q.Type,
		"Points":                q.Points,
		"MultiChoiceAnsLimit":   q.MultiChoiceAnsLimit,
		"HintExplanation":       q.HintExplanation,
		"DifficultyLevel":       q.DifficultyLevel,
		"CreationDate":          q.CreationDate,
		"LastModifiedDate":      q.LastModifiedDate,
		"Answers":               m.QuestionForms[m.CurrentFormGroup].Answers,
		"CorrectAnswer":         m.QuestionForms[m.CurrentFormGroup].CorrectAnswer,
		"answersCommaSeparated": answersCommaSeparated[m.CurrentFormGroup],
		"multiChoiceAnsLimit":   multiChoiceAnsLimit[m.CurrentFormGroup],
		"questionScorePoint":    questionScorePoint[m.CurrentFormGroup],
	})
}

// Move to the next form in the sequence
func (m *DynamicQuizModel) NextForm() {
	logger.Info("NextForm called", "CurrentFormGroup", m.CurrentFormGroup, "Focused", m.Focused, "ButtonLabel", m.Buttons[0].Label)

	// Check if the current form group is less than the total form groups
	if m.CurrentFormGroup < m.TotalFormGroups-1 {
		// If the current focus is on the table, handle the save operation
		if m.Focused == "table" {
			logger.Info("Saving Data to Backend Server", "Operation", "Update")

			// Check if the quiz exists in the backend
			quizExists, _, err := m.checkIfQuizExists()
			if err != nil {
				logger.Error("Error checking if quiz exists", "error", err)
			}

			// Save responses to backend based on whether the quiz exists
			if quizExists {
				m.saveResponsesToBackend()
			} else {
				m.postResponseToBackend()
			}
		}

		// If the current focus is on the button and the label is "Create", handle the create operation
		if m.Focused == "button" && m.Buttons[0].Label == "Create" {
			logger.Info("Posting Data to Backend Server", "Operation", "Create")
			m.postResponseToBackend()
		}

		// Move to the next form group
		m.CurrentFormGroup++
		logger.Info("Moved to next form", "CurrentFormGroup", m.CurrentFormGroup)
	} else {
		// If the current form group is the last one, log that the end of forms has been reached
		logger.Info("NextForm: Reached end of forms")
	}
}

// Move to the previous form in the sequence
func (m *DynamicQuizModel) PreviousForm() {
	logger.Info("PreviousForm called", "CurrentFormGroup", m.CurrentFormGroup)
	if m.CurrentFormGroup > 0 {
		m.CurrentFormGroup--
	}
	logger.Info("Moved to previous form", "CurrentFormGroup", m.CurrentFormGroup)
}

// Post the response to the backend
func (m *DynamicQuizModel) postResponseToBackend() {
	logger.Info("Posting responses to backend")
	for i, q := range m.QuestionForms {
		// Parse the answersCommaSeparated into individual answers
		answerTexts := strings.Split(answersCommaSeparated[i], ",")
		for _, answerText := range answerTexts {
			if strings.TrimSpace(answerText) == q.CorrectAnswer {
				q.Answers = append(q.Answers, Answer{QuestionId: q.Question.ID, Text: strings.TrimSpace(answerText), IsCorrect: true})
			} else {
				q.Answers = append(q.Answers, Answer{QuestionId: q.Question.ID, Text: strings.TrimSpace(answerText), IsCorrect: false})
			}
		}

		// Set the creation date and last modified date to the current date and time
		currentDate := time.Now().UTC()
		q.Question.CreationDate = currentDate
		q.Question.LastModifiedDate = currentDate

		// Create question response data
		questionResponse := map[string]interface{}{
			"quiz_id":                m.QuizID,
			"id":                     q.Question.ID,
			"text":                   q.Question.Text,
			"type":                   q.Question.Type,
			"points":                 q.Question.Points,
			"hint_explanation":       q.Question.HintExplanation,
			"difficulty_level":       q.Question.DifficultyLevel,
			"multi_choice_ans_limit": q.Question.MultiChoiceAnsLimit,
			"creation_date":          q.Question.CreationDate.Format(time.RFC3339),
			"last_modified_date":     q.Question.LastModifiedDate.Format(time.RFC3339),
		}
		logger.Info("Sending question to backend (POST)", "questionID", q.Question.ID, "questionResponse", questionResponse)
		questionID, err := m.postQuestionToBackend(questionResponse)
		if err != nil {
			logger.Error("Failed to send question to backend", "error", err)
		}
		q.Question.ID = questionID

		// Clear the Answers list to avoid duplication
		q.Answers = nil

		// Send answers to the backend
		for _, ans := range q.Answers {
			// Set the creation date and last modified date to the current date and time for answers
			ans.CreationDate = currentDate
			ans.LastModifiedDate = currentDate
			ans.QuestionId = questionID

			answerResponse := map[string]interface{}{
				"question_id":        questionID,
				"text":               ans.Text,
				"is_correct":         ans.IsCorrect,
				"creation_date":      ans.CreationDate.Format(time.RFC3339),
				"last_modified_date": ans.LastModifiedDate.Format(time.RFC3339),
			}
			logger.Info("Sending answer to backend", "answer", ans, "answerResponse", answerResponse)
			err := m.postAnswerToBackend(answerResponse)
			if err != nil {
				logger.Error("Failed to send answer to backend", "error", err)
			}
		}
	}
	logger.Info("All responses posted to backend")
}
