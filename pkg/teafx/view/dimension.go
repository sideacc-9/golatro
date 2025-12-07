package view

import "strings"

type Dimensions struct {
	Height int
	Width  int
}

func FitHeight(txt string, height int) string {
	if height < 0 {
		return txt
	}
	newLinesNeeded := height - (strings.Count(txt, "\n") + 1)
	if newLinesNeeded >= 0 {
		result := txt
		for range newLinesNeeded {
			result += "\n"
		}
		return result
	}
	nextIter := txt
	lastIdx := -1
	for range height {
		idx := strings.Index(nextIter, "\n")
		if lastIdx == -1 {
			lastIdx = idx + 1
		} else {
			lastIdx += idx + 1
		}
		nextIter = nextIter[idx+1:]
	}
	return txt[:lastIdx-1]
}

func FitWidth(txt string, width int) string {
	if width < 0 {
		return txt
	}
	lines := strings.Split(txt, "\n")
	result := ""
	for _, line := range lines {
		result += fitWidthLine(line, width) + "\n"
	}
	return result[:len(result)-1] // remove last \n
}

func fitWidthLine(line string, width int) string {
	if width < 0 {
		return line
	}
	spacesNeeded := width - len([]rune(line))
	if spacesNeeded >= 0 {
		return line + strings.Repeat(" ", spacesNeeded)
	}

	return line[:width]
}

func (dimensions Dimensions) Fit(txt string) string {
	return FitDimensions(txt, dimensions)
}

func FitDimensions(txt string, dimensions Dimensions) string {
	if dimensions.Height == -1 && dimensions.Width == -1 {
		return txt
	}
	fitted := txt
	if dimensions.Height != -1 {
		fitted = FitHeight(fitted, dimensions.Height)
	}
	if dimensions.Width != -1 {
		fitted = FitWidth(fitted, dimensions.Width)
	}
	return fitted
}
