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
import requests
import urllib3
import os
from utils import util_methods

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
abs_path = os.path.abspath(os.path.dirname(__file__))


def getIDSecret(gateway_protocol, gateway_host, gateway_port, endpoint, admin_username, admin_b64):
    """
    This function will send http request to obtain client id, client secret.
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port to obtain key, secret
    :param endpoint: Endpoint to obtain key, secret
    :param admin_username: Username of the tenant admin
    :param admin_b64: Base64 encrypted username:password of the tenant admin
    :return: Client ID and client secret
    """
    
    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Basic {}".format(admin_b64),
        "Content-Type": "application/json"
    }
    data = {
        "callbackUrl": "www.google.lk",
        "clientName": "rest_api_publisher",
        "owner": admin_username,
        "grantType": "password refresh_token",
        "saasApp": "true"
    }
    data = json.dumps(data)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "getIDSecret(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url
        
        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return response['clientId'], response['clientSecret']
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None, None

    except Exception as err:
        code = 521
        log_txt = "getIDSecret(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)

        return None, None


def getAccessToken(gateway_protocol, gateway_host, gateway_port, endpoint, b64_encoded_value, scope, admin_username, admin_password):
    """
    This function will obtain an access token for a given scope.
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port to obtain access token
    :param endpoint: Endpoint to obtain access token
    :param b64_encoded_value: Base64 encrypted value of client_id:client_secret
    :param scope: Scope of the token
    :param admin_username: Username of the tenant admin
    :param admin_password: Password of the tenant admin
    :return: Access token and refresh token
    """
    
    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Basic {}".format(b64_encoded_value),
        "Content-Type": "application/x-www-form-urlencoded"
    }
    data = {
        'grant_type': 'password',
        'username': admin_username,
        'password': admin_password,
        'scope': scope
    }

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "getAccessToken(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return response['access_token'], response['refresh_token']
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None, None

    except Exception as err:
        code = 521
        log_txt = "getAccessToken(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        
        return None, None


def createAPI(gateway_protocol, gateway_host, gateway_port, endpoint, create_token, name, desc, path, version, swagger, tags, throttling_tier, visibility, production_endpoint, sandbox_endpoint, admin_username):
    """
    This function will create a given API in WSO2 API Manager
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param create_token: API creation token
    :param name: API name
    :param desc: API description
    :param path: API path
    :param version: API version
    :param swagger: Swagger definition (swagger string)
    :param tags: API tag list
    :param throttling_tier: Throttling tier of the API
    :param visibility: Visibility of the API
    :param production_endpoint: Production endpoint
    :param sandbox_endpoint: Sandbox endpoint
    :param admin_username: Username of the tenant admin
    :return: ID of the create API
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Bearer {}".format(create_token),
        "Content-Type": "application/json"
    }
    data = {
        "name": name,
        "description": desc,
        "context": path,
        "version": version,
        "provider": admin_username,
        "apiDefinition": str(swagger),
        "wsdlUri": None,
        "responseCaching": "Disabled",
        "cacheTimeout": 300,
        "destinationStatsEnabled": "false",
        "isDefaultVersion": "false",
        "type": "HTTP",
        "transport": [
            "http",
            "https"
        ],
        "tags": tags,
        "tiers": [throttling_tier],
        "maxTps": None,
        "visibility": visibility,
        "visibleRoles": [],
        "endpointConfig": "{\"production_endpoints\":{\"url\":\"%s\",\"config\":null},\"sandbox_endpoints\":{\"url\":\"%s\",\"config\":null},\"endpoint_type\":\"http\"}" % (production_endpoint, sandbox_endpoint),
        "endpointSecurity":    {
            "username": "user",
            "type": "basic",
            "password": "pass"
        },
        "gatewayEnvironments": "Production and Sandbox",
        "sequences": [],
        "subscriptionAvailability": None,
        "subscriptionAvailableTenants": [],
        "businessInformation":    {
            "businessOwnerEmail": None,
            "technicalOwnerEmail": None,
            "technicalOwner": None,
            "businessOwner": None
        },
        "corsConfiguration":    {
            "accessControlAllowOrigins": ["*"],
            "accessControlAllowHeaders":       [
                "authorization",
                "Access-Control-Allow-Origin",
                "Content-Type",
                "SOAPAction"
            ],
            "accessControlAllowMethods":       [
                "GET",
                "PUT",
                "POST",
                "DELETE",
                "PATCH",
                "OPTIONS"
            ],
            "accessControlAllowCredentials": "false",
            "corsConfigurationEnabled": "false"
        }
    }
    data = json.dumps(data).replace('\\\\\\', '\\')

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "createAPI(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 201:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            api_id = response['id']

            with open(abs_path + '/../../../data/runtime_data/api_ids.csv', 'a+') as file:
                file.write(api_id+'\n')

            return api_id

        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None

    except Exception as err:
        code = 521
        log_txt = "createAPI(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        
        return None


def publishAPI(gateway_protocol, gateway_host, gateway_port, endpoint, publish_token, api_id):
    """
    This function will publish an API
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param publish_token: Token to publish
    :param api_id: ID of the API
    :return: True if API published successfully. False otherwise
    """

    url = "{}://{}:{}{}/change-lifecycle?apiId={}&action=Publish".format(gateway_protocol, gateway_host, gateway_port, endpoint, api_id)
    headers = {
        "Authorization": "Bearer {}".format(publish_token),
        "Content-Type": "application/json"
    }
    data = {}

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text

        log_txt = "publishAPI(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "publishAPI(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        
        return False


def createApplication(gateway_protocol, gateway_host, gateway_port, endpoint, subs_token, name, desc, throttling_tier):
    """
    This function will create a given application
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param subs_token: Subscribe access token
    :param name: App name
    :param desc: App description
    :param throttling_tier: Throttling tier of the app
    :return: ID of the created app
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Bearer {}".format(subs_token),
        "Content-Type": "application/json"
    }
    data = {
        "throttlingTier": throttling_tier,
        "description": desc,
        "name": name,
        "callbackUrl": "http://my.server.com/callback"
    }
    data = json.dumps(data)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "createApplication(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 201:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            app_id = response['applicationId']

            with open(abs_path + '/../../../data/runtime_data/app_ids.csv', 'a+') as file:
                file.write(app_id+'\n')

            return app_id

        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None

    except Exception as err:
        code = 521
        log_txt = "createApplication(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        
        return None


