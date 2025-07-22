package core

import "runtime"

type UserArgs struct {
	OpenDirWith string
}

func (u UserArgs) GetOpenDirArgs() string {
	if u.OpenDirWith != "" {
		return u.OpenDirWith
	}

	var openDir string
	switch runtime.GOOS {
	case "windows":
		openDir = "explorer"
	case "linux":
		openDir = "xdg-open"
	}

	return openDir
}
