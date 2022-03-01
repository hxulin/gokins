package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gokins/job"
	"gokins/model"
	"gokins/util"
	"os"
	"strconv"
	"strings"
	"time"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "执行任务",
	Long: `执行某一个任务，后面接任务ID参数，如：
gokins deploy 1005`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(`deploy 命令后面需要指定任务ID，例如：
gokins deploy 1005`)
			return
		}
		// 读取配置文件
		sysConfig, err := util.ReadConfigInfo()
		if err != nil {
			return
		}
		jobList := sysConfig.Gokins.Job
		// 将 job 任务放到 map 中，方便按 id 查找
		jobMap := make(map[int]model.Job)
		for _, job := range jobList {
			jobMap[job.Id] = job
		}
		for _, jobId := range args {
			id, err := strconv.Atoi(jobId)
			if err != nil {
				fmt.Println("无效的任务ID：" + jobId)
				continue
			}
			// 获取 map 数据，ok 表示是否存在
			task, ok := jobMap[id]
			if !ok {
				fmt.Println("无效的任务ID：" + jobId)
				continue
			}
			if task.Ack {
				ackText := "是否确认部署？(y/N)："
				if len(task.AckText) > 0 {
					ackText = task.AckText
				}
				fmt.Print(ackText)
				// 二次确认，读取用户输入信息
				in := bufio.NewReader(os.Stdin)
				ackBytes, _, err := in.ReadLine()
				if err != nil {
					continue
				}
				// 121 表示用户输入字符 y
				if len(ackBytes) == 0 || ackBytes[0] != 121 {
					fmt.Println("用户取消部署操作，任务ID为：" + jobId)
					continue
				}
			}
			baseUrl := sysConfig.Gokins.JenkinsUrl
			if !strings.HasPrefix(baseUrl, "http://") && !strings.HasPrefix(baseUrl, "https://") {
				fmt.Println("Jenkins地址配置错误")
				return
			}
			if !strings.HasSuffix(baseUrl, "/") {
				baseUrl += "/"
			}
			// 读取用户名和token配置信息
			authConfig, err := util.ReadAuthInfo()
			if err != nil {
				fmt.Println("用户名和token信息读取失败")
				return
			}
			username, err := util.Decrypt(authConfig.Auth.Username)
			if err != nil {
				fmt.Println("用户名解密失败")
				return
			}
			token, err := util.Decrypt(authConfig.Auth.Token)
			if err != nil {
				fmt.Println("token信息解密失败")
				return
			}
			// 查询任务状态
			status, err := job.QueryLastBuildStatus(baseUrl, task.Name, username, token)
			if err != nil {
				fmt.Println("查询任务最近一次部署状态失败")
				return
			}
			statusText := job.ParseBuildStatus(status)
			if statusText == job.Building {
				buildParams := job.ParseCurrentJobBuildParam(status)
				// 判断构建参数是否相等
				if buildParamIsEquals(buildParams, task.Params) {
					fmt.Print("其他终端正在部署此任务，是否继续本次部署？(y/N)：")
					in := bufio.NewReader(os.Stdin)
					ackBytes, _, err := in.ReadLine()
					if err != nil {
						continue
					}
					// 121 表示用户输入字符 y
					if len(ackBytes) == 0 || ackBytes[0] != 121 {
						fmt.Println("用户取消部署操作，任务ID为：" + jobId)
						continue
					}
				}
			}
			queueId, err := job.Build(baseUrl, task.Name, task.Params, username, token)
			if err != nil {
				fmt.Println("✘ 任务部署失败：" + err.Error())
				continue
			}
			// 查询重试次数
			retry := 0
			number := -1
			for {
				// 查询执行队列状态信息
				queueInfo, err := job.GetQueueInfo(baseUrl, queueId, username, token)
				if err != nil {
					retry += 1
					if retry >= 5 {
						fmt.Print("\033[2K")
						fmt.Println("\r✘ 查询执行队列状态失败。")
						break
					}
					time.Sleep(time.Duration(100) * time.Millisecond)
					continue
				} else {
					retry = 0
				}
				if !queueInfo.Blocked && queueInfo.Executable.URL != "" {
					number = queueInfo.Executable.Number
					fmt.Print("\033[2K")
					fmt.Println("\r✔ 任务调度器初始化完成。")
					break
				} else {
					loading("任务已提交，等待调度器执行，请稍候...")
				}
			}
			if number == -1 {
				continue
			}
			// 重置重试次数
			retry = 0
			for {
				status, err = job.QueryBuildStatus(baseUrl, task.Name, number, username, token)
				if err != nil {
					retry += 1
					if retry >= 3 {
						fmt.Print("\033[2K")
						fmt.Println("\r✘ 部署结果查询失败，请登录web页面查看")
						break
					}
					time.Sleep(time.Duration(1) * time.Second)
					continue
				} else {
					retry = 0
				}
				statusText = job.ParseBuildStatus(status)
				if statusText == job.Building {
					loading("正在部署当前任务，请稍候...")
				} else if statusText == job.Success {
					fmt.Print("\033[2K")
					fmt.Println("\r✔ 任务ID: " + jobId + " 部署完成 -> SUCCESS")
					break
				} else {
					fmt.Print("\033[2K")
					fmt.Println("\r✘ 任务ID: " + jobId + " 部署失败 -> " + statusText)
					break
				}
			}
		}
	},
}

// 加载中提示
func loading(text string) {
	interval := time.Duration(125) * time.Millisecond
	fmt.Print("\033[2K")
	fmt.Print("\r─ " + text)
	for i := 0; i < 2; i++ {
		time.Sleep(interval)
		fmt.Print("\r\\ " + text)
		time.Sleep(interval)
		fmt.Print("\r| " + text)
		time.Sleep(interval)
		fmt.Print("\r/ " + text)
		time.Sleep(interval)
		fmt.Print("\r─ " + text)
	}
}

// 判断构建参数是否相等
func buildParamIsEquals(a, b []model.BuildParamItem) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
