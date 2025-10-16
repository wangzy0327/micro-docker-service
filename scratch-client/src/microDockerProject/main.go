package main

import (
    "microDockerProject/myhttp"
    "fmt"
    "flag"
    "net"
    "net/http"
    "os"
    "sync"
    "time"
    "io/ioutil"
    "embed"
    "strings"
    "encoding/json"
)

var input string
var output string
var taskType string
var nodeIp string
var callbackPort int // 修改为int类型，因为端口是整数
var resultChan = make(chan string)
var wg sync.WaitGroup

// 嵌入时区文件
//go:embed tzdata/Asia/Shanghai
var tzData embed.FS

const defaultTimeFilePath = "/time/start-time.txt"

func Init(){
    // 通过命令行参数读取input和output，默认值分别为input/input1和output/output1
    flag.StringVar(&input, "input", "yolov3", "project input ")
    flag.StringVar(&output, "output", "output/output1", "project output ")
    flag.StringVar(&taskType, "type", "hadoop", "Task type: micro/hadoop")
}

// 获取环境变量或默认值
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// 回调接口
func callbackHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("only POST allowed"))
        return
    }

    respBody, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("read callback data failed: " + err.Error()))
        return
    }

    type CallbackData struct {
        UUID      string `json:"uuid"`
        Success   bool   `json:"success"`
        Message   string `json:"message"`
        Timestamp string `json:"timestamp"`
    }
    var callbackData CallbackData
    if err := json.Unmarshal(respBody, &callbackData); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("parse callback JSON failed: " + err.Error()))
        return
    }

    formattedMessage := strings.ReplaceAll(callbackData.Message, "\\n", "\n")

    resultChan <- fmt.Sprintf(
        "UUID: %s\nSuccess: %v\nTimestamp: %s\nMessage:\n%s",
        callbackData.UUID,
        callbackData.Success,
        callbackData.Timestamp,
        formattedMessage,
    )

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("callback received success"))
}

// 启动客户端回调服务
func startCallbackServer() (int, error) {
    listener, err := net.Listen("tcp", ":0")
    if err != nil {
        return 0, fmt.Errorf("监听端口失败: %v", err)
    }

    port := listener.Addr().(*net.TCPAddr).Port
    fmt.Printf("[回调服务] 动态分配端口: %d\n", port)

    go func() {
        http.HandleFunc("/client/callback", callbackHandler)
        fmt.Printf("[回调服务] 启动成功，接口: /client/callback（端口: %d）\n", port)
        if err := http.Serve(listener, nil); err != nil && err != http.ErrServerClosed {
            fmt.Printf("[回调服务] 异常退出: %s\n", err.Error())
        }
    }()

    return port, nil
}

// 获取本机非环回地址的IP
func getLocalIP() (string, error) {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "", err
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String(), nil
}

// 解析时间格式
func parseScriptTime(scriptTimeStr string, loc *time.Location) (time.Time, error) {
    layout := "15:04:05.000"
    parsedTime, err := time.ParseInLocation(layout, scriptTimeStr, loc)
    if err != nil {
        return time.Time{}, fmt.Errorf("解析脚本时间失败: %v（输入格式：%s）", err, scriptTimeStr)
    }

    today := time.Now().In(loc)
    parsedTime = time.Date(
        today.Year(),
        today.Month(),
        today.Day(),
        parsedTime.Hour(),
        parsedTime.Minute(),
        parsedTime.Second(),
        parsedTime.Nanosecond(),
        loc,
    )

    return parsedTime, nil
}

// 获取开始时间
func getStartTime() string {
    if envTime := os.Getenv("START_TIME"); envTime != "" {
        fmt.Printf("[时间读取] 从环境变量获取 START_TIME: %s\n", envTime)
        return envTime
    }

    fmt.Printf("[时间读取] 环境变量为空，尝试从文件 %s 读取\n", defaultTimeFilePath)
    fileContent, err := ioutil.ReadFile(defaultTimeFilePath)
    if err != nil {
        fmt.Printf("[时间读取] 读取文件失败: %v（将返回空字符串）\n", err)
        return ""
    }

    fileTime := strings.TrimSpace(string(fileContent))
    if fileTime == "" {
        fmt.Printf("[时间读取] 文件内容为空\n")
        return ""
    }

    fmt.Printf("[时间读取] 从文件获取 START_TIME: %s\n", fileTime)
    return fileTime
}

