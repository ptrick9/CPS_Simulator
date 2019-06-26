package edu.udel.ntsee.bombdetection.data;

import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;
import edu.udel.ntsee.bombdetection.io.FileManager;
import javafx.beans.property.IntegerProperty;
import javafx.beans.property.SimpleIntegerProperty;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
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

    private Node bomb;
    private List<Wall> walls;
    private List<Node> positions;
    private List<Sample> samples;
    private List<SuperNode> superNodes;
    private List<AdHoc> adhocs;
    private Grid sensorReadings;
    private Road road;

    public static Room fromFile(File file) throws IOException {

        BufferedReader reader = new BufferedReader(new FileReader(file));
        Room room = new Room();
        room.index = new SimpleIntegerProperty(0);

        String imagePath = Util.parseString(reader.readLine());
        BufferedImage image = ImageIO.read(new File(file.getParent() + "/" + imagePath));
        room.walls = Util.createWallsFromImage(image);


        room.width = Util.parseAmount(reader.readLine());
        room.height = Util.parseAmount(reader.readLine());
        room.runs = Util.parseAmount(reader.readLine());

        int bx = Util.parseAmount(reader.readLine());
        int by = Util.parseAmount(reader.readLine());
        room.bomb = new Node(-1, bx, by);
        room.fileManager = new FileManager(room, file);

        File f = new File(file.getParent() + "\\roadLog.txt");
        if (f.exists()) room.road = Road.fromFile(f);
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

    public List<Wall> getWalls() {

        return walls;
    }
    public List<Node> getPositions() {

        return positions;
    }

    public Node getNodeByID(int id) {

        for (Node node : positions) {
            if (node.getID() == id) {
                return node;
            }
        }

        return null;
    }

    public List<Node> getNodesByIDs(List<Integer> ids) {

        List<Node> nodes = new ArrayList<>();
        for (int id : ids) {
            Node node = getNodeByID(id);
            if (node != null) {
                nodes.add(node);
            }
        }

        return nodes;
    }

    public List<Sample> getSamples() {

        return samples;
    }

    public List<SuperNode> getSuperNodes() {

        return superNodes;
    }

    public List<AdHoc> getAdHocs() {
        return adhocs;
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
        this.adhocs = fileManager.getAdHocs() != null ? fileManager.getAdHocs().getData(getIndex()) : null;
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

        if (fileManager.getAdHocs() != null) {
            fileManager.getAdHocs().close();
        }

        if (fileManager.getSensorReadings() != null) {
            fileManager.getSensorReadings().close();
        }

    }

    public Node getBomb() {
        return bomb;
    }

    public Road getRoad() {
        return road;
    }
}