def subscribe(gateway_protocol, gateway_host, gateway_port, endpoint, subs_token, subs_tier, api_id, app_id):
    """
    This function will subscribe given app to given API
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param subs_token: API subscription token
    :param subs_tier: Subscription tier
    :param api_id: ID of the API
    :param app_id: ID of the app
    :return: True if subscription successful. False otherwise
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Bearer {}".format(subs_token),
        "Content-Type": "application/json"
    }
    data = {
        "tier": subs_tier,
        "apiIdentifier": api_id,
        "applicationId": app_id
    }
    data = json.dumps(data)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "subscribe(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 201:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "subscribe(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return False


def selfSignupStoreAPI(gateway_protocol, gateway_host, gateway_port, endpoint, username, password, all_fields):
    """
    This function will self signup users through store API (deprecated in version 3.0.0 onwards)
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param username: Username
    :param password: Password
    :param all_fields: Remaining fields seperated by '|' in a string
    :return: True if user signup successful. False otherwise
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept": "application/json"
    }
    data = {
        'action': 'addUser',
        'username': username,
        'password': password,
        'allFieldsValues': all_fields
    }

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text

        log_txt = "selfSignupStoreAPI(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url + ", username: " + username

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "selfSignupStoreAPI(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return False


def selfSignupIS(gateway_protocol, gateway_host, gateway_port, endpoint, admin_b64, username, password, firstname, lastname, email, country, organization, land_no, mobile_no, IM, user_url):
    """
    This function will self signup users through WSO2 Identity server
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param admin_b64: Base64 encrypted username:password of the tenant admin
    :param username: Username
    :param password: Password
    :param firstname: First name of the user
    :param lastname: Last name of the user
    :param email: Email address of the user
    :param country: Country of the user
    :param organization: Organization of the user
    :param land_no: Land phone number of the user
    :param mobile_no: Mobile phone number of the user
    :param IM: IM of the user
    :param user_url: Url of the user
    :return: True if user signup successful. False otherwise
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Basic {}".format(admin_b64),
        "Content-Type": "application/json",
    }
    data = {
        "user": {
            "username": username,
            "realm": "PRIMARY",
            "password": password,
            "claims": [
                    {"uri": "http://wso2.org/claims/givenname","value": firstname},
                    {"uri": "http://wso2.org/claims/emailaddress","value": email},
                    {"uri": "http://wso2.org/claims/lastname","value": lastname},
                    {"uri": "http://wso2.org/claims/mobile","value": mobile_no},
                    {"uri": "http://wso2.org/claims/organization","value": organization},
                    {"uri": "http://wso2.org/claims/telephone","value": land_no},
                    {"uri": "http://wso2.org/claims/country","value": country},
                    {"uri": "http://wso2.org/claims/url","value": user_url},
                    {"uri": "http://wso2.org/claims/im","value": IM}
                ]
        },
        "properties": []
    }
    data = json.dumps(data)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text

        log_txt = "selfSignupIS(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url + ", username: " + username

        if code == 201:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "selfSignupIS(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return False


def genProductionKey(gateway_protocol, gateway_host, gateway_port, endpoint, subs_token, token_validity_period):
    """
    This function will generate production keys
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param subs_token: Token in subscribe scope
    :param token_validity_period: Access token validity period
    :return: Consumer key and consumer secret
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Bearer {}".format(subs_token),
        "Content-Type": "application/json"
    }
    data = {
        "validityTime": token_validity_period,
        "keyType": "PRODUCTION",
        "accessAllowDomains": ["ALL"],
        "scopes": ["am_application_scope", "default"],
        "supportedGrantTypes": ["urn:ietf:params:oauth:grant-type:saml2-bearer", "iwa:ntlm", "refresh_token", "client_credentials", "password"]
    }
    data = json.dumps(data)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        log_txt = "genProductionKey(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return response['consumerKey'], response['consumerSecret']
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None, None

    except Exception as err:
        code = 521
        log_txt = "genProductionKey(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return None, None


