package core

type SearchTypeChangedMsg struct {
	Search string
}

type SearchResultMsg struct {
	Result []FsEntry
}

type PathMsg struct {
	Path string
}
