package cmd

import (
	"fmt"

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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run() {
	c, _ := twitch.TwitchNozzle().Pump()
	zap.S().Info("Start pumping ...")
	zap.S().Info(viper.Get("hoge"))
	for c := range c {
		zap.S().Debugf("%+v", c)
		fmt.Println(c.Msg)
	}
}
