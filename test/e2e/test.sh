#!/bin/bash
cd ./pretest
./start.sh

cd ..
python3 test.py
