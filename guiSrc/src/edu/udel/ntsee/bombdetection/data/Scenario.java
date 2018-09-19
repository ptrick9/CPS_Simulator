package edu.udel.ntsee.bombdetection.data;

import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.ui.AdvancedCanvas;
import javafx.collections.FXCollections;
import javafx.util.Pair;

import java.io.*;
import java.util.ArrayList;
import java.util.List;

public class Scenario {

    private int maxWidth;
    private int maxHeight;
    private int totalRuns;
    private int squareRow;
    private int squareCol;
    private int superNodeType;
    private int totalRandomNodes;

    private Node bomb;
    private List<Node> walls;
    private List<TimedNode> nodes;
    private List<TemporaryNode> attractions;

    public Scenario() {

        this.maxWidth = 0;
        this.maxHeight = 0;
        this.totalRuns = 0;
        this.squareRow = 0;
        this.squareCol = 0;
        this.superNodeType = 1;
        this.totalRandomNodes = 0;

        this.bomb = new Node(0, 0);
        this.walls = FXCollections.observableArrayList();
        this.nodes = FXCollections.observableArrayList();
        this.attractions = FXCollections.observableArrayList();

    }

    public static Scenario fromFile(File file)
            throws IOException {

        Scenario scenario = new Scenario();
        BufferedReader reader = new BufferedReader(new FileReader(file));

        scenario.squareRow = Util.parseVariable(reader.readLine());
        scenario.squareCol = Util.parseVariable(reader.readLine());
        scenario.totalRandomNodes = Util.parseVariable(reader.readLine());
        scenario.maxWidth = Util.parseVariable(reader.readLine());
        scenario.maxHeight = Util.parseVariable(reader.readLine());
        scenario.totalRuns = Util.parseVariable(reader.readLine());
        scenario.superNodeType = Util.parseVariable(reader.readLine());

        int bombX = Util.parseVariable(reader.readLine());
        int bombY = Util.parseVariable(reader.readLine());
        scenario.bomb = new Node(bombX, bombY);

        int totalNodes = Util.parseAmount(reader.readLine());
        for (int i = 0; i < totalNodes; i++) {
            scenario.nodes.add(TimedNode.fromString(reader.readLine()));
        }

        int totalWalls = Util.parseAmount(reader.readLine());
        for (int i = 0; i < totalWalls; i++) {
            scenario.walls.add(Node.fromString(reader.readLine()));
        }

        reader.readLine(); // s
        reader.readLine(); // p

        int totalAttractions = Util.parseAmount(reader.readLine());
        for (int i = 0; i < totalAttractions; i++) {
            scenario.attractions.add(TemporaryNode.fromString(reader.readLine()));
        }

        return scenario;
    }

    public void writeToFile(File file)
            throws IOException {

        BufferedWriter writer = new BufferedWriter(new FileWriter(file));
        writer.write("squareRow-" + squareRow); writer.newLine();
        writer.write("squareCol-" + squareCol); writer.newLine();
        writer.write("numNodes-" + totalRandomNodes); writer.newLine();
        writer.write("maxX-" + maxWidth); writer.newLine();
        writer.write("maxY-" + maxHeight); writer.newLine();
        writer.write("runs-" + totalRuns); writer.newLine();
        writer.write("superNodeType-" + superNodeType); writer.newLine();
        writer.write("bombX-" + bomb.getX()); writer.newLine();
        writer.write("bombY-" + bomb.getY()); writer.newLine();

        writer.write("N: " + nodes.size()); writer.newLine();
        for(TimedNode node : nodes) {
            writer.write(String.format("x:%d, y:%d, t:%d", node.getX(), node.getY(), node.getTime()));
            writer.newLine();
        }

        writer.write("W: " + walls.size()); writer.newLine();
        for(Node node : walls) {
            writer.write(String.format("x:%d, y:%d", node.getX(), node.getY()));
            writer.newLine();
        }

        writer.write("S: 0"); writer.newLine();
        writer.write("P: 0"); writer.newLine();

        writer.write("POIS: " + attractions.size()); writer.newLine();
        for(TemporaryNode node : attractions) {
            writer.write(String.format("x:%d, y:%d, ti:%d, to:%d", node.getX(), node.getY(), node.getStart(), node.getEnd()));
            writer.newLine();
        }

        writer.close();
    }

    public void setMaxWidth(int maxWidth) {
        this.maxWidth = maxWidth;
    }

    public void setMaxHeight(int maxHeight) {
        this.maxHeight = maxHeight;
    }

    public void setTotalRuns(int totalRuns) {
        this.totalRuns = totalRuns;
    }

    public void setSquareRow(int squareRow) {
        this.squareRow = squareRow;
    }

    public void setSquareCol(int squareCol) {
        this.squareCol = squareCol;
    }

    public void setSuperNodeType(int superNodeType) {
        this.superNodeType = superNodeType;
    }

    public void setTotalRandomNodes(int totalRandomNodes) {
        this.totalRandomNodes = totalRandomNodes;
    }

    public void setBomb(Node bomb) {
        this.bomb = bomb;
    }

    public int getMaxWidth() {
        return maxWidth;
    }

    public int getMaxHeight() {
        return maxHeight;
    }

    public int getTotalRuns() {
        return totalRuns;
    }

    public int getSquareRow() {
        return squareRow;
    }

    public int getSquareCol() {
        return squareCol;
    }

    public int getSuperNodeType() {
        return superNodeType;
    }

    public int getTotalRandomNodes() {
        return totalRandomNodes;
    }

    public Node getBomb() {
        return bomb;
    }

    public List<Node> getWalls() {
        return walls;
    }

    public List<TimedNode> getNodes() {
        return nodes;
    }

    public List<TemporaryNode> getAttractions() {
        return attractions;
    }

    public void remove(Node node) {
        walls.remove(node);
        nodes.remove(node);
        attractions.remove(node);
    }

    public boolean contains(Node node) {
        return bomb.equals(node) ||
                walls.contains(node) ||
                nodes.contains(node);
    }
}
