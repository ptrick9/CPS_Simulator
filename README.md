

# General Description
The open area CPS simulator is a program designed to simulate the detection of a bomb utilizing cheap chemical sensors carried by every member of the event. Specifically, the simulator attempts to emulate realistic detector equations for sensor drifting, energy usage, and communication overheads. The simulator utilizes a general flow which will be listed and described below. Detailed documentation can be found in the project wiki here: [Documentation](https://github.com/ptrick9/CPS_Simulator/wiki)


### [Quick Link To Tutorial](#tutorial)
### [Simulator Instructions](#simulator)
### [Data Processing Instructions](#data-processing)


## General Detection
The simulator runs with the given parameters until it detects a bomb. In order to detect a bomb, a few things must happen. A node (person) must walk near enough to a bomb for it to create a high reading on the node's sensor. The node's individual energy model must also decide to take a sample within the time period that the node is within the detection radius of the bomb. At some point in the future, again, as determined by the energy model, the node will send its data to the server.

Upon receiving the data from the node, the server will place the readings into a map of the arena, allowing for spatial correlation of multiple readings. In order to better correlate readings and decrease the number of false positives caused by sensor error, the server averages readings over areas. These areas are represented as large (user settable) squares in which the most recent readings in that area are averaged together. Once a square has a high enough average to trigger a detection, a supernode (a highly accurate node with greater sensor accuracy and detection distance) is routed to that area to ensure that there is no false alarm.

# Tutorial

The simulator project comes with a bash script that builds and runs a test project. In order for everything to compile correctly, some dependencies are required:

* [Golang](https://golang.org/dl/) must be installed
    * Choose most recent version for your system, install, and restart any terminal
    * For Mac you may need to add the go binary to your PATH. This is accomplished by adding `/usr/local/go` to the `/etc/paths` file
* [Java](http://www.oracle.com/technetwork/java/javase/downloads/jdk10-downloads-4416644.html) w/ JDK must be installed in order to run viewer. Accept agreement and install
    * If on LINUX, run ``` sudo apt-get install openjdk openjfx ``` which will install the required packages. If you are running Ubuntu 18.04, this does not seem to work.
* [Python3](https://www.python.org/downloads/) with numpy, pandas, matplotlib, and jupyter must be installed if you want to use the included statistics processing.
    * Ensure that python is on the path (checkbox on first page of installer)
    * Once python is installed, we need to install the libraries, this varies by system
        * Windows
            * Python defaults to C:\Users\\--username--\AppData\Local\Programs\Python\Python**__Number__**, so open a shell and get to that directory.
            * `cd Scripts`
            * `./pip3 install numpy pandas matplotlib jupyter`
        * Mac
            * install python3
            * `/usr/local/bin/pip3 install numpy pandas matplotlib jupyter`


Once Go and Java are installed, you can run one of the test cases.


### All Platforms
1. Please download the files from this [Google Drive](https://drive.google.com/drive/folders/1g77OqPrcu9dB5okUgC5AI4Z974BK7CmF?usp=sharing) and place them into the CPS_Simulator folder.

### Windows
 1. Navigate to CPS_Simulator folder in file viewer.
 2. `Shift+RightClick` and select `Open Powershell Window Here`.
 3. type `python test_run.py`
    * This command will run a script which will compile the simulator and run it with some basic command line arguments
4. Once completed, navigate to viewer folder and run Viewer program
5. In viewer, click file->open and select `tutorial-simulatorOutput.txt`

### Mac
1. Open terminal and navigate to CPS_Simulator folder
2. type `./test_run.sh`
    * This command will run a script which will compile the simulator and run it with some basic command line arguments
3. Once completed, navigate to viewer folder and run Viewer program
4. In viewer, click file->open and select `tutorial-simulatorOutput.txt`

# Simulator

In order to read a group of log files follow these steps:
1. Toolbar: File -> Menu -> Open
2. Navigate and open log file ending in -simulatorOutput.txt

Once the logs are loaded, you will see the simulation start to play.
You will now be able to use the following controls:
- Click the "Play" button to start/stop the simulation.
- Click the "< -" button to move back one instance in time
- Click the "- >" button to move forward one instance in time
- Click anywhere on the progress bar to skip to a certain time
- Click and drag the mouse around the main pane in order to pan
  around the room.
- Scroll up in order to zoom into the room.
- Scroll down in order to zoom out of the room.

Display Details: This is varied depending on the options selected
on Toolbar: View.
- Black lines represent walls or buildings in the scenario
- Blue squares represent node locations
- Nodes will have a yellow circle around them when taking a Sensor Reading
  if (View -> Nodes -> Sensor Coverage) is enabled.
- If (View -> Extras -> Sensor Reading) is selected the room will be overlayed
  with a subgrid of red squares. Low sensor readings will make a square closer
  to transparent, while higher sensor readings will make this square appear
  more solid.

|  Before Detection             |  Bomb Detected |
:-------------------------:|:-------------------------:
| ![Pre Detection](https://imgur.com/AomBuDe.png) |  ![Detected](https://imgur.com/61HAKHD.png) |

# Data Processing
Data processing is handled in a jupyter notebook.
Due to the large element of randomness in the simulator, meaningful data can only really be gathered by hundreds of runs. This can take a long time for big simulations, and to save you this time, an already processed dataset has been included in the tutorial_output folder.

1. To get started, open either Powershell or Terminal and navigate to the `CPS_Simulator/tutorial_output` directory. Then run: `jupyter notebook` on linux or mac and `jupyter-notebook.exe` on Windows. This will start a jupyter session and should open a web browser. If not, navigate your browser to the url given by the output of your `juptyer notebook` command.
2. In your web browser, click the `StatisticEngine.ipynb` file
3. This should open a new tab to the data processor
4. Here you will see python code to process the provided data files. In each cell on the page you can either press `ctrl+Enter` or click the run button at the top of the page.
    - Jupyter allows for us to run cells in order and pause the python interpreter, so if we ingest all of our data (sometimes a long process) in one cell and make a mistake trying to process it in another, we don't need to reingest our data, we can simply fix our processing step and rerun.
6. In the final cell, we filter the data, average similar runs, and then graph them. You can see, we end up with a graph showing that as we increase the number of validators required to make a bomb detection, we experience fewer false positives per second.

![False Positive Rejection](https://imgur.com/CP3nHqH.png)
