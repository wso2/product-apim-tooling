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

import os
import yaml
import math
import random
from scipy.stats import truncnorm
from utils import log

# variables
logger = log.setLogger('gen_scenario_yaml')
abs_path = None
no_of_users = 0
frequency_limits = None
app_list = {}
time_patterns = {}
app_api_subs = {}
api_resources = {}


def getNormal(mean, sd, lower, upper):
    """
    This function will return a normal distribution bound to a given limit
    :param mean: mean of the distribution
    :param sd: standard deviation
    :param mean: lower bound of the range
    :param mean: upper bound of the range
    :return: normal distribution
    """
    return truncnorm( (lower - mean) / sd, (upper - mean) / sd, loc=mean, scale=sd)


def genUserApps():
    """
    This function will distribute users among applications and write applications to user_details.yaml file.
    User app distribution (no of apps a particular user is using) happen according to a normal distribution.
    :return: None
    """
    global no_of_users, app_list

    app_usage = []
    no_of_apps = len(app_list)
    mn = int(math.floor(no_of_apps/2))
    user_data = None

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as f:
        user_data = yaml.load(f, Loader=yaml.FullLoader)
    users = user_data['users']

    # calculate the number of apps each user is using (a normal distribution is used)
    distribution = getNormal(mn-1, mn, 1, no_of_apps)
    app_usage = distribution.rvs(no_of_users)
    app_usage = [int(round(num, 0)) for num in app_usage]

    # distribute apps among users
    for i in range(len(users)):
        user_apps = random.sample(app_list.keys(), app_usage.pop(0))
        users[i]['applications'] = ','.join(user_apps)

        for app in user_apps:
            app_list.get(app).append(users[i]['username'])

    # update user_details.yaml file
    with open(abs_path + '/../../data/scenario/user_details.yaml', 'w') as f:
        yaml.dump(user_data, f, sort_keys=False)

    logger.info('Users distributed among apps')


def genScenarioDistribution():
    """
    This function will generate a random invoke scenario and write output to the invoke_scenario.yaml file
    :return: None
    """
    scenario_list = []
    low_freq_lower = frequency_limits['low']['lower']
    low_freq_upper = frequency_limits['low']['upper']
    medium_freq_lower = frequency_limits['medium']['lower']
    medium_freq_upper = frequency_limits['medium']['upper']
    high_freq_lower = frequency_limits['high']['lower']
    high_freq_upper = frequency_limits['high']['upper']

    for app in app_list:
        # get user count for the app
        app_user_count = len(app_list.get(app))

        # distribute users for each type based on request frequency
        mn = int(math.floor((high_freq_upper-low_freq_lower)/2))

        distribution = getNormal(mn, mn, low_freq_lower, high_freq_upper)
        user_frequencies = distribution.rvs(app_user_count)
        user_frequencies = [int(round(num, 0)) for num in user_frequencies]

        no_of_users_low_freq = len(list(fq for fq in user_frequencies if low_freq_lower <= fq < low_freq_upper))
        no_of_users_medium_freq = len(list(fq for fq in user_frequencies if medium_freq_lower <= fq < medium_freq_upper))
        no_of_users_high_freq = len(list(fq for fq in user_frequencies if high_freq_lower <= fq <= high_freq_upper))

        # generate scenario for low frequent category
        if(no_of_users_low_freq != 0):
            temp = {}
            temp['app_name'] = app
            temp['no_of_users'] = no_of_users_low_freq
            temp['time_pattern'] = random.choice(time_patterns['low'].split(',')).strip()
            api_calls = []

            for api in app_api_subs.get(app):
                for resource in api_resources.get(api):
                    temp_2 = {}
                    temp_2['api'] = api
                    temp_2['method'] = resource

                    mn = int(math.floor((low_freq_upper-low_freq_lower)/2))
                    distribution = getNormal(mn, mn, low_freq_lower, low_freq_upper)
                    no_of_reqs = int(distribution.rvs(1)[0])

                    if resource in ['POST', 'DELETE']:
                        no_of_reqs = int(no_of_reqs/2)

                    temp_2['no_of_requests'] = no_of_reqs
                    api_calls.append(temp_2)

            temp['api_calls'] = api_calls
            scenario_list.append(temp)

        # generate scenario for medium frequent category
        if(no_of_users_medium_freq != 0):
            temp = {}
            temp['app_name'] = app
            temp['no_of_users'] = no_of_users_medium_freq
            temp['time_pattern'] = random.choice(time_patterns['medium'].split(',')).strip()
            api_calls = []

            for api in app_api_subs.get(app):
                for resource in api_resources.get(api):
                    temp_2 = {}
                    temp_2['api'] = api
                    temp_2['method'] = resource

                    mn = int(math.floor((medium_freq_upper-medium_freq_lower)/2))
                    distribution = getNormal(mn, mn, medium_freq_lower, medium_freq_upper)
                    no_of_reqs = int(distribution.rvs(1)[0])

                    if resource in ['POST', 'DELETE']:
                        no_of_reqs = int(no_of_reqs/2)

                    temp_2['no_of_requests'] = no_of_reqs
                    api_calls.append(temp_2)

            temp['api_calls'] = api_calls
            scenario_list.append(temp)

        # generate scenario for high frequent category
        if(no_of_users_high_freq != 0):
            temp = {}
            temp['app_name'] = app
            temp['no_of_users'] = no_of_users_high_freq
            temp['time_pattern'] = random.choice(time_patterns['high'].split(',')).strip()
            api_calls = []

            for api in app_api_subs.get(app):
                for resource in api_resources.get(api):
                    temp_2 = {}
                    temp_2['api'] = api
                    temp_2['method'] = resource

                    mn = int(math.floor((high_freq_upper-high_freq_lower)/2))
                    distribution = getNormal(mn, mn, high_freq_lower, high_freq_upper)
                    no_of_reqs = int(distribution.rvs(1)[0])

                    if resource in ['POST', 'DELETE']:
                        no_of_reqs = int(no_of_reqs/2)

                    temp_2['no_of_requests'] = no_of_reqs
                    api_calls.append(temp_2)

            temp['api_calls'] = api_calls
            scenario_list.append(temp)

    with open(abs_path + '/../../data/scenario/invoke_scenario.yaml', 'w') as f:
        yaml.dump({'invoke_scenario': scenario_list}, f, sort_keys=False)

    logger.info('Invoke scenario generated successfully')


if __name__ == "__main__":
    abs_path = os.path.abspath(os.path.dirname(__file__))

    # load and set variables from config files
    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)
    no_of_users = int(traffic_config['tool_config']['no_of_users'])
    frequency_limits = traffic_config['tool_config']['frequency_limits']

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as file:
        apim_config = yaml.load(file, Loader=yaml.FullLoader)
    apps = apim_config['apps']

    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as file:
        api_config = yaml.load(file, Loader=yaml.FullLoader)
    apis = api_config['apis']

    with open(abs_path + '/../../data/tool_data/invoke_patterns.yaml', 'r') as file:
        pattern_config = yaml.load(file, Loader=yaml.FullLoader)
    time_patterns = pattern_config['frequency']

    for app in apps:
        app_list[app['name']] = []
        app_api_subs[app['name']] = app['api_subscriptions'].split(',')

    for api in apis:
        methods = []
        for resource in api['resources']:
            methods.append(resource['method'])
        api_resources[api['name']] = methods

    # execute the functions
    genUserApps()
    genScenarioDistribution()

    out_txt = 'Invoke scenario generated successfully. Output written to \'user_details.yaml\' and \'invoke_scenario.yaml\' files'
    logger.info(out_txt)
