package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	primaryv1 = "/primaryv1.json"
	primaryv2 = "/primaryv2.json"

	backupv1 = "/backupv1.json"
	backupv2 = "/backupv2.json"
)

func AssetsDir() string {
	_, b, _, _ := runtime.Caller(0)
	internalDir := path.Join(path.Dir(b))
	projectDir := path.Join(path.Dir(internalDir))
	return filepath.Join(filepath.Dir(projectDir), "/assets/"+strconv.Itoa(os.Getpid()))
}

func assetExists(assetsDir string, directory string, name string) bool {
	path := filepath.Join(assetsDir, directory+name)
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true

}

//selectFileNames takes in an array and selects (writeFile,deleteFile/readFile)
func selectFileNames(data interface{}, primary bool, assetsDir string, directory string) (string, string) {
	switch primary {
	case true:
		if assetExists(assetsDir, directory, primaryv1) {
			return primaryv2, primaryv1
		}
		return primaryv1, primaryv2
	case false:
		if assetExists(assetsDir, directory, backupv1) {
			return backupv2, backupv1
		}
		return backupv1, backupv2
	default:
		return "DefaultWrite.json", "DefaultDelete.json"
	}
}

func WriteLog(data interface{}, primary bool, directory string) {
	assetsDir := AssetsDir()
	result, err := json.MarshalIndent(data, "", "")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(filepath.Join(assetsDir, directory)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(assetsDir, directory), 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
	writefile, deletefile := selectFileNames(data, primary, assetsDir, directory)
	err = ioutil.WriteFile(filepath.Join(AssetsDir(), directory)+writefile, result, 0644)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(filepath.Join(AssetsDir(), directory) + deletefile)
	if err != nil {
		//fmt.Println(err)
	}
}

func ReadLogQueue(data *[]datatypes.Order, primary bool, directory string) {
	assetsDir := AssetsDir()
	_, readFile := selectFileNames(*data, primary, assetsDir, directory)
	file, _ := ioutil.ReadFile(filepath.Join(assetsDir, directory) + readFile)
	_ = json.Unmarshal([]byte(file), data)
}
