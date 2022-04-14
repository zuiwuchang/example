#!/bin/bash

if [[ "$@" == "command-default" ]];then
    exec "server"
else
    exec "$@"
fi