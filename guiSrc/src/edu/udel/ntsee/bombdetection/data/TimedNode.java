package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class TimedNode {

    private int x;
    private int y;
    private int time;

    public TimedNode(int x, int y, int time) {

        this.x = x;
        this.y = y;
        this.time = time;
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
    public int getTime() {

        return time;
    }

    public void setTime(int time) {
        this.time = time;
    }



    @Override
    public String toString() {
        return String.format("x:%d, y:%d, t:%d", getX(), getY(), time);
    }
}
