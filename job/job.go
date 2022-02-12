package job

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gokins/model"
	"gokins/util"
	"net/http"
	"strings"
)

func QueryDeployStatus(queryUrl string, username string, token string) error {

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
	req, err := http.NewRequest("POST", buildUrl, strings.NewReader("json="+util.EncodeURIComponent(string(paramStr))))
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
