package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.data.Grid;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class FileManager {

    private static final Pattern TIME_REGEX = Pattern.compile("^t=(\\d+)$");

    private int[] properties;
    private PositionFile positions;
    private SamplesFile samples;
    private RoutesFile routes;
    private GridFile sensorReadings;
    private Map<Integer, Grid> nodePath;

    public FileManager(File file) {

        this.properties = new int[3];
        String base;
        try {
            base = file.getPath().substring(0, file.getPath().lastIndexOf("-"));
        } catch (StringIndexOutOfBoundsException e) {
            throw new IllegalArgumentException("Invalid File Name");
        }

        try {
            BufferedReader reader = new BufferedReader(new FileReader(base + "-simulatorOutput.txt"));
            properties[0] = Util.parseAmount(reader.readLine());
            properties[1] = Util.parseAmount(reader.readLine());
            properties[2] = Util.parseAmount(reader.readLine());
            this.positions = new PositionFile(base + "-simulatorOutput.txt");

        } catch (IOException e) {
            return;
        }

        try {
            this.samples = new SamplesFile(base + "-node.txt");
        } catch (IOException e) {
            this.samples = null;
        }

        try {
            this.routes = new RoutesFile(base + "-path.txt");
        } catch (IOException e) {
            this.routes = null;
        }

        try {
            this.sensorReadings = new GridFile(base + "-grid.txt");
        } catch (IOException e) {
            this.sensorReadings = null;
        }

        try {
            this.nodePath = new HashMap<>();
            this.loadNodePaths(base + "-pathgrid.txt");
        } catch (IOException e) {
            this.nodePath = null;
        }
    }

    public int[] getProperties() {
        return properties;
    }

    public PositionFile getPositions() {
        return positions;
    }

    public SamplesFile getSamples() {
        return samples;
    }

    public RoutesFile getRoutes() {
        return routes;
    }

    public GridFile getSensorReadings() {
        return sensorReadings;
    }

    public Map<Integer, Grid> getNodePath() {
        return nodePath;
    }

    private void loadNodePaths(String path) throws IOException {

        Matcher m = TIME_REGEX.matcher("");
        BufferedReader reader = new BufferedReader(new FileReader(path));

        String line;
        while ((line = reader.readLine()) != null) {
            m.reset(line);
            if (m.find()) {
                reader.readLine();
                int time = Integer.parseInt(m.group(1));
                double[][] data = new double[properties[1]][properties[0]];
                for(int i=0; i<properties[1]; i++) {
                    String[] rowVals = reader.readLine().split(" ");
                    double[] vals = new double[rowVals.length];
                    for (int j = 0; j < rowVals.length; j++) {
                        vals[j] = Double.parseDouble(rowVals[j]);
                    }
                    data[i] = vals;
                }
                nodePath.put(time, new Grid(data));
            }
        }
    }

}
