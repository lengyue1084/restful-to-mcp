package util

import (
	"os"
	"time"
)

// PathExists 返回true则不存在，返回err具体分析
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Date 格式化日期
func Date(layout string) string {
	return time.Now().Format(layout)
}

// DateByStr Date 根据格式返回可读日期字符串
func DateByStr(f string) string {
	if f == "Y-m-d H:i:s" || f == "" {
		return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	}
	if f == "Y-m-d" {
		return time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	}
	if f == "Y-m-d" {
		return time.Unix(time.Now().Unix(), 0).Format("15:04:05")
	}

	return ""
}
