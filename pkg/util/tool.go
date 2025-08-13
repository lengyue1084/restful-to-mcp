package util

import (
	"os"
	"time"
)

// FileExists 检查文件是否存在
func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DirExists 检查目录是否存在
func DirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
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
