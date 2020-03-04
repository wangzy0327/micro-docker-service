current=`date "+%Y-%m-%d %H:%M:%S"`
timeStamp=`date -d "$current" +%s`
currentTimeStamp=$((10#$timeStamp*1000+10#`date "+%N"`/1000000))
beforeStartTime=`echo "scale=3;$currentTimeStamp/1000"|bc`
echo "开始时间戳:"$beforeStartTime
docker run -it -v /home/wzy/micro-docker-service/:/home/wzy/micro-docker-service scratch-client /scratch0 -input 100
