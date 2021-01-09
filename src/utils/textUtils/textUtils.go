package textUtils

import (
	"strings"
)

func PullMiddleText(input string, pre string, suf string) (result string) {
	index := strings.Index(input, pre)
	if index != -1 {
		endIndex := strings.Index(input[index+len(pre):], suf)
		if endIndex != -1 {
			return input[index+len(pre) : index+len(pre)+endIndex]
		}
	} else {
		return ""
	}

	return ""
}

func PullRightText(input string, pre string) (result string) {
	index := strings.Index(input, pre)
	if index != -1 {
		return input[index+len(pre):]
	}
	return ""
}

func PullAllMiddleText(input string, pre string, suf string) []string {

	result := make([]string, 0, 0)
	index := 0
	for index != -1 {
		index = strings.Index(input, pre)
		endIndex := strings.Index(input[index+len(pre):], suf)
		if endIndex != -1 {
			result = append(result, input[index+len(pre):index+len(pre)+endIndex])
			input = input[index+len(pre)+endIndex:]
		} else {
			break
		}
		index = strings.Index(input, pre)
	}

	return result
}
