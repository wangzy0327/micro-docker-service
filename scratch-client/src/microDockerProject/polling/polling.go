package polling

import (
	"fmt"
        "strings"
//	"strconv"
	"time"
	"microDockerProject/file"
)

const TIME_LAYOUT="2006-01-02 15:04:05"

type poll interface{

}

func count() int{
	return 10;
}

func diffMax() int64{
	return 5;
}

func diffMin() int64{
	return -5;
}

func Track(uuidStr string,output string,sublog string){
//         uuidStr="d9ef5254-09a4-11ea-87c2-3b8a8b7bbf5a"
for ;;{  
         uuidArr:=parseInfo(output,sublog)
         for i:=0;i<len(uuidArr);i++{
               if uuidStr == uuidArr[i]{
	          fmt.Println("find uuid :"+uuidStr)
                  fmt.Println("......................end...................")
                  return;
               }
         }
     }
}

func datetimeHandler(startAtTime []string){
   for i:=0;i<len(startAtTime)-1;i++{
     workflowStartAtTime:=strings.Replace(strings.Replace(startAtTime[i],"T"," ",-1),"Z","",-1)
     fmt.Println("handler time : "+workflowStartAtTime)
     loc,_:=time.LoadLocation("Local")
     //theTime,_:=time.ParseInLocation(TIME_LAYOUT,"2018-09-10 00:00:00",loc)
     theTime,_:=time.ParseInLocation(TIME_LAYOUT,workflowStartAtTime,loc)
     workflowTimeStamp:=theTime.Unix()
     fmt.Println(theTime.Format(TIME_LAYOUT))
     fmt.Println(workflowTimeStamp) 
   }
}

func parseInfo(output string,sublog string)([]string){
  var strs = file.ReadLastBlock(count(),sublog)
  time.Sleep(time.Duration(1)*time.Second)
	var uuidArr []string
	for j:=0;j<len(strs);j++{
		strarray := strings.Fields(strings.TrimSpace(strs[j]))
                if strarray!=nil && len(strarray) !=0{
                   //fmt.Println("strarray len : "+strconv.Itoa(len(strarray)))
                   //fmt.Println("uuid : "+strarray[len(strarray)-1])
		   uuidArr = append(uuidArr,strarray[len(strarray)-1])
                }
	}
        /*
	for i:=0;i<len(uuidArr);i++{
		fmt.Println(uuidArr[i])
	}
        */
        /*
	for i:=0;i<len(workflowStatus);i++{
		fmt.Println(workflowStatus[i])
	}
	for i:=0;i<len(startAtTime);i++{
		fmt.Println(startAtTime[i])
	}
	*/
	return uuidArr 
}

