package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Node {

    private int x;
    private int y;

    public Node(int x, int y) {

        this.x = x;
        this.y = y;
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

    private static final Pattern NODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+)$");
    public static Node fromString(String string) {

        Matcher m = NODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        return new Node(x, y);
    }

}
