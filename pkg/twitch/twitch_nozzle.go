package twitch

import (
	"fmt"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/spf13/viper"
	"github.com/ya0201/go-mcv/pkg/comment"
	"github.com/ya0201/go-mcv/pkg/nozzle"
)

// twitchNozzleはNozzle interfaceを実装している
var _ nozzle.Nozzle = (*twitchNozzle)(nil)

type twitchNozzle struct {
	conn *websocket.Conn
}

func TwitchNozzle() *twitchNozzle {
	channel := viper.GetString("TWITCH_CHANNEL_ID")
	if channel != "" {
		zap.S().Infof("twitch channel id: %s", channel)
	}

	url := viper.GetString("TWITCH_IRC_SSL_WEBSOCKET_ENDPOINT")
	if channel != "" && url == "" {
		zap.S().Panic("Could not get TWITCH_IRC_SSL_WEBSOCKET_ENDPOINT.")
	}

	if channel == "" || url == "" {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		zap.S().Panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("PASS %s", viper.GetString("TWITCH_OAUTH_TOKEN"))))
	if err != nil {
		zap.S().Panic(err)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("NICK %s", viper.GetString("TWITCH_NICKNAME"))))
	if err != nil {
		zap.S().Panic(err)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("JOIN %s", viper.GetString("TWITCH_CHANNEL_ID"))))
	if err != nil {
		zap.S().Panic(err)
	}

	zap.S().Info("TwitchNozzle initialized!")
	return &twitchNozzle{conn: conn}
}

func (tn *twitchNozzle) Pump() (<-chan comment.Comment, error) {
	if tn == nil {
		zap.S().Panic("TwitchNozzle is nil.")
	}

	go func() {
		zap.S().Infof("%+v", *tn)
		defer tn.conn.Close()
		for {
			_, message, err := tn.conn.ReadMessage()
			if err != nil {
				zap.S().Panic(err)
				return
			}
			zap.S().Infof("recv: %s", message)
		}
	}()
	c := make(chan comment.Comment, 50)

	comm := comment.Comment{
		StreamingPlatform: "twitch",
		Msg:               "hello, twitch nozzle!",
	}

	c <- comm
	close(c)

	return c, nil
}
