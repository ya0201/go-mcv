package youtube_nozzle

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ya0201/go-mcv/pkg/comment"
	"github.com/ya0201/go-mcv/pkg/nozzle"
	"github.com/ya0201/go-mcv/pkg/youtube_live_chat"
)

// youtubeNozzleはNozzle interfaceを実装している
var _ nozzle.Nozzle = (*youtubeNozzle)(nil)

type youtubeNozzle struct {
	client        youtube_live_chat.SimpleLiveChatClient
	commentFilter *comment.CommentFilter
}

func NewYoutubeNozzle() *youtubeNozzle {
	channel := viper.GetString("YOUTUBE_CHANNEL_ID")
	if channel != "" {
		zap.S().Infof("youtube channel id: %s", channel)
	}

	if channel == "" {
		return nil
	}

	client := youtube_live_chat.NewSimpleLiveChatClient()
	err := client.Join(channel)

	if err != nil {
		zap.S().Warn(err.Error())
		return nil
	} else {
		zap.S().Info("YoutubeNozzle initialized!")
		return &youtubeNozzle{client: client}
	}

}

func (this *youtubeNozzle) SetCommentFilter(cf *comment.CommentFilter) {
	this.commentFilter = cf
}

func (this *youtubeNozzle) Pump() (<-chan comment.Comment, error) {
	if this == nil {
		zap.S().Panic("YoutubeNozzle is nil.")
	}

	zap.S().Debugf("youtube nozzle: %+v", *this)
	c := make(chan comment.Comment, 50)

	this.client.OnMessage(func(message *youtube_live_chat.SimpleLiveChatMessage) error {
		zap.S().Debugf("received message (variable of youtube_live_chat.SimpleLiveChatMessage): %+v\n", message)
		if message.Msg == "" {
			return nil
		}

		comm := comment.Comment{
			StreamingPlatform: "youtube",
			Msg:               message.Msg,
		}
		if this.commentFilter != nil && this.commentFilter.IsInvalid(comm) {
			return nil
		} else {
			c <- comm
			return nil
		}
	})

	go func() {
		err := this.client.Connect()
		if err != nil {
			zap.S().Panic(err)
		}
	}()

	return c, nil
}
