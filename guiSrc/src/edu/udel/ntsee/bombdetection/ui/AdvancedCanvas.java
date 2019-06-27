package edu.udel.ntsee.bombdetection.ui;

import edu.udel.ntsee.bombdetection.data.Node;
import javafx.beans.property.*;
import javafx.event.Event;
import javafx.event.EventHandler;
import javafx.geometry.VPos;
import javafx.scene.canvas.Canvas;
import javafx.scene.canvas.GraphicsContext;
import javafx.scene.input.MouseEvent;
import javafx.scene.paint.Color;
import javafx.scene.text.Font;
import javafx.scene.text.TextAlignment;

import java.awt.*;

/*
 * Note: This class was designed so that the (0, 0) is top left.
 * However, it transforms x,y so that (0, 0) is bottom left in render methods.
 */
public class AdvancedCanvas extends Canvas {

    private EventHandler<MouseEvent> mouseEventHandler;
    private Drawable drawable;
    private Camera camera;

    private ObjectProperty<Node> mouseProperty;
    private BooleanProperty allowPanning;
    private int columns, rows;
    private double clickX, clickY;

    // to do:
    // copied positioning code from draw block
    // should extract into a function
    public void drawNumber(double number, int x, int y, int squares) {

        GraphicsContext gc = getGraphicsContext2D();
        gc.setTextAlign(TextAlignment.CENTER);
        gc.setTextBaseline(VPos.CENTER);
        gc.setFont(new Font(gc.getFont().getName(), camera.getBlockSize()));

        double worldX = x * squares - getStartColumn();
        double width = squares;
        if (worldX < 0) {
            worldX = 0;
            width = squares + (x * squares - getStartColumn());
        }

        double worldY = y * squares - getStartRow();
        double height = squares;
        if (worldY < 0) {
            worldY = 0;
            height = squares + (y * squares - getStartRow());
        }

        double pixelX = -camera.getOffsetX() + worldX * camera.getBlockSize() + (camera.getBlockSize() * width / 2);
        double pixelY = -camera.getOffsetY() + worldY * camera.getBlockSize() + (camera.getBlockSize() * height / 2);
        gc.strokeText(String.valueOf(number), pixelX, pixelY, width * camera.getBlockSize());
    }

    public Node getMouse() {
        return mouseProperty.get();
    }

    public ObjectProperty<Node> mouseProperty() {
        return mouseProperty;
    }

    public AdvancedCanvas(Drawable drawable) {

        this.drawable = drawable;
        this.camera = new Camera(getWidth(), getHeight());
        this.allowPanning = new SimpleBooleanProperty(true);
        this.columns = 0;
        this.rows = 0;
        this.clickX = 0;
        this.clickY = 0;
        this.mouseProperty = new SimpleObjectProperty<>(new Node(-1, 0, 0));

        widthProperty().addListener((observable, oldValue, newValue) -> {
            this.camera.setWidth(newValue.doubleValue());
            this.drawable.draw();
        });
        heightProperty().addListener(((observable, oldValue, newValue) -> {
            this.camera.setHeight(newValue.doubleValue());
            this.drawable.draw();
        }));
        setOnMousePressed(event -> {
            this.clickX = event.getX();
            this.clickY = event.getY();
            if (mouseEventHandler != null) this.mouseEventHandler.handle(event);
        });
        setOnMouseDragged(event -> {
            double dx = event.getX() - clickX;
            double dy = event.getY() - clickY;
            if (allowPanning.get()) this.camera.move(dx, dy);
            this.clickX = event.getX();
            this.clickY = event.getY();
            if (mouseEventHandler != null) this.mouseEventHandler.handle(event);
            this.drawable.draw();
        });
        setOnMouseMoved(event -> {
            int x = (int)((event.getX() + camera.getOffsetX()) / camera.getBlockSize()) + getStartColumn();
            int y = (int)((event.getY() + camera.getOffsetY()) / camera.getBlockSize()) + getStartRow();
            this.mouseProperty.set(new Node(-1, x, y));
        });
        setOnScroll(event -> {

            if (event.getDeltaY() > 0)
                camera.zoomIn();
            else if (event.getDeltaY() < 0)
                camera.zoomOut();

            GraphicsContext gc = getGraphicsContext2D();
            gc.setLineWidth(camera.getScale());
            drawable.draw();
        });

    }

    public void setMouseEventHandler(EventHandler<MouseEvent> mouseEventHandler) {
        this.mouseEventHandler = mouseEventHandler;
    }

    public Camera getCamera() {
        return camera;
    }

    public BooleanProperty allowPanningProperty() {
        return allowPanning;
    }

    public void setColumns(int columns) {

        this.columns = columns;
    }

    public void setRows(int rows) {

        this.rows = rows;
    }

    public void clear() {

        GraphicsContext gc = getGraphicsContext2D();
        gc.clearRect(0, 0, getWidth(), getHeight());
    }

