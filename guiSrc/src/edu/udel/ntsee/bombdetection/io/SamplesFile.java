package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.Sample;

import java.io.IOException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class SamplesFile extends LogFile<Sample> {

    private static final Pattern HEADER = Pattern.compile("^Amount: (\\d+).*$");
    private static final Pattern DATA = Pattern.compile("^battery: (\\d+) sensor checked: (true|false) GPS checked: (true|false).*");

    public SamplesFile(String path) throws IOException {
        super(path, HEADER, DATA);
    }

    @Override
    protected Sample parseData(Matcher m) {

        int battery = Integer.parseInt(m.group(1));
        boolean sensorChecked = Boolean.parseBoolean(m.group(2));
        boolean gpsChecked = Boolean.parseBoolean(m.group(3));
        return new Sample(battery, sensorChecked, gpsChecked);
    }
}