def generateInvokeToken(gateway_protocol, gateway_host, gateway_port, endpoint, b64_key_secret, username, password, scope):
    """
    This function will generate API invoke tokens in password grant type for a given user
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param b64_key_secret: Base64 encrypted value of consumer_key:consumer_secret
    :param username: Username of the user
    :param password: Password of the user
    :param scope: Scope to generate tokens
    :return: Access token and refresh token
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, endpoint)
    headers = {
        "Authorization": "Basic {}".format(b64_key_secret),
        "Content-Type": "application/x-www-form-urlencoded"
    }
    data = {
        "grant_type": "password",
        "username": username,
        "password": password,
        "scope": scope
    }

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text
        response = json.loads(res_txt)

        expires_in = response['expires_in']
        if int(expires_in) <= 4800:
            util_methods.log('traffic-tool.log', "WARN", "Your token expiration time is {}. It is recommended to increase expiration time to prevent unnecessary token expirations!".format(str(expires_in)))

        log_txt = "generateInvokeToken(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return response['access_token'], response['refresh_token']
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return None, None

    except Exception as err:
        code = 521
        log_txt = "generateInvokeToken(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return None, None


def deleteAppAPI(gateway_protocol, gateway_host, gateway_port, endpoint, access_token, id):
    """
    This function will delete a given API or application
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param endpoint: Endpoint
    :param access_token: Access token to delete
    :param id: ID of the API or application
    :return: True if the deletion successful. False otherwise
    """

    url = "{}://{}:{}{}/{}".format(gateway_protocol, gateway_host, gateway_port, endpoint, id)
    headers = {
        "Authorization": "Bearer {}".format(access_token),
        "Content-Type": "application/json"
    }

    try:
        response = requests.delete(url=url, headers=headers, verify=False)
        code = response.status_code
        res_txt = response.text

        log_txt = "deleteAppAPI(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: DELETE" + ", url: " + url

        if code == 200:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "deleteAppAPI(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: DELETE" + ", url: " + url
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return False


def removeUserSOAP(gateway_protocol, gateway_host, gateway_port, soap_endpoint, username, admin_b64):
    """
    This function will remove a given user account from carbon (Uses SOAP endpoint)
    :param gateway_protocol: Running protocol of the gateway
    :param gateway_host: Host IP address of the gateway
    :param gateway_port: Port
    :param soap_endpoint: SOAP Endpoint
    :param username: Username of the user/ account
    :param admin_b64: Base64 encrypted username:password of the tenant admin
    :return: True if the user account removed. False otherwise
    """

    url = "{}://{}:{}{}".format(gateway_protocol, gateway_host, gateway_port, soap_endpoint)
    headers = {
        "Authorization": "Basic {}".format(admin_b64),
        "Content-Type": "text/xml",
        "SOAPAction": "urn:deleteUser"
    }
    data = """
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://service.ws.um.carbon.wso2.org">
        <soapenv:Header/>
        <soapenv:Body>
            <ser:deleteUser>
                <!--Optional:-->
                <ser:userName>{}</ser:userName>
            </ser:deleteUser>
        </soapenv:Body>
    </soapenv:Envelope>""".format(username)

    try:
        response = requests.post(url=url, headers=headers, data=data, verify=False)
        code = response.status_code
        res_txt = response.text

        log_txt = "removeUserSOAP(). responseCode: " + str(code) + ", responseMessage: " + str(res_txt) + ", method: POST" + ", url: " + url + ", username: " + username

        if code == 202:
            util_methods.log('traffic-requests.log', "SUCCESS", log_txt)
            return True
        else:
            util_methods.log('traffic-requests.log', "FAILED", log_txt)
            return False

    except Exception as err:
        code = 521
        log_txt = "removeUserSOAP(). responseCode: " + str(code) + ", errorLog: " + str(err) + ", method: POST" + ", url: " + url + ", body: " + str(data)
        util_methods.log('traffic-requests.log', "ERROR", log_txt)
        return False
