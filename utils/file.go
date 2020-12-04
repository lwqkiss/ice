package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

/**
 * @author @codenai.com
 * @date 2020/2/20
 */

func WriteFile(content string) {
	//打开文件，新建文件
	f, err := os.OpenFile("./create.sql", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666) //传递文件路径
	if err != nil {                                                                //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s\n", content)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

func ReadFileLine() (map[string]int, []string) {
	//打开文件
	f, err := os.Open("./history")
	if err != nil {
		fmt.Println("err = ", err)
		return nil, nil
	}

	//关闭文件
	defer f.Close()

	//新建一个缓冲区，把内容先放在缓冲区
	r := bufio.NewReader(f)
	historymap := make(map[string]int, 0)
	historylist := make([]string, 0)
	for {
		//遇到'\n'结束读取, 但是'\n'也读取进入
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF { //文件已经结束
				ss := strings.Trim(string(buf), "\n")
				if ss != "" {
					historymap[ss] = 1
					historylist = append(historylist, ss)
				}
				break
			}
			//fmt.Println("err = ", err)
		}
		ss := strings.Trim(string(buf), "\n")
		if ss != "" {
			historymap[ss] = 1
			historylist = append(historylist, ss)
		}

		//fmt.Printf("buf = #%s#\n", ss)
	}
	return historymap, historylist

}

func ClearFile() {
	//打开文件，新建文件
	f, err := os.OpenFile("./history", os.O_RDWR|os.O_TRUNC, 0666) //传递文件路径
	if err != nil {                                                //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = ""
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

func WriteErrFile(str string) {
	//打开文件，新建文件
	f, err := os.OpenFile("./err.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666) //传递文件路径
	if err != nil {                                                             //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s\n", str)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

//一般这里不会出现什么错误，记录历史文件的
func WriteHistoryFile(filename string, fileStoragePath string, historyFileName string) {
	//打开文件，新建文件
	f, err := os.OpenFile(path.Join(fileStoragePath, historyFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //传递文件路径
	if err != nil {                                                                                             //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s\n", filename)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

func ClearHistoryFile(path string) {
	//打开文件，新建文件
	f, err := os.OpenFile("./"+path+"/history", os.O_RDWR|os.O_TRUNC, 0666) //传递文件路径
	if err != nil {                                                         //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = ""
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

//@param historyFileName : 历史文件的文件名
func ReadHistoryFile(filePath, historyFileName string) (map[string]int, []string) {

	if !CheckPathIsExist(filePath) {
		os.MkdirAll(filePath, 0666)
	}

	//打开文件,不存在就创建，存在打开，并追加
	f, err := os.OpenFile(path.Join(filePath, historyFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err = ", err)
		return nil, nil
	}

	//关闭文件
	defer f.Close()

	//新建一个缓冲区，把内容先放在缓冲区
	r := bufio.NewReader(f)
	historymap := make(map[string]int, 0)
	historylist := make([]string, 0)
	for {
		//遇到'\n'结束读取, 但是'\n'也读取进入
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF { //文件已经结束
				ss := strings.Trim(string(buf), "\n")
				if ss != "" {
					historymap[ss] = 1
					historylist = append(historylist, ss)
				}
				break
			}
			//fmt.Println("err = ", err)
		}

		ss := strings.Trim(string(buf), "\n")
		if ss != "" {
			historymap[ss] = 1
			historylist = append(historylist, ss)
		}

		//fmt.Printf("buf = #%s#\n", ss)
	}
	return historymap, historylist

}
