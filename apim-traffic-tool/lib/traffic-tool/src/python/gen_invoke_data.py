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
import datetime as dt
import os
import pickle
import random
import sys
import numpy as np
import yaml
from datetime import datetime
from multiprocessing import Process, Value
from collections import defaultdict
from utils import log

# variables
logger = log.setLogger('gen_invoke_data')

no_of_data_points = None
heavy_traffic = None
time_patterns = None

scenario_pool = {}
process_pool = []
current_data_points = Value('i', 0)
script_start_time = None

abs_path = os.path.abspath(os.path.dirname(__file__))


def loadConfig():
    """
    This function will load and set the configuration data
    :return: None
    """
    global no_of_data_points, heavy_traffic, time_patterns

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as config_file:
        traffic_config = yaml.load(config_file, Loader=yaml.FullLoader)

    no_of_data_points = int(traffic_config['tool_config']['no_of_data_points'])
    heavy_traffic = str(traffic_config['tool_config']['heavy_traffic']).lower()

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


def writeInvokeData(timestamp, path, access_token, method, user_ip, cookie, user_agent):
    """
    This function will write the invoke request data to a file
    :param timestamp: Timestamp of the request
    :param path: Invoke path of the request
    :param access_token: Access token of the request
    :param method: Http method of the request
    :param user_ip: User IP of the request
    :param cookie: User cookie of the request
    :param user_agent: User agent of the request
    :return: None
    """
    accept = 'application/json'
    content_type = 'application/json'
    code = '200'

    # user agent is wrapped around quotes because there are commas in the user agent and they clash with the commas in csv file
    write_string = str(
        timestamp) + "," + user_ip + "," + access_token + "," + method + "," + path + "," + cookie + "," + accept + "," + content_type + "," + user_ip + ",\"" + user_agent + "\"," + str(code) + "\n"

    with open(abs_path + '/../../../../dataset/generated-traffic/{}'.format(filename), 'a+') as dataset_file:
        dataset_file.write(write_string)


def runInvoker(user_scenario, current_data_points):
    """
    This function will take a given invoke scenario and generate data for it.
    Supposed to be executed from a process.
    :param user_scenario: User scenario as a list
    :param current_data_points: Current data point count
    :return: None
    """
    global no_of_data_points

    timestamp = datetime.now()
    appNames = list(user_scenario.keys())
    it = 0

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
            if current_data_points.value >= no_of_data_points:
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

            writeInvokeData(timestamp, path, access_token, method, user_ip, cookie, user_agent)
            current_data_points.value += 1

            if heavy_traffic != 'true':
                sleep_time = np.absolute(np.random.normal(time_pattern['mean'], time_pattern['std']))
                timestamp += dt.timedelta(seconds=sleep_time)
            else:
                timestamp += dt.timedelta(seconds=abs(int(np.random.normal())))
            it += 1

        if current_data_points.value >= no_of_data_points:
            break
        else:
            timestamp += dt.timedelta(seconds=abs(int(np.random.normal() * 10)))


if __name__ == "__main__":
    '''
        Generate the dataset according to the scenario
        Usage: python3 gen_invoke_data.py filename
        output folder: dataset/generated-traffic/
    '''

    parser = argparse.ArgumentParser("generate traffic data")
    parser.add_argument("filename", help="Enter a filename to write final output (without extension)", type=str)
    args = parser.parse_args()
    filename = args.filename + ".csv"

    # load and set tool configurations
    try:
        loadConfig()
    except FileNotFoundError as e:
        logger.exception(str(e))
        sys.exit()
    except Exception as e:
        logger.exception(str(e))
        sys.exit()

    with open(abs_path + '/../../../../dataset/generated-traffic/{}'.format(filename), 'w') as file:
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
    for key_uname, val_scenario in scenario_pool.items():
        process = Process(target=runInvoker, args=(val_scenario, current_data_points))
        process.daemon = False
        processes_list.append(process)
        process.start()

        with open(abs_path + '/../../data/runtime_data/traffic_processes.pid', 'a+') as file:
            file.write(str(process.pid) + '\n')

    logger.info("Scenario loaded successfully. Wait until data generation complete!")

    while True:
        if current_data_points.value >= no_of_data_points:
            for process in processes_list:
                process.terminate()
            with open(abs_path + '/../../data/runtime_data/traffic_processes.pid', 'w') as file:
                file.write('')

            time_elapsed = datetime.now() - script_start_time
            logger.info("Data generated successfully. Time elapsed: {} seconds".format(time_elapsed.seconds))
            break
        else:
            pass
