package main

import (
	"dbtool/pkg/command"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "備份 PostgreSQL 資料庫",
	Run: func(cmd *cobra.Command, args []string) {
		localConfig := getDBConfig(name)
		command.Backup(localConfig, filename)
	},
}
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "還原 PostgreSQL 資料庫",
	Run: func(cmd *cobra.Command, args []string) {
		localConfig := getDBConfig(name)
		command.Restore(localConfig, filename)
	},
}

var buildCmd = &cobra.Command{
	Use:   "create",
	Short: "建立 PostgreSQL 資料庫",
	Run: func(cmd *cobra.Command, args []string) {
		localConfig := getDBConfig(name)
		dirPath := dir
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 檢查是否為 .sql 檔案
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".sql") {
				command.Create(localConfig, path)
			}
			return nil
		})
		if err != nil {
			sugar.Error(err)
		}
	},
}
