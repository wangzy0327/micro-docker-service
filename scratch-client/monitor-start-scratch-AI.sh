#!/bin/bash

# 设置时区为 Asia/Shanghai
export TZ='Asia/Shanghai'

# 获取开始时间，用于计算总运行时长
start_total_time=$(date +%s)
# 输出开始信息
current_time=$(date +"%H:%M:%S.%3N")
echo "=== 脚本开始时间: $current_time ==="

# 定义日志文件
LOG_FILE="docker_run_logs_$(date +%Y%m%d_%H%M%S).txt"
echo "日志将保存到: $LOG_FILE" | tee -a "$LOG_FILE"

# 总运行时间30分钟（1800秒）
TOTAL_DURATION=$((1800+800))
# 记录当前已运行时间
elapsed_time=0

# 循环执行直到达到总运行时间
while [ $elapsed_time -lt $TOTAL_DURATION ]; do
    # 获取当前时间并记录到日志
    current_time=$(date +"%H:%M:%S.%3N")
    echo "=== 开始执行docker run - $current_time ===" | tee -a "$LOG_FILE"
    
    # 生成0或1的随机数，用于随机选择两个命令之一
    random_choice=$(( RANDOM % 2 ))
    
    # 根据随机数选择要执行的docker命令
    if [ $random_choice -eq 0 ]; then
        echo "选择执行classification命令" | tee -a "$LOG_FILE"
        docker run --privileged --net=host -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input=classification -type=micro 2>&1 | tee -a "$LOG_FILE"
    else
        echo "选择执行yolov3命令" | tee -a "$LOG_FILE"
        docker run --privileged --net=host -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input=yolov3 -type=micro 2>&1 | tee -a "$LOG_FILE"
    fi
    
    # 记录完成时间
    current_time=$(date +"%H:%M:%S.%3N")
    echo "=== docker run 执行完成 - $current_time ===" | tee -a "$LOG_FILE"
    
    # 计算下次执行前的随机等待时间（5-10分钟，转换为秒）
    # $RANDOM 生成0-32767的随机数，取模300得到0-299，加300得到300-599秒（5-10分钟）
    wait_seconds=$(( RANDOM % 60 + 180 ))
    wait_minutes=$(echo "scale=2; $wait_seconds / 60" | bc)
    
    #echo "下一次任务将在 $wait_minutes 分钟后提交执行" | tee -a "$LOG_FILE"
    echo "等待下一次任务提交执行" | tee -a "$LOG_FILE"
    
    # 等待随机时间
    sleep $wait_seconds
    
    # 更新已运行时间
    current_total_time=$(date +%s)
    elapsed_time=$(( current_total_time - start_total_time ))
    
    # 检查剩余时间是否足够进行下一次循环
    if [ $(( elapsed_time + wait_seconds )) -ge $TOTAL_DURATION ]; then
        echo "剩余时间不足，结束任务循环" | tee -a "$LOG_FILE"
        break
    fi
done

# 记录脚本结束时间
current_time=$(date +"%H:%M:%S.%3N")
echo "=== 脚本结束时间: $current_time ===" | tee -a "$LOG_FILE"
