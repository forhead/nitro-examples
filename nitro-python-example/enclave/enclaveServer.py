import base64
import json
import os
import socket
import ecdsa
import hashlib
import binascii

from Crypto.Cipher import AES
from kms import NitroKms

class nitroServer:

    # generate a string, not an implementation for crypto wallet
    def generateWallet(self, credential, keyId):
        nitro_kms = NitroKms()
        nitro_kms.set_region( os.getenv("REGION", "ap-southeast-1"))
        nitro_kms.set_credentials(credential)     

        # Generate Random as User private Key from KMS(256 bits)
        random = nitro_kms.kms_generate_random(32)  # return bytes 
        private_key_hex = binascii.hexlify(random).decode('utf-8')  # bytes to Hex

        # Convert private key to ECDSA signing key
        sk = ecdsa.SigningKey.from_string(bytes.fromhex(private_key_hex), curve=ecdsa.SECP256k1, hashfunc = hashlib.sha256)

        # Generate user Public key from ECDSA signing key
        vk = sk.get_verifying_key()
        public_key_hex = vk.to_string().hex()

        # Generate data key by KMS GenerateDataKey API with attestation
        datakey = nitro_kms.kms_generate_data_key(24, keyId)  # return bytes
        plain_datakey = base64.b64encode(datakey[0]).decode('utf-8') # bytes to string
        encrypted_datakey = datakey[1]['CiphertextBlob']

        """Server side encrypt and decrypt by KMS"""
        # Encrypt text(from GenerateRandom) by KMS Encrypt API
        # kms_encrypted = nitro_kms.kms_encrypt(random, key_id)  # Input plaintext bytes and KMS keyID

        """Client side encrypt and decrypt by AES"""
        # Encrypt User Private_Key using datakey from KMS, by client-side AES
        encrypted_privatekey = self.__encrypt(plain_datakey, private_key_hex) # Input datakey string and plaintext string

        content = {
            'encryptedPrivateKey': encrypted_privatekey,
            'publicKey': public_key_hex,
            'encryptedDatakey': encrypted_datakey
        }
        return content
    

    # return the private key's hash value, not an implement of crypty sign operation
    def sign(self,credential, encryptedPrivateKey, encryptedDatakey, message):
        nitro_kms = NitroKms()

        # Set environment variables
        nitro_kms.set_region('ap-southeast-1')
        nitro_kms.set_credentials(credential)
        
        # Decrypt encrypted data_key by KMS Decrypt API with attestation
        datakey = nitro_kms.kms_decrypt(encryptedDatakey)  # Key metadata included in Ciphertextblob, return bytes
        plain_datakey = base64.b64encode(datakey).decode('utf-8') # bytes to string

        # Decrypt private_key with plain data key
        private_key = self.__decrypt(plain_datakey, encryptedPrivateKey) # Input datakey string and encrypted string

        # Convert private key to ECDSA signing key
        sk = ecdsa.SigningKey.from_string(bytes.fromhex(private_key), curve=ecdsa.SECP256k1, hashfunc = hashlib.sha256)

        # Sign message using private_key
        bmessage = bytes(message, 'utf-8') 
        content = sk.sign(bmessage)  #Signature in bytes

        return content
    
    def __add_to_16(self, value):
        while len(value) % 16 != 0:
            value += '\0'
        return str.encode(value)  # return bytes

    def __encrypt(self, key, text):
        aes = AES.new(self.__add_to_16(key), AES.MODE_ECB)  # Initialize encryption method
        encrypt_aes = aes.encrypt(self.__add_to_16(text))  # Execute Encryption, return bytes
        encrypted_text = str(base64.encodebytes(encrypt_aes), encoding='utf-8')  # return base64 encoded string
        return encrypted_text

    def __decrypt(self, key, text):
        aes = AES.new(self.__add_to_16(key), AES.MODE_ECB)  # Initialize decryption method
        base64_decrypted = base64.decodebytes(text.encode(encoding='utf-8'))  # Execute Decryption, return bytes
        decrypted_text = str(aes.decrypt(base64_decrypted), encoding='utf-8').replace('\0', '')  # return base64 encoded string
        return decrypted_text

def main():
    print("Starting server...")

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

    # server client which call generateWallet or sign
    server = nitroServer()

    while True:
        c, addr = s.accept()
        
        # Get AWS credential sent from parent instance
        payload = c.recv(4096)
        payload_json = json.loads(payload.decode())
        print("payload json: {}".format(payload_json))

        apiCall = payload_json["apiCall"]

        if apiCall == "generateWallet":
            print("generateWallet request")
            credential = payload_json["credential"]
            keyId = payload_json["keyId"]
            result = server.generateWallet(credential, keyId)
            # send back to parent instance
            c.send(str.encode(json.dumps(result)))
            print("generateWallet finished")

        elif apiCall == "sign":
            print("sign request")
            credential = payload_json["credential"]
            message = payload_json["message"]
            encryptedPrivateKey = payload_json["encryptedPrivateKey"]
            encryptedDatakey = payload_json['encryptedDatakey']
        
            print('Message Received:  '+ message)

            signedStr = server.sign(credential, encryptedPrivateKey, encryptedDatakey, message)
            c.send(signedStr)

            print("sign fihished")
        else :
            print("nothing to do")

        c.close()


if __name__ == '__main__':
    main()