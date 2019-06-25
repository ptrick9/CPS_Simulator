package edu.udel.ntsee.bombdetection.data;

import java.util.List;

public class SuperNode {

    private int x;
    private int y;
    private List<TimedNode> points;
    private List<TimedNode> path;
    private List<TimedNode> unvisitedPoints;

    public SuperNode(int x, int y, List<TimedNode> points, List<TimedNode> path, List<TimedNode> unvisitedPoints) {
        this.x = x;
        this.y = y;
        this.points = points;
        this.path = path;
        this.unvisitedPoints = unvisitedPoints;
    }

    public int getX() {

        return x;
    }

    public int getY() {

        return y;
    }

    public List<TimedNode> getPoints() {

        return points;
    }

    public List<TimedNode> getPath() {

        return path;
    }

    public List<TimedNode> getUnvisitedPoints() {

        return unvisitedPoints;
    }

}
