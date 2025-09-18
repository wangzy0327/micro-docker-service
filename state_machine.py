import subprocess
import os
import time
import re

class State():
    count = 10
    publish_output = "publish/publish.txt"
    subscribe_output = "subscribe/subscribe.txt"
    one_count = 0
    shcmd = []
    uuid = []

    def get_shcmd(self,file_name):
        with open(file_name) as f:
            self.shcmd = []
            self.uuid = []
            for i in range(self.count):
                sstr = f.readline().strip()
                if sstr:
                    self.shcmd.append(sstr)
                    print("task id : "+sstr)
                    #self.uuid.append(sstr.split()[-1])
                    #print("uuid : "+sstr.split()[-1])

    def parse_shell(self,shcmd):
        '''
        p = subprocess.Popen(shcmd,shell=True,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
        stdout,stderr = p.communicate()
        if p.returncode != 0:
           print(str(stderr,encoding="utf-8").strip('\n'))
           return False, str(stderr,encoding="utf-8").strip('\n')
        return True, str(stdout,encoding="utf-8").strip('\n')
        '''
        # 执行命令并捕获输出
        p = subprocess.Popen(
            shcmd,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=subprocess.PIPE,
            universal_newlines=True  # 直接返回字符串，无需 decode
        )
        stdout, stderr = p.communicate()

        # 合并输出（关键！）
        combined = stdout + "\n" + stderr

        #print("=== STDOUT ===")
        #print(stdout)
        #print("=== STDERR ===")
        #print(stderr)
        #print("=== RETURN CODE:", p.returncode)

        # 检查是否执行失败
        if p.returncode != 0:
            print("Command failed with return code:",p.returncode)
            print("Output:\n",combined)
            return False, None

        # 使用正则提取模型信息块
        #pattern = r'.*?Offline model information BEGIN.*?\n.*?Offline model information END.*?'
        # 尝试匹配：模型块 + 可选的下一行 execution time
        #pattern = r'(.*?Offline model information BEGIN.*?\n.*?Offline model information END.*?)(\n\s*I\d{6}.*?execution time.*?us)?'
        #pattern = r'I\d{6} \d{2}:\d{2}:\d{2}\.\d{6}.*?execution time: .*?us'
        #pattern = r"(\*{40} Offline model information BEGIN \*{40}\n.*?\*{41} Offline model information END \*{41})"
        pattern = r"(I\d{4} \d{2}:\d{2}:\d{2}\.\d{6} +\d+ caffe\.cpp:495\] execution time: .*? us)"
        match = re.search(pattern, combined, re.DOTALL | re.IGNORECASE)

        if match:
            print("===== Matched Execution Time Line ====")
            # 提取完整的信息块（包括 BEGIN 和 END 行）
            model_info = match.group(0).strip()
            #print(model_info)  # 只打印你关心的部分
            return True, model_info
        else:
            print("Warning: Could not find offline model information in output.")
            # 如果没找到，也可以选择打印全部输出用于调试
            # print(stdout)
            return True, combined  # 即使没提取到，也算“成功”，但返回原始输出

    def update_files(self, completed_tasks):
        # 更新 publish.txt，移除已完成的任务
        with open(self.publish_output, 'r') as file:
            lines = file.readlines()
        with open(self.publish_output, 'w') as file:
            for line in lines:
                if line.strip() not in completed_tasks:
                    file.write(line)

        # 将已完成的任务写入 subscribe.txt
        with open(self.subscribe_output, 'a') as file:
            for task in completed_tasks:
                file.write(task + '\n')

    def exec_cmd(self):
        sh_num = len(self.shcmd)
        print("sh_num is " + str(sh_num))
        completed_tasks = []
        for cmd in self.shcmd:
            res, output = self.parse_shell(
                "docker exec elastic_neumann bash -c 'cd /opt/cambricon/caffe/src/caffe && bash gen_offline_model.sh' ")
            print(f"Executing: {cmd}\nResult: {res}\nOutput: {output}")
            if res:
                completed_tasks.append(cmd)
        if completed_tasks:
            self.update_files(completed_tasks)

    def go(self):
        print("\033[1;35m get_task \033[0m")
        self.get_shcmd(self.publish_output)
        print("\033[1;35m exec_cmd \033[0m")
        self.exec_cmd()

    def end_symbol(self):
        with open('symbol','w') as symbol:
            symbol.write('end\n')


if __name__ == "__main__":
    state = State()
    startTime = round(time.time(),3)
    state.go()
    state.end_symbol()
    endTime = round(time.time(),3)
    diffTime = endTime - startTime
    print("processing time : "+str(diffTime))
    print("______________end______________")
