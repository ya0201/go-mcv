package comment

import (
	"regexp"
	"strings"
)

type CommentFilter struct {
	MaxLength int
	NgWords   []string
	NgRegexps []*regexp.Regexp
}

func NewCommentFilter(maxLength int, ngWords []string, ngRegexpsString []string) *CommentFilter {
	var regexps []*regexp.Regexp

	for _, s := range ngRegexpsString {
		r := regexp.MustCompile(s)
		regexps = append(regexps, r)
	}

	return &CommentFilter{MaxLength: maxLength, NgWords: ngWords, NgRegexps: regexps}
}

func (this *CommentFilter) IsInvalid(msg Comment) bool {
	if len(msg.Msg) > this.MaxLength {
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
