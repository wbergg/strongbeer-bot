package helper

func GetColumnLetter(index int) string {
	letters := ""
	for index >= 0 {
		letters = string('A'+(index%26)) + letters
		index = index/26 - 1
	}
	return letters
}
