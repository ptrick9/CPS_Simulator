package edu.udel.ntsee.bombdetection.simulator;

import edu.udel.ntsee.bombdetection.Main;
import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.data.*;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;
import edu.udel.ntsee.bombdetection.ui.*;
import javafx.animation.Animation;
import javafx.animation.KeyFrame;
import javafx.animation.Timeline;
import javafx.fxml.FXML;
import javafx.scene.control.*;
import javafx.scene.input.MouseEvent;
import javafx.scene.layout.Pane;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;
import javafx.scene.text.Text;
import javafx.stage.FileChooser;
import javafx.stage.Stage;
import javafx.util.Duration;

import java.io.File;
import java.io.IOException;
import java.time.LocalDateTime;
import java.util.Date;
import java.util.List;

public class SimulatorController implements Drawable {

    private Room room;
    private Timeline timeline;
    @FXML private VBox root;

    // Menu
    @FXML private MenuItem menuItemClose;
    @FXML private CheckMenuItem checkMenuGridLines;
    @FXML private CheckMenuItem checkMenuQuadrants;
    @FXML private CheckMenuItem checkMenuWalls;
    @FXML private CheckMenuItem checkMenuSensorCoverage;
    @FXML private CheckMenuItem checkMenuAdHoc;

    @FXML private ToggleGroup toggleGroupNodeColor;
    @FXML private RadioMenuItem radioMenuGPSReading;
    @FXML private RadioMenuItem radioMenuBatteryLevel;

    @FXML private CheckMenuItem checkMenuLegendEnabled;

    @FXML private ToggleGroup toggleGroupExtras;
    @FXML private RadioMenuItem radioMenuNone;
    @FXML private RadioMenuItem radioMenuSensorReading;
    @FXML private RadioMenuItem radioMenuNodePathing;
    @FXML private RadioMenuItem radioMenuRoad;
    @FXML private CheckMenuItem checkMenuItemShowText;

    // Main
    private AdvancedCanvas canvas;
    @FXML private Pane containerCanvas;
    @FXML private Text textNotLoaded;

    // Control Bar
    @FXML private Text textProgress;
    @FXML private ProgressBar progressBarSimulation;
    @FXML private ToggleButton buttonPlay;

    // Mouse Positions
    @FXML private Label labelMousePosition;
    @FXML private Label labelMouseGridPosition;

    // Legend
    @FXML private VBox legendContainer;
    @FXML private Separator legendSeparator;
    @FXML private CheckMenuItem checkMenuLegendNode;
    @FXML private CheckMenuItem checkMenuLegendSuperNode;
    @FXML private CheckMenuItem checkMenuLegendAdHocLeader;
    @FXML private CheckMenuItem checkMenuLegendBattery;
    @FXML private CheckMenuItem checkMenuLegendSensorGrid;

    private SolidLegendKey nodeLegend;
    private SolidLegendKey superNodeLegend;
    private SolidLegendKey adhocLeaderLegend;
    private GradientLegendKey batteryLegend;
    private OpaqueLegendKey sensorGridLegend;

    // Extras
    private FileChooser fileChooser;

