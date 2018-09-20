#! /bin/sh

cd tutorial_output
go build ../src/
./src -inputFileName=../testScenario/Scenario_4.txt -logPosition=true -logGrid=true -logEnergy=true -logNodes=true -outputFileName=tutorial -squareRow=20 -squareCol=20
