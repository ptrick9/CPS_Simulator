package edu.udel.ntsee.bombdetection.data;

import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;
import edu.udel.ntsee.bombdetection.io.FileManager;
import javafx.beans.property.IntegerProperty;
import javafx.beans.property.SimpleIntegerProperty;
import javafx.scene.layout.Pane;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class Room {

    private FileManager fileManager;
    private IntegerProperty index;
    private int width;
    private int height;
    private int runs;

    private List<Node> positions;
    private List<Sample> samples;
    private List<SuperNode> superNodes;
    private Grid sensorReadings;

    public static Room fromFile(File file) throws IOException {

        BufferedReader reader = new BufferedReader(new FileReader(file));
        Room room = new Room();
        room.index = new SimpleIntegerProperty(0);
        room.fileManager = new FileManager(file);
        room.width = room.fileManager.getProperties()[0];
        room.height = room.fileManager.getProperties()[1];
        room.runs = room.fileManager.getProperties()[2];
        reader.close();
        return room;
    }

    public int getIndex() {

        return index.get();
    }

    public void setIndex(int index) {

        this.index.set(index);
    }

    public IntegerProperty indexProperty() {

        return index;
    }

    public int getWidth() {

        return width;
    }

    public int getHeight() {

        return height;
    }

    public int getMaxRuns() {

        return runs;
    }

    public List<Node> getPositions() {

        return positions;
    }

    public List<Sample> getSamples() {

        return samples;
    }

    public List<SuperNode> getSuperNodes() {

        return superNodes;
    }

    public Grid getSensorReadings() {

        return sensorReadings;
    }

    public Grid getNodePath() {

        int i = getIndex() + 1;
        Map<Integer, Grid> paths = fileManager.getNodePath();
        Grid grid = null;
        while(grid == null) {
            i -= 1;
            grid = paths.get(i);
        }
        return grid;
    }

    public void updateData() throws IOException, LogFormatException {

        this.positions = fileManager.getPositions() != null ? fileManager.getPositions().getData(getIndex()) : null;
        this.samples = fileManager.getSamples() != null ? fileManager.getSamples().getData(getIndex()) : null;
        this.superNodes = fileManager.getRoutes() != null ? fileManager.getRoutes().getData(getIndex()) : null;
        this.sensorReadings = fileManager.getSensorReadings() != null ? fileManager.getSensorReadings().getGrid(getIndex()) : null;
    }

    public void close() throws IOException {


        if (fileManager.getPositions() != null) {
            fileManager.getPositions().close();
        }

        if (fileManager.getSamples() != null) {
            fileManager.getSamples().close();
        }

        if (fileManager.getRoutes() != null) {
            fileManager.getRoutes().close();
        }

        if (fileManager.getSensorReadings() != null) {
            fileManager.getSensorReadings().close();
        }
    }
}
