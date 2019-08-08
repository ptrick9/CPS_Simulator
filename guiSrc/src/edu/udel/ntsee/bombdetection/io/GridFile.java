package edu.udel.ntsee.bombdetection.io;

import edu.udel.ntsee.bombdetection.Util;
import edu.udel.ntsee.bombdetection.data.Grid;

import java.io.Closeable;
import java.io.FileReader;
import java.io.IOException;
import java.io.LineNumberReader;

public class GridFile implements Closeable {

    private String path;
    private LineNumberReader lnr;

    private int width;
    private int height;

    public GridFile(String path)
        throws IOException {

        this.path = path;
        this.lnr = new LineNumberReader(new FileReader(path));
        this.width = Util.parseAmount(lnr.readLine());
        this.height = Util.parseAmount(lnr.readLine());
    }

    public Grid getGrid(int run) throws IOException {

        int targetLine = 2 + (run * height) + (2 * run);
        if(targetLine < lnr.getLineNumber() - 1) {
            this.lnr.close();
            this.lnr = new LineNumberReader(new FileReader(path));
            this.lnr.readLine();
            this.lnr.readLine();
        }

        while(lnr.getLineNumber() < targetLine) {
            this.lnr.readLine();
        }

        double[][] data = new double[height][width];
        for(int j=0; j<height; j++) {
            String[] vars = lnr.readLine().split("\t");
            for (int i=0; i<width; i++) {
                data[j][i] = Double.parseDouble(vars[i]);
            }
        }

        return new Grid(data);
    }



    @Override
    public void close() throws IOException {
        this.lnr.close();
    }
}
