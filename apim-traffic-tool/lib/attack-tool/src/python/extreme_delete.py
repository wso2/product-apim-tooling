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

import os
import random
import sys
import time
from datetime import datetime
from multiprocessing.dummy import Pool

import pandas as pd
import requests
import yaml
import numpy as np
from utils import util_methods
from utils.entity_classes import API
from utils.util_methods import generate_random_string


def request_handler(i):
    """
    Handle the requests
    :return: None
    """
    global attack_duration, protocol, host, port, payloads, user_agents, api_list, dataset_path

    up_time = datetime.now() - start_time
    if up_time.seconds < attack_duration:
        api = random.choice(api_list)
        context = api.context
        version = api.version
        resource_path = random.choice(api.resources['DELETE'])
        # random_user = random.choice(api.users)
        random_user = api.single_user
        method = "DELETE"
        accept = content_type = "application/json"

        # sleep the process for a random period of time
        time.sleep(abs(int(np.random.normal() * 10)))

        request_path = "{}://{}:{}/{}/{}/{}".format(protocol, host, port, context, version, resource_path)
        random_user_agent = random.choice(user_agents)
        token = random_user[0]
        ip = random_user[2]
        cookie = random_user[3]
        path_param = generate_random_string(10)
        try:
            response = util_methods.send_simple_request(request_path, method, token, ip, cookie, accept, content_type, random_user_agent, path_params=path_param)
            request_info = "{},{},{},{},{}/{},{},{},{},{},\"{}\",{}".format(datetime.now(), ip, token, method, request_path,path_param, cookie, accept, content_type, ip, random_user_agent,
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
    try:
        with open(os.path.abspath(os.path.join(__file__, "../../../../../config/apim.yaml")), "r") as config_file:
            config = yaml.load(config_file, Loader=yaml.FullLoader)

        with open(os.path.abspath(os.path.join(__file__, "../../../../../config/attack-tool.yaml")), "r") as attack_config_file:
            attack_config = yaml.load(attack_config_file, Loader=yaml.FullLoader)
    except FileNotFoundError as ex:
        error_string = "[ERROR] {} - {}: \'{}\'".format(datetime.now(), ex.strerror, ex.filename)
        print(error_string)
        util_methods.log(attack_tool_log_path, error_string, "a")
        sys.exit()

    # reading configurations from attack-tool.yaml
    protocol = attack_config['general_config']['api_host']['protocol']
    host = attack_config['general_config']['api_host']['ip']
    port = attack_config['general_config']['api_host']['port']
    attack_duration = attack_config['general_config']['attack_duration']
    scenario_name = attack_config['general_config']['scenario']
    payloads = attack_config['general_config']['payloads']
    user_agents = attack_config['general_config']['user_agents']
    process_count = attack_config['general_config']['number_of_processes']

    # reading api configuration from apim.yaml
    apis = config['apis']

    # reading user data (access token, api name,user ip, and cookie)
    user_details = pd.read_csv(os.path.abspath(os.path.join(__file__, "../../../../traffic-tool/data/scenario/{}/token_ip_cookie.csv".format(scenario_name))))
    user_details_groups = user_details.groupby('api_name')

    # Instantiating API objects which has delete methods and appending them to api_list
    api_list = []
    for api in apis:
        temp = API(protocol, host, port, api['context'], api['version'], api['name'])
        temp.users = user_details_groups.get_group(temp.name).values.tolist()
        temp.set_single_user()
        for resource in api['resources']:
            temp.add_resource(resource['method'], resource['path'])
        if 'DELETE' in temp.resources.keys():
            api_list.append(temp)

    start_time = datetime.now()

    # Recording column names in the dataset csv file
    dataset_path = "../../../../../../dataset/attack/extreme_delete.csv"
    util_methods.log(dataset_path, "Timestamp, Request path, Method,Access Token, IP Address, Cookie, Response Code", "w")

    if len(api_list) == 0:
        error_string = "[ERROR] {} - There are no APIs with DELETE endpoints".format(datetime.now())
        print(error_string)
        util_methods.log(attack_tool_log_path, error_string, "a")
        sys.exit()

    log_string = "[INFO] {} - Extreme delete attack started ".format(start_time)
    print(log_string)
    util_methods.log(attack_tool_log_path, log_string, "a")
    util_methods.log(dataset_path, "timestamp,ip_address,access_token,http_method,invoke_path,cookie,accept,content_type,x_forwarded_for,user_agent,response_code", "w")

    process_pool = Pool(processes=process_count)

    # Executing scenarios until the attack duration elapses
    while True:
        time_elapsed = datetime.now() - start_time
        if time_elapsed.seconds >= attack_duration:
            log_string = "[INFO] {} - Attack terminated successfully. Time elapsed: {} minutes".format(datetime.now(), time_elapsed.seconds / 60.0)
            print(log_string)
            util_methods.log(attack_tool_log_path, log_string, "a")
            break
        else:
            process_pool.map(request_handler, range(1000))

    # closes the process pool and wait for the processes to finish
    process_pool.close()
    process_pool.join()
