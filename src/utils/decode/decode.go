package decode

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func UnicodeToUTF8(text string) string {
	compile := regexp.MustCompile(`\\u([0-9a-z]{4})`)
	matches := compile.FindAllStringSubmatch(text, -1)

	for _, v := range matches {
		temp, _ := strconv.ParseInt(v[1], 16, 32)
		str := fmt.Sprintf("%c", temp)
		text = strings.Replace(text, v[0], str, -1)
	}
	return text
}

func GbkToUtf8(s string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

