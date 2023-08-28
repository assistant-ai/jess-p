#!/bin/bash

go build .
chmod +x ./jess-p 
sudo cp ./jess-p /usr/local/bin/jess-p
rm ./jess-p