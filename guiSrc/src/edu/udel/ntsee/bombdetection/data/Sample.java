package edu.udel.ntsee.bombdetection.data;

public class Sample {

    private int battery;
    private boolean sensorChecked;
    private boolean gpsChecked;


    public Sample(int battery, boolean sensorChecked, boolean gpsChecked) {
        this.battery = battery;
        this.sensorChecked = sensorChecked;
        this.gpsChecked = gpsChecked;
    }



    public int getBattery() {
        return battery;
    }

    public boolean isSensorChecked() {
        return sensorChecked;
    }

    public boolean isGpsChecked() {
        return gpsChecked;
    }

}
