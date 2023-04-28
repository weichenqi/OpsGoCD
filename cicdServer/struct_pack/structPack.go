package structPack

type DeployTaskDetails struct {
	TaskId       string `json:"taskId"`
	PushDate     string `json:"pushDate"`
	ProjectName  string `json:"projectName"`
	BranchName   string `json:"branchName"`
	Branch       string `json:"branch"`
	GitSshUrl    string `json:"gitSshUrl"`
	GitHttpUrl   string `json:"gitHttpUrl"`
	PushUserName string `json:"pushUserName"`
}

type DeployTaskState struct {
	TaskId        string  `json:"taskId"`
	StartTime     string  `json:"startTime"`
	NodeId        string  `json:"NodeId"`
	DeployState   float64 `json:"deployState"`
	DeployDetails string  `json:"deployDetails"`
	PushImageName string  `json:"pushImageName"`
}

type DeployNodeState struct {
	NodeId   string  `json:"nodeId"`
	HostName string  `json:"hostName"`
	LoadNum  float64 `json:"loadNum"`
	//ClientToken  string  `json:"clietToken"`
	NodeUpdateTs float64 `json:"nodeUpdateTs"`
}
