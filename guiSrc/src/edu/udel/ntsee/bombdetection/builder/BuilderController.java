package edu.udel.ntsee.bombdetection.builder;

import edu.udel.ntsee.bombdetection.Main;
import edu.udel.ntsee.bombdetection.data.Node;
import edu.udel.ntsee.bombdetection.data.Scenario;
import edu.udel.ntsee.bombdetection.data.TemporaryNode;
import edu.udel.ntsee.bombdetection.data.TimedNode;
import edu.udel.ntsee.bombdetection.exceptions.ScenarioFormatException;
import edu.udel.ntsee.bombdetection.ui.AdvancedCanvas;
import edu.udel.ntsee.bombdetection.ui.Drawable;
import edu.udel.ntsee.bombdetection.ui.IntegerField;
import javafx.collections.ObservableList;
import javafx.fxml.FXML;
import javafx.scene.control.*;
import javafx.scene.control.cell.TextFieldTableCell;
import javafx.scene.layout.Pane;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;
import javafx.scene.text.Text;
import javafx.stage.FileChooser;
import javafx.stage.Stage;
import javafx.util.StringConverter;
import java.io.File;
import java.io.IOException;

public class BuilderController implements Drawable {

    private Scenario scenario;
    @FXML private VBox root;

    // Menu
    @FXML private CheckMenuItem checkMenuGridLines;
    @FXML private CheckMenuItem checkMenuQuadrants;
    @FXML private CheckMenuItem checkMenuBombs;
    @FXML private CheckMenuItem checkMenuWalls;
    @FXML private CheckMenuItem checkMenuNodes;
    @FXML private CheckMenuItem checkMenuAttractions;

    // Main
    private AdvancedCanvas canvas;
    @FXML private Pane canvasContainer;

    // Controls
    @FXML private ToggleGroup toggleGroupControls;
    @FXML private ToggleButton toggleButtonCamera;
    @FXML private ToggleButton toggleButtonBomb;
    @FXML private ToggleButton toggleButtonWalls;
    @FXML private ToggleButton toggleButtonNodes;
    @FXML private ToggleButton toggleButtonAttractions;
    @FXML private ToggleButton toggleButtonDelete;
    @FXML private Text textMousePosition;

    // Properties
    @FXML private IntegerField integerFieldWidth;
    @FXML private IntegerField integerFieldHeight;
    @FXML private IntegerField integerFieldTotalRuns;
    @FXML private IntegerField integerFieldSquareCols;
    @FXML private IntegerField integerFieldSquareRows;
    @FXML private IntegerField integerFieldSuperNodeType;

    // Bomb
    @FXML private IntegerField integerFieldBombX;
    @FXML private IntegerField integerFieldBombY;

    // Nodes
    @FXML private IntegerField integerFieldNodes;
    @FXML private TableView<TimedNode> tableViewNodes;
    @FXML private TableColumn<TimedNode, Integer> colNodeX;
    @FXML private TableColumn<TimedNode, Integer> colNodeY;
    @FXML private TableColumn<TimedNode, Integer> colNodeTime;

    // Walls
    @FXML private TableView<Node> tableViewWalls;
    @FXML private TableColumn<Node, Integer> colWallX;
    @FXML private TableColumn<Node, Integer> colWallY;

    // Attractions
    @FXML private TableView<TemporaryNode> tableViewAttractions;
    @FXML private TableColumn<TemporaryNode, Integer> colAttractionX;
    @FXML private TableColumn<TemporaryNode, Integer> colAttractionY;
    @FXML private TableColumn<TemporaryNode, Integer> colAttractionStart;
    @FXML private TableColumn<TemporaryNode, Integer> colAttractionEnd;

    // Extra
    private FileChooser fileChooser;

    @FXML
    private void initialize() {

        this.initializeMenu();
        this.initializeCanvas();
        this.initializeTables();
        this.initializeFields();
        this.initializeScenario(null);
        this.fileChooser = new FileChooser();
        FileChooser.ExtensionFilter ex = new FileChooser.ExtensionFilter("Text Files", "*.txt");
        this.fileChooser.getExtensionFilters().add(ex);
    }

    private void initializeMenu() {

        this.checkMenuGridLines.selectedProperty().addListener(observable -> draw());
        this.checkMenuQuadrants.selectedProperty().addListener(observable -> draw());
        this.checkMenuBombs.selectedProperty().addListener(observable -> draw());
        this.checkMenuWalls.selectedProperty().addListener(observable -> draw());
        this.checkMenuNodes.selectedProperty().addListener(observable -> draw());
        this.checkMenuAttractions.selectedProperty().addListener(observable -> draw());
    }

