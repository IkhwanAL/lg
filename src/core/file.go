package core

type FileType int

const (
	File FileType = iota
	Dir
)

type FsEntry struct {
	Name string
	Type FileType
	Path string
}
