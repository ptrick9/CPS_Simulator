package edu.udel.ntsee.bombdetection.data;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Parser {

    /* Used to parse scenarios, see io package for parsing logs */

    private static final Pattern NODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+)$");
    public static Node nodeFromString(String string) {

        Matcher m = NODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        return new Node(-1, x, y);
    }

    private static final Pattern TEMPORARYNODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+), ti:(\\d+), to:(\\d+)$");
    public static TemporaryNode temporaryNodeFromString(String string) {

        Matcher m = TEMPORARYNODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse temporary node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        int ti = Integer.parseInt(m.group(3));
        int to = Integer.parseInt(m.group(4));
        return new TemporaryNode(x, y, ti, to);
    }

    private static final Pattern TIMEDNODE_PATTERN = Pattern.compile("^x:(\\d+), y:(\\d+), t:(\\d+)$");
    public static TimedNode timedNodeFromString(String string) {

        Matcher m = TIMEDNODE_PATTERN.matcher(string);
        if (!m.find()) throw new IllegalArgumentException("Can not parse timed node: Invalid format");

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        int t = Integer.parseInt(m.group(3));
        return new TimedNode(x, y, t);

    }

}
