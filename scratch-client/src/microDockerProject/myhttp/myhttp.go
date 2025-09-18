package myhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "unsafe"
)

// 定义请求体结构（包含原input/output + 新增callback_url）
type MicroRequest struct {
	Input       string `json:"input"`        // 原输入路径
	Output      string `json:"output"`       // 原输出路径
	CallbackUrl string `json:"callback_url"` // 客户端回调地址（新增）
}


// 改造原HttpPost：新增回调地址参数，发送带回调的请求
func HttpPostWithCallback(serverUrl, input, output, callbackUrl string) (string, error) {
	// 1. 构建请求体（携带回调地址）
	reqData := MicroRequest{
		Input:       input,
		Output:      output,
		CallbackUrl: callbackUrl,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("构建请求JSON失败: %s", err.Error())
	}

	// 2. 发送POST请求（保留原Header格式）
	req, err := http.NewRequest("POST", serverUrl, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// 3. 执行请求并读取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %s", err.Error())
	}
	defer resp.Body.Close()

	// 4. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("服务端响应异常: 状态码=%d", resp.StatusCode)
	}

	// 关键修复：用 ioutil.ReadAll 替换 http.ReadAll（兼容低版本Go）
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %s", err.Error())
	}
	uuidStr := string(respBody)
	return uuidStr, nil
}
