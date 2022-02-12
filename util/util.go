package util

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gokins/model"
	"gokins/rsa"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const (
	workspace          = "~/.gokins"
	publicKeyFilePath  = workspace + "/public.bin"
	privateKeyFilePath = workspace + "/private.bin"
	configFilePath     = workspace + "/config.yaml"
)

// FileIsExist 判断文件是否存在，存在返回true，不存在返回false
func FileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// KeyFileIsExist 判断密钥文件是否存在，存在返回true，不存在返回false
func KeyFileIsExist() (bool, error) {
	publicKeyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return false, err
	}
	privateKeyFile, err := homedir.Expand(privateKeyFilePath)
	if err != nil {
		return false, err
	}
	return FileIsExist(publicKeyFile) && FileIsExist(privateKeyFile), nil
}

// InitWorkspace 初始化工作目录
func InitWorkspace() {
	dir, err := homedir.Expand(workspace)
	if err == nil && !FileIsExist(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}

// CreateSecretKey 创建密钥并保存到文件
func CreateSecretKey() {
	publicKeyFile, err1 := homedir.Expand(publicKeyFilePath)
	privateKeyFile, err2 := homedir.Expand(privateKeyFilePath)
	if err1 != nil || err2 != nil {
		fmt.Println("密钥文件存储位置初始化失败")
		return
	}
	// 初始化工作空间目录
	InitWorkspace()
	publicKey, privateKey, err := rsa.GenerateRsaKey(256)
	if err != nil {
		fmt.Println("密钥生成失败")
		return
	}
	err1 = ioutil.WriteFile(publicKeyFile, publicKey, os.ModePerm)
	err2 = ioutil.WriteFile(privateKeyFile, privateKey, os.ModePerm)
	if err1 != nil || err2 != nil {
		fmt.Println("密钥文件创建失败")
	} else {
		fmt.Println("密钥初始化完成，已保存到文件：")
		fmt.Println(publicKeyFile)
		fmt.Println(privateKeyFile)
	}
}

// ReadPublicKey 读取公钥内容
func ReadPublicKey() (publicKey []byte, err error) {
	keyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyFile)
	return keyBytes, err
}

// EncodeURIComponent URL 参数编码，实现和 JS 通用
func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	return strings.Replace(r, "+", "%20", -1)
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
		fmt.Println("系统配置文件读取失败")
		return model.SysConfig{}, err
	}
	var conf model.SysConfig
	if err = yaml.Unmarshal(fileBytes, &conf); err != nil {
		fmt.Println("系统配置文件解析失败")
		return model.SysConfig{}, err
	}
	return conf, err
}

// DecodeConfig 解密配置信息
func DecodeConfig(ptr interface{}) {

	fmt.Println(ptr)

	// 获取结构体实例的反射类型对象
	v := reflect.ValueOf(ptr)
	// 必须是指针类型
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	t := v.Elem()
	// 遍历结构体所有成员
	for i := 0; i < t.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		field := t.Field(i)
		//val := v.Field(i).Interface()

		//fmt.Println(field)
		//fmt.Println(val)

		if field.Type().Kind() == reflect.Struct {
			fmt.Printf("%T", v)
			DecodeConfig(v)
		} else if field.Type().Kind() == reflect.Slice {
			//rv := reflect.ValueOf(val)
			//for j := 0; j < rv.Len(); j++ {
			//	DecodeConfig(rv.Index(j).Interface())
			//}
		} else if field.Type().Kind() == reflect.String {

			//fmt.Println(field.Name)
			//reflect.ValueOf(ptr).Elem().FieldByName(field.Name).SetString("7")

			//reflect.ValueOf(val).FieldByName(field.Name).SetString("7")

			//reflect.ValueOf(&ptr).Elem().FieldByName(field.Name).SetString("7")

			//fmt.Println(field.Name)
			//
			//t.Elem().FieldByName().Set(reflect.ValueOf(name))

			//fmt.Println(field, " : ", val)
			//fmt.Println("---------")
			//fieldValue := reflect.ValueOf(ptr).FieldByName(field.Name)

			//fmt.Printf("name: %v , %v\n", field.Name, fieldValue)

		}
	}
}
