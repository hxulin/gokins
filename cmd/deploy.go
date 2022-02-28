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
			status, err := job.QueryBuildStatus(baseUrl, task.Name, username, token)
			if err != nil {
				fmt.Println("查询任务状态失败")
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
				fmt.Println("任务部署失败：" + err.Error())
				continue
			}
			// 是否排队等待过、是否需要Println
			queue, ln := false, false
			for {
				status, err = job.QueryBuildStatus(baseUrl, task.Name, username, token)
				if status.QueueId < queueId {
					queue = true
					loading("有其他任务正在部署，排队等待中，请稍候...")
					ln = true
				} else if status.QueueId == queueId {
					if queue {
						queue = false
						fmt.Print("\033[2K")
						fmt.Println("\r✔ 其他任务部署完成。")
						ln = false
					}
					statusText = job.ParseBuildStatus(status)
					if statusText == job.Building {
						loading("正在部署当前任务，请稍候...")
						ln = true
					} else if statusText == job.Success {
						fmt.Print("\033[2K")
						fmt.Println("\r✔ 当前任务部署完成 -> SUCCESS")
						break
					} else {
						fmt.Print("\033[2K")
						fmt.Println("\r✘ 当前任务部署失败 -> " + statusText)
						break
					}
				} else if status.QueueId > queueId {
					if ln {
						fmt.Println()
					}
					fmt.Println("✘ 部署结果查询失败，请登录Web页面查看。")
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
