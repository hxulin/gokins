package model

// Auth 权限信息
type Auth struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

// AuthConfig 权限配置信息
type AuthConfig struct {
	Auth Auth `yaml:"auth"`
}

// BuildParamItem 构建参数项
type BuildParamItem struct {
	Name  string      `json:"name" yaml:"name"`
	Value interface{} `json:"value" yaml:"value"`
}

// Job 任务描述
type Job struct {
	Id      int              `yaml:"id"`
	Name    string           `yaml:"name"`
	Ack     bool             `yaml:"ack"`
	AckText string           `yaml:"ack-text"`
	Params  []BuildParamItem `yaml:"params"`
	Columns []BuildParamItem `yaml:"columns"`
}

// Gokins gokins 配置
type Gokins struct {
	Username   string `yaml:"username"`
	Token      string `yaml:"token"`
	JenkinsUrl string `yaml:"jenkins-url"`
	Job        []Job  `yaml:"job"`
}

// SysConfig 系统配置
type SysConfig struct {
	Gokins Gokins `yaml:"gokins"`
}

// BuildParam 构建参数
type BuildParam struct {
	Parameter []BuildParamItem `json:"parameter"`
}

// Action 构建信息
type Action struct {
	ClassName  string           `json:"_class"`
	Parameters []BuildParamItem `json:"parameters"`
}

// BuildStatus 部署状态
type BuildStatus struct {
	Actions  []Action `json:"actions"`
	Building bool     `json:"building"`
	Result   string   `json:"result"`
}

// QueueInfo 队列信息
type QueueInfo struct {
	Blocked    bool `json:"blocked"`
	Executable struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"executable"`
}
