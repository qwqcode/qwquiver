package cmd

import (
	"github.com/qwqcode/qwquiver/app"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Version: rootCmd.Version,
	Aliases: []string{"server"},
	Short:   "服务器",
	Long:    rootCmd.Long,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunIris()
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Global configs
	flagPV(serveCmd, "name", "n", "qwquiver", "网站标题")
	flagPV(serveCmd, "address", "a", "", "网站地址 (例如: qwqaq.com)")
	flagPV(serveCmd, "port", "p", 8087, "网站端口")
	flagPV(serveCmd, "dbFile", "d", "./data/qwquiver.db", "数据文件")
	flagPV(serveCmd, "logFile", "l", "./data/qwquiver.log", "日志文件")
}
