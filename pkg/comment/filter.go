package comment

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type CommentFilter struct {
	MaxLength int
	NgWords   []string
	NgRegexps []*regexp.Regexp
}

func NewCommentFilter(maxLength int, ngWords []string, ngRegexpStrings []string) *CommentFilter {
	var regexps []*regexp.Regexp

	for _, s := range ngRegexpStrings {
		r := regexp.MustCompile(s)
		regexps = append(regexps, r)
	}

	return &CommentFilter{MaxLength: maxLength, NgWords: ngWords, NgRegexps: regexps}
}

func (this *CommentFilter) IsInvalid(msg Comment) bool {
	if utf8.RuneCountInString(msg.Msg) > this.MaxLength {
		return true
	}

	for _, word := range this.NgWords {
		if strings.Contains(msg.Msg, word) {
			return true
		}
	}

	for _, r := range this.NgRegexps {
		if r.Match([]byte(msg.Msg)) {
			return true
		}
	}

	return false
}
