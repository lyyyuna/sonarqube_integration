package internal

import ()

const (
	TASK_SUCCESS     = "SUCCESS"
	TASK_PENDING     = "PENDING"
	TASK_FAILED      = "FAILED"
	TASK_IN_PROGRESS = "IN_PROGRESS"
)

type SonarTask struct {
	AnalysisId string `json:"analysisId"`
	Status     string `json:"status"`
}

type SonarIssue struct {
	Severity  string `json:"severity"`
	Component string `json:"component"`
}

type SonarQualityGate struct {
	Status            string `json:"status"`
	IgnoredConditions bool   `json:"ignoredConditions,omitempty"`
}
