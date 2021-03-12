package http_client

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestGet(t *testing.T) {
	client := g.Client()
	if response, err := client.Get("http://localhost:1937/get"); err != nil {
		panic(err)
	} else {
		defer response.Close()
		fmt.Println(response.ReadAllString())
	}
}

func TestPost(t *testing.T) {
	client := g.Client()

	//if response,err:=client.Post("http://localhost:1937/post_map",g.Map{
	//	"name":"小飞鱼骚骚",
	//	"age":25,
	//});err != nil {
	//	panic(err)
	//}else {
	//	defer response.Close()
	//}

	girl := &girl{
		Name:      "电气鼠小骚",
		Age:       21,
		FaceScore: 7,
		IsBeauty:  true,
	}
	marshal, _ := json.Marshal(girl)
	js := string(marshal)

	if response, err := client.Post("http://localhost:1937/post_json", js); err != nil {
		panic(err)
	} else {
		defer response.Close()
	}
}

/*
GetContent GetBytes GetVar
*/

func TestGetValue(t *testing.T) {
	client := g.Client()
	//url := "http://localhost:1937/get"
	//data := "name=电气鼠&age=20&face_score=7.0&is_beauty=true"

	//content := client.GetContent(url, data)
	//fmt.Println(content)
	//
	//bytes := client.GetBytes(url, data)
	//fmt.Println(bytes)
	//
	//girl := new(girl)
	//variable := client.GetVar(url, data)
	//_ = variable.Scan(girl)
	//fmt.Println(*girl)

	resp := client.GetVar("http://localhost:1937/getIntList")

	var ints *[]int
	_ = resp.Scan(ints)
	fmt.Println(resp)
}
