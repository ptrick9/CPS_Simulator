package edu.udel.ntsee.bombdetection.statistics;

import edu.udel.ntsee.bombdetection.Main;
import edu.udel.ntsee.bombdetection.data.DriftLog;
import javafx.embed.swing.SwingFXUtils;
import javafx.fxml.FXML;
import javafx.scene.SnapshotParameters;
import javafx.scene.chart.CategoryAxis;
import javafx.scene.chart.LineChart;
import javafx.scene.chart.NumberAxis;
import javafx.scene.chart.XYChart;
import javafx.scene.control.ComboBox;
import javafx.scene.image.WritableImage;
import javafx.scene.input.Clipboard;
import javafx.scene.input.ClipboardContent;
import javafx.scene.layout.VBox;
import javafx.stage.DirectoryChooser;
import javafx.stage.FileChooser;
import javafx.stage.Stage;

import javax.imageio.ImageIO;
import java.io.File;
import java.io.IOException;
import java.util.*;
import java.util.stream.Collectors;

public class StatisticsController {

    private List<DriftLog> driftLogs;

    @FXML private VBox root;
    @FXML private ComboBox<String> comboBoxInputFile;
    @FXML private ComboBox<String> comboBoxXAxis;
    @FXML private ComboBox<String> comboBoxReason;
    @FXML private ComboBox<String> comboBoxCategory;

    @FXML private LineChart<String, Double> chart;
    @FXML private CategoryAxis xAxis;
    @FXML private NumberAxis yAxis;

    private DirectoryChooser directoryChooser;
    private FileChooser txtFileChooser;
    private FileChooser pngFileChooser;

    public void initialize() {

        this.driftLogs = new ArrayList<>();
        this.directoryChooser = new DirectoryChooser();
        this.txtFileChooser = new FileChooser();
        FileChooser.ExtensionFilter ex = new FileChooser.ExtensionFilter("Text Files", "*.txt");
        this.txtFileChooser.getExtensionFilters().add(ex);
        this.pngFileChooser = new FileChooser();
        ex = new FileChooser.ExtensionFilter("PNG", "*.png");
        this.pngFileChooser.getExtensionFilters().add(ex);

        this.comboBoxXAxis.getSelectionModel().selectFirst();
        this.comboBoxReason.getSelectionModel().selectFirst();
        this.comboBoxCategory.getSelectionModel().selectFirst();

        this.comboBoxInputFile.getSelectionModel().selectedIndexProperty().addListener(observable -> {
            updateChart();
        });
        this.comboBoxXAxis.getSelectionModel().selectedIndexProperty().addListener(observable -> {
            updateChart();
        });
        this.comboBoxReason.getSelectionModel().selectedIndexProperty().addListener(observable -> {
            updateChart();
        });
        this.comboBoxCategory.getSelectionModel().selectedIndexProperty().addListener(observable -> {
            updateChart();
        });

        this.xAxis.labelProperty().bind(comboBoxXAxis.valueProperty());
        this.yAxis.labelProperty().bind(comboBoxCategory.valueProperty());
    }

    @FXML
    private void onMenuItemOpen() {

        File directory = directoryChooser.showDialog(root.getScene().getWindow());
        if (directory == null) return;

        directoryChooser.setInitialDirectory(directory.getParentFile());
        driftLogs.clear();
        File[] files = directory.listFiles((dir, name) -> name.endsWith("_drift.txt"));
        for(File file : files) {
            try { driftLogs.add(DriftLog.fromFile(file)); }
            catch (Exception e) { Main.showErrorDialog(e); }
        }

        comboBoxInputFile.getItems().clear();
        Set<DriftLog> inputFiles = new TreeSet<>(DriftLog.COMPARE_BY_INPUT_FILE);
        inputFiles.addAll(driftLogs);
        for(DriftLog log : inputFiles) {
            comboBoxInputFile.getItems().add(log.getInputFileName());
        }
        comboBoxInputFile.getSelectionModel().selectFirst();
    }

    @FXML
    private void onMenuItemSimulator() {

        Main.openSimulator((Stage)root.getScene().getWindow());
    }

    @FXML
    private void onMenuItemBuilder() {

        Main.openBuilder((Stage)root.getScene().getWindow());
    }

    @FXML
    private void onMenuItemSave() {

        File file = pngFileChooser.showSaveDialog(root.getScene().getWindow());
        WritableImage image = chart.snapshot(new SnapshotParameters(), null);
        try {
            ImageIO.write(SwingFXUtils.fromFXImage(image, null), "png", file);
        }
        catch (IOException e) { Main.showErrorDialog(e); }
        catch (NullPointerException npe) {} // No File
    }

    @FXML
    private void onMenuItemCopy() {

        Clipboard clipboard = Clipboard.getSystemClipboard();
        ClipboardContent content = new ClipboardContent();
        StringBuilder data = new StringBuilder();
        for (XYChart.Series<String, Double> series : chart.getData()) {
            data.append(series.getName()).append("\n");
            for (XYChart.Data<String, Double> point : series.getData()) {
                data.append(point.getXValue()).append("\t").append(point.getYValue()).append("\n");
            }
        }
        content.putString(data.toString());
        clipboard.setContent(content);
    }

    private void updateChart() {

        if (driftLogs.isEmpty()) return;

        List<DriftLog> logs = driftLogs.stream()
                .filter(driftLog -> driftLog.getInputFileName().equals(comboBoxInputFile.getValue()))
                .collect(Collectors.toList());

        Comparator<DriftLog> comparator = DriftLog.getComparator(comboBoxXAxis.getValue());
        Set<DriftLog> categories = new TreeSet<>(comparator);
        categories.addAll(logs);

        List<String> xPoints = categories.stream()
                .map(driftLog -> driftLog.getProperty(comboBoxXAxis.getValue()))
                .collect(Collectors.toList());

        List<Double> yPoints = new ArrayList<>();
        for(DriftLog category : categories) {
            double y = 0;
            double total = 0;
            for(DriftLog log : categories) {
                if (category.equals(log)) {
                    double logY = log.getStatistics().get(comboBoxReason.getValue()).getProperty(comboBoxCategory.getValue());
                    y +=  logY / log.getTotal();
                    total++;
                }
            }
            y /= total;
            yPoints.add(y);
        }

        XYChart.Series series = new XYChart.Series();
        for (int i = 0; i < xPoints.size(); i++) {
            series.getData().add(new XYChart.Data<>(xPoints.get(i), yPoints.get(i)));
        }

        chart.getData().clear();
        chart.setTitle(comboBoxInputFile.getValue() + " - " + comboBoxReason.getValue() + " " + comboBoxCategory.getValue());
        chart.getData().add(series);
    }

}
