#! /bin/sh

cd run_here
go build ../src/
./src -inputFileName=../testScenario/Scenario_4.txt -logPosition=true -logGrid=true -logEnergy=true -logNodes=true -outputFileName=testing
