package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class TemporaryNode extends Node {

    private int start;
    private int end;

    public TemporaryNode(int x, int y, int start, int end) {
        super(x, y);
        this.start = start;
        this.end = end;
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

    private static final Pattern TEMPORARYNODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+), ti:(\\d+), to:(\\d+)$");
    public static TemporaryNode fromString(String string) {

        Matcher m = TEMPORARYNODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse temporary node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        int ti = Integer.parseInt(m.group(3));
        int to = Integer.parseInt(m.group(4));
        return new TemporaryNode(x, y, ti, to);
    }

    @Override
    public String toString() {
        return String.format("x:%d, y:%d, ti:%d, to:%d",
                getX(), getY(), getStart(), getEnd());
    }
}
