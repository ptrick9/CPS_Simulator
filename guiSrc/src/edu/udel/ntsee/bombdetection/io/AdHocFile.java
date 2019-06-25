package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.data.AdHoc;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class AdHocFile extends LogFile<AdHoc>  {

    private static final Pattern HEADER = Pattern.compile("^Amount: (\\d+).*$");
    public static final Pattern DATA = Pattern.compile("^(\\d+): \\[((\\d+,? ?)*)\\]$");
    public AdHocFile(String path) throws IOException {
        super(path, HEADER, DATA);
    }

    @Override
    protected AdHoc parseData(Matcher m) {

        int leaderID = Integer.parseInt(m.group(1));
        List<Integer> childrenIDs = new ArrayList<>();
        String[] rawChildren = m.group(2).split(", ");
        for(String rawChild : rawChildren) {
            childrenIDs.add(Integer.parseInt(rawChild));
        }

        return new AdHoc(leaderID, childrenIDs);
    }
}