func main(){
    data, err := tzData.ReadFile("tzdata/Asia/Shanghai")
    if err != nil {
        panic(fmt.Sprintf("读取时区文件失败: %v", err))
    }

    shanghaiLoc, err := time.LoadLocationFromTZData("Asia/Shanghai",data)
    if err != nil {
        panic(fmt.Sprintf("加载时区失败: %v", err))
    }

    scriptStartTimeStr := getStartTime()
    if scriptStartTimeStr == "" {
        fmt.Printf("[警告] 未获取到启动时间（环境变量和文件均为空），跳过耗时计算\n")
    }

    fmt.Printf("容器启动.....\n")
    programStartTime := time.Now().In(shanghaiLoc)
    programStartTimeStr := programStartTime.Format("15:04:05.000")
    fmt.Printf("容器内程序启动时间: %s\n", programStartTimeStr)

    if scriptStartTimeStr != "" {
        scriptStartTime, err := parseScriptTime(scriptStartTimeStr, shanghaiLoc)
        if err != nil {
            fmt.Printf("[警告] 计算启动耗时失败: %v\n", err)
        } else {
            duration := programStartTime.Sub(scriptStartTime)
            ms := duration.Milliseconds()
            fmt.Printf("========================================\n")
            fmt.Printf("容器启动耗时: %d 毫秒\n", ms)
            fmt.Printf("========================================\n")
        }
    }

    // 初始化命令行参数并解析
    Init()
    flag.Parse()
    // 打印当前使用的input和output路径
    //fmt.Printf("使用的输入路径: %s\n", input)
    //fmt.Printf("使用的输出路径: %s\n", output)

    callbackPort, err := startCallbackServer()
    if err != nil {
        fmt.Printf("[启动失败] 回调服务无法启动: %v\n", err)
        return
    }

    NODE_NAME := getEnv("NODE_NAME", "localhost")
    PORT := getEnv("PORT", "8800")
    var serverUrl string
    if taskType == "hadoop" {
        serverUrl = "http://" + NODE_NAME + ":" + PORT + "/hadoop"
    } else {
        serverUrl = "http://" + NODE_NAME + ":" + PORT + "/micro"
    }
    fmt.Printf("[请求地址] 服务端: %s\n", serverUrl)

    clientIP, err := getLocalIP()
    if err != nil {
        clientIP = getEnv("NODE_NAME", "localhost")
        fmt.Printf("[IP获取警告] 自动检测失败，使用默认IP: %s\n", clientIP)
    }
    fmt.Printf("[客户端信息] 真实IP: %s，回调端口: %d\n", clientIP, callbackPort)

    callbackUrl := fmt.Sprintf("http://%s:%d/client/callback", clientIP, callbackPort)
    fmt.Printf("[回调地址] 发送给服务端: %s\n", callbackUrl)

    fmt.Printf("发送请求给代理服务器......\n")
    uuidStr, err := myhttp.HttpPostWithCallback(serverUrl, input, output, callbackUrl)
    if err != nil {
        fmt.Printf("[请求失败] %s\n", err.Error())
        return
    }
    fmt.Printf("[请求成功] 服务端返回UUID: %s\n", uuidStr)

    wg.Add(1)
    go func() {
        defer wg.Done()
        select {
        case result := <-resultChan:
            fmt.Printf("\n========================================\n")
            fmt.Printf("[任务完成] UUID: %s\n", uuidStr)
            fmt.Printf("[回调结果] %s\n", result)
            fmt.Printf("========================================\n")
        case <-time.After(30 * time.Minute):
            fmt.Printf("\n[任务超时] 等待回调超过30分钟，任务可能未完成\n")
        }
    }()

    wg.Wait()
    fmt.Printf("%s 任务结束......\n", uuidStr)
}


