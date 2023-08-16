package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

// PathExist 检查路径是否不存在
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateNewFile 创建新的文件 文件如果存在的话 会删除旧文件创建新文件并返回
func CreateNewFile(path string) (*os.File, error) {
	isExist, err := PathExist(path)
	if err != nil {
		return nil, err
	}
	if isExist {
		if err = os.Remove(path); err != nil {
			return nil, err
		}
	}
	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	return newFile, nil
}

// CheckFileMD5 检查文件的 md5 值
func CheckFileMD5(file string, tagMd5Value string) (bool, error) {
	md5er := md5.New()
	checkFile, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer checkFile.Close()
	if _, err := io.Copy(md5er, checkFile); err != nil {
		return false, err
	}
	fileMd5Value := hex.EncodeToString(md5er.Sum(nil))
	return fileMd5Value == tagMd5Value, nil
}

func FileName(file string) string {
	return path.Base(file)[:len(path.Base(file))-len(path.Ext(file))-1]
}

func ReplaceFileDir(srcFile, replaceDir string) string {
	return fmt.Sprintf("%s/%s", replaceDir, path.Base(srcFile))
}
