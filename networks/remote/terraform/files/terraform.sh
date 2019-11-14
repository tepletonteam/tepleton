#!/bin/bash
# Script to initialize a testnet settings on a server

#Usage: terraform.sh <testnet_name> <testnet_node_number>

#Add tond node number for remote identification
echo "$2" > /etc/tond-nodeid

#Create tond user
useradd -m -s /bin/bash tond
#cp -r /root/.ssh /home/tond/.ssh
#chown -R tond.tond /home/tond/.ssh
#chmod -R 700 /home/tond/.ssh

#Reload services to enable the tond service (note that the tond binary is not available yet)
systemctl daemon-reload
systemctl enable tond