    public void initialize() {

        this.checkMenuGridLines.selectedProperty().addListener(event -> draw());
        this.checkMenuQuadrants.selectedProperty().addListener(event -> draw());
        this.checkMenuWalls.selectedProperty().addListener(event -> draw());
        this.checkMenuSensorCoverage.selectedProperty().addListener(event -> {
            room.getFileManager().loadSamplesFile();
            draw();
        });
        this.toggleGroupNodeColor.selectedToggleProperty().addListener(event -> {
            room.getFileManager().loadSamplesFile();
            draw();
        });
        this.toggleGroupExtras.selectedToggleProperty().addListener(event -> {
            if (toggleGroupExtras.getSelectedToggle() == radioMenuSensorReading) {
                room.getFileManager().loadGridFile();
            }
            draw();
        });
        this.checkMenuItemShowText.selectedProperty().addListener(event-> draw());
        this.checkMenuAdHoc.selectedProperty().addListener(event-> {
            room.getFileManager().loadAdHocFile();
            draw();
        });

        this.canvas = new AdvancedCanvas(this);
        this.canvas.widthProperty().bind(containerCanvas.widthProperty());
        this.canvas.heightProperty().bind(containerCanvas.heightProperty());
        this.canvas.requestFocus();
        this.containerCanvas.getChildren().add(canvas);

        this.timeline = new Timeline(new KeyFrame(Duration.millis(100), event -> {
            if (buttonPlay.isSelected() && room.getIndex() < room.getMaxRuns() - 1){
                room.setIndex(room.getIndex() + 1);
            } else { timeline.stop(); }
        }));
        this.timeline.setCycleCount(Animation.INDEFINITE);
        this.buttonPlay.selectedProperty().addListener((observable, oldValue, newValue) -> {
            if (newValue && room != null) { timeline.play(); }
            else { timeline.stop(); }
        });

        this.canvas.mouseProperty().addListener((observable, oldValue, newValue) -> {
            if (room != null) {

                AdvancedCanvas.Camera camera = canvas.getCamera();

                int x = (int)((newValue.getX() + camera.getOffsetX()) / camera.getBlockSize()) + canvas.getStartColumn();
                x = Math.max(0, x);
                x = Math.min(x, room.getWidth() - 1);

                int y = (int)((newValue.getY() + camera.getOffsetY()) / camera.getBlockSize()) + canvas.getStartRow();
                y = Math.max(0, this.room.getHeight() - y - 1);
                y = Math.min(y, room.getHeight() - 1);

                this.labelMousePosition.setText(String.format("Simulation Position - (%d, %d)", x, y));
                if (this.room.getSensorReadings() != null) {
                    int squares = room.getHeight() / room.getSensorReadings().getValues().length;
                    this.labelMouseGridPosition.setText(String.format("Sensor Reading Position - (%d, %d)", x/squares, y/squares));
                }
            }
        });
        this.initializeLegend();

        this.fileChooser = new FileChooser();
        this.fileChooser.setInitialDirectory(new File("../tutorial_output"));
        //Create filter to ease readability
        FileChooser.ExtensionFilter extFilter = new FileChooser.ExtensionFilter("PositionFiles (*-simulatorOutput.txt)", "*-simulatorOutput.txt");
        this.fileChooser.getExtensionFilters().add(extFilter);
    }

    private void initializeLegend() {

        this.legendSeparator.visibleProperty().bind(checkMenuLegendEnabled.selectedProperty());
        this.legendSeparator.managedProperty().bind(checkMenuLegendEnabled.selectedProperty());
        this.legendContainer.visibleProperty().bind(checkMenuLegendEnabled.selectedProperty());
        this.legendContainer.managedProperty().bind(checkMenuLegendEnabled.selectedProperty());

        this.nodeLegend = new SolidLegendKey("Node", Color.BLUE);
        this.nodeLegend.visibleProperty().bind(checkMenuLegendNode.selectedProperty());
        this.nodeLegend.managedProperty().bind(checkMenuLegendNode.selectedProperty());

        this.superNodeLegend = new SolidLegendKey("Super Node", Color.PLUM);
        this.superNodeLegend.visibleProperty().bind(checkMenuLegendSuperNode.selectedProperty());
        this.superNodeLegend.managedProperty().bind(checkMenuLegendSuperNode.selectedProperty());

        this.adhocLeaderLegend = new SolidLegendKey("Ad Hoc Leader Node", Color.DARKGREEN);
        this.adhocLeaderLegend.visibleProperty().bind(checkMenuLegendAdHocLeader.selectedProperty());
        this.adhocLeaderLegend.managedProperty().bindBidirectional(checkMenuLegendAdHocLeader.selectedProperty());

        this.batteryLegend = new GradientLegendKey("Battery", Color.RED, Color.GREEN, 0, 100);
        this.batteryLegend.visibleProperty().bind(checkMenuLegendBattery.selectedProperty());
        this.batteryLegend.managedProperty().bind(checkMenuLegendBattery.selectedProperty());

        this.sensorGridLegend = new OpaqueLegendKey("Sensor Grid", Color.RED, 0, 0);
        this.sensorGridLegend.visibleProperty().bind(checkMenuLegendSensorGrid.selectedProperty());
        this.sensorGridLegend.managedProperty().bind(checkMenuLegendSensorGrid.selectedProperty());

        this.legendContainer.getChildren().addAll(nodeLegend, superNodeLegend, adhocLeaderLegend, batteryLegend, sensorGridLegend);

    }
    @Override
    public void draw() {

        if (room == null) return;
        canvas.clear();


        if (radioMenuSensorReading.isSelected()) {
            drawSensorGrid(room.getSensorReadings());
        } else if (radioMenuNodePathing.isSelected()) {
            drawNodePathing(room.getNodePath());
        } else if (radioMenuRoad.isSelected()) {
            drawRoad(room.getRoad());
        }

        if (checkMenuAdHoc.isSelected()) {
            drawAdHocs(room.getAdHocs());
        }

        if (checkMenuWalls.isSelected()) {
            drawWalls(room.getWalls());
        }

        drawNodes(room.getPositions(), room.getSamples());
        drawSuperNodes(room.getSuperNodes());

        if (checkMenuGridLines.isSelected()) {
            canvas.drawGrid();
        }

        if (checkMenuQuadrants.isSelected()) {
            canvas.drawQuadrants();
        }

        canvas.outline();
    }

