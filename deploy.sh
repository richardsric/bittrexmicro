#!/bin/bash
eval "$(aws ecr get-login --no-include-email --region us-east-2)"&&\
cd $GITPATH/bittrexmicro &&\
GOOS=linux GOARCH=386 go build &&\
docker build -t 375749533262.dkr.ecr.us-east-2.amazonaws.com/bittrexmicrox:latest .&&\
docker rmi -f $(docker images | grep 'bittrexmicro' | tr -s ' ' | cut -d ' ' -f 3)
GOOS=linux GOARCH=386 go build &&\
docker build -t 375749533262.dkr.ecr.us-east-2.amazonaws.com/bittrexmicrox:latest .&&\
docker push 375749533262.dkr.ecr.us-east-2.amazonaws.com/bittrexmicrox:latest &&\
cd ~ &&\
clear &&\
ssh -i "awsec2key.pem" ec2-user@ec2-18-221-72-6.us-east-2.compute.amazonaws.com /bin/bash /home/ec2-user/restart_itc.sh &&\
cd $GITPATH/bittrexmicro &&clear
