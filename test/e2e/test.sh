#!/bin/bash
currentdir=$(pwd)
cd ../../
workspaceBaseDir=$(pwd)
cd ${currentdir}


#set up environments
cd ./pretest
./start.sh
cd ..

#start test
cd ${currentdir}
python3 test.py
