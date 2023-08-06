package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var vFlag *bool
var Version string

const baseURL = "https://www.data.jma.go.jp/mscweb/data/himawari/list_r2w.html"

var rootCmd = &cobra.Command{
	Use:   "himawari-dl",
	Short: "Download Himawari satellite images from data.jma.go.jp",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *vFlag {
			fmt.Println(Version)
			return
		}

		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	vFlag = rootCmd.Flags().BoolP("version", "v", false, "Show program version")
}
