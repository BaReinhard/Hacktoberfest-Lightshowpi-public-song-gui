#!/bin/bash


RUNNING=`ps -ef | grep "lame" | grep -v "grep" | wc -l`

if [ "$RUNNING" = "0" ]
then
	echo false > /tmp/show-running
	sleep 5
	sudo systemctl stop rpi-server
else
	echo true > /tmp/show-running
fi