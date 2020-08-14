package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.Infection;

import java.io.IOException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class InfectionFile extends LogFile<Infection> {

    private static final Pattern HEADER = Pattern.compile("^Amount: (\\d+).*$");
    private static final Pattern DATA = Pattern.compile("^id: (\\d+) infection: (\\d+) mask: (true|false).*");

    public InfectionFile(String path) throws IOException {
        super(path, HEADER, DATA);
    }

    @Override
    protected Infection parseData(Matcher m) {
        int id = Integer.parseInt(m.group(1));
        Infection.Type type = Infection.Type.valueOf(Integer.parseInt(m.group(2)));
        boolean mask = Boolean.parseBoolean(m.group(3));
        return new Infection(id, type, mask);
    }
}
