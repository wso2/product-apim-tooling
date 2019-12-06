
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


import rstr
from faker import Faker
import argparse
import os
import yaml
from datetime import datetime
import sys

# global variables
faker = Faker()
usernames = []
scenario_name = None
no_of_users = 0
abs_path = os.path.abspath(os.path.dirname(__file__))

# load and set configurations
try:
    with open(abs_path+'/../../../../config/traffic-tool.yaml', 'r') as file:
        traffic_config = yaml.load(file, Loader=yaml.FullLoader)
    scenario_name = traffic_config['scenario_name']
    no_of_users = int(traffic_config['tool_config']['no_of_users'])

    if no_of_users <= 0:
        print('[ERROR] {} gen_user_details.py: Invalid user count: {}'.format(str(datetime.now()), str(no_of_users)))
        sys.exit()

except FileNotFoundError as e:
    print('[ERROR] {} gen_user_details.py: {}: {}'.format(str(datetime.now()), e.strerror, e.filename))
    sys.exit()


'''
    This function will return a username and password for a given user (username and password are considered as the same)
'''
def genUnPw(firstname:str, num:int):
    username = firstname.lower() + str(num)
    if len(username) < 5:
        username += '123'
    usernames.append(username)
    return username


'''
    This function will generate random user details (for a single user)
'''
def generateUser(num:int):
    user = []
    firstname = faker.first_name()
    username = genUnPw(firstname, num)
    user.append(username)
    user.append(username)
    user.append(firstname)
    user.append(faker.last_name())
    user.append(faker.company())
    user.append(faker.country())
    user.append(firstname.lower() + str(num) + '@gmail.com')
    user.append(faker.phone_number())
    user.append(faker.phone_number())
    user.append(firstname.lower() + str(num))
    user.append('http://{0}.{1}.com/{2}/?{3}'.format(rstr.domainsafe(), rstr.letters(3), rstr.urlsafe(), rstr.urlsafe()))

    return user


'''
    This function will generate app name, username pattern according to the scenario
'''
def app_userScenario():
    finalArr = []
    finalStr = ""

    individual_app_users = int(no_of_users * 3/5)
    only_onlineShopping = int(individual_app_users*1/4)
    only_cricScore = int(individual_app_users*1/6)
    only_taxi = individual_app_users - (only_onlineShopping + only_cricScore)

    all_app = int((no_of_users - individual_app_users) * 1/4)
    shopping_taxi = int((no_of_users - individual_app_users) * 1/4)
    shopping_cricScore = int((no_of_users - individual_app_users) * 1/8)
    taxi_cricScore = no_of_users - individual_app_users - (all_app + shopping_taxi + shopping_cricScore)

    finalArr.append([ usernames[i]+",Online Shopping\n" for i in range(0, only_onlineShopping) ])     # only online shopping app users
    finalArr.append([ usernames[i]+",CricScore\n" for i in range(only_onlineShopping, only_onlineShopping+only_cricScore) ])     # only cricscore app users
    finalArr.append([ usernames[i]+",Taxi\n" for i in range(only_onlineShopping+only_cricScore, only_onlineShopping+only_cricScore+only_taxi) ])     # only taxi app users

    v1 = individual_app_users + shopping_taxi
    finalArr.append([ usernames[i]+",Online Shopping\n" for i in range(individual_app_users, v1) ])  # both shopping and taxi app users
    finalArr.append([ usernames[i]+",Taxi\n" for i in range(individual_app_users, v1) ])

    v2 = v1 + shopping_cricScore
    finalArr.append([ usernames[i]+",Online Shopping\n" for i in range(v1, v2) ])  # both shopping and cricscore app users
    finalArr.append([ usernames[i]+",CricScore\n" for i in range(v1, v2) ])

    v3 = v2 + taxi_cricScore
    finalArr.append([ usernames[i]+",Taxi\n" for i in range(v2, v3) ])  # both taxi and cricscore app users
    finalArr.append([ usernames[i]+",CricScore\n" for i in range(v2, v3) ])

    v4 = v3 + all_app
    finalArr.append([ usernames[i]+",Online Shopping\n" for i in range(v3, v4) ])  # all 3 app users
    finalArr.append([ usernames[i]+",Taxi\n" for i in range(v3, v4) ])
    finalArr.append([ usernames[i]+",CricScore\n" for i in range(v3, v4) ])

    for outer in finalArr:
        for inner in outer:
            finalStr += inner

    file = open(abs_path+'/../../data/scenario/{}/data/user_app_pattern.csv'.format(scenario_name), 'w')
    file.write(finalStr)
    file.close()

    print('[INFO] {}: User app pattern generation successful!'.format(str(datetime.now())))


'''
    This function will generate given number of users and write data to a csv file
    data format: <username>, <password>, <first_name>, <last_name>, <organization>, <country>, <email>, <no(land)>, <no(mobile)>, <IM>, <url>
    delimiter : '$$ '
'''
def genUsersCSV():
    csvString = ""
    for i in range(no_of_users):
        userArr = generateUser(i+1)
        for ele in userArr:
            csvString += ele + '$$ '
        csvString += '\n'

    file = open(abs_path+'/../../data/scenario/{}/data/user_generation.csv'.format(scenario_name), 'w')
    file.write(csvString)
    file.close()
    print('[INFO] {}: User generation successful!'.format(str(datetime.now())))


# execute
parser = argparse.ArgumentParser("generate user details")
parser.add_argument("option", help="Pass 0 to generate only user details. Pass 1 to generate user details and the scenario distribution", type=int)
args = parser.parse_args()

if args.option == 0:
    genUsersCSV()
elif args.option == 1:
    genUsersCSV()
    app_userScenario()
else:
    print("[INFO] {}: Invalid argument value {}!".format(str(datetime.now()), args.option))
