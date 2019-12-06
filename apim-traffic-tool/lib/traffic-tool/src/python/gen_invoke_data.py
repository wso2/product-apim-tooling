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

import csv
import random
import string
import time
from datetime import datetime
import datetime as dt
import sys
import argparse
import urllib3
import pickle
import yaml
import os
import json
import math
from multiprocessing import Process, Value
import numpy as np
from scipy.stats import norm


parser = argparse.ArgumentParser("generate traffic data")
parser.add_argument("filename", help="Enter a filename to write final output", type=str)
args = parser.parse_args()
filename = args.filename + ".csv"

# variables
no_of_data_points = None
heavy_traffic = None
time_patterns = None

scenario_pool = {}
process_pool = []
current_data_points = Value('i', 0)
script_starttime = None

abs_path = os.path.abspath(os.path.dirname(__file__))


'''
    This function will load and set the configuration data
'''
def loadConfig():
    global no_of_data_points, heavy_traffic, time_patterns

    with open(abs_path+'/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)

    no_of_data_points = int(traffic_config['tool_config']['no_of_data_points'])
    heavy_traffic = str(traffic_config['tool_config']['heavy_traffic']).lower()

    with open(abs_path+'/../../data/access_pattern/invoke_patterns.yaml') as file:
        invoke_patterns = yaml.load(file, Loader=yaml.FullLoader)

    time_patterns = invoke_patterns['time_patterns']


'''
    This function will write the given log output to the log.txt file
'''
def log(tag, write_string):
    with open(abs_path+'/../../../../logs/traffic-tool.log', 'a+') as file:
        file.write("[{}] ".format(tag) + str(datetime.now()) + ": " + write_string + "\n")


'''
    This function will write the invoke request data to a file
'''
def writeInvokeData(timestamp, path, access_token, method, user_ip, cookie, app_name, username, user_agent):
    accept = 'application/json'
    content_type = 'application/json'
    code = '200'

    # user agent is wrapped around quotes because there are commas in the user agent and they clash with the commas in csv file
    write_string = str(timestamp) + "," + user_ip + "," + access_token + "," + method + "," + path + "," + cookie + "," + accept + "," + content_type + "," + user_ip + ",\"" + user_agent + "\"," + str(code) + "\n"

    with open(abs_path+'/../../../../dataset/generated-traffic/{}'.format(filename), 'a+') as file:
        file.write(write_string)


'''
    This function will take a given invoke scenario and generate data for it.
    Supposed to be executed from a process.
'''
def runInvoker(username, user_scenario, current_data_points):
    global no_of_data_points

    timestamp = datetime.now()
    appNames = list(user_scenario.keys())
    it = 0

    while(True):
        app_name = appNames[random.randint(0,len(appNames)-1)]
        app_scenario_list = user_scenario.get(app_name)
        time_pattern = None

        iterations = 0
        probability_list = []
        invoke_pattern_indices = None

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
            api_name = scenario[1]
            path = scenario[2]
            access_token = scenario[3]
            method = scenario[4]
            user_ip = scenario[5]
            cookie = scenario[6]
            user_agent = scenario[7]

            if time_pattern == None:
                time_pattern = scenario[8]
                time_pattern = time_patterns.get(time_pattern)
                if type(time_pattern) is str:
                    time_pattern = [int(t) for t in time_pattern.split(',')]
                else:
                    time_pattern = [time_pattern]

            writeInvokeData(timestamp, path, access_token, method, user_ip, cookie, app_name, username, user_agent)
            current_data_points.value += 1

            if heavy_traffic != 'true':
                timestamp += dt.timedelta(seconds=it%len(time_pattern))
            else:
                timestamp += dt.timedelta(seconds=abs(int(np.random.normal())))
            it += 1

        if current_data_points.value >= no_of_data_points:
            break
        else:
            timestamp += dt.timedelta(seconds=abs(int(np.random.normal() * 10)))


'''
    Generate the dataset according to the scenario
    Usage: python3 gen_invoke_data.py filename
    output folder: dataset/generated-traffic/
'''

# load and set tool configurations
try:
    loadConfig()
except FileNotFoundError as e:
    log('ERROR', '{}: {}'.format(e.strerror, e.filename))
    sys.exit()
except Exception as e:
    log('ERROR', '{}'.format(str(e)))
    sys.exit()

with open(abs_path+'/../../../../dataset/generated-traffic/{}'.format(filename), 'w') as file:
    file.write("timestamp,ip_address,access_token,http_method,invoke_path,cookie,accept,content_type,x_forwarded_for,user_agent,response_code\n")

try:
    # load and set the scenario pool
    scenario_pool = pickle.load(open(abs_path+"/../../data/runtime_data/scenario_pool.sav", "rb"))
except FileNotFoundError as e:
    log('ERROR', '{}: {}'.format(e.strerror, e.filename))
    sys.exit()

# record script starttime
script_starttime = datetime.now()

processes_list = []

# create and start a process for each user
for key_uname, val_scenario in scenario_pool.items():
    process = Process(target=runInvoker, args=(key_uname, val_scenario, current_data_points))
    process.daemon = False
    processes_list.append(process)
    process.start()

    with open(abs_path+'/../../data/runtime_data/traffic_processes.pid', 'a+') as file:
        file.write(str(process.pid)+'\n')

print("[INFO] Scenario loaded successfully. Wait until data generation complete!")
log("INFO", "Scenario loaded successfully. Wait until data generation complete!")

while True:
    if current_data_points.value >= no_of_data_points:
        for process in processes_list:
            process.terminate()
        with open(abs_path+'/../../data/runtime_data/traffic_processes.pid', 'w') as file:
            file.write('')

        time_elapsed = datetime.now() - script_starttime
        print("[INFO] Data generated successfully. Time elapsed: {} seconds".format(time_elapsed.seconds))
        log("INFO", "Data generated successfully. Time elapsed: {} seconds".format(time_elapsed.seconds))
        break
    else:
        pass
