package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	primaryv1 = "/primaryv1.json"
	primaryv2 = "/primaryv2.json"
)

func AssetsDir() string {
	_, b, _, _ := runtime.Caller(2) //2 To get root dir of project
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), "/assets")
}

func assetExists(assetsDir string, directory string, name string) bool {
	path := filepath.Join(assetsDir, directory+name)
	fmt.Println(path)
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Asset does not exist:", name)
		return false
	}
	fmt.Println("Asset exists:", name)
	return true

}

//selectFileNames takes in an array and selects (writeFile,deleteFile/readFile)
func selectFileNames(data interface{}, assetsDir string, directory string) (string, string) {
	switch data.(type) {
	case []datatypes.SWOrder:
		if assetExists(assetsDir, directory, primaryv1) {
			fmt.Println("Primaryv1 exists, want to write primaryv2 and delete primary v1")
			return primaryv2, primaryv1
		}
		return primaryv1, primaryv2
	default:
		return "DefaultWrite.json", "DefaultDelete.json"
	}
}

func WriteLog(data interface{}, directory string) {
	assetsDir := AssetsDir()
	result, err := json.MarshalIndent(data, "", "")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("jsonInfo: %s\n", result)
	if _, err := os.Stat(filepath.Join(assetsDir, directory)); os.IsNotExist(err) {
		fmt.Println("Making dir: ", filepath.Join(assetsDir, directory))
		err := os.MkdirAll(filepath.Join(assetsDir, directory), 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
	writefile, deletefile := selectFileNames(data, assetsDir, directory)
	err = ioutil.WriteFile(filepath.Join(AssetsDir(), directory)+writefile, result, 0644)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(filepath.Join(AssetsDir(), directory) + deletefile)
	if err != nil {
		fmt.Println(err)
	}
}

func ReadLogQueue(data *[]datatypes.SWOrder, directory string) {
	assetsDir := AssetsDir()
	_, readFile := selectFileNames(*data, assetsDir, directory)
	file, _ := ioutil.ReadFile(filepath.Join(assetsDir, directory) + readFile)
	_ = json.Unmarshal([]byte(file), data)
}
