package cmd

import (
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ya0201/go-mcv/pkg/comment"
	"github.com/ya0201/go-mcv/pkg/logging"
	"github.com/ya0201/go-mcv/pkg/nozzle"
	"github.com/ya0201/go-mcv/pkg/twitch_nozzle"
	"github.com/ya0201/go-mcv/pkg/youtube_nozzle"
	"go.uber.org/zap"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs go-mcv.",
	Long:  `Runs go-mcv. Currently, only twitch is supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	runCmd.PersistentFlags().String("twitch-channel-id", "", "the twitch channel id you want to get comments")
	viper.BindPFlag("TWITCH_CHANNEL_ID", runCmd.PersistentFlags().Lookup("twitch-channel-id"))
	runCmd.PersistentFlags().String("youtube-channel-id", "", "the youtube channel id you want to get comments")
	viper.BindPFlag("YOUTUBE_CHANNEL_ID", runCmd.PersistentFlags().Lookup("youtube-channel-id"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run() {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			textView.ScrollToEnd()
			return nil
		}
		return event
	})
	logging.SetLoggerOutputToTview(textView)

	tn := twitch_nozzle.NewTwitchNozzle()
	yn := youtube_nozzle.NewYoutubeNozzle()
	var tc, yc <-chan comment.Comment

	if tn == nil {
		zap.S().Info("TwitchNozzle is not initialized, and does not panic. So, ignoring twitch...")
	} else {
		initCommentFilter(tn)
		tc, _ = tn.Pump()
	}
	if yn == nil {
		zap.S().Info("YoutubeNozzle is not initialized, and does not panic. So, ignoring youtube...")
	} else {
		initCommentFilter(yn)
		yc, _ = yn.Pump()
	}

	if tn == nil && yn == nil {
		zap.S().Info("All nozzles are nil. Do nothing.")
		zap.S().Info("Press Ctrl+C to exit...")
	} else {
		zap.S().Info("Start pumping ...")
		go func() {
			for {
				select {
				case msg := <-tc:
					printMsg(textView, msg)
				case msg := <-yc:
					printMsg(textView, msg)
				}
			}
		}()
	}

	frame := tview.NewFrame(textView).
		AddText("Enter: scroll to latest comments", false, tview.AlignCenter, tcell.ColorWhite).
		AddText("Ctrl+C: exit", false, tview.AlignCenter, tcell.ColorWhite)
	if err := app.SetRoot(frame, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func printMsg(w io.Writer, msg comment.Comment) {
	zap.S().Debugf("%+v", msg)

	var msgPrefix string
	switch msg.StreamingPlatform {
	case "twitch":
		msgPrefix = "[blue]TW[-]"
	case "youtube":
		msgPrefix = "[red]YT[-]"
	}

	fmt.Fprintf(w, "%s %s\n\n", msgPrefix, msg.Msg)
}

func initCommentFilter(nozzle nozzle.Nozzle) {
	maxLength := viper.GetInt("COMMENT_FILTER_MAX_LENGTH")
	ngWords := viper.GetStringSlice("COMMENT_FILTER_NG_WORDS")
	ngRegexpStrings := viper.GetStringSlice("COMMENT_FILTER_NG_REGEXP_STRINGS")

	filter := comment.NewCommentFilter(maxLength, ngWords, ngRegexpStrings)
	nozzle.SetCommentFilter(filter)
}
