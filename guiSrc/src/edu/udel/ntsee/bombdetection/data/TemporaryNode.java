package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class TemporaryNode {

    private int x;
    private int y;
    private int start;
    private int end;

    public TemporaryNode(int x, int y, int start, int end) {

        this.x = x;
        this.y = y;
        this.start = start;
        this.end = end;
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

    public int getStart() {
        return start;
    }

    public void setStart(int start) {
        this.start = start;
    }

    public int getEnd() {
        return end;
    }

    public void setEnd(int end) {
        this.end = end;
    }

    @Override
    public String toString() {
        return String.format("x:%d, y:%d, ti:%d, to:%d",
                getX(), getY(), getStart(), getEnd());
    }
}
