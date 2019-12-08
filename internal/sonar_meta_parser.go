package internal

import (
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

// SonarProjectProperties defines the struct of sonar-project.properties
type SonarProjectProperties struct {
	HostUrl    string `properties:"sonar.host.url,omitempty"`
	ProjectKey string `properties:"sonar.projectKey,omitempty"`
	Token      string `properties:"sonar.login,omitempty"`
}

type SonarReportTask struct {
	ProjectKey    string `properties:"projectKey"`
	ServerUrl     string `properties:"serverUrl"`
	ServerVersion string `properties:"serverVersion"`
	Branch        string `properties:"branch"`
	DashboardUrl  string `properties:"dashboardUrl"`
	CeTaskId      string `properties:"ceTaskId"`
	CeTaskUrl     string `properties:"ceTaskUrl"`
}

func NewSonarProperty(path string) *SonarProjectProperties {
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		log.Fatalf("Fail to open the sonar-project.properties, the error is %v", err)
	}
	var cfg SonarProjectProperties
	if err = p.Decode(&cfg); err != nil {
		log.Fatalf("Fail to decode the sonar-project.properties, the error is %v", err)
	}
	return &cfg
}

func NewSonarReportTask(path string) *SonarReportTask {
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		log.Fatalf("Fail to open the report-task.txt, the error is %v", err)
	}
	var cfg SonarReportTask
	if err = p.Decode(&cfg); err != nil {
		log.Fatalf("Fail to decode the report-task.txt, the error is %v", err)
	}
	return &cfg
}
