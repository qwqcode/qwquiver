package cmd

import (
	"fmt"

	"github.com/qwqcode/qwquiver/http"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Version: rootCmd.Version,
	Aliases: []string{"server"},
	Short:   "启动服务器",
	Long:    rootCmd.Long,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Banner)
		fmt.Println("-------------------")
		http.Run()
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Global configs
	flagPV(serveCmd, "name", "n", "qwquiver", "网站标题")
	flagPV(serveCmd, "address", "a", "", "网站地址 (例如: qwqaq.com)")
	flagPV(serveCmd, "port", "p", 8087, "网站端口")
}
