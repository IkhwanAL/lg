package core

type SearchTypeChangedMsg struct {
	Search string
}

type SearchResultMsg struct {
	Result []string
}
