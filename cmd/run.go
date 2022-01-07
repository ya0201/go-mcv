package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ya0201/go-mcv/pkg/twitch"
	"go.uber.org/zap"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs go-mcv.",
	Long:  `Runs go-mcv. Currently, twitch is only supported.`,
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
	runCmd.PersistentFlags().String("twitch-channel-id", "", "channel id you want to get comments")
	viper.BindPFlag("TWITCH_CHANNEL_ID", runCmd.PersistentFlags().Lookup("twitch-channel-id"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run() {
	tn := twitch.TwitchNozzle()
	if tn == nil {
		zap.S().Info("TwitchNozzle is not initialized, and does not panic. So, ignoring twitch...")
		return
	}

	c, _ := tn.Pump()
	time.Sleep(time.Second * 10)
	zap.S().Info("Start pumping ...")
	for c := range c {
		zap.S().Debugf("%+v", c)
		fmt.Println(c.Msg)
	}
}
