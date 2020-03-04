package file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
//	"strconv"
	"strings"
)

const BUFFER_SIZE int64 = 2048
//const sublog = "/home/wzy/hadoop-cluster-docker/subscribe/subscribe.txt"
func ReadLastBlock(lastNumber int,sublog string)[] string{
        file, err := os.Open(sublog)
        //fmt.Println("hadoop subscribe path : "+sublog)
	if err != nil {
		log.Println(err)
		return nil
	}
	fileInfo, _ := file.Stat()
	//var offset = fileInfo.Size()%BUFFER_SIZE
	//fmt.Println("offset : "+strconv.FormatInt(offset,10))
	buf := bufio.NewReader(file)
	//i := int64(0)
	//readLines := int64(0)
	buffer:=make([]byte,BUFFER_SIZE)
	//i := fileInfo.Size()/BUFFER_SIZE
	var n int
	file.Seek(fileInfo.Size() - BUFFER_SIZE,io.SeekStart)
        n, err = buf.Read(buffer)
	//fmt.Println("本次读取字节数 n : "+strconv.Itoa(n))
	if err == io.EOF {
		fmt.Println(" 文件为空！")
	}
	fmt.Println("reading ...")
	strs := strings.Split(strings.TrimSpace(string(buffer[:n])), "\n")
        var start int
        var temp = (len(strs) - lastNumber)
	if (temp > 0) {
		start = temp
	} else {
		start = 0
	}
        var length = len(strs)
        /*
	fmt.Println("start : "+strconv.Itoa(start))
        fmt.Println("end : "+strconv.Itoa(length))
         
        for j := start; j < len(strs); j++ {
		for _,ch:=range []rune(strs[j]){
			fmt.Printf("%c",  ch)
		}
		fmt.Println()
	}
	*/
	file.Close()
       return strs[start:length]
}
