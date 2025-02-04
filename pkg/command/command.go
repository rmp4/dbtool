package command

import (
	"dbtool/pkg/configs"
	"dbtool/pkg/logger"
	"fmt"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	logger.InitLogger(true)
	sugar = logger.GetSugar()
}

func Backup(c configs.DatabaseConfig, output string) {
	sugar.Info("執行 pg_dump")
	cmd := exec.Command("pg_dump",
		"-h", c.Address,
		"-p", c.Port,
		"-U", c.Account,
		"-Fc", // 使用自訂格式
		"-f", output,
		c.DBName,
	)
	// 設定環境變數
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", c.Password))
	o, err := cmd.CombinedOutput()
	if err != nil {
		sugar.Errorf("備份失敗: %v\n", err)
		sugar.Error(string(o))
		return
	}
	sugar.Info("備份完成")
}
func isTextFormat(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// 讀取前幾個字節
	buffer := make([]byte, 5)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	// 檢查是否以 "-- " 開頭（SQL 文件的常見特徵）
	return string(buffer[:3]) == "-- "
}
func Restore(c configs.DatabaseConfig, output string) {
	sugar.Info("執行 pg_restore")
	// 檢查檔案格式
	var cmd *exec.Cmd
	if isTextFormat(output) {
		// 如果是文字格式，使用 psql
		cmd = exec.Command("psql",
			"-h", c.Address,
			"-p", c.Port,
			"-U", c.Account,
			"-d", c.DBName,
			"-f", output,
		)
	} else {
		// 如果是自訂格式，使用 pg_restore
		cmd = exec.Command("pg_restore",
			"--clean",
			"--if-exists",
			"--no-owner",
			"--no-acl",
			"-h", c.Address,
			"-p", c.Port,
			"-U", c.Account,
			"-d", c.DBName,
			output,
		)
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", c.Password))
	o, err := cmd.CombinedOutput()
	if err != nil {
		sugar.Errorf("還原失敗: %v\n", err)
		sugar.Error(string(o))
		return
	}
	sugar.Info("還原完成")
}

func Create(c configs.DatabaseConfig, sqlfile string) {
	sugar.Info("執行 createdb")
	cmd := exec.Command("psql",
		"-h", c.Address,
		"-p", c.Port,
		"-U", c.Account,
		"-d", c.DBName,
		"-f", sqlfile)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", c.Password))
	o, err := cmd.CombinedOutput()
	if err != nil {
		sugar.Errorf("建立資料庫失敗: %v\n", err)
		sugar.Error(string(o))
		return
	}
	sugar.Info("建立資料庫完成")
}
