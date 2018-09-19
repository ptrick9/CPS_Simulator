package edu.udel.ntsee.bombdetection.data;

import edu.udel.ntsee.bombdetection.Util;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.Comparator;
import java.util.HashMap;
import java.util.Map;
import java.util.function.BiConsumer;

public class DriftLog {

    public static final String GRID = "Grid";
    public static final String NODE_ENERGY = "Node (Energy)";
    public static final String NODE_DRIFTING = "Node (Drifting)";
    public static final String NODE_BOTH = "Node (Both)";

    private String inputFileName;
    private String outputFileName;
    private String squareSize;
    private int storedNodeSamples;
    private int storedGridSamples;
    private double detectionThreshold;
    private Map<String, Statistics> statistics;

    private static final String INPUT_FILE = "Input File";
    private static final String DETECTION_THRESHOLD = "Detection Threshold";
    private static final String STORED_NODE_SAMPLES = "Samples Stored by Node";
    private static final String STORED_GRID_SAMPLES = "Samples Stored by Grid";
    private static final String SQUARE_DIMENSIONS = "Square Dimensions";

    private DriftLog() {

        this.statistics = new HashMap<>();
        this.statistics.put(GRID, new Statistics());
        this.statistics.put(NODE_ENERGY, new Statistics());
        this.statistics.put(NODE_DRIFTING, new Statistics());
        this.statistics.put(NODE_BOTH, new Statistics());
    }

    public static DriftLog fromFile(File file) throws IOException {

        DriftLog log = new DriftLog();
        BufferedReader reader = new BufferedReader(new FileReader(file));
        reader.readLine(); // Nodes
        reader.readLine(); // Rows
        reader.readLine(); // Columns
        log.storedNodeSamples = (int) Util.parseDouble(reader.readLine());
        log.storedGridSamples = (int) Util.parseDouble(reader.readLine());
        reader.readLine(); // Width
        reader.readLine(); // Height
        reader.readLine(); // Bomb X
        reader.readLine(); // Bomb Y
        reader.readLine(); // Iterations
        log.squareSize = Util.parseString(reader.readLine());
        log.detectionThreshold = Util.parseDouble(reader.readLine());
        log.inputFileName = Util.parseString(reader.readLine());
        log.outputFileName = Util.parseString(reader.readLine()); // output file
        reader.readLine(); // battery natural loss
        reader.readLine(); // sensor/gps/server loss
        reader.readLine(); // printing
        reader.readLine(); // super nodes

        String line;
        while ((line = reader.readLine()) != null) {

            if (line.isEmpty() || line.equals("----------------")) { continue;
            } else if (line.startsWith("Grid True Positive")) {
                log.statistics.get(GRID).truePositive++;
            } else if (line.startsWith("Grid False Positive")) {
                log.statistics.get(GRID).falsePositive++;
            } else if (line.startsWith("Grid True Negative")) {
                log.statistics.get(GRID).trueNegative++;
            } else if (line.startsWith("Grid False Negative")) {
                log.statistics.get(GRID).falseNegative++;
            } else if (line.startsWith("Node True Positive (energy)")) {
                log.statistics.get(NODE_ENERGY).truePositive++;
            } else if (line.startsWith("Node False Positive (energy)")) {
                log.statistics.get(NODE_ENERGY).falsePositive++;
            } else if (line.startsWith("Node True Negative")) {
                log.statistics.get(NODE_ENERGY).trueNegative++;
            } else if (line.startsWith("Node False Negative (energy)")) {
                log.statistics.get(NODE_ENERGY).falseNegative++;
            }
        }

        reader.close();
        return log;
    }

    public String getProperty(String name) {

        switch (name) {
            case DETECTION_THRESHOLD:
                return String.valueOf(detectionThreshold);
            case STORED_NODE_SAMPLES:
                return String.valueOf(storedNodeSamples);
            case STORED_GRID_SAMPLES:
                return String.valueOf(storedGridSamples);
            case SQUARE_DIMENSIONS:
                return squareSize;
            default:
                return null;
        }
    }

    public String getInputFileName() {
        return inputFileName;
    }

    public String getOutputFileName() {
        return outputFileName;
    }

    public String getSquareSize() {
        return squareSize;
    }

    public int getStoredNodeSamples() {
        return storedNodeSamples;
    }

    public int getStoredGridSamples() {
        return storedGridSamples;
    }

    public double getDetectionThreshold() {
        return detectionThreshold;
    }

    public int getTotal() {

        int total = 0;
        for (Statistics statistic : statistics.values()) {
            total += statistic.truePositive + statistic.trueNegative +
                    statistic.falsePositive + statistic.falseNegative;
        }

        return total;
    }


    public Map<String, Statistics> getStatistics() {
        return statistics;
    }

    @Override
    public String toString() {
        return "StatisticsLog{" +
                "inputFileName='" + inputFileName + '\'' +
                ", outputFileName='" + outputFileName + '\'' +
                ", storedNodeSamples=" + storedNodeSamples +
                ", storedGridSamples=" + storedGridSamples +
                ", detectionThreshold=" + detectionThreshold +
                ", squareSize='" + squareSize + '\'' +
                '}';
    }

    public static final Comparator<DriftLog> COMPARE_BY_INPUT_FILE = Comparator.comparing(o -> o.inputFileName);
    public static final Comparator<DriftLog> COMPARE_BY_DETECTION_THRESHOLD = Comparator.comparingDouble(o -> o.detectionThreshold);
    public static final Comparator<DriftLog> COMPARE_BY_STORED_NODE = Comparator.comparingInt(o -> o.storedNodeSamples);
    public static final Comparator<DriftLog> COMPARE_BY_STORED_GRID = Comparator.comparingInt(o -> o.storedGridSamples);
    public static final Comparator<DriftLog> COMPARE_BY_SQUARE_SIZE = Comparator.comparing(o -> o.squareSize);
    public static Comparator<DriftLog> getComparator(String name) {

        switch (name) {
            case INPUT_FILE:
                return COMPARE_BY_INPUT_FILE;
            case DETECTION_THRESHOLD:
                return COMPARE_BY_DETECTION_THRESHOLD;
            case STORED_NODE_SAMPLES:
                return COMPARE_BY_STORED_NODE;
            case STORED_GRID_SAMPLES:
                return COMPARE_BY_STORED_GRID;
            case SQUARE_DIMENSIONS:
                return COMPARE_BY_SQUARE_SIZE;
            default:
                return null;
        }
    }

    public class Statistics {

        private int truePositive;
        private int falsePositive;
        private int trueNegative;
        private int falseNegative;

        public Statistics() {

            this.truePositive = 0;
            this.falsePositive = 0;
            this.trueNegative = 0;
            this.falseNegative = 0;
        }

        public int getTruePositive() {
            return truePositive;
        }

        public int getFalsePositive() {
            return falsePositive;
        }

        public int getTrueNegative() {
            return trueNegative;
        }

        public int getFalseNegative() {
            return falseNegative;
        }

        public int getProperty(String name) {
            switch (name) {
                case "True Positive":
                    return truePositive;
                case "False Positive":
                    return falsePositive;
                case "True Negative":
                    return trueNegative;
                case "False Negative":
                    return trueNegative;
                default:
                    return 0;
            }
        }

        public int getTotal() {
            return truePositive + trueNegative + falsePositive + falseNegative;
        }

        @Override
        public String toString() {
            return "Statistics{" +
                    "truePositive=" + truePositive +
                    ", falsePositive=" + falsePositive +
                    ", trueNegative=" + trueNegative +
                    ", falseNegative=" + falseNegative +
                    '}';
        }
    }

}
