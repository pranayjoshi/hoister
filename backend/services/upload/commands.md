* Docker run: `docker build -t hoister-upload-service .`
* docker tag: `docker tag hoister-upload-service:latest 891376924010.dkr.ecr.us-east-2.amazonaws.com/hoister-upload-service:latest`
* docker push: `docker push 891376924010.dkr.ecr.us-east-2.amazonaws.com/hoister-upload-service:latest`
*  aws login: `aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 891376924010.dkr.ecr.us-east-2.amazonaws.com`