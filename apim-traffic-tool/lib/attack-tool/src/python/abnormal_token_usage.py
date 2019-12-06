# Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
import atexit
import sys
from multiprocessing.dummy import Pool
import os
import pickle
import time

import requests
import yaml
from datetime import datetime
from utils import util_methods
import random
import numpy as np
from multiprocessing import Process, Value
from utils.util_methods import generate_biased_random


def execute_scenario(scenario):
    """
    Execute scenarios from the scenario pool to simulate abnormal token usage
    :param scenario: A list containing a scenario
    :return: none
    """
    global attack_duration, protocol, host, port, payloads, user_agents, start_time, max_request_multiplier, min_request_multiplier, dataset_path

    up_time = datetime.now() - start_time

    if up_time.seconds < attack_duration:
        # multiply normal request count by a random value between user defined min and max value
        request_target = scenario[0] * random.randint(min_request_multiplier, max_request_multiplier)
        context = scenario[1]
        version = scenario[2]
        resource_path = scenario[3]
        token = scenario[4]
        method = scenario[5]
        ip = scenario[6]
        cookie = scenario[7]
        user_agent = scenario[10]
        accept = content_type = "application/json"

        request_path = "{}://{}:{}/{}/{}/{}".format(protocol, host, port, context, version, resource_path)
        random_payload = random.choice(payloads)

        # sending requests until the request target achieved or attack duration elapses
        for i in range(request_target):
            up_time = datetime.now() - start_time
            if up_time.seconds >= attack_duration:
                break
            try:
                response = util_methods.send_simple_request(request_path, method, token, ip, cookie, accept, content_type, user_agent, payload=random_payload)
                request_info = "{},{},{},{},{},{},{},{},{},\"{}\",{}".format(datetime.now(), ip, token, method, request_path, cookie, accept, content_type, ip, user_agent,
                                                                             response.status_code,
                                                                             )
                util_methods.log(dataset_path, request_info, "a")
            except requests.exceptions.RequestException:
                msg_string = "[Error] {} - Request Failure\n\t {}".format(datetime.now(), str(ex))
                print(msg_string)
                util_methods.log(attack_tool_log_path, msg_string, "a")

            # sleep the process for a random period of time
            time.sleep(abs(int(np.random.normal() * 10)))


def simulate_user(user_data):
    global attack_duration, protocol, host, port, payloads, user_agents, start_time, dataset_path, invoke_patterns, max_request_multiplier, min_request_multiplier

    up_time = datetime.now() - start_time

    if up_time.seconds < attack_duration:
        for app in user_data.values():
            invoke_pattern_iterator = 0

            for scenario in app:
                scenario[0] *= random.randint(min_request_multiplier, max_request_multiplier)

            invoke_pattern_indices = util_methods.generate_method_invoke_pattern(app)
            for i in invoke_pattern_indices:
                up_time = datetime.now() - start_time

                if up_time.seconds >= attack_duration:
                    break

                scenario = app[i]
                path = scenario[2]
                token = scenario[3]
                method = scenario[4]
                pattern = invoke_patterns[random.choice(list(invoke_patterns.keys()))].split(',')
                request_path = "{}://{}:{}/{}".format(protocol, host, port, path)
                random_user_agent = scenario[7]
                ip = scenario[5]
                cookie = scenario[6]
                random_payload = random.choice(payloads)
                accept = content_type = "application/json"

                try:
                    response = util_methods.send_simple_request(request_path, method, token, ip, cookie, accept, content_type, random_user_agent, payload=random_payload)
                    request_info = "{},{},{},{},{},{},{},{},{},\"{}\",{}".format(datetime.now(), ip, token, method, request_path, cookie, accept, content_type, ip,
                                                                                 random_user_agent,
                                                                                 response.status_code,
                                                                                 )
                    util_methods.log(dataset_path, request_info, "a")
                except requests.exceptions.RequestException:
                    msg_string = "[Error] {} - Request Failure\n\t {}".format(datetime.now(), str(ex))
                    print(msg_string)
                    util_methods.log(attack_tool_log_path, msg_string, "a")

                time.sleep(int(pattern[invoke_pattern_iterator % len(pattern)]))
                invoke_pattern_iterator += 1


