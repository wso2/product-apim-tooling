# Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
#
# WSO2 Inc. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.

import argparse
import json
import os
import pickle
import random
import sys
import time
import numpy as np
import requests
import urllib3
import yaml
from collections import defaultdict
from datetime import datetime
from multiprocessing import Process, Value
from utils import util_methods, log


urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# variables
logger = log.setLogger('invoke_API')

max_connection_refuse_count = None
host_protocol = None
host_ip = None
host_port = None
heavy_traffic = None
post_data = []
delete_data = []
time_patterns = None

script_start_time = None
scenario_pool = {}
connection_refuse_count = Value('i', 0)
process_pool = []

abs_path = os.path.abspath(os.path.dirname(__file__))


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global max_connection_refuse_count, host_protocol, host_ip, host_port, heavy_traffic, post_data, delete_data, time_patterns

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as config_file:
        traffic_config = yaml.load(config_file, Loader=yaml.FullLoader)

    max_connection_refuse_count = int(traffic_config['tool_config']['max_connection_refuse_count'])
    heavy_traffic = str(traffic_config['tool_config']['heavy_traffic']).lower()
    host_protocol = traffic_config['api_host']['protocol']
    host_ip = traffic_config['api_host']['ip']
    host_port = traffic_config['api_host']['port']
    post_data = traffic_config['payloads']['post']
    delete_data = traffic_config['payloads']['delete']

    with open(abs_path + '/../../data/tool_data/invoke_patterns.yaml') as pattern_file:
        invoke_patterns = yaml.load(pattern_file, Loader=yaml.FullLoader)

    time_patterns = process_time_patterns(invoke_patterns['time_patterns'])


def process_time_patterns(patterns: dict) -> defaultdict:
    """
    Process time patterns to obtain mean and standard deviation to be used with distributions.
    :param patterns: Patterns dictionary.
    :return: Dictionary with mean and std for each pattern.
    """
    processed_patterns = defaultdict()

    for key, pattern in patterns.items():
        pattern = list(map(int, pattern.split(',')))
        mean = np.mean(pattern)
        std = np.std(pattern)
        processed_patterns[key] = {'mean': mean, 'std': std}
    return processed_patterns


def sendRequest(url_protocol, url_ip, url_port, path, access_token, method, user_ip, cookie, user_agent):
    """
    This function will send http requests to the given address
    :param url_protocol: Protocol of the URL
    :param url_ip: IP of the URL
    :param url_port: Port of the URL
    :param path: Invoke path for the request
    :param access_token: Access token for the request
    :param method: HTTP method of the request
    :param user_ip: IP of the user
    :param cookie: Cookie of the user
    :param user_agent: User agent for the user
    :return: Response code and response text
    """
    global post_data, delete_data

    url = "{}://{}:{}/{}".format(url_protocol, url_ip, url_port, path)
    accept = 'application/json'
    content_type = 'application/json'
    headers = {
        'accept': '{}'.format(accept),
        'Content-Type': '{}'.format(content_type),
        'Authorization': 'Bearer {}'.format(access_token),
        'client-ip': '{}'.format(user_ip),
        'x-forwarded-for': '{}'.format(user_ip),
        'cookie': '{}'.format(cookie),
        'User-Agent': '{}'.format(user_agent)
    }
    res_txt = ""

    try:
        if method == "GET":
            response = requests.get(url=url, headers=headers, verify=False)
            code = response.status_code
            res_txt = response.text

        elif method == "POST":
            data = json.dumps(random.choice(post_data))
            response = requests.post(url=url, headers=headers, data=data, verify=False)
            code = response.status_code
            res_txt = response.text

        elif method == "DELETE":
            if delete_data is not None:
                data = json.dumps(random.choice(delete_data))
                response = requests.delete(url=url, headers=headers, data=data, verify=False)
                code = response.status_code
                res_txt = response.text
            else:
                url = url + '/' + str(random.randint(0, 1000))
                response = requests.delete(url=url, headers=headers, verify=False)
                code = response.status_code
                res_txt = response.text

        elif method == "PUT":
            response = requests.put(url=url, headers=headers, verify=False)
            code = response.status_code
            res_txt = response.text

        elif method == "PATCH":
            response = requests.patch(url=url, headers=headers, verify=False)
            code = response.status_code
            res_txt = response.text

        else:
            code = '400'
            res_txt = 'Invalid type'

    except Exception as err:
        code = '521'
        log_txt = "sendRequest(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: " + method + ", url: " + url
        util_methods.log('traffic-requests.log', 'ERROR', log_txt)
        util_methods.log('traffic-tool.log', 'ERROR', str(err))

    # user agent is wrapped around quotes because there are commas in the user agent and they clash with the commas in csv file
    write_string = str(
        datetime.now()) + "," + user_ip + "," + access_token + "," + method + "," + path + "," + cookie + "," + accept + "," + content_type + "," + user_ip + ",\"" + user_agent + "\"," + str(
        code) + "\n"

    with open(abs_path + '/../../../../dataset/traffic/{}'.format(filename), 'a+') as dataset_file:
        dataset_file.write(write_string)

    return code, res_txt


