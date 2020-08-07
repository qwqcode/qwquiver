package cmd

import (
	"fmt"
	"strings"

	"github.com/qwqcode/qwquiver/lib/tools"
	"github.com/qwqcode/qwquiver/lib/utils"
	"github.com/qwqcode/qwquiver/model"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import [Excel 文件路径]",
	Version: rootCmd.Version,
	Aliases: []string{"excel"},
	Short:   "数据导入工具",
	Long: `qwquiver 数据导入工具 - 快速导入 excel 成绩数据

表头可选字段名：` + getOptionalFieldNames(),
	Run: func(cmd *cobra.Command, args []string) {
		// 导入 Excel
		dataName := ""      // default is "", will use filename as the dataName
		if len(args) <= 1 { // be effective on importing single file
			flagDataName, err := cmd.Flags().GetString("name") // read the flag “name”
			if err == nil && strings.TrimSpace(flagDataName) != "" {
				dataName = flagDataName
			}
		}

		for _, filename := range args {
			tools.ImportExcel(dataName, filename)
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(importCmd)

	flagP(importCmd, "name", "n", "", "数据名")
}

func getOptionalFieldNames() (s string) {
	for _, fn := range utils.GetStructFields(&model.Score{}) {
		s += fmt.Sprintf("%s (%s), ", model.ScoreFieldTransMap[fn], fn)
	}
	return
}