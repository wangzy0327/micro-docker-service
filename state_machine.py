import subprocess

import os
import time

class State():
    count = 10
    publish_output = "publish/publish.txt"
    subscribe_output = "subscribe/subscribe.txt"
    one_count = 0
    shcmd = []
    uuid = []

    def get_shcmd(self,file_name):
        f = open(file_name)
        self.shcmd = []
        self.uuid = []
        for i in range(self.count):
            sstr = f.readline().strip()
            if sstr != None and sstr.strip()!='':
                self.shcmd.append(sstr)
                print("shcmd : "+sstr)
                self.uuid.append(sstr.split()[-1])
                print("uuid : "+sstr.split()[-1])
        pass

    def parse_shell(self,shcmd):
        p = subprocess.Popen(shcmd,shell=True,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE)	
        stdout,stderr = p.communicate()
        if p.returncode != 0:
           print(str(stderr,encoding="utf-8").strip('\n'))
           return str(stderr,encoding="utf-8").strip('\n')
        return str(stdout,encoding="utf-8".strip('\n'))
        pass



    def get_last_lines(self,file_name,count):
        file_size = os.path.getsize(file_name)
        block_size = 1024
        file = open(file_name,'r')
        # last_line = ""
        if file_size > block_size:
            # maxseekpoint = (file_size//block_size)
            maxseekpoint = file_size - block_size if file_size >= block_size else 0
            # remainder = (file_size%block_size)
            file.seek(maxseekpoint)
        elif file_size:
            file.seek(0,0)
        lines = file.readlines()
        if lines :
            last_line = [ s.strip() for s in lines[-1*count:]]
            #print("last_line : ",last_line)
        file.close()
        time.sleep(1.5)
        return lines
        pass


    def exec_cmd(self):
        sh_num = len(self.shcmd)
        print("sh_num is " + str(sh_num))
        i = 0
        while i < sh_num:
            res = self.parse_shell(
                "docker run -v /home/wzy/micro-docker-service/:/root/ sum-java " + self.shcmd[i])
            print("docker run -v /home/wzy/micro-docker-service/:/root/ sum-java "+self.shcmd[i])
            print("exec business docker "+self.shcmd[i])
            print(res)
            i+=1
        pass        

    def go(self):
        i = 0
        #while True:
        while i < 1:
            print("\033[1;35m get_shcmd \033[0m")
            self.get_shcmd(self.publish_output)
            print("\033[1;35m exec_cmd \033[0m")
            self.exec_cmd()
            i+=1
        #pass

    def end_symbol(self):
        symbol = open('symbol','w')
        symbol.write('end\n')
        symbol.close()
        pass


if __name__ == "__main__":
    state = State()
    #state.get_shcmd(State.publish_output)
    #state.extend_group(3)
    #state.reduce_group()
    #lines = state.get_last_lines(state.subscribe_output,10)
    #state.shcmd=['sh run-wordcount2.sh input/input1 output/output1 908f7dee-0a6e-11ea-84ff-35f681938c05','sh run-wordcount2.sh input/input1 output/output1 a2d1d326-0a6e-11ea-84ff-35f681938c05','sh run-wordcount2.sh input/input2 output/output2 aed091c6-0a6e-11ea-84ff-35f681938c05','sh run-wordcount2.sh input/input3 output/output3 ab970ad0-0a6e-11ea-84ff-35f681938c05']
    #state.hadoop_list=['hadoop-master-0','hadoop-master-1','hadoop-master-2']
    #state.uuid=['908f7dee-0a6e-11ea-84ff-35f681938c05','a2d1d326-0a6e-11ea-84ff-35f681938c05','aed091c6-0a6e-11ea-84ff-35f681938c05','ab970ad0-0a6e-11ea-84ff-35f681938c05']
    #state.exec_cmd()
    #state.parse_shell("docker exec -d hadoop-master-0 bash -c 'sh run-wordcount2.sh input/input1 output/output1 908f7dee-0a6e-11ea-84ff-35f681938c05' ")
    startTime = round(time.time(),3)
    print('\033[1;35m 开始时间戳:'+str(startTime)+" \033[0m")
    state.go()
    state.end_symbol()
    endTime = round(time.time(),3)
    print('\033[1;35m 结束时间戳:'+str(endTime)+" \033[0m")
    diffTime = endTime - startTime
    print("processing time : "+str(diffTime))
    print("______________end______________")
