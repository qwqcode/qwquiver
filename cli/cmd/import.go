package cmd

import (
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import",
	Version: rootCmd.Version,
	Aliases: []string{"excel"},
	Short:   "数据导入工具",
	Long:    `qwquiver 数据导入工具 - 快速导入 excel 成绩数据`,
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(importCmd)

}
