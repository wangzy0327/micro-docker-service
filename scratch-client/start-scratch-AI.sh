# 设置时区为 Asia/Shanghai
export TZ='Asia/Shanghai'

# 获取当前时间，格式：HH:MM:SS.mmm（毫秒）
current_time=$(date +"%H:%M:%S.%3N")

# 输出
echo "开始时间: $current_time"

# 运行 Docker 命令
docker run --privileged --net=host -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input=classification -type=micro
#docker run --privileged --net=host -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:arm64 /scratch1 -input=yolov3 -type=micro
