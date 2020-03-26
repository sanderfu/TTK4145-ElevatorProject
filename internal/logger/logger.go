package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	primaryv1 = "/primaryv1.json"
	primaryv2 = "/primaryv2.json"

	backupv1 = "/backupv1.json"
	backupv2 = "/backupv2.json"

	assetDir = "./assets/"

	permissionRW = 0755
)

func fileExists(dir string, filename string) bool {
	path := filepath.Join(dir, filename)
	fmt.Println("fileExist path: ", path)
	_, err := os.Stat(path)

	return err == nil
}

func selectFileNames(primary bool) (string, string) {
	processAssetsDir := filepath.Join(assetDir, strconv.Itoa(os.Getpid()))

	if primary {
		if fileExists(processAssetsDir, primaryv1) {
			return primaryv2, primaryv1
		} else {
			return primaryv1, primaryv2
		}
	} else {
		if fileExists(processAssetsDir, backupv1) {
			return backupv2, backupv1
		} else {
			return backupv1, backupv2
		}
	}
}

func SaveQueue(queue []datatypes.QueueOrder, primary bool) {
	processAssetsDir := filepath.Join(assetDir, strconv.Itoa(os.Getpid()))
	result, err := json.MarshalIndent(queue, "", "")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := os.Stat(processAssetsDir); os.IsNotExist(err) {
		err := os.MkdirAll(processAssetsDir, permissionRW)
		if err != nil {
			fmt.Println(err)
		}
	}
	writefile, deletefile := selectFileNames(primary)
	err = ioutil.WriteFile(processAssetsDir+writefile, result, permissionRW)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(processAssetsDir + deletefile)
	if err != nil {
		//fmt.Println(err)
	}
}

func LoadQueue(queue *[]datatypes.QueueOrder, primary bool, pid string) {
	_, readFile := selectFileNames(primary)
	file, err := ioutil.ReadFile(filepath.Join(assetDir, pid) + readFile)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	_ = json.Unmarshal([]byte(file), queue)
}
