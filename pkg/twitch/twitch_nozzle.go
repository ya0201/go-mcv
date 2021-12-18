package twitch

import (
	"github.com/ya0201/go-mcv/pkg/comment"
	"github.com/ya0201/go-mcv/pkg/nozzle"
)

// twitchNozzleはNozzle interfaceを実装している
var _ nozzle.Nozzle = (*twitchNozzle)(nil)

func TwitchNozzle() *twitchNozzle {
	return &twitchNozzle{}
}

type twitchNozzle struct {
}

func (tn *twitchNozzle) Pump() (<-chan comment.Comment, error) {
	c := make(chan comment.Comment, 50)

	comm := comment.Comment{
		StreamingPlatform: "twitch",
		Msg:               "hello, twitch nozzle!",
	}

	c <- comm
	close(c)

	return c, nil
}
