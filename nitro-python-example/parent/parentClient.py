import socket
import requests
import json
import boto3
from botocore.config import Config


"""
dynamodb design
table name: AccountTable

colume:
keyId: kms alias id which used for encryption for the private key
Name: account name for this Account
encryptedPrivateKey: encrypted Account private key
address: the address of the Account
encryptedDataKey: the data key used to encrypt the private key
"""

class AccountClient:

    """
    region: region to deploy
    ddbTableName: the dynamodb table name which used to store and retrive the encrypted
    keyId: the kms alais id, used to encrypt the plaintext and decrypt the ciphertext
    cid: cid for vsock client to connect
    port: port for vsock client to connect
    """
    def __init__(self, region, ddbTableName, keyId, cid, port):
        self.__region = region
        self.__ddbTableName = ddbTableName
        self.__keyId = keyId
        self.__cid = cid
        self.__port = port

    """
    generateAccount which will send credential and kms-key-id to enclave server, it waits util enclave
    server sends the encrypted private key back, then save to dynamodb
    name: name for the Account
    """
    def generateAccount(self, name):
        credential = self.__getIAMToken()
        # Create a vsock socket object
        s = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)
        # Connect to the server
        s.connect((self.__cid, self.__port))

        playload = {}
        playload['apiCall'] = "generateAccount"
        playload['credential'] = credential
        playload['keyId'] = self.__keyId
        # Send AWS credential and KMS keyId to the server running in enclave
        s.send(str.encode(json.dumps(playload)))
        # receive data from the server and save to dynamodb with the name
        response = s.recv(65536).decode()
        self.__saveEncryptAccountToDDB(name, response, self.__keyId)
        s.close()

    def sign(self, keyId, name, transaction):
        # Get item from DynamoDB
        dynamodb = boto3.resource('dynamodb', region_name=self.__region)
        table = dynamodb.Table(self.__ddbTableName)
        try:
            response = table.get_item(Key={
                'keyId': keyId,
                'name': name
            })
            if 'Item' not in response:
                print( name + ' not found in DynamoDB')
                return
        except Exception as error:
            print(error)
            return

        encryptedPrivateKey = response['Item']['encryptedPrivateKey']
        # address = response['Item']['address']
        encryptedDataKey = response['Item']['encryptedDataKey']

        credential = self.__getIAMToken()

        playload = {}
        playload['apiCall'] = "sign"
        playload['credential'] = credential
        playload['encryptedPrivateKey'] = encryptedPrivateKey
        playload['encryptedDataKey'] = encryptedDataKey
        playload['transaction'] = transaction

        s = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)
        s.connect((self.__cid, self.__port))
        s.send(str.encode(json.dumps(playload)))
        response = s.recv(65536)
        s.close()
        return response

    def __getIAMToken(self):
        """
        Get the AWS credential from EC2 instance metadata
        """
        r = requests.get(
            "http://169.254.169.254/latest/meta-data/iam/security-credentials/")
        instance_profile_name = r.text

        r = requests.get(
            "http://169.254.169.254/latest/meta-data/iam/security-credentials/%s" % instance_profile_name)
        response = r.json()

        credential = {
            'aws_access_key_id': response['AccessKeyId'],
            'aws_secret_access_key': response['SecretAccessKey'],
            'aws_session_token': response['Token']
        }
        return credential

    def __saveEncryptAccountToDDB(self, name, response, keyId):
        dynamodb = boto3.resource('dynamodb', self.__region)
        table = dynamodb.Table(self.__ddbTableName)
        response_json = json.loads(response)
        print("saved account value to ddb:",response_json)
        table.put_item(Item={
            'name': name,
            'keyId': keyId,
            'encryptedPrivateKey': response_json['encryptedPrivateKey'],
            'address': response_json['address'],
            'encryptedDataKey': response_json['encryptedDataKey'],
        })


def main():
    # check dynamodb AccountTable exists or not, if not exists, create it
    my_config = Config(
    region_name = 'ap-southeast-1'
    )

    client = boto3.client('dynamodb', config=my_config)
    try:
        client.describe_table(TableName='AccountTable')
    except:
         client.create_table(
            TableName = "AccountTable",
            KeySchema = [
                {'AttributeName': 'keyId', 'KeyType': 'HASH'},
                {'AttributeName': 'name', 'KeyType': 'RANGE'}
            ],
            AttributeDefinitions = [
                {'AttributeName': 'keyId', 'AttributeType': 'S'},
                {'AttributeName': 'name', 'AttributeType': 'S'}
            ],
            ProvisionedThroughput = {
                'ReadCapacityUnits': 10,
                'WriteCapacityUnits': 10
            }
        )
    # generate a client and demo it
    client = AccountClient("ap-southeast-1", "AccountTable", "c314a998-b78f-44f8-8c12-6f426ccf89fd", 16, 5000)
    client.generateAccount("Account2")

    # test transaction
    transaction = {
        'value': 1000000000,
        'to': '0xF0109fC8DF283027b6285cc889F5aA624EaC1F55',
        'nonce': 0,
        'chainId': 4,
        'gas': 100000,
        'gasPrice' :234567897654321
    }
    # with defined kms id, you should replace it with yours
    signedValue = client.sign('c314a998-b78f-44f8-8c12-6f426ccf89fd', "Account2", transaction)
    print("signed with value: ", signedValue)

if __name__ == '__main__':
    main()
