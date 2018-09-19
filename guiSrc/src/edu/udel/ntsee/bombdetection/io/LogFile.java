package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.exceptions.LogFormatException;

import java.io.*;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public abstract class LogFile<T> implements Closeable {

    protected String path;
    protected Pattern headerRegex;
    protected Pattern dataRegex;

    protected int lastRun;
    protected LineNumberReader lnr;
    protected Map<Integer, Integer> offsets;

    public LogFile(String path, Pattern headerRegex, Pattern dataRegex)
            throws IOException {

        this.path = path;
        this.headerRegex = headerRegex;
        this.dataRegex = dataRegex;

        this.lastRun = -1;
        this.lnr = new LineNumberReader(new FileReader(new File(path)));
        this.offsets = new HashMap<>();
        this.offsets.put(0, 0);
    }

    protected int parseHeader()
            throws IOException, LogFormatException {

        Matcher h = headerRegex.matcher(lnr.readLine());
        if (!h.matches()) throw new LogFormatException(getClass().getSimpleName() + " Invalid header on line " + lnr.getLineNumber());
        return Integer.parseInt(h.group(1));
    }

    protected abstract T parseData(Matcher m);

    public List<T> getData(int run)
            throws IOException, LogFormatException {

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
        List<T> data = new ArrayList<>();
        Matcher m = dataRegex.matcher("");
        for (int i=0; i<amount; i++) {
            m.reset(lnr.readLine());
            if (!m.matches()) throw new LogFormatException("Invalid data on line " + lnr.getLineNumber());
            data.add(parseData(m));
        }

        this.lastRun = run;
        return data;
    }

    protected void cacheNextIteration()
            throws IOException, LogFormatException {
;
        int lastRunLine = Collections.max(offsets.values());
        this.goToLine(lastRunLine);
        int amount = parseHeader();
        for (int i = 0; i < amount; i++) {
            lnr.readLine();
        }
        this.offsets.put(Collections.max(offsets.keySet()) + 1, lnr.getLineNumber());
    }

    protected void goToLine(int targetLine)
            throws IOException {

        if (lnr.getLineNumber() > targetLine) {
            this.close();
            this.lnr = new LineNumberReader(new FileReader(new File(path)));
        }

        for (int i = lnr.getLineNumber(); i < targetLine; i++) {
            this.lnr.readLine();
        }
    }

    @Override
    public void close() throws IOException  {

        this.lnr.close();
    }
}