    public void drawWalls(List<Wall> walls) {

        if (walls == null) return;

        for (Wall wall : walls) {
            canvas.drawBlock(Color.BLACK, false, wall.getX(), wall.getY());
        }
    }

    public void drawAdHocs(List<AdHoc> adhocs) {

        if (adhocs == null) return;

        for(AdHoc adhoc : adhocs) {
            Node leader = room.getNodeByID(adhoc.getLeaderID());
            int lY = room.getHeight() - leader.getY() - 1;
            List<Node> children = room.getNodesByIDs(adhoc.getChildrenIDs());
            for(Node child : children) {
                int cY = room.getHeight() - child.getY() - 1;
                canvas.drawLine(leader.getX(), lY, child.getX(), cY);
            }
        }
    }

    public void drawNodes(List<Node> nodes, List<Sample> samples) {

        if (nodes == null) return;

        // Sensor Coverage
        if(checkMenuSensorCoverage.isSelected()) {
            if (samples == null) {
                checkMenuSensorCoverage.setSelected(false);
                Main.showErrorDialog(new LogFormatException("Samples log is unavailable."));
                return;
            }
            for (int i = 0; i < nodes.size(); i++) {
                Node node = nodes.get(i);
                Sample sample = samples.get(i);
                if (sample.isSensorChecked()) {
                    int y = room.getHeight() - node.getY() - 1;
                    canvas.drawCircle(Color.YELLOW,(2.34 / .5) * canvas.getCamera().getBlockSize(),
                            node.getX(), y);
                }
            }
        }

        // Bomb
        int by = room.getHeight() - (int)room.getBomb().getY() - 1;
        canvas.drawCircle(Color.RED, canvas.getCamera().getBlockSize(), (int)room.getBomb().getX(), by);

        // Nodes
        for (int i = 0; i < nodes.size(); i++) {

            Node node = nodes.get(i);
            int y = room.getHeight() - node.getY() - 1;
            Color color = Color.BLUE;
            if(samples != null && radioMenuBatteryLevel.isSelected()) {
                Sample sample = samples.get(i);
                color = Util.gradient(Color.RED, Color.GREEN, (double)sample.getBattery()/100);
            }

            canvas.drawBlock(color, true, node.getX(), y);
        }

        if (room.getAdHocs() != null) {
            for (AdHoc adHoc : room.getAdHocs()) {
                Node leader = room.getNodeByID(adHoc.getLeaderID());
                canvas.drawBlock(Color.DARKGREEN, true, leader.getX(), room.getHeight() - leader.getY() - 1);
            }
        }

    }

