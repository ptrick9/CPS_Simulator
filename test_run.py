#! /bin/sh

import os

#cd run_here
os.system("go build .\\src\\")
os.system('src.exe -inputFileName="testScenario\\Scenario_4.txt" -logPosition=true -logGrid=true -logEnergy=true -logNodes=true -outputFileName=run_here\\testing -squareRow=20 -squareCol=20')
