#!/bin/bash
set -euo pipefail  # 严格模式：报错立即退出，避免空变量，捕获管道错误

# ==================== 配置参数 ====================
NUM_PROCESSES=5                 # 同时启动的Docker进程数
IMAGE_NAME="scratch-client:amd64"      # Docker镜像名
CONTAINER_PREFIX="scratch-hadoop-client-amd64-" # 容器名前缀（区分不同容器）
LOG_DIR="./logs"                  # 日志根目录
TZ="Asia/Shanghai"               # 时区（确保容器内时间与宿主机一致）

# ==================== 函数定义 ====================
# 函数：打印带时间戳的脚本日志（区分脚本自身日志和容器日志）
script_log() {
    local current_time=$(TZ="$TZ" date +"%H:%M:%S.%3N")
    echo "[$current_time] [脚本] $1"
}

# 函数：启动单个Docker进程（后台执行，日志独立输出，实时计算启动时间）
start_single_container() {
    local index=$1  # 进程序号（1~10，用于日志文件名和容器名）
    local container_name="${CONTAINER_PREFIX}${index}"  # 容器名：micro-client-1
    local container_log_path="${LOG_DIR}/${container_name}.log"  # 容器日志路径

    # 1. 确保日志目录存在（不存在则创建）
    if [ ! -d "$LOG_DIR" ]; then
        mkdir -p "$LOG_DIR"
        script_log "日志目录 $LOG_DIR 不存在，已自动创建"
    fi

    # 2. 实时计算当前容器的启动时间（每个容器启动前重新计算，避免间隙影响）
    local script_start_time=$(date +"%H:%M:%S.%3N")
    script_log "准备启动第 $index 个容器，实时启动时间：$script_start_time"

    # 3. 启动Docker容器：后台执行+日志输出到指定文件+容器退出自动清理
    docker run  \
        --privileged \
        --net=host \
        --rm \
        --name "$container_name" \
        -v /dev:/dev \
        -e START_TIME="$script_start_time" \
        "$IMAGE_NAME" \
        /scratch0 > "$container_log_path" 2>&1 &  # 日志重定向：stdout/stderr都写入文件

    # 4. 记录容器后台PID（用于后续等待执行完成）
    local container_pid=$!
    script_log "第 $index 个容器启动成功！容器名：$container_name，后台PID：$container_pid，日志文件：$container_log_path"
    echo "$container_pid $index" >> .container_pids.tmp  # 记录PID和序号（便于日志关联）
}

# ==================== 主逻辑 ====================
# 1. 初始化：清理上一次的PID临时文件
rm -f .container_pids.tmp
script_log "===== 开始启动 $NUM_PROCESSES 个Docker容器 ====="

# 2. 循环启动10个Docker进程（并行执行，每个容器独立计算启动时间）
for ((i=1; i<=NUM_PROCESSES; i++)); do
    start_single_container "$i"  # 传入进程序号，启动单个容器
    # （可选）轻微延迟避免Docker瞬时资源竞争，根据宿主机性能调整（0.05~0.2秒）
    sleep 0.05
done

# 3. 等待所有后台容器执行完成（阻塞直到所有容器退出）
script_log "===== 所有 $NUM_PROCESSES 个容器已启动，等待执行完成 ====="
while read -r pid index; do
    if wait "$pid"; then
        container_name="${CONTAINER_PREFIX}${index}"
        container_log_path="${LOG_DIR}/${container_name}.log"
        script_log "第 $index 个容器（PID：$pid）执行完成！完整日志：$container_log_path"
    else
        script_log "警告：第 $index 个容器（PID：$pid）执行失败！可查看日志排查：${LOG_DIR}/${CONTAINER_PREFIX}${index}.log"
    fi
done < .container_pids.tmp

# 4. 清理临时文件，提示执行完成
rm -f .container_pids.tmp
script_log "===== 所有 $NUM_PROCESSES 个Docker容器执行完成 ====="
