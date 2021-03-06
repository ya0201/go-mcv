package nozzle

import (
	"github.com/ya0201/go-mcv/pkg/comment"
)

type Nozzle interface {
	Pump() (<-chan comment.Comment, error)
	SetCommentFilter(*comment.CommentFilter)
}

// helloNozzleはNozzle interfaceを実装している
var _ Nozzle = (*helloNozzle)(nil)

func HelloNozzle() *helloNozzle {
	return &helloNozzle{}
}

type helloNozzle struct {
}

func (hn *helloNozzle) Pump() (<-chan comment.Comment, error) {
	c := make(chan comment.Comment, 50)

	comm := comment.Comment{
		StreamingPlatform: "hoge",
		Msg:               "hello, nozzle!",
	}

	c <- comm
	close(c)

	return c, nil
}

func (this *helloNozzle) SetCommentFilter(cf *comment.CommentFilter) {
}
