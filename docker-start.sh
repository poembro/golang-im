#!/bin/bash

nohup /webser/go_wepapp/golang-im/cmd/logic/logic > /webser/logs/logic.log 2>&1 &
 
echo "启动logic服务"

sleep 2
/webser/go_wepapp/golang-im/cmd/connect/connect > /webser/logs/connect.log 2>&1
echo "启动connect服务"
 