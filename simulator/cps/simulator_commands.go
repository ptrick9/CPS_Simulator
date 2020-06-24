package cps

import (
	"os"
)

type Params struct {
	Events      PriorityQueue
	CurrentTime int

	NegativeSittingStopThresholdCM int     // This is a negative number for the sitting to be set to when map is reset
	SittingStopThresholdCM         int     // This is the threshold for the longest Time a node can sit before no longer moving
	GridCapacityPercentageCM       float64 // This is the percent of a subgrid that can be filled with nodes, between 0.0 and 1.0
	ErrorModifierCM                float64 // Multiplier for error model
	OutputFileNameCM               string  // This is the prefix of the output text file
	InputFileNameCM                string  // This must be the name of the input text file with ".txt"
	NaturalLossCM                  float64 // This can be any number n: 0 < n < .1

	WifiOr4G					   bool	   //True: nodes speak to server over wifi, False: nodes speak to server over 4G

	ThresholdBatteryToHaveCM       int     // This can be any number n: 0 < n < 50
	ThresholdBatteryToUseCM        int     // This can be any number n: 0 < n < 20 < 100-thresholdBatteryToHaveCM
	MovementSamplingSpeedCM        int     // This can be any number n: 0 < n < 100
	MovementSamplingPeriodCM       int     // This can be any int number n: 1 <= n <= 100
	MaxBufferCapacityCM            int     // This can be aby int number n: 10 <= n <= 100
	EnergyModelCM                  string  // This can be "custom", "2StageServer", or other string will result in dynamic
	NoEnergyModelCM                bool    // If set to true, all energy model values ignored
	SensorSamplingPeriodCM         int     // This can be any int n: 1 <= n <= 100
	GPSSamplingPeriodCM            int     // This can be any int n: 1 <= n < sensorSamplingPeriodCM <=  100
	ServerSamplingPeriodCM         int     // This can be any int n: 1 <= n < GPSSamplingPeriodCM <= 100
	NumStoredSamplesCM             int     // This can be any int n: 5 <= n <= 25
	GridStoredSamplesCM            int     // This can be any int n: 5 <= n <= 25
	DetectionThresholdCM           float64 //This is whatever Value 1-1000 we determine should constitute a "detection" of a bomb
	PositionPrintCM                bool    //This is either true or false for whether to print positions to log file
	EnergyPrintCM                  bool    //This is either true or false for whether to print energy info to log file
	NodesPrintCM                   bool    //This is either true or false for whether to print node readings/averages to log file
	GridPrintCM                    bool    //This is either true or false for whether to print grid readings to log file
	SquareRowCM                    int     //This is an int 1 through maxX representing how many rows of squares there are
	SquareColCM                    int     //This is an int 1 through maxY representing how many columns of squares there are
	StdDevThresholdCM			   float64 //Detection Threshold based on standard deviations from mean
	CalibrationThresholdCM		   float64
	DetectionDistance 			   float64
	RandomBomb					   bool


	StimFileNameCM        string
	ImageFileNameCM       string
	OutRoutingStatsNameCM string
	OutRoutingNameCM      string
	CPUProfile            string
	MemProfile            string

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

	GridWidth 		int
	GridHeight 		int

	MoveReadingsFile *os.File
	DriftFile      *os.File
	NodeFile       *os.File
	PositionFile   *os.File
	GridFile       *os.File
	EnergyFile     *os.File
	RoutingFile    *os.File
	ServerFile	   *os.File
	DetectionFile  *os.File
	BatteryFile    *os.File
	RunParamFile   *os.File
	DriftExploreFile *os.File
	DistanceFile 	*os.File
	OutputLog       *os.File
	ZipFiles 		bool
	Files 			[]string
	NodeDataFile   *os.File
	ClusterFile	   *os.File
	ClusterStatsFile *os.File
	ClusterDebug	*os.File
	ClusterReadings *os.File
	ClusterMessages *os.File
	AliveValidNodes	*os.File
	SamplingData    *os.File
	SampleRates     *os.File

	SensorPath  string
	FineSensorPath  string
	MovementPath  string
	WindRegionPath string

	SensorTimes []int
	TimeStep    int
	MaxTimeStep int

	FoundBomb bool
	Err       error

	BoardMap [][]int

	BoolGrid [][]bool
	Grid     [][]*Square

	SensorReadings [][][]float64
	FineSensorReadings [][][]float64
	NodeMovements  [][]Tuple

	SquareCapacity int

	XDiv int
	YDiv int

	MaxX  int
	MaxY  int
	BombX int
	BombY int
	BombXCM int
	BombYCM int
	CommBomb bool

	ThreshHoldBatteryToHave  float32
	TotalPercentBatteryToUse float32
	IterationsCM		     int
	Iterations_used          int
	Iterations_of_event      int
	EstimatedPingsNeeded     int

	B *Bomb

	Tau1 float64
	Tau2 float64


	FileName string

	RegionRouting bool
	AStarRouting  bool
	CSVMovement   bool
	CSVSensor     bool
	SuperNodes     bool
	IsSense		  bool

	CurrentNodes               int
	NumWallNodes               int
	NumPoints                  int
	NumPointsOfInterestKinetic int
	NumPointsOfInterestStatic  int

	NodeEntryTimes [][]int // node positions
	Wpos           [][]int // wall positions
	Spos           [][]int // super node positions
	Ppos           [][]int // super node points of interest positions
	Poikpos        [][]int // points of interest kinetic
	Poispos        [][]int // points of interest static

	DetectionThreshold float64

	//SquareRow        int
	//SquareCol        int
	TotalNodes       int
	NumStoredSamples int
	NumGridSamples   int

	WallNodeList []WallNodes
	NodeList     []*NodeImpl

	BatteryCharges []float32
	BatteryLosses  []float32

	BatteryLossesSensor				  []float32
	BatteryLossesGPS 			      []float32
	BatteryLossesServer 			  []float32
	BatteryLossesBT					  []float32
	BatteryLossesWiFi				  []float32
	BatteryLosses4G					  []float32
	BatteryLossesAccelerometer		  []float32

	NumAtt      int
	Attractions []*Attraction
	BombSquare  *Square
	XLoc        int
	YLoc        int

	Width 			int
	Height 			int
	Server 			FusionCenter //Server object

	MaxRaw 			float32
	EdgeRaw 		float32
	MaxADC 			float32
	EdgeADC 		float32
	ADCWidth 		float32
	ADCOffset		float32


	NodePositionMap			map[Tuple]*NodeImpl
	ValidationThreshold	int
	WindRegion [][]Coord

	FineWidth 		int
	FineHeight		int
	FineScale		int
	Scale 			int

	NodeTree            * Quadtree
	ClusterNetwork      * AdHocNetwork
	NodeBTRange         float64
	ClusterMaxThreshold int
	ClusterMinThreshold	int
	ClusteringOn        bool
	RedundantClustering bool
	DegreeWeight        float64
	BatteryWeight       float64
	Penalty             float64
	ReclusterThreshold  float64
	ReclusterPeriod		int

	DriftExplorer 	bool
	NumNodeMovements 	int
	MovementOffset 		int
	MovementSize 		int

	ReadingHistorySize	int
	ServerRecal 	bool
	MinDistance 	int

	ValidationType string
	RecalReject 	bool
	DensityThreshold int  // number of nodes that must be in a square for it to be considered dense and have the sampling rate decreased
	SamplingPeriodMS	 int
	BatteryCapacity				int
	AverageBatteryLevel			float64
	SampleLossPercentage		float64
	BluetoothLossPercentage		float64
	WifiLossPercentage			float64
}

// returns the amount of battery drained when a sampling event occurs
func (p *Params) SampleLossAmount() int {
	return int(float64(p.BatteryCapacity) * p.SampleLossPercentage)
}

// returns the amount of battery drained when a bluetooth event occurs
func (p *Params) BluetoothLossAmount() int {
	return int(float64(p.BatteryCapacity) * p.BluetoothLossPercentage)
}

func (p *Params) WifiLossAmount() int {
	return int(float64(p.BatteryCapacity) * p.WifiLossPercentage)
}