    public void drawSuperNodes(List<SuperNode> superNodes) {

        if (superNodes == null) { return; }
        for(SuperNode superNode : superNodes) {

            for(TimedNode node : superNode.getPath()) {
                int y = room.getHeight() - node.getY() - 1;
                canvas.drawBlock(Color.WHITE, true, node.getX(), y);
            }

            for(TimedNode node : superNode.getUnvisitedPoints()) {
                int y = room.getHeight() - node.getY() - 1;
                Color color = Util.gradient(Color.RED, Color.GREEN, (double)node.getTime()/120);
                canvas.drawBlock(color, true, node.getX(), y);
            }

            for(TimedNode node : superNode.getPoints()) {
                int y = room.getHeight() - node.getY() - 1;
                Color color = Util.gradient(Color.RED, Color.GREEN, (double)node.getTime()/120);
                canvas.drawBlock(color, true, node.getX(), y);
            }

            int y = room.getHeight() - superNode.getY();
            canvas.drawBlock(Color.PLUM, true, superNode.getX(), y);
        }
    }

    public void drawSensorGrid(Grid grid) {

        if (grid == null) {
            radioMenuSensorReading.setSelected(false);
            radioMenuNone.setSelected(true);
            Main.showErrorDialog(new LogFormatException("Sensor Reading log is unavailable."));
            return;
        }

        canvas.getGraphicsContext2D().save();
        int squares = room.getHeight() / grid.getValues().length;

        int yStart = canvas.getStartRow() / squares;
        int yEnd = Math.min(canvas.getEndRow() / squares, room.getHeight() / squares - 1);
        int xStart = canvas.getStartColumn() / squares;
        int xEnd = Math.min(canvas.getEndColumn() / squares, room.getWidth() / squares - 1);
        for (int y=yStart; y<=yEnd; y++) {

            int flipY = Math.max(0, grid.getValues().length - y - 1);
            for (int x=xStart; x<=xEnd; x++) {

                double reading = grid.getValues()[flipY][x];
                if (reading > 0) {
                    double amount = reading - grid.getMinValue();
                    double max = Math.max(grid.getMaxValue(), 20);
                    canvas.getGraphicsContext2D().setGlobalAlpha(amount / max);
                    canvas.drawBlock(Color.RED, true, x, y, squares);
                }
            }
        }

        canvas.getGraphicsContext2D().restore();
        sensorGridLegend.setMax(grid.getMaxValue());

    }

    private void drawNodePathing(Grid grid) {

        if (grid == null) {
            radioMenuNone.setSelected(true);
            Main.showErrorDialog(new LogFormatException("Path Grid log is unavailable."));
            return;
        }
        canvas.getGraphicsContext2D().save();
        int nodes = (canvas.getEndColumn() - canvas.getStartColumn())
                * (canvas.getEndRow() - canvas.getStartRow());
        if (nodes > 100000) grid = grid.getAveragedValues(2);
        int squares = room.getWidth() / grid.getValues().length;
        int yStart = canvas.getStartRow()/squares;
        int yEnd = (int)Math.ceil((double)canvas.getEndRow() / squares);
        yEnd = Math.min(yEnd, grid.getValues().length - 1);
        for(int y=yStart; y<=yEnd; y++) {
            int xStart = canvas.getStartColumn()/squares;
            int xEnd = (int)Math.ceil((double)canvas.getEndColumn() / squares);
            xEnd = Math.min(xEnd, grid.getValues()[y].length - 1);
            for(int x=xStart; x<=xEnd; x++) {
                Color color;
                if(grid.getValues()[y][x] == -1) {
                    color = Color.BLACK;
                } else if (grid.getValues()[y][x] == 1) {
                    color = Color.WHITE;
                } else {
                    double percentage = grid.getValues()[y][x] / grid.getMaxValue();
                    color = Util.gradient(Color.GREEN, Color.RED, percentage);
                }
                canvas.drawBlock(color, false, x, y, squares);
            }
        }

        canvas.getGraphicsContext2D().restore();
    }

