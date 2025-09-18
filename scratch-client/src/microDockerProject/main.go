package main

import (
    "microDockerProject/myhttp"
    // "microDockerProject/polling"
    "fmt"
    "flag"
	"net"       // 新增：网络操作（动态端口、IP获取）
	"net/http"
	"os"
	// "strconv"   // 新增：整数转字符串（端口处理
	"sync"
    "time"
	"io/ioutil" // 新增这行
    "embed"  //嵌入时区
)


var input string
var output string
var nodeIp string
// var callbackPort = "8082" // 客户端回调服务端口（需确保与服务端网络互通）
var callbackPort string // 客户端回调服务端口（需确保与服务端网络互通）
var resultChan = make(chan string) // 接收服务端回调结果的通道
var wg sync.WaitGroup              // 等待回调结果，避免程序提前退出

// 嵌入时区文件（只包含需要的时区）
//go:embed tzdata/Asia/Shanghai
var tzData embed.FS


func Init(){
	//flag.StringVar(&nodeIp,"nodeIp","10.18.127.4","cluster node ip")
	flag.StringVar(&input,"input","input/input1","project input")
	flag.StringVar(&output,"output","output/output1","project output")
}

// 获取环境变量或默认值
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// 回调接口：服务端处理完成后会调用此接口返回结果
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// 只允许POST请求
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only POST allowed"))
		return
	}

	// 读取服务端回调的JSON数据
	// 关键修复：用 ioutil.ReadAll 替换 http.ReadAll
	// respBody, err := http.ReadAll(r.Body)
	respBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("read callback data failed: " + err.Error()))
		return
	}

	// 将结果发送到通道，通知主逻辑
	resultChan <- string(respBody)
	// 响应服务端：回调接收成功
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("callback received success"))
}

// 启动客户端回调服务（独立goroutine，不阻塞主逻辑）
func startCallbackServer() (int, error) {
	// 用端口0让系统分配空闲端口
    listener, err := net.Listen("tcp", ":0")
    if err != nil {
        return 0, fmt.Errorf("监听端口失败: %v", err)
    }

    // 获取实际分配的端口（从listener中提取）
    port := listener.Addr().(*net.TCPAddr).Port
    fmt.Printf("[回调服务] 动态分配端口: %d\n", port)

    // 启动HTTP服务（使用已绑定的listener，避免二次端口冲突）
    go func() {
        http.HandleFunc("/client/callback", callbackHandler)
        fmt.Printf("[回调服务] 启动成功，接口: /client/callback（端口: %d）\n", port)
        // 使用http.Serve而非ListenAndServe，避免重复绑定端口
        if err := http.Serve(listener, nil); err != nil && err != http.ErrServerClosed {
            fmt.Printf("[回调服务] 异常退出: %s\n", err.Error())
        }
    }()

    return port, nil
}

// 获取本机非环回地址的IP（适用于多网卡场景）
func getLocalIP() (string, error) {
    // 连接外部地址（不实际建立连接，仅用于获取本地出口IP）
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "", err
    }
    defer conn.Close()

    // 从本地地址中提取IP
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String(), nil
}

