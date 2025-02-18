package main

import (
	"dbtool/pkg/configs"
	"dbtool/pkg/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile  string
	name     string
	dir      string
	filename string
	sugar    *zap.SugaredLogger
)
var rootCmd = &cobra.Command{
	Use:   "dbtool",
	Short: "An example app reading configs",
}

// Initial logger and cobra command
func init() {
	logger.InitLogger(true)
	sugar = logger.GetSugar()
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "configs/configs.yaml", "config file (default is configs/configs.yaml)")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "local", "config name")
	buildCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "sql", "sql file directory")
	backupCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "backup.dump", "backup file name")
	restoreCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "backup.dump", "restore file name")
}
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("configs") // 添加配置檔案所在的路徑
		viper.SetConfigName("configs") // 設定配置檔案名稱(無需副檔名)
	}

	viper.AutomaticEnv() // 讀取匹配的環境變數

	if err := viper.ReadInConfig(); err == nil {
		sugar.Debugf("Using config file: %s\n\r", viper.ConfigFileUsed())
	} else {
		sugar.Debugf("Error reading config file: %s \n", err)
	}
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func getDBConfig(environment string) configs.DatabaseConfig {
	return configs.DatabaseConfig{
		Address:  viper.GetString(fmt.Sprintf("%s.address", environment)),
		Account:  viper.GetString(fmt.Sprintf("%s.account", environment)),
		Password: viper.GetString(fmt.Sprintf("%s.password", environment)),
		DBName:   viper.GetString(fmt.Sprintf("%s.dbname", environment)),
		Port:     viper.GetString(fmt.Sprintf("%s.port", environment)),
	}
}
