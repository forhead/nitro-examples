## setup Nitro Enclave Environment on EC2

1) install the packages and dependencies
```
sudo amazon-linux-extras install aws-nitro-enclaves-cli -y
sudo yum install aws-nitro-enclaves-cli-devel -y
sudo yum install docker -y
sudo yum install jq -y
```

2) add ec2-user to the ne and docker user group 
```
sudo usermod -aG ne ec2-user
sudo usermod -aG docker ec2-user
```

3) verify nitro-cli version, check the installation status
```
nitro-cli --version
```

4) start services
```
# nitro-enclaves-allocator.service
sudo systemctl start nitro-enclaves-allocator.service
sudo systemctl enable nitro-enclaves-allocator.service
 
# nitro vsock proxy
sudo systemctl enable nitro-enclaves-vsock-proxy.service
sudo systemctl start nitro-enclaves-vsock-proxy.service

# Docker service 
sudo systemctl start docker && sudo systemctl enable docker
```

5) logout and re-login to activate the ec2-user in group docker and ne
