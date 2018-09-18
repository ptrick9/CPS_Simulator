#! /bin/bash

cd run_here
go build ../src/
./src -inputFileName=../testScenario/Scenario_1.txt -logPosition=true -outputFileName=testing
