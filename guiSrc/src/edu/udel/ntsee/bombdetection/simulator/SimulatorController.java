package edu.udel.ntsee.bombdetection.simulator;

import edu.udel.ntsee.bombdetection.Main;
import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.data.*;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;
import edu.udel.ntsee.bombdetection.exceptions.RoomLoadException;
import edu.udel.ntsee.bombdetection.ui.AdvancedCanvas;
import edu.udel.ntsee.bombdetection.ui.Drawable;
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
import java.util.List;

public class SimulatorController implements Drawable {

    private Room room;
    private Timeline timeline;
    @FXML private VBox root;

    // Menu
    @FXML private MenuItem menuItemClose;
    @FXML private CheckMenuItem checkMenuGridLines;
    @FXML private CheckMenuItem checkMenuQuadrants;
    @FXML private CheckMenuItem checkMenuSensorCoverage;

    @FXML private ToggleGroup toggleGroupNodeColor;
    @FXML private RadioMenuItem radioMenuGPSReading;
    @FXML private RadioMenuItem radioMenuBatteryLevel;

    @FXML private ToggleGroup toggleGroupExtras;
    @FXML private RadioMenuItem radioMenuNone;
    @FXML private RadioMenuItem radioMenuSensorReading;
    @FXML private RadioMenuItem radioMenuNodePathing;

    // Main
    private AdvancedCanvas canvas;
    @FXML private Pane containerCanvas;
    @FXML private Text textNotLoaded;

    // Control Bar
    @FXML private Text textProgress;
    @FXML private ProgressBar progressBarSimulation;
    @FXML private ToggleButton buttonPlay;

    // Extras
    private FileChooser fileChooser;

    public void initialize() {

        this.checkMenuGridLines.selectedProperty().addListener(event -> draw());
        this.checkMenuQuadrants.selectedProperty().addListener(event -> draw());
        this.checkMenuSensorCoverage.selectedProperty().addListener(event -> draw());
        this.toggleGroupNodeColor.selectedToggleProperty().addListener(event -> draw());
        this.toggleGroupExtras.selectedToggleProperty().addListener(event -> draw());

        this.canvas = new AdvancedCanvas(this);
        this.canvas.widthProperty().bind(containerCanvas.widthProperty());
        this.canvas.heightProperty().bind(containerCanvas.heightProperty());
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

        this.fileChooser = new FileChooser();
    }

    @Override
    public void draw() {

        if (room == null) return;
        canvas.clear();

        if (radioMenuSensorReading.isSelected()) {
            drawSensorGrid(room.getSensorReadings());
        }

        if (radioMenuNodePathing.isSelected()) {
            drawNodePathing(room.getNodePath());
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
                    canvas.drawCircle(Color.YELLOW, node.getX(), node.getY());
                }
            }
        }

        // Nodes
        for (int i = 0; i < nodes.size(); i++) {

            Node node = nodes.get(i);

            Color color = Color.BLUE;
            if(samples != null && radioMenuBatteryLevel.isSelected()) {
                Sample sample = samples.get(i);
                color = Util.gradient(Color.RED, Color.GREEN, (double)sample.getBattery()/100);
            }

            canvas.drawBlock(color, true, node.getX(), node.getY());
        }
    }

    public void drawSuperNodes(List<SuperNode> superNodes) {

        if (superNodes == null) { return; }
        for(SuperNode superNode : superNodes) {

            for(Node node : superNode.getPath()) {
                canvas.drawBlock(Color.WHITE, true, node.getX(), node.getY());
            }

            for(TimedNode node : superNode.getUnvisitedPoints()) {
                Color color = Util.gradient(Color.RED, Color.GREEN, (double)node.getTime()/120);
                canvas.drawBlock(color, true, node.getX(), node.getY());
            }

            for(TimedNode node : superNode.getPoints()) {
                Color color = Util.gradient(Color.RED, Color.GREEN, (double)node.getTime()/120);
                canvas.drawBlock(color, true, node.getX(), node.getY());
            }

            canvas.drawBlock(Color.PLUM, true, superNode.getX(), superNode.getY());
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
        int squares = room.getWidth() / grid.getValues().length;
        int yStart = canvas.getStartRow()/squares;
        int yEnd = (int)Math.ceil((double)canvas.getEndRow() / squares);
        yEnd = Math.min(yEnd, grid.getValues().length - 1);
        for(int y=yStart; y<=yEnd; y++) {
            int xStart = canvas.getStartColumn()/squares;
            int xEnd = (int)Math.ceil((double)canvas.getEndColumn() / squares);
            xEnd = Math.min(xEnd, grid.getValues()[y].length - 1);
            for(int x=xStart; x<=xEnd; x++) {
                double percentage = grid.getValues()[y][x] / grid.getMaxValue();
                canvas.getGraphicsContext2D().setGlobalAlpha(percentage);
                canvas.drawBlock(Color.RED, true, x, y, squares);
            }
        }

        canvas.getGraphicsContext2D().restore();
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
    private void onMenuItemZoomIn() {
        canvas.getCamera().zoomIn();
        draw();
        int nodes = (canvas.getEndColumn() - canvas.getStartColumn())
                * (canvas.getEndRow() - canvas.getStartRow());
        System.out.println("Draw: - " + nodes);
    }

    @FXML
    private void onMenuItemZoomOut() {
        canvas.getCamera().zoomOut();
        draw();
        int nodes = (canvas.getEndColumn() - canvas.getStartColumn())
                * (canvas.getEndRow() - canvas.getStartRow());
        System.out.println("Draw: - " + nodes);
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
