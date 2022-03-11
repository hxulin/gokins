package cmd

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"gokins/util"
	"os"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "查看任务列表",
	Long:  "查看任务列表信息",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取配置文件
		sysConfig, err := util.ReadConfigInfo()
		if err != nil {
			return
		}
		jobList := sysConfig.Gokins.Job
		var tableHeader table.Row
		var tableRows []table.Row
		for i, job := range jobList {
			row := table.Row{}
			for _, col := range job.Columns {
				if i == 0 {
					tableHeader = append(tableHeader, col.Name)
				}
				row = append(row, col.Value)
			}
			tableRows = append(tableRows, row)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(tableHeader)
		t.AppendRows(tableRows)
		t.SetStyle(table.StyleLight)
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
