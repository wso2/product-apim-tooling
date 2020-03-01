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
import pandas as pd
import yaml
from datetime import datetime
from utils import log
from faker import Factory
from scipy.stats import norm

# variables
logger = log.setLogger('gen_invoke_scenario')

tenant_name = ''
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
multi_tenancy_enabled = 'false'
app_tenant = {}

fake_generator = Factory.create()

# setup configurations
abs_path = os.path.abspath(os.path.dirname(__file__))

try:
    with open(abs_path + '/../../../../config/apim.yaml', 'r') as config_file:
        apim_config = yaml.load(config_file, Loader=yaml.FullLoader)
    multi_tenancy = apim_config['multi_tenancy']['enabled']

    if str(multi_tenancy).lower() == 'false':
        tenant_name = apim_config['main_tenant']['tenant_name']
    else:       # read and set tenant names for each app
        with open(abs_path + '/../../data/scenario/tenant_details.yaml', 'r') as f:
            tenant_apps = yaml.load(f, Loader=yaml.FullLoader)['tenant_apps']

        for tenant in tenant_apps:
            for app in tenant_apps.get(tenant):
                app_tenant[app] = tenant

    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)
    user_agents = traffic_config['user_agents']

    with open(abs_path + '/../../../../config/api_details.yaml', 'r') as file:
        api_config = yaml.load(file, Loader=yaml.FullLoader)
    apis = api_config['apis']

    with open(abs_path + '/../../../../config/user-settings.yaml', 'r') as file:
        user_settings = yaml.load(file, Loader=yaml.FullLoader)
    ip_dataset_name = user_settings['resources']['ip_database']
    dataset_column_order = user_settings['resources']['column_order']

except FileNotFoundError as e:
    out_txt = 'FileNotFoundError in gen_invoke_scenario.py: {}: {}'.format(e.strerror, e.filename)
    logger.exception(out_txt)
    sys.exit()


def getPath(api_name, method):
    """
    This function will return the invoke path for a given api and http method (for a single tenant environment)
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
                    if tenant_name.lower() == 'super' or tenant_name.lower() == "carbon.super":
                        return context + '/' + version + '/' + str(resource.get('path'))
                    else:
                        return 't/' + tenant_name + '/' + context + '/' + version + '/' + str(resource.get('path'))


def getMultiTenantPath(api_name, method, app_name):
    """
    This function will return the invoke path for a given api and http method (for a multi tenant setup)
    :param api_name: Name of the API
    :param method: Method of the API
    :param app_name: Name of the application
    :return: Invoke path
    """
    global apis, app_tenant

    for api in apis:
        if api.get('name') == api_name:
            context = str(api.get('context'))
            version = str(api.get('version'))
            resources = api.get('resources')
            tenant_name = app_tenant.get(app_name)
            for resource in resources:
                if resource.get('method') == method:
                    return 't/' + tenant_name + '/' + context + '/' + version + '/' + str(resource.get('path'))


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
        temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(int(ip_decs[0]), int(ip_decs[1])))
        while temp_ip in used_ips:
            temp_ip = ipaddress.IPv4Address._string_from_ip_int(random.randint(int(ip_decs[0]), int(ip_decs[1])))

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
        with open(abs_path + '/../../data/scenario/user_details.yaml', 'r') as user_file:
            userlist = yaml.load(user_file, Loader=yaml.FullLoader)

        cookie_list = genUniqueCookieList(len(userlist['users']))

        for user in userlist['users']:
            username = user['username']
            country = user['country']
            user_country.update({username: country})
            user_ip.update({username: ipGen(country)})
            user_cookie.update({username: cookie_list.pop()})

        # read user token csv file
        user_token = pd.read_csv(abs_path + '/../../data/scenario/api_invoke_tokens.csv')

    except FileNotFoundError as e:
        out_txt = 'FileNotFoundError in gen_invoke_scenario.py: {}: {}'.format(e.strerror, e.filename)
        logger.exception(out_txt)
        sys.exit()
    except pd.errors.EmptyDataError as e:
        out_txt = 'EmptyDataError in gen_invoke_scenario.py: {}'.format(str(e))
        logger.exception(out_txt)
        sys.exit()
    except Exception as e:
        out_txt = 'Exception in gen_invoke_scenario.py: {}'.format(e)
        logger.exception(out_txt)
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
        with open(abs_path + '/../../data/scenario/invoke_scenario.yaml') as file:
            invoke_scenario = yaml.load(file, Loader=yaml.FullLoader)
        scenario_data = invoke_scenario['invoke_scenario']

    except FileNotFoundError as e:
        out_txt = 'FileNotFoundError in gen_invoke_scenario.py: {}: {}'.format(e.strerror, e.filename)
        logger.exception(out_txt)
        sys.exit()

    # check whether the user count is valid (not more than the created number of users)
    for item in scenario_data:
        user_count = int(item.get('no_of_users'))
        total_no_of_user_combinations += user_count
    
        if total_no_of_user_combinations > existing_no_of_user_combinations:
            # invalid no of users (cannot execute the scenario)
            logger.error(
                "Invalid number of user count declared in 'invoke_scenario.yaml'. Expected {} user combinations. Found {}.".format(existing_no_of_user_combinations, total_no_of_user_combinations))
            logger.critical("Cannot run the tool!")
            sys.exit()

    # generate the scenario pool
    for item in scenario_data:
        app_name = item.get('app_name')
        user_count = int(item.get('no_of_users'))
        time_pattern = item.get('time_pattern')
        invokes = item.get('api_calls')

        users = []
        for i in range(user_count):
            users.append(users_apps.get(app_name).pop())

        for invoke in invokes:
            api_name = invoke.get('api')
            method = invoke.get('method')
            call_median = int(invoke.get('no_of_requests'))

            if str(multi_tenancy).lower() == 'false':
                full_path = getPath(api_name, method)
            else:
                full_path = getMultiTenantPath(api_name, method, app_name)

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

    with open(abs_path + '/../../data/scenario/token_ip_cookie.csv', 'w') as file:
        file.write(write_str)

    # saving scenario pool to a pickle file
    pickle.dump(scenario_pool, open(abs_path + "/../../data/runtime_data/scenario_pool.sav", "wb"))

    logger.info("User scenario distribution generated successfully")