def runInvoker(user_scenario, connection_refuse_count):
    """
    This function will take a given invoke scenario and execute it.
    Supposed to be executed from a process.
    :param user_scenario: User scenario data
    :param connection_refuse_count: Current connection refuse count
    :return: None
    """
    global script_start_time, script_runtime

    appNames = list(user_scenario.keys())

    while True:
        app_name = appNames[random.randint(0, len(appNames) - 1)]
        app_scenario_list = user_scenario.get(app_name)
        time_pattern = None

        iterations = 0
        probability_list = []

        # prepare probabilities for the scenario
        for scenario in app_scenario_list:
            iterations += scenario[0]
            probability_list.append(scenario[0])

        if iterations == 0:
            continue

        for i in range(len(probability_list)):
            probability_list[i] = probability_list[i] / iterations

        # increase probabilities if it's too small compared to max value
        for i in range(len(probability_list)):
            max_pro = max(probability_list)
            if max_pro - probability_list[i] >= 0.5:
                probability_list[i] = probability_list[i] + 0.075
                probability_list[probability_list.index(max_pro)] = max_pro - 0.075

        # prepare request pattern from list indices
        invoke_pattern_indices = np.random.choice(len(app_scenario_list), size=iterations, p=probability_list)

        for i in invoke_pattern_indices:
            up_time = datetime.now() - script_start_time
            if up_time.seconds >= script_runtime:
                break

            scenario = app_scenario_list[i]
            path = scenario[2]
            access_token = scenario[3]
            method = scenario[4]
            user_ip = scenario[5]
            cookie = scenario[6]
            user_agent = scenario[7]

            # set time pattern if not set
            if time_pattern is None:
                time_pattern = scenario[8]
                time_pattern = time_patterns.get(time_pattern)

            # send the request
            try:
                if heavy_traffic != 'true':
                    sleep_time = np.absolute(np.random.normal(time_pattern['mean'], time_pattern['std']))
                    time.sleep(sleep_time)
                res_code = sendRequest(host_protocol, host_ip, host_port, path, access_token, method, user_ip, cookie, user_agent)[0]
                if res_code == '521':
                    connection_refuse_count.value += 1

            except Exception as err:
                util_methods.log('traffic-tool.log', 'ERROR', str(err))
                connection_refuse_count.value += 1

        up_time = datetime.now() - script_start_time
        if up_time.seconds >= script_runtime:
            break
        else:
            time.sleep(abs(int(np.random.normal() * 10)))


if __name__ == "__main__":
    """
        Execute the scenario and generate the dataset
        Usage: python3 invoke_API.py filename exec_time
        output folder: dataset/traffic/
    """

    parser = argparse.ArgumentParser("run traffic tool")
    parser.add_argument("filename", help="Enter a filename to write final output (without extension)", type=str)
    parser.add_argument("runtime", help="Enter the script execution time in minutes", type=float)
    args = parser.parse_args()
    filename = args.filename + ".csv"
    script_runtime = args.runtime * 60  # in seconds

    # load and set tool configurations
    try:
        loadConfig()
    except FileNotFoundError as e:
        logger.exception(str(e))
        sys.exit()
    except Exception as e:
        logger.exception(str(e))
        sys.exit()

    with open(abs_path + '/../../../../dataset/traffic/{}'.format(filename), 'w') as file:
        file.write("timestamp,ip_address,access_token,http_method,invoke_path,cookie,accept,content_type,x_forwarded_for,user_agent,response_code\n")

    try:
        # load and set the scenario pool
        scenario_pool = pickle.load(open(abs_path + "/../../data/runtime_data/scenario_pool.sav", "rb"))
    except FileNotFoundError as e:
        logger.exception(str(e))
        sys.exit()

    # record script start_time
    script_start_time = datetime.now()

    processes_list = []

    # create and start a process for each user
    for scenario in scenario_pool.values():
        process = Process(target=runInvoker, args=(scenario, connection_refuse_count))
        process.daemon = False
        processes_list.append(process)
        process.start()

        with open(abs_path + '/../../data/runtime_data/traffic_processes.pid', 'a+') as file:
            file.write(str(process.pid) + '\n')

    logger.info("Scenario loaded successfully. Wait {} minutes to complete the script!".format(str(script_runtime / 60)))

    while True:
        time_elapsed = datetime.now() - script_start_time

        if time_elapsed.seconds >= script_runtime:
            for process in processes_list:
                process.terminate()
            with open(abs_path + '/../../data/runtime_data/traffic_processes.pid', 'w') as file:
                file.write('')

            logger.info("Script terminated successfully. Time elapsed: {} minutes".format(time_elapsed.seconds / 60.0))
            break

        elif connection_refuse_count.value > max_connection_refuse_count:
            for process in processes_list:
                process.terminate()
            with open(abs_path + '/../../data/runtime_data/traffic_processes.pid', 'w') as file:
                file.write('')
            logger.error("Terminating the program due to maximum no of connection refuses!")
            break

        else:
            pass
