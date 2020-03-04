package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"
)

type JsonPostSample struct{

}


func HttpPost(url string,input string,output string) (string){

	data:=make(map[string]interface{})
	data["input"] = input
	data["output"] = output
	bytesData,err := json.Marshal(data)

	if err != nil{
		fmt.Println(err.Error())
		return err.Error()
	}

	reader:= bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST",url, reader)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	request.Header.Set("Content-Type","application/json;charset=UTF-8")
	client:=http.Client{}
	resp,err := client.Do(request)
	if err != nil{
		fmt.Println(err.Error())
		return err.Error()
	}
	respBytes,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err.Error())
		return err.Error()
	}

	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
        
        return *str

}

