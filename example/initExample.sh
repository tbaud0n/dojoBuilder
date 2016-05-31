#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/client


for process in {"dojo","dijit","dojox","util"}; do
    if [ ! -d $process ]; then
        git clone https://github.com/dojo/$process.git $process
    else
        echo $process already installed
    fi
done


# Local variables:
# coding: utf-8
# End:
