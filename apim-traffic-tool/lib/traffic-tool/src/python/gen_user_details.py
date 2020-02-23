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


import rstr
from faker import Faker
import argparse
import os
import yaml
from datetime import datetime
import sys
from utils import log

# global variables
logger = log.setLogger('gen_user_details')
faker = Faker()
scenario_name = None
no_of_users = 0
abs_path = os.path.abspath(os.path.dirname(__file__))

# load and set configurations
try:
    with open(abs_path + '/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)
    no_of_users = int(traffic_config['tool_config']['no_of_users'])

    if no_of_users <= 0:
        logger.error('User creation Failed!. Invalid user count: {}'.format(str(no_of_users)))
        sys.exit()

except FileNotFoundError as e:
    logger.exception('User creation Failed!. errLog: {}. Filename: {}'.format(e.strerror, e.filename))
    sys.exit()


def genUsernamePassword(firstname: str, num: int):
    """
    This function will return a username and password for a given user (username and password are considered as the same)
    :param firstname: User's first name
    :param num: User number
    :return: User name for the user
    """
    username = firstname.lower() + str(num)
    if len(username) < 5:
        username += '123'
    return username


def generateUser(num: int):
    """
    This function will generate random user details (for a single user)
    :param num: User number
    :return: Dictionary containing details of the user
    """
    user = {}
    firstname = faker.first_name()
    username = genUsernamePassword(firstname, num)
    user['username'] = username
    user['password'] = username
    user['firstname'] = firstname
    user['lastname'] = faker.last_name()
    user['organization'] = faker.company()
    user['country'] = faker.country()
    user['email'] = firstname.lower() + str(num) + '@gmail.com'
    user['no_land'] = faker.phone_number()
    user['no_mobile'] = faker.phone_number()
    user['IM'] = firstname.lower() + str(num)
    user['url'] = 'http://{0}.{1}.com/{2}/?{3}'.format(rstr.domainsafe(), rstr.letters(3), rstr.urlsafe(), rstr.urlsafe())
    user['applications'] = None
    
    return user


def appUserScenario():
    """
    This function will generate all users and divide them among applications according to the example scenario
    :return: None
    """
    user_list = []
    for i in range(no_of_users):
        user = generateUser(i + 1)
        user_list.append(user)

    individual_app_users = int(no_of_users * 3 / 5)
    only_onlineShopping = int(individual_app_users * 1 / 4)
    only_cricScore = int(individual_app_users * 1 / 6)
    only_taxi = individual_app_users - (only_onlineShopping + only_cricScore)

    all_app = int((no_of_users - individual_app_users) * 1 / 4)
    shopping_taxi = int((no_of_users - individual_app_users) * 1 / 4)
    shopping_cricScore = int((no_of_users - individual_app_users) * 1 / 8)
    taxi_cricScore = no_of_users - individual_app_users - (all_app + shopping_taxi + shopping_cricScore)

    # only online shopping app users
    for i in range(0, only_onlineShopping):
        user_list[i]['applications'] = 'Online Shopping'
    # only cricscore app users
    for i in range(only_onlineShopping, only_onlineShopping + only_cricScore):
        user_list[i]['applications'] = 'CricScore'
    # only taxi app users
    for i in range(only_onlineShopping + only_cricScore, only_onlineShopping + only_cricScore + only_taxi):
        user_list[i]['applications'] = 'Taxi'

    # both shopping and taxi app users
    v1 = individual_app_users + shopping_taxi
    for i in range(individual_app_users, v1):
        user_list[i]['applications'] = 'Online Shopping, Taxi'

    # both shopping and cricscore app users
    v2 = v1 + shopping_cricScore
    for i in range(v1, v2):
        user_list[i]['applications'] = 'Online Shopping, CricScore'

    # both taxi and cricscore app users
    v3 = v2 + taxi_cricScore
    for i in range(v2, v3):
        user_list[i]['applications'] = 'Taxi, CricScore'
    
    # all 3 app users
    v4 = v3 + all_app
    for i in range(v3, v4):
        user_list[i]['applications'] = 'Online Shopping, Taxi, CricScore'

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'w') as file:
        yaml.dump({'users': user_list}, file, sort_keys=False)

    logger.info('User details generated. Users divided among applications according to the example scenario. No of users: {}'.format(str(no_of_users)))


def genUsers():
    """
    This function will generate given number of users and write data to the user_details.yaml file
    :return: None
    """
    user_list = []
    for i in range(no_of_users):
        user = generateUser(i + 1)
        user_list.append(user)

    with open(abs_path + '/../../data/scenario/user_details.yaml', 'w') as file:
        yaml.dump({'users': user_list}, file, sort_keys=False)

    logger.info('User details generated successfully. No of users: {}'.format(str(no_of_users)))


if __name__ == "__main__":
    # execute
    parser = argparse.ArgumentParser("generate user details")
    parser.add_argument("option", help="Pass 0 to generate only user details. Pass 1 to generate user details and the scenario distribution", type=int)
    args = parser.parse_args()

    if args.option == 0:
        genUsers()
    elif args.option == 1:
        appUserScenario()
    else:
        logger.error("Invalid argument {}!".format(args.option))
