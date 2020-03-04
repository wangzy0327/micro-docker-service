package main

import (
    "microDockerProject/http"
    "microDockerProject/polling"
    "fmt"
	"flag"
	"os"
    "time"
)


var input string
var output string
var nodeIp string

func Init(){
	//flag.StringVar(&nodeIp,"nodeIp","10.18.127.4","cluster node ip")
	flag.StringVar(&input,"input","input/input1","project input")
	flag.StringVar(&output,"output","output/output1","project output")
}

func main(){
    Init()
    flag.Parse()
	timeUnixNano := time.Now().UnixNano()
	timeUnixMirco := float64(timeUnixNano)/1000000000
	//fmt.Printf("纳秒时间为：%d\n",timeUnixNano)
	fmt.Printf("容器启动后时间戳为：%.3f\n",timeUnixMirco)
	//var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	fmt.Printf("容器启动后时间为: %s\n",time.Now().Format("2006-01-02 15:04:05"))
        var NODE_NAME string
	NODE_NAME = os.Getenv("NODE_NAME")
	var host string = "http://"+NODE_NAME+":8081/micro"
	fmt.Printf(host+"  \n")
        var SUBSCRIBE_PATH string
        SUBSCRIBE_PATH = os.Getenv("SUBSCRIBE_PATH")
        fmt.Print("subscribe path: "+SUBSCRIBE_PATH+" \n")
	uuidStr := http.HttpPost(host,input,output)
	//var output = "/mnt/xfs/pipeline_server/output/workflowArr.log"
	//fmt.Printf("发送请求后时间戳(调整后)为:%d\n",time.Now().UTC().Unix()-ZONE_DIFF)
	fmt.Printf("发送请求后时间戳为:%.3f\n",float64(time.Now().UnixNano())/1000000000)
	//fmt.Println(time.Now().UTC().Format(polling.TIME_LAYOUT))
	//fmt.Println(time.Now().UTC().Unix())
        fmt.Println("uuid str : "+uuidStr)
	polling.Track(uuidStr, output, SUBSCRIBE_PATH)
	fmt.Printf("任务结束后时间戳为:%.3f\n",float64(time.Now().UnixNano())/1000000000)
	fmt.Printf("任务结束后时间为: %s\n",time.Now().Format("2006-01-02 15:04:05"))
}



