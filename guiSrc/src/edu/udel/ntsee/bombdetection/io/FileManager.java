package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.Grid;
import edu.udel.ntsee.bombdetection.data.Room;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;

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

    private String base;
    private Room room;
    private PositionFile positions;
    private SamplesFile samples;
    private RoutesFile routes;
    private AdHocFile adhocs;
    private GridFile sensorReadings;
    private Map<Integer, Grid> nodePath;

    public FileManager(Room room, File file) {

        this.room = room;

        try {
            base = file.getPath().substring(0, file.getPath().lastIndexOf("-"));
        } catch (StringIndexOutOfBoundsException e) {
            throw new IllegalArgumentException("Invalid File Name");
        }

        try {
            this.positions = new PositionFile(base + "-simulatorOutput.txt");

        } catch (IOException e) {
            return;
        }

        try {
            this.routes = new RoutesFile(base + "-path.txt");
        } catch (IOException e) {
            this.routes = null;
        }

        try {
            this.nodePath = new HashMap<>();
            this.loadNodePaths(base + "-pathgrid.txt");
        } catch (IOException e) {
            this.nodePath = null;
        }
    }

    public void loadSamplesFile() {

        if (samples != null) {
            return;
        }

        try {
            this.samples = new SamplesFile(base + "-node.txt");
            room.updateData();
        } catch (IOException | LogFormatException e) {
            this.samples = null;
        }


    }

    public void loadAdHocFile() {

        if (adhocs != null) {
            return;
        }

        try {
            this.adhocs = new AdHocFile(base + "-adhoc.txt");
            room.updateData();
        } catch (IOException | LogFormatException e) {
            this.adhocs = null;
        }
    }

    public void loadGridFile() {

        if (sensorReadings == null) {
            return;
        }

        try {
            this.sensorReadings = new GridFile(base + "-grid.txt");
            room.updateData();
        } catch (IOException | LogFormatException e) {
            this.sensorReadings = null;
        }
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

    public AdHocFile getAdHocs() {
        return adhocs;
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
                double[][] data = new double[room.getHeight()][room.getWidth()];
                for(int i=0; i<room.getHeight(); i++) {
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
