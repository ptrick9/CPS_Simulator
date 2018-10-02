package edu.udel.ntsee.bombdetection.ui;

import edu.udel.ntsee.bombdetection.Util;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.paint.Color;

public class OpaqueLegendKey extends GradientLegendKey{

    public OpaqueLegendKey(String text, Color color, double min, double max) {
        super(text, Color.WHITE, Color.WHITE, min, max);
        this.initializeCanvas(color);
    }

    private void initializeCanvas(Color color) {

        GraphicsContext gc = getCanvas().getGraphicsContext2D();
        gc.clearRect(0, 0, getCanvas().getWidth(), getCanvas().getHeight());
        gc.setFill(color);
        for(int i=0; i<=4; i++) {
            gc.setGlobalAlpha((double)i/4);
            gc.fillRect(i * 20, 0, 20, 20);
        }

        gc.setGlobalAlpha(1);
        for(int i=0; i<=4; i++) {
            gc.strokeRect(i * 20, 0, 20, 20);
        }
    }
}
