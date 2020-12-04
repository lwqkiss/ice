package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//将PolyGon 和lineStr 字符串 转换为 二维切片
func ConvertPolyGonLineStr2Arr(str string) [][]float64 {
	result := [][]float64{}
	//pxx((1 3, 1 3,),())
	firstBrack := strings.Index(str, "(")
	lastBrack := strings.LastIndex(str, ")")
	str = str[firstBrack+1 : lastBrack]
	//secodeTrim
	//(),()
	outArr := strings.Split(str, "),(")

	for _, oe := range outArr {
		oe = strings.Trim(oe, string('('))
		oe = strings.Trim(oe, ")")
		fmt.Println("oe..", oe)
		// 1 2,1 3 ,1 2
		marr := strings.Split(oe, ",")
		for _, ie := range marr {
			//1 2
			numArr := strings.Split(ie, " ")
			fa := []float64{}
			for _, ele := range numArr {
				nu, err := strconv.ParseFloat(ele, 64)
				if err != nil {
					log.Println(err, ele)
				}
				fa = append(fa, nu)
			}
			result = append(result, fa)
		}
	}
	return result

}

//point字符串 转为切片
func ConverPoint2Arr(ss string) []float64 {
	var str = `Point(1 2)`
	str = str[strings.Index(str, "(")+1 : strings.Index(str, ")")]
	arr := make([]float64, 2)
	strArr := strings.Split(str, " ")
	for i, ele := range strArr {
		num, err := strconv.ParseFloat(ele, 64)
		if err != nil {
			log.Println(err)
		}
		arr[i] = num
	}
	return arr
}

//切片转 GeomFromText(Point)
func ParsePoint(arr []float64) (string, error) {
	if len(arr) != 2 {
		return "", errors.New("长度不符合,PrasePoint 失败")
	}
	return fmt.Sprintf("GeomFromText('Point(%f %f)')", arr[0], arr[1]), nil
}

//todo
/*func ParsePolyGonFromArr(arr [][]float64)(string,error){

}*/

func GenGeostr(lon, lat float64) string {
	return "(" + strconv.FormatFloat(lon, 'f', -1, 64) + " " + strconv.FormatFloat(lat, 'f', -1, 64) + ")"
}

func DealPoint(str string) (lon, lat float64, geo string) {
	sl := strings.Split(str, ",")

	lon, err := strconv.ParseFloat(sl[0], 64)
	if err != nil {
		fmt.Println(err)
	}
	lat, err = strconv.ParseFloat(sl[1], 64)
	if err != nil {
		fmt.Println(err)
	}
	geo = "(" + strconv.FormatFloat(lon, 'f', -1, 64) + " " + strconv.FormatFloat(lat, 'f', -1, 64) +
		")"
	return lon, lat, geo

}
