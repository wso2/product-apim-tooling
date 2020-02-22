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

import ipaddress
import os
import pickle
import random
import string
import sys
from datetime import datetime

import pandas as pd
import yaml
from faker import Factory
from scipy.stats import norm

# variables
scenario_name = None
ip_dataset_name = None
apis = None
invoke_scenario = None
user_country = {}
user_ip = {}
user_cookie = {}
users_apps = {}
scenario_pool = {}
scenario_distribution = []
existing_no_of_user_combinations = 0  # to validate the user count
total_no_of_user_combinations = 0
used_ips = []
ip_dataset = None

fake_generator = Factory.create()

# setup configurations
abs_path = os.path.abspath(os.path.dirname(__file__))

try:
    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)
    scenario_name = traffic_config['scenario_name']
    user_agents = traffic_config['user_agents']

    with open(abs_path + '/../../../../config/apim.yaml', 'r') as file:
        apim_config = yaml.load(file, Loader=yaml.FullLoader)
    apis = apim_config['apis']

    with open(abs_path + '/../../../../config/user-settings.yaml', 'r') as file:
        user_settings = yaml.load(file, Loader=yaml.FullLoader)
    ip_dataset_name = user_settings['resources']['ip_database']
    dataset_column_order = user_settings['resources']['column_order']

except FileNotFoundError as e:
    print('[ERROR] {} gen_invoke_scenario.py: {}: {}'.format(str(datetime.now()), e.strerror, e.filename))
    sys.exit()


def log(tag, write_string):
    """
    This function will write the given log output to the log.txt file
    :param tag: Log tag
    :param write_string: Message to be written
    :return: None
    """
    with open(abs_path + '/../../../../logs/traffic-tool.log', 'a+') as log_file:
        log_file.write("[{}] ".format(tag) + str(datetime.now()) + ": " + write_string + "\n")


def getPath(api_name, method):
    """
    This function will return the invoke path for a given api and http method
    :param api_name: Name of the API
    :param method: Method of the API
    :return: Invoke path
    """
    global apis

    for api in apis:
        if api.get('name') == api_name:
            context = str(api.get('context'))
            version = str(api.get('version'))
            resources = api.get('resources')
            for resource in resources:
                if resource.get('method') == method:
                    return context + '/' + version + '/' + str(resource.get('path'))


def varySlightly(median):
    """
    This function will return an integer slightly varied to the given median
    :param median: An integer
    :return: An integer slightly varied to the median
    """
    st_div = 1
    return abs(int(norm.rvs(loc=median, scale=st_div)))


def ipGen(country):
    """
    This function will return a randomly generated ipv4 address for a given country
    :param country: A country name
    :return: IPv4 address for a given country
    """
    global used_ips, ip_dataset, dataset_column_order

    country_name_index = int(dataset_column_order['country_name']) - 1
    ip_from_index = int(dataset_column_order['ip_from']) - 1
    ip_to_index = int(dataset_column_order['ip_to']) - 1

    ip_range = ip_dataset.loc[ip_dataset.iloc[:, country_name_index] == country].iloc[:, ip_from_index:ip_to_index + 1]

    if len(ip_range) > 0:
        ip_decs = ip_range.sample(n=1).values[0]

        random.seed()
        temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(ip_decs[0], ip_decs[1]))
        while temp_ip in used_ips:
            temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(ip_decs[0], ip_decs[1]))

        used_ips.append(temp_ip)

        return temp_ip

    else:
        return fake_generator.ipv4()


def getCookie():
    """
    This function will return a randomly generated cookie
    :return: A cookie
    """
    lettersAndDigits = string.ascii_lowercase + string.digits
    cookie = 'JSESSIONID='
    cookie += ''.join(random.choice(lettersAndDigits) for ch in range(31))
    return cookie


def genUniqueCookieList(count: int):
    """
    This function will return a list of unique cookies
    :param count: Number of cookies to be generated
    :return: A list of unique cookies
    """
    cookie_list = set()
    while len(cookie_list) != count:
        cookie_list.add(getCookie())

    return list(cookie_list)


