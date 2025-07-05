package util

// Need test
func ShiftLeftArray(files []string) []string {
	for i := 1; i < len(files); i++ {
		files[i-1] = files[i]
		files[i] = ""
	}

	return files
}