    public void drawRoad(Road road) {

        if (road == null) return;
        for(TimedNode node : road.getNodes()) {
            if (node.getTime() == 0) {

            } else if (node.getTime() == -1) {
                Color color = Color.BLACK;//Util.gradient(Color.GREEN, Color.RED, (double)node.getTime() / road.getMax());
                canvas.drawBlock(color, true, node.getX(), node.getY());
            } else {
                Color color = Util.gradient(Color.GREEN, Color.RED, (double)node.getTime() / road.getMax());
                canvas.drawBlock(color, true, node.getX(), node.getY());
            }
        }
    }

    private void onRoomLoaded() {

        this.menuItemClose.setDisable(false);
        this.textNotLoaded.setVisible(false);
        this.room.indexProperty().addListener((observable, oldValue, newValue) -> {
            textProgress.setText(String.format("%d / %d", newValue.intValue() + 1, room.getMaxRuns()));
            progressBarSimulation.setProgress((newValue.doubleValue() + 1) / room.getMaxRuns());
            if (newValue.intValue() == room.getMaxRuns() - 1) {
                timeline.stop();
            }

            draw();
        });

        this.canvas.setRows(room.getHeight());
        this.canvas.setColumns(room.getWidth());
        this.canvas.center();
        this.timeline.play();
    }

    @FXML
    private void onMenuItemOpen() {

        File f = fileChooser.showOpenDialog(root.getScene().getWindow());
        if (f == null) return;
        try {
            this.fileChooser.setInitialDirectory(f.getParentFile());
            if (room != null) onMenuItemClose();
            try {
                this.room = Room.fromFile(f);
                this.room.indexProperty().addListener(((observable, oldValue, newValue) -> {
                    try {
                        this.room.updateData();
                    } catch (IOException | LogFormatException e) {
                        e.printStackTrace();
                    }
                }));
            } catch (IllegalArgumentException e) {
                Main.showErrorDialog(e);//new RoomLoadException("Invalid log file name"));
                return;
            }
            this.onRoomLoaded();
        }
        catch (IOException e) { Main.showErrorDialog(e); }
    }

    @FXML
    private void onMenuItemClose() {

        this.labelMousePosition.setText(" Simulation Position - N/A");
        this.labelMouseGridPosition.setText("Grid Position - N/A");
        this.sensorGridLegend.setMax(0);
        this.timeline.stop();
        try {
            this.room.close();
            this.room = null;
        }
        catch (IOException e) {Main.showErrorDialog(e); }
        catch (NullPointerException npe) {}
        finally {
            this.textNotLoaded.setVisible(true);
            this.textProgress.setText("0 / 0");
            this.progressBarSimulation.setProgress(0);
            this.menuItemClose.setDisable(true);
            this.canvas.clear();
        }
    }

    @FXML
    private void onMenuItemZoomFit() {

        canvas.center();
        draw();
    }

    @FXML
    private void onMenuItemFullscreen() {

        Stage stage = ((Stage)this.root.getScene().getWindow());
        stage.setFullScreen(!stage.isFullScreen());
    }

    @FXML
    private void onMenuItemToolBuilder() {

        Stage stage = (Stage)root.getScene().getWindow();
        Main.openBuilder(stage);
    }

    @FXML
    private void onMenuItemStatistics() {

        Main.openStatistics((Stage)root.getScene().getWindow());
    }

    @FXML
    private void onButtonBack() {

        if (room == null) return;
        if (room.getIndex() > 0) {
            timeline.stop();
            room.setIndex(room.getIndex() - 1);
        }
    }

    @FXML
    private void onButtonForward() {

        if (room == null) return;
        if (room.getIndex() < room.getMaxRuns() - 1) {
            timeline.stop();
            room.setIndex(room.getIndex() + 1);
        }
    }

    @FXML
    private void onClickProgressBar(MouseEvent event) {

        if (room == null) return;
        double percentage = Math.max(0, event.getX() / progressBarSimulation.getWidth());
        percentage = Math.min(percentage, 1);
        double index = (room.getMaxRuns() - 1) * percentage;
        room.setIndex((int)index);
    }

}
