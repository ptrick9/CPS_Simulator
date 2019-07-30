package edu.udel.ntsee.bombdetection.data;

import java.util.ArrayList;
import java.util.List;

public class Node {

    private int id;
    private int x;
    private int y;
    private Sample sample;
    private List<Integer>  children;

    public Node(int id, int x, int y) {

        this.id = id;
        this.x = x;
        this.y = y;
        this.sample = null;
        this.children = null;
    }

    public int getID() {
        return id;
    }

    public int getX() {
        return x;
    }

    public void setX(int x) {
        this.x = x;
    }

    public int getY() {
        return y;
    }

    public void setY(int y) {
        this.y = y;
    }

    public Sample getSample() {
        return sample;
    }

    public void setSample(Sample sample) {
        this.sample = sample;
    }

    public boolean hasSample() {
        return sample != null;
    }

    public void setChildren(List<Integer> children) {
        this.children = children;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null) return false;

        Node node = (Node) o;

        if (x != node.x) return false;
        return y == node.y;
    }

    @Override
    public int hashCode() {
        int result = x;
        result = 31 * result + y;
        return result;
    }

    @Override
    public String toString() {
        return String.format("x:%d, y:%d", x, y);
    }



}
