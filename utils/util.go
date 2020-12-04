package utils

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"strings"
)

/**
 * @author @codenai.com
 * @date 2020/2/20
 */

func GetFileList(filePath string) []string {
	fileInfoList, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println(err.Error())
		logs.Error(err.Error())
		return []string{}
	}
	//fmt.Println(len(fileInfoList))
	filelist := make([]string, 0)
	for i := range fileInfoList {
		name := fileInfoList[i].Name()
		/*last := strings.LastIndex(name, ".")

		if name[last+1:] != "ready" {
			continue
		}*/

		filelist = append(filelist, name) //打印当前文件或目录下的文件或目录名
	}
	return filelist
}

func CheckPathIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetFileListByPath(fileStoragePath string) []string {
	fileInfoList, err := ioutil.ReadDir(fileStoragePath)
	if err != nil {
		fmt.Println(err.Error())
		logs.Error(err.Error())
		return []string{}
	}
	//fmt.Println(len(fileInfoList))
	filelist := make([]string, 0)
	for i := range fileInfoList {
		name := fileInfoList[i].Name()
		last := strings.LastIndex(name, ".")

		if name[last+1:] != "xlsx" {
			continue
		}

		filelist = append(filelist, name) //打印当前文件或目录下的文件或目录名
	}
	return filelist
}

func GetStaticFileListByPath(fileStoragePath string) []string {
	fileInfoList, err := ioutil.ReadDir(fileStoragePath)
	if err != nil {
		fmt.Println(err.Error())
		logs.Error(err.Error())
		return []string{}
	}
	//fmt.Println(len(fileInfoList))
	filelist := make([]string, 0)
	for i := range fileInfoList {
		name := fileInfoList[i].Name()
		if strings.Contains(name, "history_zhxx") {
			continue
		}
		filelist = append(filelist, name) //打印当前文件或目录下的文件或目录名
	}
	return filelist
}
