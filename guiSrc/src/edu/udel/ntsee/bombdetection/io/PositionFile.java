package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.Node;
import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;

import java.io.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class PositionFile extends LogFile<Node> {

    private static final Pattern HEADER = Pattern.compile("^t=  \\d+  amount=  (\\d+).*$");
    private static final Pattern DATA = Pattern.compile("^ID: \\d+ x: (\\d+) y: (\\d+).*$");

    public PositionFile(String path) throws IOException {
        super(path, HEADER, DATA);

        this.lastRun = -1;
        this.lnr = new LineNumberReader(new FileReader(new File(path)));
        this.lnr.readLine();
        this.lnr.readLine();
        this.lnr.readLine();

        this.offsets = new HashMap<>();
        this.offsets.put(0, 3);
    }


    @Override
    public List<Node> getData(int run) throws IOException, LogFormatException {
        int targetLine;
        if (lastRun + 1 == run) {
            targetLine = lnr.getLineNumber();
        } else if (offsets.containsKey(run)) {
            targetLine = offsets.get(run);
        } else {
            this.cacheNextIteration();
            return getData(run);
        }

        this.goToLine(targetLine);

        // parse
        int amount = parseHeader();
        List<Node> data = new ArrayList<>();
        Matcher m = dataRegex.matcher("");
        for (int i=0; i<amount; i++) {
            m.reset(lnr.readLine());
            if (!m.matches()) throw new LogFormatException("Invalid data on line " + lnr.getLineNumber());
            data.add(parseData(m));
        }

        this.lastRun = run;
        return data;
    }

    @Override
    protected Node parseData(Matcher m) {

        int x = Integer.parseInt(m.group(1));
        int y = Integer.parseInt(m.group(2));
        return new Node(x,  y);
    }
}