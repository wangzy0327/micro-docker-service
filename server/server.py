#coding:utf-8
import subprocess
import requests
import threading
import logging
import sys
from concurrent.futures import ThreadPoolExecutor
from flask import Flask, request, jsonify
import uuid
import time
import re
import traceback

# 修复中文编码问题：强制stdout/stderr使用UTF-8
try:
    sys.stdout.buffer.write('\ufffd'.encode('utf-8'))  # 测试编码支持
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer)
    sys.stderr = codecs.getwriter('utf-8')(sys.stderr.buffer)
except:
    pass  # 忽略不支持的环境

# 配置日志（移除中文，使用英文避免编码问题）
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("server.log", encoding='utf-8'),  # 日志文件强制UTF-8
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

# 初始化Flask应用
app = Flask(__name__)
# 创建线程池控制并发任务数
executor = ThreadPoolExecutor(max_workers=10)

def parse_shell(shcmd):
    """Execute command and parse output"""
    try:
        logger.info("Executing command: %s" % shcmd)
        p = subprocess.Popen(
            shcmd,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=subprocess.PIPE,
            universal_newlines=True
        )
        stdout, stderr = p.communicate(timeout=300)
        combined = stdout + "\n" + stderr

        if p.returncode != 0:
            logger.error("Command failed (return code: %s), output: %s" % (p.returncode, combined))
            return False, combined
        else:
            logger.info("Command executed successfully, return code: %s" % p.returncode)

        # Extract model information
        pattern = r"(I\d{4} \d{2}:\d{2}:\d{2}\.\d{6} +\d+ caffe\.cpp:495\] execution time: .*? us)"
        match = re.search(pattern, combined, re.DOTALL | re.IGNORECASE)

        if match:
            model_info = match.group(0).strip()
            logger.info("Extracted model info: %s" % model_info)
            return True, model_info
        else:
            logger.warning("No execution time found, returning full output")
            return True, combined

    except subprocess.TimeoutExpired:
        p.kill()
        error_msg = "Command timed out (5 minutes)"
        logger.error(error_msg)
        return False, error_msg
    except Exception as e:
        error_msg = "Command execution error: %s" % str(e)
        logger.error(error_msg)
        logger.error("Traceback: %s" % traceback.format_exc())
        return False, error_msg


def async_process_task(input_path, output_path, uuid_str, callback_url):
    """Process task asynchronously"""
    try:
        logger.info("[Async Task] Starting processing UUID: %s" % uuid_str)
        logger.info("[Async Task] Input path: %s, Output path: %s" % (input_path, output_path))

        # Execute inference command
        inference_cmd = "cd /opt/cambricon/caffe/src/caffe && bash gen_offline_model.sh"
        success, result_output = parse_shell(inference_cmd)
        
        # Build task result
        task_result = {
            "status": "success" if success else "failed",
            "cmd": inference_cmd,
            "output": result_output,
            "input_path": input_path,
            "output_path": output_path,
            "execution_time": time.strftime("%Y-%m-%d %H:%M:%S")
        }

        # Callback to client
        if callback_url:
            logger.info("[Async Task] Calling back client: %s" % callback_url)
            callback_data = {
                "uuid": uuid_str,
                "task_status": task_result["status"],
                "result": task_result
            }

            response = requests.post(
                url=callback_url,
                json=callback_data,
                headers={"Content-Type": "application/json"},
                timeout=10
            )
            logger.info("[Async Task] Callback completed, status code: %s, response: %s" % (response.status_code, response.text))
        else:
            logger.warning("[Async Task] No callback URL provided, skipping callback")

    except Exception as e:
        error_msg = "Task processing error: %s" % str(e)
        logger.error(error_msg)
        logger.error("Traceback: %s" % traceback.format_exc())
        if callback_url:
            try:
                requests.post(
                    url=callback_url,
                    json={
                        "uuid": uuid_str,
                        "task_status": "error",
                        "error_msg": error_msg
                    },
                    headers={"Content-Type": "application/json"},
                    timeout=10
                )
            except Exception as ce:
                logger.error("Failed to send error callback: %s" % str(ce))
    finally:
        logger.info("[Async Task] Processing finished UUID: %s" % uuid_str)


@app.route('/pipeline', methods=['POST'])
def handle_pipeline():
    try:
        data = request.get_json()
        logger.info("Received /pipeline request: %s" % data)

        ipy_path = data["ipy_path"]
        pipeline_name = data["pipeline"]
        logger.info("ipy_path: %s, pipeline_name: %s" % (ipy_path, pipeline_name))

        root_path = "/home/pipeline_server/shells/"
        cmd = "%sstart.sh %s" % (root_path, ipy_path)
        output = subprocess.run(
            cmd,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            universal_newlines=True
        )
        return jsonify({
            "status": "success" if output.returncode == 0 else "failed",
            "output": output.stdout + output.stderr
        })

    except KeyError as e:
        error_msg = "Missing parameter: %s" % e
        logger.error(error_msg)
        return error_msg, 400
    except Exception as e:
        error_msg = "Server error: %s" % str(e)
        logger.error(error_msg)
        logger.error("Traceback: %s" % traceback.format_exc())
        return error_msg, 500


@app.route('/micro', methods=['POST'])
def handle_micro():
    try:
        # Log request details
        logger.info("Received /micro request, headers: %s" % request.headers)
        data = request.get_json()
        if not data:
            data = request.form  # Compatible with form-data
        logger.info("Received /micro request data: %s" % data)

        # Validate required parameters
        required_params = ["input", "callback_url"]
        for param in required_params:
            if param not in data:
                error_msg = "Missing parameter: %s" % param
                logger.error(error_msg)
                return error_msg, 400

        input_path = data["input"]
        output_path = data.get("output", "")
        callback_url = data["callback_url"]

        # Generate UUID
        uuid_str = str(uuid.uuid1())
        logger.info("Generated UUID: %s" % uuid_str)

        # Submit async task
        executor.submit(
            async_process_task,
            input_path,
            output_path,
            uuid_str,
            callback_url
        )

        return uuid_str

    except KeyError as e:
        error_msg = "Missing parameter: %s" % e
        logger.error(error_msg)
        return error_msg, 400
    except Exception as e:
        error_msg = "Server error: %s" % str(e)
        logger.error(error_msg)
        logger.error("Traceback: %s" % traceback.format_exc())
        return error_msg, 500


if __name__ == "__main__":
    logger.info("Starting MLU inference server...")
    app.run(
        debug=False,
        threaded=True,
        host="0.0.0.0",
        port=8800
    )