if __name__ == "__main__":
    '''
        Execute the script and generate the user scenario distribution
        Usage: python3 gen_invoke_scenario.py
        output folders: lib/traffic-tool/data/scenario/ and lib/traffic-tool/data/runtime_data/
    '''

    try:
        # read and load the ip database
        ip_dataset = pd.read_csv(abs_path + '/../../../../resources/libraries/{}'.format(ip_dataset_name), header=None)

        # generate a set of ips and cookies for each user
        with open(abs_path + '/../../data/scenario/{}/data/user_generation.csv'.format(scenario_name)) as file:
            userlist = file.readlines()

            cookie_list = genUniqueCookieList(len(userlist))

            for user in userlist:
                username = user.split('$$ ')[0]
                country = user.split('$$ ')[5]
                user_country.update({username: country})
                user_ip.update({username: ipGen(country)})
                user_cookie.update({username: cookie_list.pop()})

        # read user token csv file
        user_token = pd.read_csv(abs_path + '/../../data/scenario/{}/api_invoke_tokens.csv'.format(scenario_name))

    except FileNotFoundError as e:
        print('[ERROR] {} gen_invoke_scenario.py: {}: {}'.format(str(datetime.now()), e.strerror, e.filename))
        log('ERROR', '{}: {}'.format(e.strerror, e.filename))
        sys.exit()
    except pd.errors.EmptyDataError as e:
        print('[ERROR] {} gen_invoke_scenario.py: {}'.format(str(datetime.now()), e))
        log('ERROR', '{}'.format(str(e)))
        sys.exit()
    except Exception as e:
        print('[ERROR] {} gen_invoke_scenario.py: {}'.format(str(datetime.now()), e))
        log('ERROR', '{}'.format(str(e)))
        sys.exit()

    # filter out unique app names and prepare dictionary
    appNames = user_token.drop_duplicates(subset=['app_name'], keep='first')[['app_name']]

    for row in appNames.itertuples():
        users_apps.update({row.app_name: []})

    # set ips with username, access tokens and append to relevant lists
    for row in user_token.itertuples():
        username = row.username
        app_name = row.app_name
        token = row.access_token
        ip = user_ip.get(username)
        cookie = user_cookie.get(username)
        user_agent = random.choice(user_agents)
        (users_apps[app_name]).append([username, token, ip, cookie, user_agent])
        existing_no_of_user_combinations += 1

        if username not in scenario_pool:
            scenario_pool.update({username: {}})
        scenario_pool.get(username).update({app_name: []})

    try:
        # generate scenario data according to the script and append to the pool
        with open(abs_path + '/../../data/scenario/{}/data/invoke_scenario.yaml'.format(scenario_name)) as file:
            invoke_scenario = yaml.load(file, Loader=yaml.FullLoader)
        scenario_data = invoke_scenario['invoke_scenario']

    except FileNotFoundError as e:
        print('[ERROR] {} gen_invoke_scenario.py: {}: {}'.format(str(datetime.now()), e.strerror, e.filename))
        log('ERROR', '{}: {}'.format(e.strerror, e.filename))
        sys.exit()

    for item in scenario_data:
        app_name = item.get('app_name')
        user_count = int(item.get('no_of_users'))
        time_pattern = item.get('time_pattern')
        invokes = item.get('api_calls')

        # check whether the user count is valid (not more than the created number of users)
        total_no_of_user_combinations += user_count
        if total_no_of_user_combinations > existing_no_of_user_combinations:
            # invalid no of users (cannot execute the scenario)
            log("ERROR",
                "Invalid number of user count declared in 'invoke_scenario.yaml'. Expected {} user combinations. Found {} or more.".format(existing_no_of_user_combinations, total_no_of_user_combinations))
            raise ArithmeticError(
                "Invalid number of user count declared in 'invoke_scenario.yaml'. Expected {} user combinations. Found {} or more.".format(existing_no_of_user_combinations, total_no_of_user_combinations))

        users = []
        for i in range(user_count):
            users.append(users_apps.get(app_name).pop())

        for invoke in invokes:
            api_name = invoke.get('api')
            method = invoke.get('method')
            call_median = int(invoke.get('no_of_requests'))
            full_path = getPath(api_name, method)

            for user in users:  # user[username,token,ip,cookie,user_agent]
                no_of_requests = varySlightly(call_median)
                scenario_pool.get(user[0]).get(app_name).append([no_of_requests, api_name, full_path, user[1], method, user[2], user[3], user[4], time_pattern])
                scenario_distribution.append([api_name, user[1], user[2], user[3]])

    # save scenario data
    write_str = "access_token,api_name,ip_address,user_cookie\n"

    for row in scenario_distribution:
        api_name = row[0]
        access_token = row[1]
        ip_address = row[2]
        user_cookie = row[3]
        write_str += access_token + ',' + api_name + ',' + ip_address + ',' + user_cookie + "\n"

    with open(abs_path + '/../../data/scenario/{}/token_ip_cookie.csv'.format(scenario_name), 'w') as file:
        file.write(write_str)

    # saving scenario pool to a pickle file
    pickle.dump(scenario_pool, open(abs_path + "/../../data/runtime_data/scenario_pool.sav", "wb"))

    log("INFO", "User scenario distribution generated successfully")
    print('[INFO] {}: User scenario distribution generated successfully'.format(str(datetime.now())))
