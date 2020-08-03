package edu.udel.ntsee.bombdetection.data;

public class Infection {

    private int id;
    private Type type;

    public Infection(int id, Type type) {
        this.id = id;
        this.type = type;
    }

    public int getID() {
        return this.id;
    }

    public Type getType() {
        return this.type;
    }

    public enum Type {
        NONE,
        HOST,
        INFECTED;

        public static Type valueOf(int id) {
            if (id < 0 || id > Type.values().length) {
                throw new IllegalArgumentException("Invalid infection type id - " + id);
            }

            return Type.values()[id];
        }
    }
}
