package egis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type EgisAddress struct {
	Msg    string `json:"msg"`
	Result Result `json:"result"`
	Status int    `json:"status"`
}
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Result struct {
	Comprehension int      `json:"comprehension"`
	Level         string   `json:"level"`
	Location      Location `json:"location"`
}

// {"msg":"OK","result":{"coords":[{"lat":32.10638487479034,"lng":121.4220299717268}]},"status":0}
func GetLonAndLat(address string) EgisAddress {
	// egisUrl := "https://zhxx.ntyjgl.cn:7777/autoroutenoauth/zbzh/api/v1/transform?" + "_type=geocoding&address="

	egisUrl := "http://ntyj.codenai.com:9001/api/v1/transform?" + "_type=geocoding&address="
	qryStr := address
	egisUrl = egisUrl + url.QueryEscape(qryStr)
	// fmt.Println(egisUrl)
	// client := http.Client{}
	resp, err := http.Get(egisUrl)
	if err != nil {
		fmt.Println(err)
	}
	db, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return EgisAddress{}
	}
	defer resp.Body.Close()

	data := EgisAddress{}

	err = json.Unmarshal(db, &data)
	fmt.Println(data)
	return data

}
