package config

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Instance 配置实例
var Instance *Config

// Config 配置
// @link https://godoc.org/github.com/mitchellh/mapstructure
type Config struct {
	Name      string `mapstructure:"name"`      // 网站标题
	Copyright string `mapstructure:"copyright"` // 网站版权
	Address   string `mapstructure:"address"`   // 地址
	Port      int    `mapstructure:"port"`      // 端口

	DbFile  string `mapstructure:"dbFile"`  // 数据库文件路径
	LogFile string `mapstructure:"logFile"` // 日志文件路径
}

// Init 初始化配置
func Init(cfgFile string) {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logrus.Error(err)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/qwquiver/")
		viper.SetConfigName(".qwquiver")
	}

	viper.SetEnvPrefix("QWQU")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	Instance = &Config{}
	err := viper.Unmarshal(&Instance)
	if err != nil {
		logrus.Errorf("unable to decode into struct, %v", err)
	}
}

var (
	// CMDInjections 将在 Cmd 初始化时被执行（用于功能注入）
	CMDInjections = [](func(rootCmd *cobra.Command)){}
	// HTTPInjections 将在 Http 初始化时被执行（用于功能注入）
	HTTPInjections = [](func(e *echo.Echo, api *echo.Group)){}
)
