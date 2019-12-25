package internal

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (c *ClientAnalysis) ceLatestTask() *SonarTask {
	log.Infof("CE engine's task url is %v", c.Task.CeTaskUrl)
	resp, err := c.webApi.Get(c.Task.CeTaskUrl)
	if err != nil {
		log.Fatalf("Cannot connect to SonarQube server, the error is %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("SonarQube Server response not 200, the response code is: %v, response body is %v", resp.StatusCode(), resp.Body())
	}

	body := resp.Body()
	ret := struct {
		Task SonarTask `json:"task"`
	}{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Fatalf("Fail to decode the json body, the error is %v", err)
	}

	return &ret.Task
}

func (c *ClientAnalysis) searchOpenIssues() []SonarIssue {
	targetUrl := fmt.Sprintf("%v/api/issues/search", c.Task.ServerUrl)

	var allIssues []SonarIssue
	for i := 1; i < 10; i++ {
		var payload map[string]string
		if c.Task.Branch == "" {
			payload = map[string]string{
				"componentKeys": c.Task.ProjectKey,
				"statuses":      "OPEN",
				"ps":            "500",
				"p":             strconv.Itoa(i),
			}
		} else {
			payload = map[string]string{
				"componentKeys": c.Task.ProjectKey,
				"statuses":      "OPEN",
				"ps":            "500",
				"p":             strconv.Itoa(i),
				"branch":        c.Task.Branch,
			}
		}
		resp, err := c.webApi.
			SetQueryParams(payload).Get(targetUrl)
		if err != nil {
			log.Fatalf("Cannot connect to SonarQube server, the error is %v", err)
		}
		if resp.StatusCode() != http.StatusOK {
			log.Fatalf("SonarQube Server response not 200, the response code is: %v, response body is %v", resp.StatusCode(), resp.Body())
		}

		ret := struct {
			Issues []SonarIssue `json:"issues"`
		}{}
		err = json.Unmarshal(resp.Body(), &ret)
		if err != nil {
			log.Fatalf("Fail to decode the json body, the error is %v", err)
		}
		if len(ret.Issues) > 0 {
			allIssues = append(allIssues, ret.Issues...)
		} else {
			break
		}
	}

	return allIssues
}

func (c *ClientAnalysis) qualityGatesProjectStatus() *SonarQualityGate {
	targetUrl := fmt.Sprintf("%v/api/qualitygates/project_status", c.Task.ServerUrl)
	resp, err := c.webApi.
		SetQueryParams(map[string]string{
			"analysisId": c.analysisId,
		}).Get(targetUrl)
	if err != nil {
		log.Fatalf("Cannot connect to SonarQube server, the error is %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("SonarQube Server response not 200, the response code is: %v, response body is %v", resp.StatusCode(), string(resp.Body()))
	}

	ret := struct {
		ProjectStatus SonarQualityGate `json:"projectStatus"`
	}{}
	err = json.Unmarshal(resp.Body(), &ret)
	if err != nil {
		log.Fatalf("Fail to decode the json body, the error is %v", err)
	}
	return &ret.ProjectStatus
}
