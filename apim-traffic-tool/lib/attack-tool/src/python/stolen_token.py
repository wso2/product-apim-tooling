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

import atexit
import ipaddress
import os
import pickle
import random
import string
import sys
import time
from datetime import datetime
from multiprocessing import Process

import numpy as np
import requests
import yaml

from constants import *
from utils import util_methods


def generate_unique_ip():
    """
    Returns a unique ip address
    :return: an unique ip
    """
    global used_ips

    random.seed()
    MAX_IPV4 = ipaddress.IPv4Address._ALL_ONES
    temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(0, MAX_IPV4))
    while temp_ip in used_ips:
        temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(0, MAX_IPV4))

    used_ips.append(temp_ip)
    return temp_ip


def generate_cookie():
    """
    generates a random cookie
    :return: a randomly generated cookie
    """
    letters_and_digits = string.ascii_lowercase + string.digits
    cookie = 'JSESSIONID='
    cookie += ''.join(random.choice(letters_and_digits) for ch in range(31))
    return cookie


def simulate_user(user_data):
    """
      Simulate the behaviour of a user during the attack duration.
      :param user_data: A dictionary containing the user data
      :return: None
      """
    global attack_duration, protocol, host, port, payloads, user_agents, start_time, dataset_path, invoke_patterns

    up_time = datetime.now() - start_time
    sleep_pattern = invoke_patterns[random.choice(list(invoke_patterns.keys()))]

    if up_time.seconds < attack_duration:
        for app in user_data.values():

            invoke_pattern_indices = util_methods.generate_method_invoke_pattern(app)

            for i in invoke_pattern_indices:
                up_time = datetime.now() - start_time

                if up_time.seconds >= attack_duration:
                    break

                sleep_time = np.absolute(np.random.normal(sleep_pattern['mean'], sleep_pattern['std']))
                time.sleep(sleep_time)

                scenario = app[i]
                path = scenario[2]
                token = scenario[3]
                method = scenario[4]
                request_path = "{}://{}:{}/{}".format(protocol, host, port, path)
                random_user_agent = random.choice(user_agents)
                random_ip = generate_unique_ip()
                random_cookie = generate_cookie()
                random_payload = random.choice(payloads)
                accept = content_type = "application/json"

                try:
                    response = util_methods.send_simple_request(request_path, method, token, random_ip, random_cookie, accept, content_type, random_user_agent, payload=random_payload)
                    request_info = "{},{},{},{},{},{},{},{},{},\"{}\",{}".format(datetime.now(), random_ip, token, method, request_path, random_cookie, accept, content_type, random_ip,
                                                                                 random_user_agent,
                                                                                 response.status_code,
                                                                                 )
                    util_methods.log(dataset_path, request_info, "a")
                except requests.exceptions.RequestException:
                    msg_string = "[Error] {} - Request Failure\n\t {}".format(datetime.now(), str(ex))
                    print(msg_string)
                    util_methods.log(attack_tool_log_path, msg_string, "a")


# Program Execution
if __name__ == '__main__':

    attack_tool_log_path = "../../../../../../logs/attack-tool.log"

    # Constants
    STOLEN_TOKEN = 'stolen_token'

    try:
        with open(os.path.abspath(os.path.join(__file__, "../../../../traffic-tool/data/runtime_data/scenario_pool.sav")), "rb") as scenario_file:
            scenario_pool = pickle.load(scenario_file, )

        with open(os.path.abspath(os.path.join(__file__, "../../../../../config/attack-tool.yaml")), "r") as attack_config_file:
            attack_config = yaml.load(attack_config_file, Loader=yaml.FullLoader)

    except FileNotFoundError as ex:
        error_string = "[ERROR] {} - {}: \'{}\'".format(datetime.now(), ex.strerror, ex.filename)
        print(error_string)
        util_methods.log(attack_tool_log_path, error_string, "a")
        sys.exit()

    # Reading configurations from attack-tool.yaml
    protocol = attack_config[GENERAL_CONFIG][API_HOST][PROTOCOL]
    host = attack_config[GENERAL_CONFIG][API_HOST][IP]
    port = attack_config[GENERAL_CONFIG][API_HOST][PORT]
    attack_duration = attack_config[GENERAL_CONFIG][ATTACK_DURATION]
    payloads = attack_config[GENERAL_CONFIG][PAYLOADS]
    user_agents = attack_config[GENERAL_CONFIG][USER_AGENTS]
    process_count = attack_config[GENERAL_CONFIG][NUMBER_OF_PROCESSES]
    compromised_user_count = attack_config[ATTACKS][STOLEN_TOKEN][COMPROMISED_USER_COUNT]
    invoke_patterns = util_methods.process_time_patterns(attack_config[GENERAL_CONFIG][TIME_PATTERNS])

    # Recording column names in the dataset csv file
    dataset_path = "../../../../../../dataset/attack/stolen_token.csv"
    util_methods.log(dataset_path, "timestamp,ip_address,access_token,http_method,invoke_path,cookie,accept,content_type,x_forwarded_for,user_agent,response_code", "w")

    used_ips = []
    start_time = datetime.now()

    log_string = "[INFO] {} - Stolen token attack started ".format(start_time)
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

    # cleaning up the processes at exit
    atexit.register(util_methods.cleanup, process_list=process_list)
