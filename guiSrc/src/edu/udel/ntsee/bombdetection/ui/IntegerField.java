package edu.udel.ntsee.bombdetection.ui;

import javafx.beans.property.IntegerProperty;
import javafx.beans.property.SimpleIntegerProperty;
import javafx.scene.control.TextField;

public class IntegerField extends TextField {

    private IntegerProperty value;

    public IntegerField() {
        super();
        this.value = new SimpleIntegerProperty(0);
        textProperty().addListener(((observable, oldValue, newValue) -> {
            boolean isNumeric = newValue.matches("^(\\d)*$");
            if (!isNumeric) {
                this.setText(oldValue);
                return;
            }

            if (newValue.isEmpty()) return;
            this.value.set(Integer.parseInt(newValue));
        }));
    }

    public IntegerProperty valueProperty() {
        return value;
    }

    public int getValue() {
        return value.get();
    }
}
