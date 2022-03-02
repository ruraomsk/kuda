#!/bin/bash
echo 'Compiling'
GOOS=linux GOARCH=arm go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy kuda to device'
scp kuda admin@192.168.115.29:/home/admin
#scp test.bin admin@192.168.115.29:/home/admin