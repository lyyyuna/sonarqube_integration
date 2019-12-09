package internal

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ceTaskUrl = "%v/api/ce/activity"
)

type ClientAnalysis struct {
	Property        *SonarProjectProperties
	Task            *SonarReportTask
	webApi          *resty.Request
	analysisId      string
	GithubTokenPath string
	WechatToken     string
}

func NewClientAnalysis(token string) *ClientAnalysis {
	return &ClientAnalysis{
		webApi: resty.New().SetDisableWarn(true).
			R().
			SetHeader("Content-Type", "application/json").
			SetBasicAuth(token, ""),
	}
}

// wait until the analysis finished on SonarQube server side
// and return the analysis ID of this task
func (c *ClientAnalysis) WaitUntilFinished() {
	for i := 1; i < 240; i++ {
		sonarTask := c.ceLatestTask()
		log.Infof("Latest ce task status is %v", sonarTask.Status)
		if sonarTask.Status == TASK_IN_PROGRESS || sonarTask.Status == TASK_PENDING {
			time.Sleep(time.Second)
		} else {
			if sonarTask.Status == TASK_FAILED {
				log.Fatal("Status is failed. Maybe two analysis for one project run in the same time!!!")
			}
			c.analysisId = sonarTask.AnalysisId
			log.Infof("Analysis ID is %v", c.analysisId)
			return
		}
	}
	log.Fatalf("The task didn't finished in 240s")
}

// if not pass, it will panic
func (c *ClientAnalysis) FeedbackToCICheck() {
	qgStatus := c.qualityGatesProjectStatus()
	if qgStatus.Status == "ERROR" {
		log.Fatalln("Why ERROR? It means the project did not pass the Quality Gate, \nyou have to check it in the SonarQube UI.")
		panic("Why ERROR? It means the project did not pass the Quality Gate, \nyou have to check it in the SonarQube UI.")
	} else {
		log.Info("Your project passed the SonarQube scan!")
	}
}

func (c *ClientAnalysis) FeedbackToPRComment() {
	githubc, err := NewGitHubClient(c.GithubTokenPath)
	if err != nil {
		log.Error("Skip comment to GitHub PR, as no github token found, the error is ", err)
		return
	}

	openIssues := c.searchOpenIssues()
	if len(openIssues) == 0 {
		log.Info("No issues found, will not comment in the PR.")
		//return
	}
	dashboardUrl := c.Task.DashboardUrl
	link := fmt.Sprintf("\nClick [%v](%v) to view issue details\n", dashboardUrl, dashboardUrl)
	rerunCmd := fmt.Sprintf("\nSay `/test %v` to re-run this static analysis\n", githubc.JobName)

	body := fmt.Sprintf("\n**%v** issues found during static analysis\n", len(openIssues))
	body += link + rerunCmd

	err = githubc.deletePreviousComments()
	if err != nil {
		log.Error(err)
		return
	}
	err = githubc.postCommentsToPR(body)
	if err != nil {
		log.Error(err)
		return
	}
}
