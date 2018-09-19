package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class TimedNode extends Node {

    private int time;

    public TimedNode(int x, int y, int time) {

        super(x, y);
        this.time = time;
    }

    public int getTime() {

        return time;
    }

    public void setTime(int time) {
        this.time = time;
    }

    private static final Pattern TIMEDNODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+), t:(\\d+)$");
    public static TimedNode fromString(String string) {

        Matcher m = TIMEDNODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse timed node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        int t = Integer.parseInt(m.group(3));
        return new TimedNode(x, y, t);

    }

    @Override
    public String toString() {
        return String.format("x:%d, y:%d, t:%d", getX(), getY(), time);
    }
}