    public void outline() {

        GraphicsContext gc = getGraphicsContext2D();
        gc.save();

        gc.setLineWidth(camera.getScale() * 3);

        if(getStartRow() <= 0) {
            double startX = -camera.getOffsetX();
            double startY = -camera.getOffsetY();
            double endX = -camera.getOffsetX() + getTotalVisibleColumns() * camera.getBlockSize();
            double endY = startY;
            gc.strokeLine(startX, startY, endX, endY);
        }

        if (getEndRow() >= rows) {
            double startX = -camera.getOffsetX();
            double startY = -camera.getOffsetY() + getTotalVisibleRows() * camera.getBlockSize();
            double endX = -camera.getOffsetX() + getTotalVisibleColumns() * camera.getBlockSize();
            double endY = startY;
            gc.strokeLine(startX, startY, endX, endY);
        }

        if (getStartColumn() <= 0) {
            double startX = -camera.getOffsetX();
            double startY = -camera.getOffsetY();
            double endX = startX;
            double endY = -camera.getOffsetY() + getTotalVisibleRows() * camera.getBlockSize();
            gc.strokeLine(startX, startY, endX, endY);
        }

        if (getEndColumn() >= columns) {
            double startX = -camera.getOffsetX() + getTotalVisibleColumns() * camera.getBlockSize();
            double startY = -camera.getOffsetY();
            double endX = startX;
            double endY = -camera.getOffsetY() + getTotalVisibleRows() * camera.getBlockSize();
            gc.strokeLine(startX, startY, endX, endY);
        }

        gc.restore();
    }

    public void drawGrid() {

        GraphicsContext gc = getGraphicsContext2D();
        for(int i=getStartRow(); i<=getEndRow(); i++) {
            double startX = -camera.getOffsetX();
            double startY = -camera.getOffsetY() + (i - getStartRow()) * camera.getBlockSize();
            double endX = -camera.getOffsetX() + (getEndColumn() - getStartColumn()) * camera.getBlockSize();
            double endY = startY;
            gc.strokeLine(startX, startY, endX, endY);
        }

        for(int i=getStartColumn(); i<=getEndColumn(); i++) {
            double startX = -camera.getOffsetX() + (i - getStartColumn()) * camera.getBlockSize();
            double startY = -camera.getOffsetY();
            double endX = startX;
            double endY = -camera.getOffsetY() + (getEndRow() - getStartRow()) * camera.getBlockSize();
            gc.strokeLine(startX, startY, endX, endY);
        }
    }

    public void drawQuadrants() {

        GraphicsContext gc = getGraphicsContext2D();
        gc.save();

        gc.setStroke(Color.RED);
        gc.setLineWidth(camera.getScale() * 3);

        double centerColumn = (double)columns / 2;
        if (getStartColumn() < centerColumn && getEndColumn() > centerColumn
                && getStartRow() <= rows && getEndRow() >= 0) {
            double startX = -camera.getOffsetX() + (centerColumn - getStartColumn()) * camera.getBlockSize();
            double startY = -camera.getOffsetY();
            double endX = startX;
            double endY = -camera.getOffsetY() + getTotalVisibleRows() * camera.getBlockSize();
            gc.strokeLine(startX, startY, endX, endY);
        }

        double centerRow = (double)rows / 2;
        if (getStartRow() < centerRow && getEndRow() > centerRow
                && getStartColumn() <= columns && getEndColumn() >= 0) {
            double startX = -camera.getOffsetX();
            double startY = -camera.getOffsetY() + (centerRow - getStartRow()) * camera.getBlockSize();
            double endX = -camera.getOffsetX() + getTotalVisibleColumns() * camera.getBlockSize();
            double endY = startY;
            gc.strokeLine(startX, startY, endX, endY);
        }

        gc.restore();
    }

    public void drawCircle(Color color, double diameter, int x, int y) {

        y = rows - y - 1;

        GraphicsContext gc = getGraphicsContext2D();
        gc.setFill(color);

        double pixelX = -camera.getOffsetX() + ((x - getStartColumn()) * camera.getBlockSize() - diameter  / 2 + camera.getBlockSize() / 2);
        double pixelY = -camera.getOffsetY() + ((y - getStartRow()) * camera.getBlockSize() - diameter / 2 + camera.getBlockSize() / 2);
        gc.fillOval(pixelX, pixelY, diameter, diameter);

    }

    public void drawLine(int x1, int y1, int x2, int y2) {

        y1 = rows - y1 - 1;
        y2 = rows - y2 - 1;

        GraphicsContext gc = getGraphicsContext2D();
        gc.setLineWidth(.1);

        double worldX1 = x1 - getStartColumn();
        if (worldX1 < 0) {
            worldX1 = 0;
        }

        double worldX2 = x2 - getStartColumn();
        if (worldX2 < 0) {
            worldX2 = 0;
        }

        double worldY1 = y1 - getStartRow();
        if (worldY1 < 0) {
            worldY1 = 0;
        }

        double worldY2 = y2 - getStartRow();
        if (worldY2 < 0) {
            worldY2 = 0;
        }

        double pixelX1 = -camera.getOffsetX() + worldX1 * camera.getBlockSize();
        double pixelY1 = -camera.getOffsetY() + worldY1 * camera.getBlockSize();
        double pixelX2 = -camera.getOffsetX() + worldX2 * camera.getBlockSize();
        double pixelY2 = -camera.getOffsetY() + worldY2 * camera.getBlockSize();
        gc.strokeLine(pixelX1, pixelY1, pixelX2, pixelY2);
        gc.setLineWidth(1);
    }