# Program Execution
if __name__ == '__main__':

    attack_tool_log_path = "../../../../../../logs/attack-tool.log"

    try:
        with open(os.path.abspath(os.path.join(__file__, "../../../../traffic-tool/data/runtime_data/scenario_pool.sav")), "rb") as scenario_file:
            scenario_pool = pickle.load(scenario_file, )

        with open(os.path.abspath(os.path.join(__file__, "../../../../traffic-tool/data/access_pattern/invoke_patterns.yaml")), "rb") as pattern_file:
            invoke_patterns = yaml.load(pattern_file, Loader=yaml.FullLoader)['time_patterns']

        with open(os.path.abspath(os.path.join(__file__, "../../../../../config/attack-tool.yaml")), "r") as attack_config_file:
            attack_config = yaml.load(attack_config_file, Loader=yaml.FullLoader)
    except FileNotFoundError as ex:
        error_string = "[ERROR] {} - {}: \'{}\'".format(datetime.now(), ex.strerror, ex.filename)
        print(error_string)
        util_methods.log(attack_tool_log_path, error_string, "a")
        sys.exit()

    # Reading configurations from attack-tool.yaml
    protocol = attack_config['general_config']['api_host']['protocol']
    host = attack_config['general_config']['api_host']['ip']
    port = attack_config['general_config']['api_host']['port']
    attack_duration = attack_config['general_config']['attack_duration']
    payloads = attack_config['general_config']['payloads']
    user_agents = attack_config['general_config']['user_agents']
    process_count = attack_config['general_config']['number_of_processes']
    max_request_multiplier = attack_config['attacks']['abnormal_token_usage']['max_request_scalar']
    min_request_multiplier = attack_config['attacks']['abnormal_token_usage']['min_request_scalar']
    compromised_user_count = attack_config['attacks']['abnormal_token_usage']['compromised_user_count']
    start_time = datetime.now()

    # Recording column names in the dataset csv file
    dataset_path = "../../../../../../dataset/attack/abnormal_token.csv"
    util_methods.log(dataset_path, "timestamp,ip_address,access_token,http_method,invoke_path,cookie,accept,content_type,x_forwarded_for,user_agent,response_code", "w")

    log_string = "[INFO] {} - Abnormal token usage attack started ".format(start_time)
    print(log_string)
    util_methods.log(attack_tool_log_path, log_string, "a")

    if compromised_user_count > len(scenario_pool):
        error_string = "[ERROR] {} - More compromised users than the total users".format(datetime.now())
        print(error_string)
        util_methods.log(attack_tool_log_path, error_string, "a")
        sys.exit()

    compromised_users = np.random.choice(list(scenario_pool.values()), size=compromised_user_count, replace=False)
    process_list = []
    for user in compromised_users:
        process = Process(target=simulate_user, args=(user,))
        process.daemon = False
        process_list.append(process)
        process.start()

    while True:
        time_elapsed = datetime.now() - start_time
        if time_elapsed.seconds >= attack_duration:
            for process in process_list:
                process.terminate()
            log_string = "[INFO] {} - Attack terminated successfully. Time elapsed: {} minutes".format(datetime.now(), time_elapsed.seconds / 60.0)
            print(log_string)
            util_methods.log(attack_tool_log_path, log_string, "a")
            break

    atexit.register(util_methods.cleanup)

    # process_pool = Pool(processes=process_count)

    # Executing scenarios until the attack duration elapses
    # while True:
    #     time_elapsed = datetime.now() - start_time
    #     if time_elapsed.seconds >= attack_duration:
    #         log_string = "[INFO] {} - Attack terminated successfully. Time elapsed: {} minutes".format(datetime.now(), time_elapsed.seconds / 60.0)
    #         print(log_string)
    #         util_methods.log(attack_tool_log_path, log_string, "a")
    #         break
    #     else:
    #         process_pool.map(execute_scenario, scenario_pool)
    #
    # # closes the process pool and wait for the processes to finish
    # process_pool.close()
    # process_pool.join()
