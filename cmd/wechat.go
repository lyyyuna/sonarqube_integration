package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "sonarqube-bot/internal"
)

var (
	wechatToken string
)

var wechatCmd = &cobra.Command{
	Use:   "wechat",
	Short: "Post SonarQube result to Wechat",
	Run:   wechat,
}

func init() {
	wechatCmd.Flags().StringVarP(&sonarqubeToken, "token", "t", "", "Specify the SonarQube token to get authorized")
	wechatCmd.Flags().StringVarP(&sonarPropertiesPath, "sonarproperty", "", "./sonar-project.properties", "Specify the sonar-project.properties path")
	wechatCmd.Flags().StringVarP(&sonarReportTaskPath, "sonartask", "", "./.scannerwork/report-task.txt", "Specify the report-task.txt path")
	wechatCmd.Flags().StringVarP(&wechatToken, "wxtoken", "w", "", "Specify the wechat token to post to wechat")
}

func wechat(cmd *cobra.Command, args []string) {
	log.Infof("sonar-project.properties path is %v", sonarPropertiesPath)
	log.Infof("report-task.txt path is %v", sonarReportTaskPath)
	p := NewSonarProperty(sonarPropertiesPath)
	if sonarqubeToken == "" {
		sonarqubeToken = p.Token
	}
	clientAnalysis := NewClientAnalysis(sonarqubeToken)
	clientAnalysis.Property = p
	clientAnalysis.Task = NewSonarReportTask(sonarReportTaskPath)
	//
	// wait until the analysis is finished
	clientAnalysis.WaitUntilFinished()
	clientAnalysis.PostToWechat(wechatToken)
}
