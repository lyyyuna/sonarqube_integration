package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type WechatClient struct {
	wechatApi  *resty.Request
	token      string
	PullNumber string
	RepoOwner  string
	RepoName   string
	JobName    string
	JobSpec    *ProwJobSpec
}

func newWechatClient(token string) *WechatClient {
	return &WechatClient{
		wechatApi: resty.New().SetDisableWarn(true).
			R(),
		token: token,
	}
}

func (c *WechatClient) prepareWechatContent(dashboardUrl string, cnt int) string {
	content := fmt.Sprintf(`## %v 代码静态扫描结果
[点击登录 sonarqube](%v)

#### PR 信息
提交者: %v
PR 编号: %v
GitHub PR链接: [%v](%v)

#### 问题数量
%v
`, c.RepoName, dashboardUrl, c.JobSpec.Refs.Pulls[0].Author, c.PullNumber, c.RepoName, c.JobSpec.Refs.Pulls[0].Link, cnt)
	return content
}

func (c *ClientAnalysis) PostToWechat(token string) {
	wxClient := newWechatClient(token)
	err := wxClient.getEnvironmentVariables()
	if err != nil {
		return
	}
	body := wxClient.prepareWechatContent(c.Task.DashboardUrl, len(c.searchOpenIssues()))
	wxClient.postToWechat(body)
}

func (c *WechatClient) getEnvironmentVariables() error {
	c.PullNumber = os.Getenv("PULL_NUMBER")
	c.RepoOwner = os.Getenv("REPO_OWNER")
	c.RepoName = os.Getenv("REPO_NAME")
	c.JobName = os.Getenv("JOB_NAME")
	jobSpec := os.Getenv("JOB_SPEC")

	if c.PullNumber == "" ||
		c.RepoName == "" ||
		c.RepoOwner == "" ||
		c.JobName == "" ||
		jobSpec == "" {
		log.Errorf("Fail to get environment variables, check the following variables to see if they are right. "+
			"PULL_NUMBER: %v, REPO_OWNER: %v, REPO_NAME： %v, JOB_NAME: %v, JOB_SPEC: %v", c.PullNumber, c.RepoOwner, c.RepoName, c.JobName, jobSpec)
		return errors.New("Some environment variables missing.")
	}

	var prowJobSpec ProwJobSpec
	err := json.Unmarshal([]byte(jobSpec), &prowJobSpec)
	if err != nil {
		log.Fatalf("Fail to decode the json body, the error is %v", err)
	}
	c.JobSpec = &prowJobSpec
	return nil
}

func (c *WechatClient) postToWechat(content string) {
	targetUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%v", c.token)

	type markdown struct {
		Content string `json:"content"`
	}
	req := struct {
		MsgType  string   `json:"msgtype"`
		Markdown markdown `json:"markdown"`
	}{
		MsgType: "markdown",
		Markdown: markdown{
			Content: content,
		},
	}

	resp, err := c.wechatApi.
		SetBody(req).
		Post(targetUrl)
	if err != nil {
		log.Errorf("Cannot connect to Wechat, the error is %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Wechat server response not 200, the response code is: %v, response body is %v", resp.StatusCode(), resp.Body())
	}
}
