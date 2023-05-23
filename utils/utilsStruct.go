package utils

var GlobalJobs []JobStruct

type JobStruct struct {
	Id         string
	Name       string
	Expression string
	Comment    string
}
