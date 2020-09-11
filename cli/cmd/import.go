package cmd

import (
	"fmt"
	"strings"

	"github.com/qwqcode/qwquiver/lib/exd"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import [Excel 文件路径]",
	Version: rootCmd.Version,
	Aliases: []string{"excel"},
	Short:   "数据导入",
	Long: `qwquiver 数据导入工具 - 快速导入 excel 成绩数据

表头可选字段名：` + getOptionalFieldNames(),
	Run: func(cmd *cobra.Command, args []string) {
		// 导入 Excel
		examName := ""      // default is "", will use filename as the examName
		if len(args) <= 1 { // be effective on importing single file
			flagExamName, _ := cmd.Flags().GetString("exam-name") // read the flag “name”
			if strings.TrimSpace(flagExamName) != "" {
				examName = flagExamName
			}
		}

		examConfJSON, _ := cmd.Flags().GetString("exam-conf")

		// 导入多个文件
		for _, filename := range args {
			exd.ImportExcel(examName, filename, examConfJSON)
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(importCmd)

	flagP(importCmd, "exam-name", "n", "", "Exam 名称")
	flagP(importCmd, "exam-conf", "c", "", "Exam 配置 (JSON格式)")
}

func getOptionalFieldNames() (s string) {
	for _, fn := range utils.GetStructFields(&model.Score{}) {
		s += fmt.Sprintf("%s (%s), ", model.ScoreFieldTransMap[fn], fn)
	}
	return
}
