package structPack

type DeployBasicInfo struct {
	TaskId          string   `json:"taskId"`
	ProjectName     string   `json:"projectName"`
	BranchName      string   `json:"branchName"`
	Branch          string   `json:"branch"`
	GitSshUrl       string   `json:"gitSshUrl"`
	GithttpUrl      string   `json:"githttpUrl"`
	CodeType        string   `json:"codeType"`
	DeployCmd       []string `json:"deployCmd"`
	DeployBaseImage string   `json:"deployImage"`
	Language        string   `json:"language"`
	BuildScripts    string   `json:"buildScripts"`
	WorkDir         string   `json:"work_dir"`
}

type AgentDeployInfo struct {
	TaskId        string `json:"taskId"`
	StartTime     string `json:"startTime"`
	TaskInfo      string `json:"taskInfo"`
	TaskStats     int    `json:"taskStats"`
	PushImageName string `json:"pushImageName"`
}

type ClientInfo struct {
	ClientToken      string          `json:"clientToken"`
	NodeId           string          `json:"nodeId"`
	HostName         string          `json:"hostName"`
	LastUpdate       int64           `json:"lastUpdate"`
	Category         string          `json:"category"`
	DeployRunningNum float64         `json:"deployRunningNum"`
	AgentDeployInfo  AgentDeployInfo `json:"agentDeployInfo"`
}

type Ci struct {
	Build struct {
		Stage    string `yaml:"stage"`
		Language string `yaml:"language"`
		Image    string `yaml:"image"`
		Workdir  string `yaml:"workdir"`
		Script   string `yaml:"script"`
	}
}
