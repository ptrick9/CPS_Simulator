package edu.udel.ntsee.bombdetection.ui;

import javafx.geometry.Pos;
import javafx.scene.canvas.Canvas;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.control.Label;
import javafx.scene.layout.HBox;
import javafx.scene.paint.Color;

public class SolidLegendKey extends HBox {

    private Label label;
    private Canvas canvas;

    public SolidLegendKey(String text, Color color) {

        this.setSpacing(5);
        this.setAlignment(Pos.CENTER_LEFT);

        this.label = new Label(text);
        this.canvas = new Canvas(20, 20);
        this.initializeCanvas(color);

        this.getChildren().add(canvas);
        this.getChildren().add(label);

    }

    private void initializeCanvas(Color color) {

        GraphicsContext gc = canvas.getGraphicsContext2D();
        gc.setFill(color);
        gc.fillRect(0, 0, 20, 20);
        gc.setLineWidth(2);
        gc.strokeRect(0, 0, 20, 20);
    }

    public void setText(String text) {
        this.label.setText(text);
    }

    public void setColor(Color color) {
        this.initializeCanvas(color);
    }

}
