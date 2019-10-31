package starGo

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// 文件夹是否存在(obsolete)
func IsDirExists(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return file.IsDir()
	}
}

// 文件是否存在
func IsFileExists(path string) (bool, error) {
	file, err := os.Stat(path)
	if err == nil {
		return file.IsDir() == false, nil
	} else {
		if os.IsNotExist(err) {
			return false, nil
		}
	}

	return true, err
}

// 读取json文件
func ReadJsonFile(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
