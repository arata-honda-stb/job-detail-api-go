package internal

type Job struct {
	JobID      string `json:"jobId" dynamodbav:"jobId"`
	JobTitle   string `json:"jobTitle" dynamodbav:"jobTitle"`
	JobContent string `json:"jobContent" dynamodbav:"jobContent"`
}
