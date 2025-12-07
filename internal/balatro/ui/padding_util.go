package ui

func PadRight(str string, length int, char string) string {
	strRunes := []rune(str)
	if len(strRunes) >= length {
		return string(strRunes[:length])
	}
	strRes := str
	for i := len(strRunes); i < length; i++ {
		strRes += char
	}
	return strRes
}

func PadLeft(str string, length int, char string) string {
	strRunes := []rune(str)
	if len(strRunes) >= length {
		return string(strRunes[:length])
	}
	strRes := ""
	for i := 0; i < length-len(strRunes); i++ {
		strRes += char
	}
	strRes += str
	return strRes
}

func CenterString(str string, length int, filler string) string {
	runes := []rune(str)
	if len(runes) >= length {
		return string(runes[:length])
	}
	whiteSpaces := length - len(runes)
	startIdx := whiteSpaces / 2
	padded := ""
	for i := 0; i < length; i++ {
		if i >= startIdx && i < startIdx+len(runes) {
			padded += str
			i += len(runes) - 1
		} else {
			padded += " "
		}
	}
	return padded
}

func RepeatString(str string, times int) string {
	return CenterString("", times, str)
}
