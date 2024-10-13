package common

// InitializeChoices sets up the initial choices for a given screen
func InitializeChoices(screen string) Model {
	var choices []string
	switch screen {
	case "login":
		choices = []string{"Login via SSO", "Signup via SSO"}
	case "menu":
		choices = []string{
			"Setup (for Admins Only)",
			"Select category & Start Quiz",
			"Create/Edit Questionnaire & Answers",
			"View Leaderboard score",
			"Submit Enhancement Request",
			"Exit",
		}
	}

	return Model{
		CurrentScreen: screen,
		Choices:       choices,
		ButtonPos:     make([]Rect, len(choices)),
	}
}
