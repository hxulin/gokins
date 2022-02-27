package job

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gokins/model"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Building 任务构建中
// Success 任务构建成功
// Failure 任务构建失败
// Aborted 用户终止构建
const (
	Building = "BUILDING"
	Success  = "SUCCESS"
	//Failure  = "FAILURE"
	//Aborted  = "ABORTED"
)

// QueryBuildStatus 查询任务的构建状态
func QueryBuildStatus(baseUrl, jobName, username, token string) (model.BuildStatus, error) {
	queryUrl := baseUrl + "job/" + jobName + "/lastBuild/api/json"
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

// Build 启动任务构建
func Build(baseUrl, jobName string, params []model.BuildParamItem, username, token string) (int, error) {
	buildUrl := baseUrl + "job/" + jobName + "/build"
	if len(params) > 0 {
		query := url.Values{}
		for _, param := range params {
			query.Add(param.Name, fmt.Sprintf("%v", param.Value))
		}
		buildUrl = baseUrl + "job/" + jobName + "/buildWithParameters?" + query.Encode()
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodPost, buildUrl, nil)
	if err != nil {
		return -1, err
	}
	//auth := []byte(username + ":" + token)
	//req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
	req.SetBasicAuth(username, token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	status := resp.StatusCode
	if status != http.StatusCreated {
		return -1, errors.New("failed to start job, http status: " + strconv.Itoa(status))
	}
	location := resp.Header["Location"][0]
	splitUrl := strings.Split(location, "/")
	return strconv.Atoi(splitUrl[len(splitUrl)-2])
}

// 查询任务构建状态
func queryBuildStatus(queryUrl, username, token string) (model.BuildStatus, error) {
	buildStatus := model.BuildStatus{}
	//client := http.Client{}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
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