// 关键函数：解析 "HH:MM:SS.mmm" 格式时间为 time.Time（带时区）
func parseScriptTime(scriptTimeStr string, loc *time.Location) (time.Time, error) {
    // 解析格式：HH:MM:SS.mmm（mmm为毫秒，3位）
    layout := "15:04:05.000"
    parsedTime, err := time.ParseInLocation(layout, scriptTimeStr, loc)
    if err != nil {
        return time.Time{}, fmt.Errorf("解析脚本时间失败: %v（输入格式：%s）", err, scriptTimeStr)
    }

    // 补充年份、月份、日期（与当前日期一致，因为脚本时间是当天的）
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

func main(){
    // 读取嵌入的时区数据
    data, err := tzData.ReadFile("tzdata/Asia/Shanghai")
    if err != nil {
        panic(fmt.Sprintf("读取时区文件失败: %v", err))
    }

    shanghaiLoc, err := time.LoadLocationFromTZData("Asia/Shanghai",data)
    if err != nil {
        panic(fmt.Sprintf("加载时区失败: %v", err))
    }

    //2. 读取脚本传递的 START_TIME 环境变量
    scriptStartTimeStr := getEnv("START_TIME", "")
    if scriptStartTimeStr == "" {
        fmt.Printf("[警告] 未收到脚本传递的启动时间，跳过耗时计算\n")
    }

    // 3. 打印启动信息
    fmt.Printf("容器启动.....\n")
    programStartTime := time.Now().In(shanghaiLoc)  // 程序启动完成时间（容器内）
    programStartTimeStr := programStartTime.Format("15:04:05.000")
    fmt.Printf("容器内程序启动时间: %s\n", programStartTimeStr)

    // 如果收到脚本时间，计算耗时
    if scriptStartTimeStr != "" {
        scriptStartTime, err := parseScriptTime(scriptStartTimeStr, shanghaiLoc)
        if err != nil {
            fmt.Printf("[警告] 计算启动耗时失败: %v\n", err)
        } else {
            // 计算时间差（毫秒级）
            duration := programStartTime.Sub(scriptStartTime)
            ms := duration.Milliseconds()  // 转换为毫秒
            fmt.Printf("========================================\n")
            fmt.Printf("容器启动耗时: %d 毫秒\n", ms)
            fmt.Printf("========================================\n")
        }
    }

	// 1. 初始化命令行参数
    Init()
    flag.Parse()

    // 2. 启动回调服务（仅调用一次！动态分配端口）
    callbackPort, err := startCallbackServer()
    if err != nil {
        fmt.Printf("[启动失败] 回调服务无法启动: %v\n", err)
        return
    }


    // 3. 构建服务端请求地址（保留原环境变量逻辑）
    NODE_NAME := getEnv("NODE_NAME", "localhost")
    PORT := getEnv("PORT", "8800")
    serverUrl := "http://" + NODE_NAME + ":" + PORT + "/micro"
    fmt.Printf("[请求地址] 服务端: %s\n", serverUrl)

    // 4. 获取客户端真实IP（自动检测，失败时降级）
    clientIP, err := getLocalIP()
    if err != nil {
        clientIP = getEnv("NODE_NAME", "localhost")
        fmt.Printf("[IP获取警告] 自动检测失败，使用默认IP: %s\n", clientIP)
    }
    fmt.Printf("[客户端信息] 真实IP: %s，回调端口: %d\n", clientIP, callbackPort)

    // 5. 构建唯一回调地址（IP+动态端口，确保50个进程不冲突）
    callbackUrl := fmt.Sprintf("http://%s:%d/client/callback", clientIP, callbackPort)
    fmt.Printf("[回调地址] 发送给服务端: %s\n", callbackUrl)

    // 6. 发送任务请求到服务端
    fmt.Printf("发送请求给代理服务器......\n")
    uuidStr, err := myhttp.HttpPostWithCallback(serverUrl, input, output, callbackUrl)
    if err != nil {
        fmt.Printf("[请求失败] %s\n", err.Error())
        return
    }
    fmt.Printf("[请求成功] 服务端返回UUID: %s\n", uuidStr)

    // 7. 等待服务端回调结果（30分钟超时）
    wg.Add(1)
    go func() {
        defer wg.Done()
        select {
        case result := <-resultChan:
            // 打印回调结果（格式化显示，便于查看）
            fmt.Printf("\n========================================\n")
            fmt.Printf("[任务完成] UUID: %s\n", uuidStr)
            fmt.Printf("[回调结果] %s\n", result)
            fmt.Printf("========================================\n")
        case <-time.After(30 * time.Minute):
            fmt.Printf("\n[任务超时] 等待回调超过30分钟，任务可能未完成\n")
        }
    }()

    // 8. 等待回调逻辑结束（避免主进程提前退出）
    wg.Wait()
    fmt.Printf("%s 任务结束......\n", uuidStr)
    
}



