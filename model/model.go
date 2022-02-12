package model

// BuildParamItem 构建参数项
type BuildParamItem struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// BuildParam 构建参数
type BuildParam struct {
	Parameter []BuildParamItem `json:"parameter"`
}

// DeployStatus 部署状态
type DeployStatus struct {
	Building bool   `json:"building"`
	Result   string `json:"result"`
}

// Job 任务描述
type Job struct {
	Id      int              `yaml:"id"`
	Tag     string           `yaml:"tag"`
	Name    string           `yaml:"name"`
	Ack     bool             `yaml:"ack"`
	AckText string           `yaml:"ack-text"`
	Remark  string           `yaml:"remark"`
	Params  []BuildParamItem `yaml:"params"`
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
