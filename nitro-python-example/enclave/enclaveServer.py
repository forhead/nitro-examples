import base64
import json
import os
import socket
import hashlib
from Crypto import Random
from Crypto.Cipher import AES
from eth_account import Account
import secrets
from web3.auto import w3

from kms import nitroKms

class nitroServer:

    def __init__(self, region):
        self.__region = region
        pass

    def generateAccount(self, credential, keyId):
        nitro_kms = nitroKms()

        # generate eth private key and calculate the account address
        priv = secrets.token_hex(32)
        private_key = "0x" + priv
        print ("SAVE BUT DO NOT SHARE THIS:", private_key)
        acct = Account.from_key(private_key)
        print("Address:", acct.address)

        # Generate data key by KMS GenerateDataKey API with attestation
        datakeyTextBase64 = nitro_kms.call_kms_generate_datakey(credential, keyId)  
        
        datakey_split = datakeyTextBase64.split("\n")
        datakeyCipherText = datakey_split[0].split(":")[1].strip()
        datakeyPlaintext = datakey_split[1].split(":")[1].strip()

        print("datakeyCipherText: ", datakeyCipherText)
        print("datakeyPlaintext: ", datakeyPlaintext)

        # Encrypt User Private_Key using datakey from KMS, by client-side AES
        # Input datakey string and plaintext string
        aesclient = AESCipher(datakeyPlaintext)

        #convert bytes[] to string
        encrypted_privatekey = str(aesclient.encrypt(private_key),encoding='utf-8')
        content = {
            'encryptedPrivateKey': encrypted_privatekey,
            'address': acct.address,
            'encryptedDataKey': datakeyCipherText.strip()
        }
        print(content)
        return content

    # return the private key's hash value, not an implement of crypty sign operation
    def sign(self, credential, encryptedPrivateKey, encryptedDataKey, transaction):
        nitro_kms = nitroKms()

        # Decrypt encrypted data_key by KMS Decrypt API with attestation
        # Key metadata included in Ciphertextblob, return bytes 
        datakeyPlainTextBase64 = nitro_kms.call_kms_decrypt(credential,encryptedDataKey)
        print("decrypted datakey",datakeyPlainTextBase64)
        datakeyPlaintextBase64String = datakeyPlainTextBase64.split(":")[1].strip()

        # Decrypt private_key with plain data key
        # Input datakey string and encrypted string
        aesclient = AESCipher(datakeyPlaintextBase64String)
        private_key = aesclient.decrypt(encryptedPrivateKey)

        # sign
        transaction_signed = w3.eth.account.sign_transaction(transaction, private_key)
        response_plaintext = {"transaction_signed": transaction_signed.rawTransaction.hex(),
                                      "transaction_hash": transaction_signed.hash.hex()}
        return response_plaintext


class AESCipher(object):

    def __init__(self, key):
        self.bs = AES.block_size
        self.key = hashlib.sha256(key.encode()).digest()

    def encrypt(self, raw):
        raw = self._pad(raw)
        iv = Random.new().read(AES.block_size)
        cipher = AES.new(self.key, AES.MODE_CBC, iv)
        return base64.b64encode(iv + cipher.encrypt(raw.encode()))

    def decrypt(self, enc):
        enc = base64.b64decode(enc)
        iv = enc[:AES.block_size]
        cipher = AES.new(self.key, AES.MODE_CBC, iv)
        return self._unpad(cipher.decrypt(enc[AES.block_size:])).decode('utf-8')

    def _pad(self, s):
        return s + (self.bs - len(s) % self.bs) * chr(self.bs - len(s) % self.bs)

    @staticmethod
    def _unpad(s):
        return s[:-ord(s[len(s)-1:])]
    

def main():
    print("nitro server started ...")

    # Create a vsock socket object
    s = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)
    # Listen for connection from any CID
    cid = socket.VMADDR_CID_ANY
    # The port should match the client running in parent EC2 instance
    port = 5000
    # Bind the socket to CID and port
    s.bind((cid, port))
    # Listen for connection from client
    s.listen()

    # read region from environment variable
    region = os.getenv("REGION")
    # server client which call generateWallet or sign
    server = nitroServer(region)

    while True:
        c, addr = s.accept()

        # Get AWS credential sent from parent instance
        payload = c.recv(4096)
        payload_json = json.loads(payload.decode())
        print("payload json: {}".format(payload_json))

        apiCall = payload_json["apiCall"]

        if apiCall == "generateAccount":
            print("generateWallet request")
            credential = payload_json["credential"]
            keyId = payload_json["keyId"]
            result = server.generateAccount(credential, keyId)
            # send back to parent instance
            c.send(str.encode(json.dumps(result)))
            print("generateWallet finished")

        elif apiCall == "sign":
            print("sign request")
            credential = payload_json["credential"]
            transaction = payload_json["transaction"]
            encryptedPrivateKey = payload_json["encryptedPrivateKey"]
            encryptedDataKey = payload_json['encryptedDataKey']

            signedStr = server.sign(
                credential, encryptedPrivateKey, encryptedDataKey, transaction)
            c.send(str.encode(json.dumps(signedStr)))

            print("sign fihished")
        else:
            print("nothing to do")

        c.close()

if __name__ == '__main__':
    main()
