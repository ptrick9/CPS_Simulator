package edu.udel.ntsee.bombdetection;

import edu.udel.ntsee.bombdetection.io.AdHocFile;
import javafx.application.Application;
import javafx.fxml.FXMLLoader;
import javafx.scene.Group;
import javafx.scene.Parent;
import javafx.scene.Scene;
import javafx.scene.control.Alert;
import javafx.stage.Stage;
import java.io.IOException;
import java.util.regex.Matcher;

public class Main extends Application {

    public static void main(String[] args) {

        launch(args);
    }

    @Override
    public void start(Stage stage)  {

        stage.setScene(new Scene(new Group()));
        openSimulator(stage);
        stage.show();
    }

    public static void openSimulator(Stage stage) {



        try {
            Parent root = FXMLLoader.load(Main.class.getResource("simulator/view.fxml"));
            stage.setTitle("Bomb Detection: Simulator");
            stage.getScene().setRoot(root);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void openBuilder(Stage stage) {

        try {
            Parent root = FXMLLoader.load(Main.class.getResource("builder/view.fxml"));
            stage.setTitle("Bomb Detection: Builder");
            stage.getScene().setRoot(root);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void openStatistics(Stage stage) {

        try {
            Parent root = FXMLLoader.load(Main.class.getResource("statistics/view.fxml"));
            stage.setTitle("Bomb Detection: Statistics");
            stage.getScene().setRoot(root);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static void showErrorDialog(Exception exception) {

        Alert alert = new Alert(Alert.AlertType.ERROR);
        alert.setTitle("Error");
        alert.setHeaderText(exception.getClass().getSimpleName());
        alert.setContentText(exception.getMessage());
        alert.show();
    }

}
