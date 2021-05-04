package teams

type TeamPayload struct {
	Summary     string          `json:"summary"`
	Title       string          `json:"title,omitempty"`
	Text        string          `json:"text"`
	Color       string          `json:"themeColor,omitempty"`
	Sections    []Section      `json:"sections"`
}

type Section struct {
	ActivityTitle   string      `json:"activityTitle"`
	Facts           []Fact     `json:"facts"`
}

type Fact struct {
	Name    string    `json:"name"`
	Value   string    `json:"value"`
}