    public void drawBlock(Color color, boolean outline, int x, int y) {
        this.drawBlock(color, outline, x, y, 1);
    }

    public void drawBlock(Color color, boolean outline, int x, int y, int squares) {


        if (x < getStartColumn() || x >= getEndColumn() || y < getStartRow() || y >= getEndRow())
            return;

        //y = rows - 1 - y; // transform coordinates so orgin is bottom left

        GraphicsContext gc = getGraphicsContext2D();
        gc.setFill(color);

        double worldX = x * squares - getStartColumn();
        double width = squares;
        if (worldX < 0) {
            worldX = 0;
            width = squares + (x * squares - getStartColumn());
        }

        double worldY = y * squares - getStartRow();
        double height = squares;
        if (worldY < 0) {
            worldY = 0;
            height = squares + (y * squares - getStartRow());
        }

        double pixelX = -camera.getOffsetX() + worldX * camera.getBlockSize();
        double pixelY = -camera.getOffsetY() + worldY * camera.getBlockSize();
        gc.fillRect(pixelX, pixelY, width * camera.getBlockSize(), height * camera.getBlockSize());
        if (outline){
            gc.strokeRect(pixelX, pixelY, width * camera.getBlockSize(), height * camera.getBlockSize());
        }
    }

    public int getStartColumn() {

        return (int)Math.max(0, camera.getX() / camera.getBlockSize());
    }

    public int getStartRow() {

        return (int)Math.max(0, (camera.getY() / camera.getBlockSize()));
    }

    public int getEndColumn() {
        double visibleColumns = Math.ceil((camera.getWidth() + camera.getOffsetX()) / camera.getBlockSize()) + getStartColumn();
        return (int)Math.min(visibleColumns + 1, columns);
    }

    public int getEndRow() {
        double visibleRows = Math.ceil((camera.getHeight() + camera.getOffsetY()) / camera.getBlockSize()) + getStartRow();
        return (int)Math.min(visibleRows + 1, rows);
    }

    public int getTotalVisibleColumns() {

        return getEndColumn() - getStartColumn();
    }

    public int getTotalVisibleRows() {

        return getEndRow() - getStartRow();
    }

    public void center() {

        double scale;
        if (getWidth() <= getHeight())
            scale = getWidth() / columns / Camera.BASE_BLOCK_SIZE;
        else
            scale = getHeight() / rows / Camera.BASE_BLOCK_SIZE;
        this.camera.setScale(scale);



        double x = (getWidth() - columns *  camera.getBlockSize()) / 2;
        double y = (getHeight() - rows * camera.getBlockSize()) / 2;

        camera.setX(-x);
        camera.setY(-y);
    }

    public class Camera {

        private static final double MINIMUM_BLOCK_SIZE = -1;
        private static final double MAXIMUM_BLOCK_SIZE = 100;
        private static final double BASE_BLOCK_SIZE = 10;
        private static final double SCALE_MULTIPLIER = 1.1;

        private double x;
        private double y;
        private double scale;
        private double width;
        private double height;

        public Camera(double width, double height) {

            this.x = 0;
            this.y = 0;
            this.scale = 1;
            this.width = width;
            this.height = height;
        }

        public double getX() {
            return x;
        }

        public void setX(double x) {
            this.x = x;
        }

        public double getY() {
            return y;
        }

        public void setY(double y) {
            this.y = y;
        }

        public double getScale() {
            return scale;
        }

        public void setScale(double scale) {
            this.scale = scale;
        }

        public double getWidth() {
            return width;
        }

        public void setWidth(double width) {
            this.width = width;
        }

        public double getHeight() {
            return height;
        }

        public void setHeight(double height) {
            this.height = height;
        }

        public void move(double dx, double dy) {

            this.x -= dx;
            this.y -= dy;
        }

        public void zoomIn() {

            if(getBlockSize() > MAXIMUM_BLOCK_SIZE) return;
            this.scale *= SCALE_MULTIPLIER;
        }

        public void zoomOut() {

            if (getBlockSize() < MINIMUM_BLOCK_SIZE) return;
            this.scale /= SCALE_MULTIPLIER;
        }

        public double getBlockSize() {
            return BASE_BLOCK_SIZE * scale;
        }

        public double getOffsetX() {
            return Math.min(x, x % getBlockSize());
        }

        public double getOffsetY() {
            return Math.min(y, y % getBlockSize());
        }

        @Override
        public String toString() {
            return String.format("---Camera---\n" +
                            "Position: (%f, %f)\n" +
                            "Offset:   (%f, %f)",
                    x, y, getOffsetX(), getOffsetY());
        }
    }
}