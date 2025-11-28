package profile

// UserProfile represents the user's language profile and learning goals.
type UserProfile struct {
	Language  string `json:"language" required:"true" description:"Language the user speaks"`
	ReadKana  bool   `json:"readKana" required:"true" description:"True if the user can read hiragana and katakana"`
	LevelFrom string `json:"levelFrom" required:"true" description:"none, N1, N2, N3, N4, N5, JLPT N1, JLPT N2, JLPT N3, JLPT N4, JLPT N5, etc"`
	LevelTo   string `json:"levelTo" required:"true" description:"N1, N2, N3, N4, N5, JLPT N1, JLPT N2, JLPT N3, JLPT N4, JLPT N5, etc"`
}
