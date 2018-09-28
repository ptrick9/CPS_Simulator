package edu.udel.ntsee.bombdetection.data;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Road {

    private int max;
    private List<TimedNode> nodes;

    private Road(int max, List<TimedNode> nodes) {

        this.max = max;
        this.nodes = nodes;
    }

    public static Road fromFile(File file) throws IOException {

        Pattern maxRegex = Pattern.compile(".*max (\\d+).*");
        Pattern dataRegex = Pattern.compile(".*(\\d+) (\\d+) (\\d+).*");

        int max = 0;
        List<TimedNode> nodes = new ArrayList<>();

        BufferedReader reader = new BufferedReader(new FileReader(file));

        Matcher m = maxRegex.matcher(reader.readLine());
        if (!m.find())
            return null;

        max = Integer.parseInt(m.group(1));

        String line;
        m = dataRegex.matcher("");
        while ((line = reader.readLine()) != null) {
            m.reset(line);
            if (m.find()) {
                int x = Integer.parseInt(m.group(1));
                int y = Integer.parseInt(m.group(2));
                int t = Integer.parseInt(m.group(3));
                nodes.add(new TimedNode(x, y, t));
            }
        }

        reader.close();
        return new Road(max, nodes);
    }

    public int getMax() {
        return max;
    }

    public List<TimedNode> getNodes() {
        return nodes;
    }


}
