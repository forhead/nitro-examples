go build appClient.go

docker build -t enclaveclientk8s -f Dockerfile .

aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com

aws ecr create-repository --repository-name enclaveclientk8s 

docker tag enclaveclientk8s:latest 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com/enclaveclientk8s:latest

docker push 774209043150.dkr.ecr.ap-northeast-1.amazonaws.com/enclaveclientk8s:latest

