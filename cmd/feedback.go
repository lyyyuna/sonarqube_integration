package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	. "sonarqube-bot/internal"
)

var (
	nofeedbackComment   bool
	sonarqubeURL        string
	sonarqubeToken      string
	githubTokenPath     string
	sonarPropertiesPath string
	sonarReportTaskPath string
)

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Feedback the SonarQube analysis result to GitHub CI checks",
	Run:   feedback,
}

func init() {
	feedbackCmd.Flags().BoolVarP(&nofeedbackComment, "nocomment", "", false, "Toggle comment, whether to feedback details to GitHub PR comment or not, default is feedback.")
	feedbackCmd.Flags().StringVarP(&sonarqubeURL, "server", "s", "", "Specify the SonarQube server url, if not specified, bot will get url from report-task.txt")
	feedbackCmd.Flags().StringVarP(&sonarqubeToken, "token", "t", "", "Specify the SonarQube token to get authorized")
	feedbackCmd.Flags().StringVarP(&githubTokenPath, "githubtokenpath", "", "/etc/github/oauth", "Specify the Github token path")
	feedbackCmd.Flags().StringVarP(&sonarPropertiesPath, "sonarproperty", "", "./sonar-project.properties", "Specify the sonar-project.properties path")
	feedbackCmd.Flags().StringVarP(&sonarReportTaskPath, "sonartask", "", "./.scannerwork/report-task.txt", "Specify the report-task.txt path")
}

func feedback(cmd *cobra.Command, args []string) {
	if nofeedbackComment == false {
		log.Info("Bot will feeback to GitHub PR comment")
	} else {
		log.Info("Bot will not feeback to GitHub PR comment")
	}
	log.Infof("GitHub token path is %v", githubTokenPath)
	log.Infof("sonar-project.properties path is %v", sonarPropertiesPath)
	log.Infof("report-task.txt path is %v", sonarReportTaskPath)

	p := NewSonarProperty(sonarPropertiesPath)
	if sonarqubeToken == "" {
		sonarqubeToken = p.Token
	}
	clientAnalysis := NewClientAnalysis(sonarqubeToken)
	clientAnalysis.Property = p
	clientAnalysis.Task = NewSonarReportTask(sonarReportTaskPath)
	clientAnalysis.GithubTokenPath = githubTokenPath
	//
	// wait until the analysis is finished
	clientAnalysis.WaitUntilFinished()
	// feedback comment or not
	if nofeedbackComment == false {
		clientAnalysis.FeedbackToPRComment()
	}
	// feedback check status
	// this will panic, so it must be the last step
	clientAnalysis.FeedbackToCICheck()
}
