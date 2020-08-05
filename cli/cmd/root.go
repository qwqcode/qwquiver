package cmd

import (
	"os"

	"github.com/qwqcode/qwquiver/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "qwquiver",
		Short: "A web-based exam results explorer.",
		Long: `      
    ____ __      ______ ___  __(_)   _____  _____
   / __  / | /| / / __  / / / / / | / / _ \/ ___/
  / /_/ /| |/ |/ / /_/ / /_/ / /| |/ /  __/ /    
  \__, / |__/|__/\__, /\__,_/_/ |___/\___/_/     
    /_/            /_/                           

A website for exploring and analyzing exam results.

More detail on https://github.com/qwqcode/qwquiver

(c) 2020 qwqaq.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

// Execute is execute cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)

	rootCmd.SetVersionTemplate("qwquiver {{printf \"version %s\" .Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件 (defaults are './.qwquiver', '$HOME/.qwquiver' or '/etc/qwquiver/.qwquiver')")
}

func initConfig() {
	config.Init(cfgFile)
}

func initLog() {
	// 初始化日志
	if file, err := os.OpenFile(config.Instance.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Error(err)
	}
}

//// 捷径函数 ////

func flag(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.Bool(name, y, usage)
	case int:
		f.Int(name, y, usage)
	case string:
		f.String(name, y, usage)
	}
	v.SetDefault(name, defaultVal)
}

func flagP(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	f := cmd.PersistentFlags()
	switch y := defaultVal.(type) {
	case bool:
		f.BoolP(name, shorthand, y, usage)
	case int:
		f.IntP(name, shorthand, y, usage)
	case string:
		f.StringP(name, shorthand, y, usage)
	}
	v.SetDefault(name, defaultVal)
}

func flagV(cmd *cobra.Command, name string, defaultVal interface{}, usage string) {
	flag(cmd, name, defaultVal, usage)
	v.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}

func flagPV(cmd *cobra.Command, name, shorthand string, defaultVal interface{}, usage string) {
	flagP(cmd, name, shorthand, defaultVal, usage)
	v.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
}
