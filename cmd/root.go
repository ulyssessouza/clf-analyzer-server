package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const AppName = "clf-analyzer-server"

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: fmt.Sprintf("%s is a tool for serving  analytics about the HTTP logs", AppName),
	Long:  fmt.Sprintf(`%s is a tool for serving  analytics about the HTTP logs coming from an external software. Eg: apache-httpd or nginx`, AppName),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("On RootCmd!")
	},
}

func Init() {
	var echoTimes int

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [# times] [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println(echoTimes)
				fmt.Println("Echo: " + strings.Join(args, " " ))
			}
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	rootCmd.AddCommand(PrintCmd, cmdEcho)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}