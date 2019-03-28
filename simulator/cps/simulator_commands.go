package cps

import "os"

type Params struct{
	NegativeSittingStopThresholdCM int     // This is a negative number for the sitting to be set to when map is reset
	SittingStopThresholdCM         int     // This is the threshold for the longest time a node can sit before no longer moving
	GridCapacityPercentageCM       float64 // This is the percent of a subgrid that can be filled with nodes, between 0.0 and 1.0
	ErrorModifierCM				   float64 // Multiplier for error model
	OutputFileNameCM               string  // This is the prefix of the output text file
	InputFileNameCM                string  // This must be the name of the input text file with ".txt"
	NaturalLossCM                  float64 // This can be any number n: 0 < n < .1
	SensorSamplingLossCM           float64 // This can be any number n: 0 < n < .1
	GPSSamplingLossCM              float64 // This can be any number n: 0 < n < GPSSamplingLossCM < .1
	ServerSamplingLossCM           float64 // This can be any number n: 0 < n < serverSamplingLossCM < .1
	ThresholdBatteryToHaveCM       int     // This can be any number n: 0 < n < 50
	ThresholdBatteryToUseCM        int     // This can be any number n: 0 < n < 20 < 100-thresholdBatteryToHaveCM
	MovementSamplingSpeedCM        int     // This can be any number n: 0 < n < 100
	MovementSamplingPeriodCM       int     // This can be any int number n: 1 <= n <= 100
	MaxBufferCapacityCM            int     // This can be aby int number n: 10 <= n <= 100
	EnergyModelCM                  string  // This can be "custom", "2StageServer", or other string will result in dynamic
	NoEnergyModelCM				   bool    // If set to true, all energy model values ignored
	SensorSamplingPeriodCM         int     // This can be any int n: 1 <= n <= 100
	GPSSamplingPeriodCM            int     // This can be any int n: 1 <= n < sensorSamplingPeriodCM <=  100
	ServerSamplingPeriodCM         int     // This can be any int n: 1 <= n < GPSSamplingPeriodCM <= 100
	NumStoredSamplesCM             int     // This can be any int n: 5 <= n <= 25
	GridStoredSamplesCM            int     // This can be any int n: 5 <= n <= 25
	DetectionThresholdCM           float64 //This is whatever value 1-1000 we determine should constitute a "detection" of a bomb
	PositionPrintCM                bool    //This is either true or false for whether to print positions to log file
	EnergyPrintCM                  bool    //This is either true or false for whether to print energy info to log file
	NodesPrintCM                   bool    //This is either true or false for whether to print node readings/averages to log file
	GridPrintCM                    bool    //This is either true or false for whether to print grid readings to log file
	SquareRowCM                    int     //This is an int 1 through maxX representing how many rows of squares there are
	SquareColCM                    int     //This is an int 1 through maxY representing how many columns of squares there are

	NumSuperNodes  int
	SuperNodeType  int
	SuperNodeSpeed int
	DoOptimize     bool
	//superNodeVariation int
	SuperNodeRadius int

	CenterCoord Coord

	Center Coord

	PositionPrint bool
	EnergyPrint   bool
	NodesPrint    bool
	GridPrint     bool

	DriftFile    *os.File
	NodeFile     *os.File
	PositionFile *os.File
	FoundBomb    bool
	Err 		 error




	BoolGrid       [][]bool
	Grid           [][]*Square
	SquareCapacity int


	XDiv			int
	YDiv			int

	MaxX             int
	MaxY             int
	BombX            int
	BombY            int

	ThreshHoldBatteryToHave  float32
	TotalPercentBatteryToUse float32
	Iterations_used          int
	Iterations_of_event      int
	EstimatedPingsNeeded     int

	b              *Bomb

	Tau1 			float64
	Tau2 		 	float64

	Recalibrate    bool

}