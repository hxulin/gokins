package util

import (
	"crypto/tls"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gokins/model"
	"gokins/rsa"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	workspace          = "~/.gokins"
	authFilePath       = workspace + "/auth.yaml"
	publicKeyFilePath  = workspace + "/public.bin"
	privateKeyFilePath = workspace + "/private.bin"
	configFilePath     = workspace + "/config.yaml"
)

// ReadAuthInfo 读取用户名、token等用户鉴权信息
func ReadAuthInfo() (model.AuthConfig, error) {
	authConfig := model.AuthConfig{}
	authFile, err := homedir.Expand(authFilePath)
	if err != nil {
		return authConfig, err
	}
	if !fileIsExist(authFile) {
		err := SaveAuthInfo(authConfig)
		if err != nil {
			return authConfig, err
		}
	}
	authBytes, err := ioutil.ReadFile(authFile)
	if err != nil {
		return authConfig, err
	}
	err = yaml.Unmarshal(authBytes, &authConfig)
	if err != nil {
		return authConfig, err
	}
	return authConfig, nil
}

// SaveAuthInfo 保存用户名、token等用户鉴权信息
func SaveAuthInfo(authConfig model.AuthConfig) error {
	authFile, err := homedir.Expand(authFilePath)
	if err != nil {
		return err
	}
	// 初始化工作空间目录
	initWorkspace()
	authBytes, _ := yaml.Marshal(authConfig)
	err = ioutil.WriteFile(authFile, authBytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Encrypt 加密文本内容
func Encrypt(plainText string) (string, error) {
	publicKey, err := readPublicKey()
	if err != nil {
		return "", err
	}
	cipherText, err := rsa.Encrypt(publicKey, plainText)
	if err != nil {
		return "", err
	}
	return cipherText, nil
}

// Decrypt 解密文本内容
func Decrypt(cipherText string) (string, error) {
	privateKey, err := readPrivateKey()
	if err != nil {
		return "", err
	}
	plainText, err := rsa.Decrypt(privateKey, cipherText)
	if err != nil {
		return "", err
	}
	return plainText, nil
}

// ReadConfigInfo 读取系统配置文件信息
func ReadConfigInfo() (model.SysConfig, error) {
	configFile, err := homedir.Expand(configFilePath)
	if err != nil {
		fmt.Println("配置文件存储位置初始化失败")
		return model.SysConfig{}, err
	}
	fileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("系统配置文件读取失败，请尝试使用load命名加载配置文件")
		return model.SysConfig{}, err
	}
	var conf model.SysConfig
	if err = yaml.Unmarshal(fileBytes, &conf); err != nil {
		fmt.Println("系统配置文件解析失败")
		return model.SysConfig{}, err
	}
	return conf, err
}

// LoadConfig 加载配置文件
func LoadConfig(url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	// 发起 http Get 请求读取远程配置文件
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	configFile, err := homedir.Expand(configFilePath)
	if err != nil {
		return err
	}
	// 初始化工作空间目录
	initWorkspace()
	err = ioutil.WriteFile(configFile, body, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 读取公钥内容
func readPublicKey() ([]byte, error) {
	keyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return nil, err
	}
	err = initSecretKey()
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// 读取私钥内容
func readPrivateKey() ([]byte, error) {
	keyFile, err := homedir.Expand(privateKeyFilePath)
	if err != nil {
		return nil, err
	}
	err = initSecretKey()
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

// 初始化密钥
func initSecretKey() error {
	publicKeyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return err
	}
	privateKeyFile, err := homedir.Expand(privateKeyFilePath)
	if err != nil {
		return err
	}
	if fileIsExist(publicKeyFile) && fileIsExist(privateKeyFile) {
		return nil
	}
	// 初始化工作空间目录
	initWorkspace()
	publicKey, privateKey, err := rsa.GenerateRsaKey(256)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(publicKeyFile, publicKey, os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(privateKeyFile, privateKey, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 初始化工作目录
func initWorkspace() {
	dir, err := homedir.Expand(workspace)
	if err == nil && !fileIsExist(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}

// 判断文件是否存在，存在返回true，不存在返回false
func fileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
