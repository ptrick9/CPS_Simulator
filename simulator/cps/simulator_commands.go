package cps

import (
	"os"
)

type Params struct {
	Events      PriorityQueue
	CurrentTime int

	SittingStopThresholdCM         int     // This is the threshold for the longest Time a node can sit before no longer moving
	ErrorModifierCM                float64 // Multiplier for error model
	OutputFileNameCM               string  // This is the prefix of the output text file
	InputFileNameCM                string  // This must be the name of the input text file with ".txt"

	MaxBufferCapacityCM            int     // This can be aby int number n: 10 <= n <= 100
	NumStoredSamplesCM             int     // This can be any int n: 5 <= n <= 25
	GridStoredSamplesCM            int     // This can be any int n: 5 <= n <= 25
	DetectionThresholdCM           float64 //This is whatever Value 1-1000 we determine should constitute a "detection" of a bomb
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

	OutputPrint			bool	//This is either true or false for whether to print OutputLog
	PositionPrint		bool	//This is either true or false for whether to print positions to log file
	EnergyPrint			bool	//This is either true or false for whether to print energy info to log file
	BatteryPrint		bool	//Similar to EnergyPrint. Log files have less words.
	NodesPrint			bool	//This is either true or false for whether to print node readings/averages to log file
	GridPrint			bool	//This is either true or false for whether to print grid readings to log file
	ClusterPrint		bool	//This is either true or false for whether to print cluster statistics to log file
	ClusterDebug		bool	//This is either true or false for whether to print cluster debug info to log file
	ServerClusterDebug	bool
	ReportBTAverages	bool	//No effect if ClusterPrint is false. Either true or false for whether to print average
								//number of nodes within bluetooth range, which slows the simulation

	GridWidth 		int
	GridHeight 		int

	MoveReadingsFile *os.File
	DriftFile      *os.File
	NodeFile       *os.File
	PositionFile   *os.File
	GridFile       *os.File
	EnergyFile     *os.File
	RoutingFile    *os.File
	ServerFile       *os.File
	DetectionFile    *os.File
	BatteryFile      *os.File
	RunParamFile     *os.File
	DriftExploreFile *os.File
	DistanceFile     *os.File
	OutputLog        *os.File
	ZipFiles         bool
	Files            []string
	NodeDataFile     *os.File
	ClusterStatsFile *os.File
	ClusterFile      *os.File
	ClusterDebugFile *os.File
	ServerClusterDebugFile	*os.File
	ClusterReadings  *os.File
	ClusterMessages  *os.File
	AliveValidNodes  *os.File
	SamplingData     *os.File
	SampleRates      *os.File

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

	IterationsCM		     int
	Iterations_used          int
	Iterations_of_event      int

	B *Bomb

	Tau1 float64
	Tau2 float64


	FileName string

	RegionRouting bool
	CSVMovement   bool
	CSVSensor     bool
	SuperNodes     bool
	IsSense		  bool

	CurrentNodes               int
	NumWallNodes               int

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
	AliveNodes   map[*NodeImpl]bool
	AliveValNodes   map[*NodeImpl]bool

	NumAtt      int
	Attractions []*Attraction
	BombSquare  *Square
	XLoc        int
	YLoc        int

	Width 			int
	Height 			int
	Server 			*FusionCenter //Server object

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
	ClusterMaxThreshold      int
	ClusterMinThreshold      int
	MaxClusterHeads          int
	ClusteringOn             bool
	RedundantClustering      bool
	DegreeWeight             float64
	BatteryWeight            float64
	Penalty                  float64
	GlobalRecluster          int
	AloneOrClusterRatio		 bool
	LocalRecluster           int
	ExpansiveRatio			 float64
	AloneThreshold           float64
	ReclusterPeriod          float64
	SmallImprovement         float64
	SmallImprovementRatio	 float64
	LargeImprovement         float64
	GlobalReclusterIncrement float64
	GlobalReclusterDecrement float64
	DisableGRThreshold       float64
	DisableCSThreshold       float64
	ServerReadyThreshold     float64
	InitClusterTime          int
	ClusterSearchThreshold   int
	ClusterHeadTimeThreshold        int
	ClusterHeadBatteryDropThreshold float64
	AdaptiveClusterSearch			bool
	ACSReset							bool
	AloneNodeClusterSearch bool


	DriftExplorer 	bool
	NumNodeMovements 	int
	MovementOffset 		int
	MovementSize 		int

	ReadingHistorySize	int
	ServerRecal 	bool
	MinDistance 	int
	MaxMoveMeters   float64
	CounterThreshold int

	ValidationType string
	RecalReject 	bool
	DensityThreshold int  // number of nodes that must be in a square for it to be considered dense and have the sampling rate decreased
	SamplingPeriodMS	 int

	BatteryCapacity				int
	AverageBatteryLevel			float64
	BatteryDeadThreshold		float64
	BatteryLowThreshold			float64
	BatteryMediumThreshold		float64
	BatteryHighThreshold		float64
	SampleLossPercentage		float64
	BluetoothLossPercentage		float64
	WifiLossPercentage			float64
	TotalAdaptations            int
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