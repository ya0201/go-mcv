package twitch_nozzle

import (
	gti "github.com/gempir/go-twitch-irc/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ya0201/go-mcv/pkg/comment"
	"github.com/ya0201/go-mcv/pkg/nozzle"
)

// twitchNozzleはNozzle interfaceを実装している
var _ nozzle.Nozzle = (*twitchNozzle)(nil)

type twitchNozzle struct {
	client        *gti.Client
	commentFilter *comment.CommentFilter
}

func NewTwitchNozzle() *twitchNozzle {
	channel := viper.GetString("TWITCH_CHANNEL_ID")
	if channel != "" {
		zap.S().Infof("twitch channel id: %s", channel)
	}

	if channel == "" {
		return nil
	}

	client := gti.NewAnonymousClient()
	client.Join(channel)

	zap.S().Info("TwitchNozzle initialized!")
	return &twitchNozzle{client: client}
}

func (this *twitchNozzle) SetCommentFilter(cf *comment.CommentFilter) {
	this.commentFilter = cf
}

func (tn *twitchNozzle) Pump() (<-chan comment.Comment, error) {
	if tn == nil {
		zap.S().Panic("TwitchNozzle is nil.")
	}

	zap.S().Debugf("twitch nozzle: %+v", *tn)
	c := make(chan comment.Comment, 50)

	tn.client.OnPrivateMessage(func(message gti.PrivateMessage) {
		zap.S().Debugf("received message (variable of go-twitch-irc.PrivateMessage): %+v\n", message)
		if message.Message == "" {
			return
		}

		comm := comment.Comment{
			StreamingPlatform: "twitch",
			Msg:               message.Message,
		}

		if tn.commentFilter != nil && tn.commentFilter.IsInvalid(comm) {
			return
		} else {
			c <- comm
		}
	})

	go func() {
		err := tn.client.Connect()
		if err != nil {
			zap.S().Panic(err)
		}
	}()

	return c, nil
}
