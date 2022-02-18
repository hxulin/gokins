package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gokins/model"
	"gokins/util"
	"os"
	"strconv"
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
			job, ok := jobMap[id]
			if !ok {
				fmt.Println("无效的任务ID：" + jobId)
				continue
			}
			if job.Ack {
				ackText := "是否确认部署？(y/N)："
				if len(job.AckText) > 0 {
					ackText = job.AckText
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
			fmt.Println(job)
			//fmt.Println("无效的任务ID：" + jobId)
			//fmt.Println(int())
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
