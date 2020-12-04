package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
)

/**
 * @author @codenai.com
 * @date 2020/3/3
 */

func WriteSchemaFile(str, filename string) {
	//打开文件，新建文件
	f, err := os.OpenFile("./genfile/"+filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) //传递文件路径
	if err != nil {                                                                      //有错误
		fmt.Println("err = ", err)
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s", str)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
	}

}

func WriteFile(str, filename string) error {
	//打开文件，新建文件
	f, err := os.OpenFile("./genfile/"+filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) //传递文件路径
	if err != nil {                                                                      //有错误
		fmt.Println("err = ", err.Error())
		logs.Error("err = ", err.Error())
		return err
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s", str)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err)
		logs.Error("err = ", err.Error())
		return err
	}
	return nil

}

func WriteFlagFile(str, filename string) {
	//打开文件，新建文件
	f, err := os.OpenFile("./"+filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) //传递文件路径
	if err != nil {                                                              //有错误
		fmt.Println("err = ", err.Error())
		logs.Error("err = ", err.Error())
		return
	}

	//使用完毕，需要关闭文件
	defer f.Close()

	var buf string
	buf = fmt.Sprintf("%s", str)
	_, err = f.WriteString(buf)
	if err != nil {
		fmt.Println("err = ", err.Error())
		logs.Error("err = ", err.Error())
	}

}
