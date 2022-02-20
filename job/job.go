package job

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gokins/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Building 任务构建中
// Success 任务构建成功
// Failure 任务构建失败
// Aborted 用户终止构建
const (
	Building = "BUILDING"
	Success  = "SUCCESS"
	Failure  = "FAILURE"
	Aborted  = "ABORTED"
)

// QueryBuildStatus 查询任务的构建状态
func QueryBuildStatus(queryUrl string, username string, token string) (model.BuildStatus, error) {
	buildStatus, err := queryBuildStatus(queryUrl, username, token)
	if err != nil {
		return model.BuildStatus{}, err
	}
	return buildStatus, nil
}

// ParseBuildStatus 解析任务的构建状态
func ParseBuildStatus(buildStatus model.BuildStatus) string {
	if buildStatus.Building {
		return Building
	}
	return buildStatus.Result
}

// ParseCurrentJobBuildParam 解析任务的构建参数
func ParseCurrentJobBuildParam(buildStatus model.BuildStatus) []model.BuildParamItem {
	for _, action := range buildStatus.Actions {
		if action.ClassName == "hudson.model.ParametersAction" {
			return action.Parameters
		}
	}
	return nil
}

func Build(buildUrl string, param model.BuildParam, username string, token string) error {
	client := &http.Client{}
	//param := buildParam{Parameter: []buildParamItem{
	//	{"branch", "dev1.1.6"},
	//	{"MAVEN_CMD", "clean install"},
	//	{"project_name", "dxm-assist"},
	//	{"namespace", "--spring.profiles.active=test"},
	//	{"dest_ip", "172.19.170.126"},
	//	{"group", "rdts/backend"},
	//}}
	paramStr, _ := json.Marshal(param)
	fmt.Println(string(paramStr))
	req, err := http.NewRequest(http.MethodPost, buildUrl, strings.NewReader("json="+encodeURIComponent(string(paramStr))))
	if err != nil {
		return err
	}
	auth := []byte(username + ":" + token)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// URL 参数编码，实现和 JS 通用
func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	return strings.Replace(r, "+", "%20", -1)
}

// 查询任务构建状态
func queryBuildStatus(queryUrl string, username string, token string) (model.BuildStatus, error) {
	buildStatus := model.BuildStatus{}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, queryUrl, nil)
	if err != nil {
		return buildStatus, err
	}
	// 添加请求头
	auth := []byte(username + ":" + token)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
	req.Header.Set("Content-type", "application/json;charset=utf-8")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return buildStatus, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return buildStatus, err
	}
	err = json.Unmarshal(respBytes, &buildStatus)
	if err != nil {
		return buildStatus, err
	}
	return buildStatus, nil
}
