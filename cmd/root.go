package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sonarqube-bot",
	Short: "Get SonarQube analysis result from server",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			line := strconv.Itoa(f.Line)
			return "[" + funcname + "]", "[" + filename + ":" + line + "]"
		},
	})
	log.SetReportCaller(true)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(feedbackCmd)
	rootCmd.AddCommand(wechatCmd)
}
