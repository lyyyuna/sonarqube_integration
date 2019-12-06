package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	feedbackComment bool
	sonarqubeURL    string
	sonarqubeToken  string
	githubTokenPath string
)

var feedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Feedback the SonarQube analysis result to GitHub CI checks",
	Run:   feedback,
}

func init() {
	feedbackCmd.Flags().BoolVarP(&feedbackComment, "comment", "c", false, "Toggle comment, whether to feedback details to GitHub PR comment or not, default is feedback.")
	feedbackCmd.Flags().StringVarP(&sonarqubeURL, "server", "s", "http://sonarqube.dev.qiniu.io", "Specify the SonarQube server url, default is http://sonarqube.dev.qiniu.io")
	feedbackCmd.Flags().StringVarP(&sonarqubeToken, "token", "t", "xxxxxxxxxxxxxxxxxx", "Specify the SonarQube token to get authorized")
	feedbackCmd.Flags().StringVarP(&githubTokenPath, "githubtokenpath", "g", "/etc/github/oauth", "Specify the Github token path")
}

func feedback(cmd *cobra.Command, args []string) {
	log.Infof("SonarQube server address is %v", sonarqubeURL)
	if feedbackComment == false {
		log.Info("Bot will feeback to GitHub PR comment")
	} else {
		log.Info("Bot will not feeback to GitHub PR comment")
	}
	log.Infof("GitHub token path is %v", githubTokenPath)
}
