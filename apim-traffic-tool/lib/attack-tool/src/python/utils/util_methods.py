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

import json
import math
import os
from collections import defaultdict
import random
import requests
import numpy as np
# disabling warnings
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def generate_random_json():
    """
    Generates a random json string
    :return: A random json string
    """
    temp = defaultdict()
    for i in range(random.randint(0, 10)):
        temp[generate_random_string(random.randint(0, 10))] = generate_random_string(random.randint(0, 10))
    return json.dumps(temp)


def generate_random_string(size):
    """
    Generates random ascii string of given size.

    :param size: integer value
    :return: a string of given size
    """
    out_str = ''
    Range = list(range(65, 90)) + list(range(97, 122))
    for i in range(size):
        # a = random.randint(65, 90)
        a = random.choice(Range)
        out_str += chr(a)
    return out_str


def write_to_file(path, data, mode):
    """
    writes the data to a file in the given path
    :param path: file path
    :param data: data to be written
    :param mode: mode to open the file (ex: w, w+, a)

    """
    temp_path = os.path.abspath(os.path.join(__file__, path))

    with open(temp_path, mode) as file:
        file.write(data + "\n")


def send_simple_request(request_path, method, token, ip, cookie, accept, content_type, user_agent, path_params=None, query_params=None, payload=None):
    """
    Sending a http request using the given parameters
    :param request_path: path used to send the request
    :param method: The request method
    :param token: An access token to be included in the header
    :param ip: An IP address to be included in the header
    :param cookie: An user cookie to be included in the header
    :param accept:
    :param content_type: Content-type of the request
    :param user_agent: An user agent to be included in the header
    :param path_params: If there are any path parameters, default value is none
    :param query_params: If there are any query parameters, default value is none
    :param payload: If there is a payload to be attached to the request body, default value is none
    :return: A response object
    """

    # append query/path parameters
    if path_params is not None:
        request_path += "/{}".format(path_params)
    elif query_params is not None:
        request_path += "?{}".format(query_params)

    if payload is not None:
        request_body = payload
    else:
        request_body = generate_random_json()

    header_data = {
        'User-Agent': user_agent,
        'accept': accept,
        'client-ip': '{}'.format(ip),
        'x-forwarded-for': '{}'.format(ip),
        'cookie': '{}'.format(cookie),
        'Authorization': 'Bearer {}'.format(token),
        'Content-Type': content_type
    }

    # default response object
    r = requests.Response()
    r.status_code = 405

    try:
        if method == 'GET':
            r = requests.get(url=request_path, headers=header_data, timeout=(15, 30), verify=False)
        elif method == 'POST':
            r = requests.post(url=request_path, headers=header_data, data=request_body, timeout=(15, 30), verify=False)
        elif method == 'PUT':
            r = requests.put(url=request_path, headers=header_data, data=request_body, timeout=(15, 30), verify=False)
        elif method == 'DELETE':
            r = requests.delete(url=request_path, headers=header_data, timeout=(15, 30), verify=False)
        elif method == 'PATCH':
            r = requests.patch(url=request_path, headers=header_data, data=request_body, timeout=(15, 30), verify=False)

        return r
    except Exception:
        raise


def generate_biased_random(minimum, maximum, exp):
    """
    Generates a random number with a bias to either minimum or maximum.
    :param minimum: Lower limit of number generation
    :param maximum: Upper limit of number generation
    :param exp: exp = 0 : number is the maximum number
                0 < exp < 1 : number is closer to maximum;
                exp = 1 : number is unbiased;
                exp > 1 : number is closer to minimum;

    :return: A biased random number
    """
    return math.floor(minimum + (maximum - minimum) * pow(random.random(), exp))


def generate_method_invoke_pattern(app):
    """
    Generates a method invoke pattern according to the likelihood of each scenario getting called.
    :param app: A list containing scenarios
    :return: Method invoke pattern list with scenario indexes as elements
    """
    probability_list = []
    iterations = 0
    DIFF_THRESHOLD = 0.5
    PROB_ADJUSTMENT = 0.075

    for scenario in app:
        iterations += scenario[0]
        probability_list.append(scenario[0])

    probability_list = list(map(lambda x: x / iterations, probability_list))

    # increase probabilities if it's too small compared to max value
    for i in range(len(probability_list)):
        max_pro = max(probability_list)
        if max_pro - probability_list[i] >= DIFF_THRESHOLD:
            probability_list[i] = probability_list[i] + PROB_ADJUSTMENT
            probability_list[probability_list.index(max_pro)] = max_pro - PROB_ADJUSTMENT

    # prepare request pattern from list indices
    method_pattern = np.random.choice(len(app), size=iterations, p=probability_list)
    return method_pattern


def cleanup(process_list):
    """
    Terminates a given list of processes
    :param process_list: List of processes that need to be terminated
    :return: None
    """
    for p in process_list:
        p.terminate()


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
