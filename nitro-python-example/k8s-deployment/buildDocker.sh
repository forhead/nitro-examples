docker build -t enclaveserverk8s -f Dockerfile .

aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com

aws ecr create-repository --repository-name enclaveserverk8s 

docker tag enclaveserverk8s:latest 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com/enclaveserverk8s:latest

docker push 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com/enclaveserverk8s:latest

kubectl delete deploy enclaveserver-deployment
kubectl apply -f enclaveserver_deployment.yaml 
systemctl restart nitro-enclaves-allocator.service

# 查看hugepage设置
sudo cat /proc/sys/vm/nr_hugepages

# 开启HugePages
sudo sysctl -w vm.nr_hugepages=4096
sudo echo "vm.nr_hugepages=4096" >> /etc/sysctl.conf

# 重启节点
reboot

# 再次检查hugepage状态及大小
sudo cat /proc/sys/vm/nr_hugepages
sudo grep Huge /proc/meminfo