    private void initializeCanvas() {

        this.canvas = new AdvancedCanvas(this);
        this.canvas.setMouseEventHandler(event -> { onCanvasInteraction(); });
        this.canvas.widthProperty().bind(canvasContainer.widthProperty());
        this.canvas.heightProperty().bind(canvasContainer.heightProperty());
        this.canvas.allowPanningProperty().bind(toggleButtonCamera.selectedProperty());
        this.canvas.mouseProperty().addListener((observable, old, mouse) -> {
            textMousePosition.setText(String.format("(%d, %d)", mouse.getX(), mouse.getY()));
        });
        this.canvasContainer.getChildren().add(canvas);
    }

    private void initializeTables() {

        this.colNodeX.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colNodeX.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setX(event.getNewValue());
                draw();
            }
        });
        this.colNodeY.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colNodeY.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setY(event.getNewValue());
                draw();
            }
        });
        this.colNodeTime.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colNodeTime.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setTime(event.getNewValue());
            }
        });
        this.colWallX.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colWallX.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setX(event.getNewValue());
                draw();
            }
        });
        this.colWallY.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colWallY.setOnEditCommit(event -> {
            event.getRowValue().setY(event.getNewValue());
            draw();
        });
        this.colAttractionX.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colAttractionX.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setX(event.getNewValue());
                draw();
            }
        });
        this.colAttractionY.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colAttractionY.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setY(event.getNewValue());
                draw();
            }
        });
        this.colAttractionStart.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colAttractionStart.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setStart(event.getNewValue());
            }
        });
        this.colAttractionEnd.setCellFactory(TextFieldTableCell.forTableColumn(NUMBER_CONVERTER));
        this.colAttractionEnd.setOnEditCommit(event -> {
            if (event.getNewValue() != null) {
                event.getRowValue().setEnd(event.getNewValue());
            }
        });
    }

    private void initializeFields() {

        this.integerFieldWidth.valueProperty().addListener((observable, oldValue, width) -> {
            scenario.setMaxWidth(width.intValue());
            canvas.setColumns(width.intValue());
            draw();
        });
        this.integerFieldHeight.valueProperty().addListener((observable, oldValue, height) -> {
            scenario.setMaxHeight(height.intValue());
            canvas.setRows(height.intValue());
            draw();
        });
        this.integerFieldTotalRuns.valueProperty().addListener((observable, oldValue, totalRuns) -> {
            scenario.setTotalRuns(totalRuns.intValue());
        });
        this.integerFieldSquareCols.valueProperty().addListener((observable, oldValue, squareCols) -> {
            scenario.setSquareCol(squareCols.intValue());
        });
        this.integerFieldSquareRows.valueProperty().addListener((observable, oldValue, squareRows) -> {
            scenario.setSquareRow( squareRows.intValue());
        });
        this.integerFieldSuperNodeType.valueProperty().addListener((observable, oldValue, superNodeType) -> {
            scenario.setSuperNodeType(superNodeType.intValue());
        });
        this.integerFieldBombX.valueProperty().addListener((observable, oldValue, x) -> {
            scenario.getBomb().setX(x.intValue());
            draw();
        });
        this.integerFieldBombY.valueProperty().addListener((observable, oldValue, y) -> {
            scenario.getBomb().setY(y.intValue());
            draw();
        });
    }

    private void initializeScenario(Scenario scenario) {

        if (scenario == null) scenario = new Scenario();
        this.scenario = scenario;

        this.integerFieldWidth.setText(String.valueOf(scenario.getMaxWidth()));
        this.integerFieldHeight.setText(String.valueOf(scenario.getMaxHeight()));
        this.integerFieldTotalRuns.setText(String.valueOf(scenario.getTotalRuns()));
        this.integerFieldSquareCols.setText(String.valueOf(scenario.getSquareCol()));
        this.integerFieldSquareRows.setText(String.valueOf(scenario.getSquareRow()));
        this.integerFieldSuperNodeType.setText(String.valueOf(scenario.getSuperNodeType()));
        this.integerFieldBombX.setText(String.valueOf(scenario.getBomb().getX()));
        this.integerFieldBombY.setText(String.valueOf(scenario.getBomb().getY()));
        this.integerFieldNodes.setText(String.valueOf(scenario.getTotalRandomNodes()));
        this.tableViewWalls.setItems((ObservableList<Node>)scenario.getWalls());
        this.tableViewNodes.setItems((ObservableList<TimedNode>)scenario.getNodes());
        this.tableViewAttractions.setItems((ObservableList<TemporaryNode>)scenario.getAttractions());
        this.draw();
    }

    @Override
    public void draw() {

        canvas.clear();

        if (checkMenuQuadrants.isSelected()) {
            canvas.drawQuadrants();
        }

        if (checkMenuBombs.isSelected()) {
            Node bomb = scenario.getBomb();
            canvas.drawBlock(Color.RED, !checkMenuGridLines.isSelected(), bomb.getX(), bomb.getY());
        }

        if (checkMenuWalls.isSelected()) {
            for(Node node : scenario.getWalls()) {
                canvas.drawBlock(Color.BLACK, !checkMenuGridLines.isSelected(), node.getX(), node.getY());
            }
        }

        if (checkMenuNodes.isSelected()) {
            for(Node node : scenario.getNodes()) {
                canvas.drawBlock(Color.BLUE, !checkMenuGridLines.isSelected(), node.getX(), node.getY());
            }
        }

        if (checkMenuAttractions.isSelected()) {
            for(Node node: scenario.getAttractions()) {
                canvas.drawBlock(Color.PLUM, !checkMenuGridLines.isSelected(), node.getX(), node.getY());
            }
        }

        if (checkMenuGridLines.isSelected()) {
            canvas.drawGrid();
        }

        canvas.outline();
    }

    @FXML
    private void onMenuItemNew() {

        this.initializeScenario(null);
    }

    @FXML
    private void onMenuItemOpen() {

        File file = fileChooser.showOpenDialog(root.getScene().getWindow());
        try {
            Scenario scenario = Scenario.fromFile(file);
            this.initializeScenario(scenario);
        } catch (IOException ioe) {
            Main.showErrorDialog(ioe);
        } catch (NumberFormatException nfe) {
            Main.showErrorDialog(new ScenarioFormatException("Invalid scenario file"));
        } catch (NullPointerException npe) {}
    }

    @FXML
    private void onMenuItemSaveAs() {

        File file = fileChooser.showSaveDialog(root.getScene().getWindow());
        try {
            scenario.writeToFile(file);
        } catch (IOException ioe) {
            Main.showErrorDialog(ioe);
        } catch (NullPointerException npe) {}
    }

    @FXML
    private void onMenuItemZoomIn() {

        canvas.getCamera().zoomIn();
        draw();
    }

    @FXML
    private void onMenuItemZoomOut() {

        canvas.getCamera().zoomOut();
        draw();
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
    private void onMenuItemSimulator() {

        Stage stage = (Stage)root.getScene().getWindow();
        Main.openSimulator(stage);
    }

    @FXML
    private void onMenuItemStatistics() {

        Stage stage = (Stage)root.getScene().getWindow();
        Main.openStatistics(stage);
    }

    private void onCanvasInteraction() {

        int x = canvas.getMouse().getX();
        int y = canvas.getMouse().getY();
        if (x < 0 || x >= scenario.getMaxWidth() || y  < 0 || y  >= scenario.getMaxHeight()) return;

        Node node = new Node(x, y);
        if (toggleButtonDelete.isSelected() && scenario.contains(node)) {
            scenario.remove(node);
        }
        else if (toggleButtonBomb.isSelected() && !scenario.contains(node)) {
            scenario.getBomb().setX(x);
            scenario.getBomb().setY(y);
            integerFieldBombX.setText(String.valueOf(x));
            integerFieldBombY.setText(String.valueOf(y));
        }
        else if (toggleButtonWalls.isSelected() && !scenario.contains(node)) {
            scenario.getWalls().add(node);
        }
        else if (toggleButtonNodes.isSelected() && !scenario.contains(node)) {
            scenario.getNodes().add(new TimedNode(x, y, 0));
        }
        else if (toggleButtonAttractions.isSelected()) {
            scenario.getAttractions().add(new TemporaryNode(x, y, 0, 0));
        }

        draw();
    }

    private static final StringConverter<Integer> NUMBER_CONVERTER = new StringConverter<Integer>() {
        @Override
        public String toString(Integer object) {
            return String.valueOf(object);
        }

        @Override
        public Integer fromString(String string) {
            try { return Integer.parseInt(string); }
            catch (NumberFormatException e) { return null; }
        }
    };

}
