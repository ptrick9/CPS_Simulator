package edu.udel.ntsee.bombdetection;

import javafx.scene.paint.Color;

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

}
