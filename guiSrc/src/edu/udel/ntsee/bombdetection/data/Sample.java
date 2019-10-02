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

    public void setBattery(int battery) {
        this.battery = battery;
    }

    public boolean isSensorChecked() {
        return sensorChecked;
    }

    public void setSensorChecked(boolean sensorChecked) {
        this.sensorChecked = sensorChecked;
    }

    public boolean isGpsChecked() {
        return gpsChecked;
    }

    public void setGpsChecked(boolean gpsChecked) {
        this.gpsChecked = gpsChecked;
    }

}
