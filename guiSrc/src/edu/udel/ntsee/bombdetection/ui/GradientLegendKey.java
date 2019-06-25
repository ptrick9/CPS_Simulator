package edu.udel.ntsee.bombdetection.ui;

import edu.udel.ntsee.bombdetection.Util;
import javafx.geometry.Pos;
import javafx.scene.canvas.Canvas;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.control.Label;
import javafx.scene.layout.HBox;
import javafx.scene.layout.VBox;
import javafx.scene.paint.Color;

public class GradientLegendKey extends VBox {
    //--module-path=/usr/share/java/ --add-modules=javafx.controls,javafx.fxml
    private Label label;
    private HBox hbox;
    private Label minLabel;
    private Label maxLabel;
    private Canvas canvas;

    public GradientLegendKey(String text, Color start, Color end, double min, double max) {

        this.setAlignment(Pos.CENTER);
        this.label = new Label(text);
        this.canvas = new Canvas(100, 40);

        this.initializeCanvas(start, end);
        this.minLabel = new Label(String.valueOf(min));
        this.maxLabel = new Label(String.valueOf(max));
        this.hbox = new HBox(5);
        this.hbox.getChildren().addAll(minLabel, canvas, maxLabel);
        this.getChildren().addAll(label, hbox);
    }

    private void initializeCanvas(Color start, Color end) {

        GraphicsContext gc = canvas.getGraphicsContext2D();
        for(int i=0; i<=4; i++) {
            Color color = Util.gradient(start, end, (double)i/4);
            gc.setFill(color);
            gc.fillRect(i * 20, 0, 20, 20);
            gc.strokeRect(i * 20, 0, 20, 20);
        }
    }

    public void setMin(double min) {
        this.minLabel.setText(String.valueOf(min));
    }

    public void setMax(double max) {
        this.maxLabel.setText(String.valueOf(max));
    }

    protected Canvas getCanvas() {
        return canvas;
    }
}
