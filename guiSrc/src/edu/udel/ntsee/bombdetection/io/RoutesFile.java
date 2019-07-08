package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.SuperNode;
import edu.udel.ntsee.bombdetection.data.TimedNode;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RoutesFile extends LogFile<SuperNode> {

    private static final Pattern HEADER = Pattern.compile("^Amount: (\\d+).*$");
    private static final Pattern DATA = Pattern.compile("^x: (\\d+) y: (\\d+) RoutePoints: \\[((\\{\\d+ \\d+ \\d+\\} ?)*+)\\] Path: \\[((\\{\\d+ \\d+ \\d+\\} ?)*+)\\] UnvisitedPoints: \\[((\\{\\d+ \\d+ \\d+\\} ?)*)\\].*$");

    public RoutesFile(String path) throws IOException {
        super(path, HEADER, DATA);
    }

    @Override
    protected SuperNode parseData(Matcher m) {

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        List<TimedNode> points = parseList(m.group(3));
        List<TimedNode> path = parseList(m.group(5));
        List<TimedNode> unvisited = parseList(m.group(7));
        return new SuperNode(x, y, points, path, unvisited);
    }

    private static List<TimedNode> parseList(String s) {

        List<TimedNode> list = new ArrayList<>();
        if (s.isEmpty()) return list;

        String[] groups = s.substring(1, s.length()-1).split("\\} \\{");
        for (String group : groups) {
            String[] vals = group.split(" ");
            int x = Integer.parseInt(vals[0]);
            int y = Integer.parseInt(vals[1]);
            int t = Integer.parseInt(vals[2]);
            list.add(new TimedNode(x, y, t));
        }

        return list;
    }
}
