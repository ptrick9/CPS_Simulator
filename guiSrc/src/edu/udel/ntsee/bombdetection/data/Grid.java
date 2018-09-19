package edu.udel.ntsee.bombdetection.data;

import java.util.Arrays;

public class Grid {

    private double maxValue;
    private double[][] values;

    public Grid(double[][] values) {

        this.values = values;
        this.maxValue = calculateMaxValue();
    }

    private double calculateMaxValue() {

        double maxValue = values[0][0];
        for(int y=0; y<values.length; y++) {
            for(int x=1; x<values[y].length; x++) {
                maxValue = Math.max(values[y][x], maxValue);
            }
        }

        return maxValue;
    }

    public double getMaxValue() {
        return maxValue;
    }

    public double[][] getValues() {
        return values;
    }

    public Grid getAveragedValues(int squareSize) {

        if (values.length % squareSize != 0)
            throw new IllegalArgumentException("Invalid square size");

        int newHeight = values.length / squareSize;
        int newWidth = values[0].length / squareSize;
        double[][] averaged = new double[newHeight][newWidth];
        for(int y=0; y<newHeight; y++) {
            for(int x=0; x<newWidth; x++) {

                double average = 0;
                for(int j=0; j<squareSize; j++) {
                    for(int i=0; i<squareSize; i++) {
                        average += values[y * squareSize + j][x * squareSize + i];
                    }
                }

                average /= (squareSize * squareSize);
                averaged[y][x] = average;
            }
        }

        return new Grid(averaged);
    }

    @Override
    public String toString() {
        String ret = String.format("--- Grid ---\nMax Vale: %f\n", maxValue);
        for(double[] row : values) {
            ret += Arrays.toString(row) + "\n";
        }
        return ret;
    }
}
