# 设置时区为 Asia/Shanghai
export TZ='Asia/Shanghai'

# 获取当前时间，格式：HH:MM:SS.mmm（毫秒）
current_time=$(date +"%H:%M:%S.%3N")

# 输出
echo "开始时间: $current_time"

# 运行 Docker 命令
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input input/input1 -output output/output1 
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input grep -output output/output1 
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input pi -output output/output1 
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input randomwriter -output output/output1 
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input sort -output output/output1 
#docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input wordmean -output output/output1 
docker run --privileged --net=host --rm -it -v /dev:/dev   -v/time:/time -e START_TIME="$current_time" scratch-client:amd64 /scratch0 -input wordmedian -output output/output1 

