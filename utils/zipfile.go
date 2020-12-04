package utils

import (
	"archive/zip"
	"github.com/astaxie/beego/logs"
	"io"
	"log"
	"os"
	"strings"
)

/**
 * @author @codenai.com
 * @date 2020/4/19
 */

/**
文件解压缩到单独文件夹
*/
func DeCompress(zipFile, dest, dirPathStr, truefilename string) (err error) {
	//目标文件夹不存在则创建
	if _, err = os.Stat(dirPathStr); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dirPathStr, 0666)
		}
	}
	//fmt.Println("dest = ", dest)
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()

	//判断压缩包解压后直接能看到的文件和文件夹数量  大于1了，说明里面有多个文件夹

	//fmt.Println(len(reader.File), dest)

	for i, file := range reader.File {
		//检查解压出来的第一个文件夹是什么，如果和虚拟目录一样就不再处理
		if i == 0 && reader.File[i].Name == truefilename+"/" {
			//fmt.Println(reader.File[0].Name)
			continue
		} else if i == 0 {
			dest = dirPathStr + "/"
		}
		//fmt.Println(dest, truefilename, reader.File[0].Name)
		if file.FileInfo().IsDir() {
			//fmt.Println(dest + file.Name)
			err := os.MkdirAll(dest+file.Name, 0777)
			if err != nil {
				log.Println(err)
			}
			continue
		} else {

			err = os.MkdirAll(getDir(dest+file.Name), 0777)
			if err != nil {
				return err
			}
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filename := dest + file.Name
		//err = os.MkdirAll(getDir(filename), 0755)
		//if err != nil {
		//    return err
		//}

		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}

	}
	return
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		logs.Error("start is wrong")
		return ""
	}

	if end < start || end > length {
		logs.Error("end is wrong")
		return ""
	}

	return string(rs[start:end])
}
