package edu.udel.ntsee.bombdetection;

import edu.udel.ntsee.bombdetection.data.Wall;
import javafx.scene.paint.Color;

import java.awt.image.BufferedImage;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Util {

    private static final Pattern VAr_REGEX = Pattern.compile(".*?:.*?(\\d+).*");

    private Util() {}

    public static Color gradient(Color from, Color to, double percentage) {

        double range = to.getHue() - from.getHue();
        percentage = Math.min(percentage, 1.0);
        return Color.hsb(from.getHue() + range * percentage, 1.0, 1.0);
    }

    // "name-(value)" -> value
    public static int parseVariable(String string) {

        return Integer.parseInt(string.substring(string.lastIndexOf("-") + 1));
    }
    // "name: (value)" -> value
    public static int parseAmount(String string) {

        Matcher m = VAr_REGEX.matcher(string);
        if (m.matches())
            return Integer.parseInt(m.group(1));

        throw new IllegalArgumentException(string);
    }

    public static String parseString(String string) {

        return string.substring(string.lastIndexOf(": ") + 2);
    }

    public static double parseDouble(String string) {
        return Double.parseDouble(parseString(string));
    }

    public static List<Wall> createWallsFromImage(BufferedImage image) {

        List<Wall> walls = new ArrayList<>();
        for (int y=0; y<image.getHeight(); y++) {
            for (int x=0; x<image.getWidth(); x++) {
                if (isBlackPixel(image.getRGB(x, y))) {

                    // Check if corner pixel of adjacent pixel is not black
                    if (x - 1 < 0
                        || x + 1 >= image.getWidth()
                        || y - 1 < 0
                        || y + 1 >= image.getHeight()
                        || !isBlackPixel(image.getRGB(x - 1, y))
                        || !isBlackPixel(image.getRGB(x + 1, y))
                        || !isBlackPixel(image.getRGB(x, y - 1))
                        || !isBlackPixel(image.getRGB(x, y + 1))) {
                        walls.add(new Wall(x, y));
                    }

                }

            }
        }

        return walls;
    }

    private static boolean isBlackPixel(int rgb) {

        java.awt.Color color = new java.awt.Color(rgb);
        return color.equals(java.awt.Color.BLACK);
    }


}
