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
)

func AssetsDir() string {
	_, b, _, _ := runtime.Caller(0)
	internalDir := path.Join(path.Dir(b))
	projectDir := path.Join(path.Dir(internalDir))
	return filepath.Join(filepath.Dir(projectDir), "/assets/"+strconv.Itoa(os.Getpid()))
}

func assetExists(assetsDir string, directory string, name string) bool {
	path := filepath.Join(assetsDir, directory+name)
	fmt.Println(path)
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true

}

//selectFileNames takes in an array and selects (writeFile,deleteFile/readFile)
func selectFileNames(data interface{}, assetsDir string, directory string) (string, string) {
	fmt.Printf("%T\n", data)
	switch data.(type) {
	case []datatypes.PrimaryOrder:
		if assetExists(assetsDir, directory, primaryv1) {
			fmt.Println("Primaryv1 exists, want to write primaryv2 and delete primary v1")
			fmt.Println("Writefile filename: ", primaryv2)
			return primaryv2, primaryv1
		}
		fmt.Println("Writefile filename: ", primaryv1)
		return primaryv1, primaryv2
	default:
		fmt.Println("Running default case")
		return "DefaultWrite.json", "DefaultDelete.json"
	}
}

func WriteLog(data interface{}, directory string) {
	assetsDir := AssetsDir()
	fmt.Println("Assets dir: ", assetsDir)
	result, err := json.MarshalIndent(data, "", "")
	if err != nil {
		fmt.Println(err)
	}
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
