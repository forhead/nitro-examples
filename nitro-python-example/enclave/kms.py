import subprocess
import os

class nitroKms:
    def call_kms_generate_datakey(self,credential,keyId):
        aws_access_key_id = credential['aws_access_key_id']
        aws_secret_access_key = credential['aws_secret_access_key']
        aws_session_token = credential['aws_session_token']

        subprocess_args = [
            "/app/kmstool_enclave_cli",
            "genkey",
            "--region", os.getenv("REGION", "us-east-1"),
            "--proxy-port", "8000",
            "--aws-access-key-id", aws_access_key_id,
            "--aws-secret-access-key", aws_secret_access_key,
            "--aws-session-token", aws_session_token,
            "--key-id", keyId,
            "--key-spec","AES-256",
        ]

        print("subprocess args: {}".format(subprocess_args))

        proc = subprocess.Popen(
            subprocess_args,
            stdout=subprocess.PIPE
        )
        
        # base64-encoded datakey, 0 ciphertext,1 plaintext
        # CIPHERTEXT: ciphertext \n PLAINTEXT: plaintext
        # datakey_split = datakeyText.split("\n")
        # ciphertext = datakey_split[0]
        # plaintext = datakey_split[1]
        datakeyText = proc.communicate()[0].decode()
        return datakeyText

    def call_kms_decrypt(self,credential, ciphertext):
        aws_access_key_id = credential['aws_access_key_id']
        aws_secret_access_key = credential['aws_secret_access_key']
        aws_session_token = credential['aws_session_token']

        subprocess_args = [
            "/app/kmstool_enclave_cli",
            "decrypt",
            "--region", os.getenv("REGION", "us-east-1"),
            "--proxy-port", "8000",
            "--aws-access-key-id", aws_access_key_id,
            "--aws-secret-access-key", aws_secret_access_key,
            "--aws-session-token", aws_session_token,
            "--ciphertext", ciphertext,
        ]

        print("subprocess args: {}".format(subprocess_args))

        proc = subprocess.Popen(
            subprocess_args,
            stdout=subprocess.PIPE
        )
        # returns b64 encoded plaintext
        plaintext = proc.communicate()[0].decode()
        print('kms decrypted the datakey: ',plaintext)
        return plaintext
