module github.com/ya0201/go-mcv

replace github.com/ya0201/go-mcv/pkg/comment => ./pkg/comment

replace github.com/ya0201/go-mcv/pkg/nozzle => ./pkg/nozzle

go 1.16

require (
	github.com/gempir/go-twitch-irc/v2 v2.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.10.0
	go.uber.org/zap v1.17.0
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
)
