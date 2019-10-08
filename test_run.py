#! /bin/sh

import os

os.system("go build .\\simulator\\Simulation.go")
os.system('Simulation.exe -logNodes=true -logPosition=true -logGrid=true -logEnergy=true -regionRouting=true -noEnergy=true -csvMove=true -zipFiles=false -windRegionPath=hull_fine_bomb_fixWind_9x9.txt -inputFileName=Scenario_3.txt -imageFileName=marathon_street_map.png -iterations=1000 -superNodes=false -doOptimize=false -movementPath=marathon2_2000_1.scb -totalNodes=2000 -detectionThreshold=5 -detectionDistance=6 -GridCapacityPercentage=0.900000 -naturalLoss=0.005000 -sensorSamplingLoss=0.001000 -GPSSamplingLoss=0.005000 -serverSamplingLoss=0.010000 -SamplingLossBTCM=0.000100 -SamplingLossWifiCM=0.001000 -SamplingLoss4GCM=0.005000 -SamplingLossAccelCM=0.001000 -thresholdBatteryToHave=30 -thresholdBatteryToUse=10 -movementSamplingSpeed=20 -movementSamplingPeriod=20 -maxBufferCapacity=25 -sensorSamplingPeriod=1000 -GPSSamplingPeriod=1000 -serverSamplingPeriod=1000 -nodeStoredSamples=10 -GridStoredSamples=10 -errorMultiplier=0.60000 -numSuperNodes=4 -RecalibrationThreshold=3 -StandardDeviationThreshold=1.700000 -SuperNodeSpeed=3 -SquareRowCM=5 -SquareColCM=5 -validationThreshold=2 -serverRecal=true -driftExplorer=false -fineSensorPath=fine_bomb9x9.csv -csvSensor=false -OutputFileName=tutorial_output/tutorial -detectionWindow=59 -moveSize=4000 -commandBomb=true -bombX=577 -bombY=213')
