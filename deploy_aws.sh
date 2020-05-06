#!/usr/bin/env bash

env GOOS=linux GOARCH=amd64 go build -o sosmed_be

scp  -i /home/deddy/Documents/pos/aws/aws-lightsail.pem sosmed_be ubuntu@52.221.255.231:/home/ubuntu

