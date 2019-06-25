package edu.udel.ntsee.bombdetection.data;

import java.util.List;

public class AdHoc {

    private int leaderID;
    private List<Integer> childrenIDs;

    public AdHoc(int leaderID, List<Integer> childrenIDs) {

        this.leaderID = leaderID;
        this.childrenIDs = childrenIDs;
    }

    public int getLeaderID() {
        return leaderID;
    }

    public List<Integer> getChildrenIDs() {
        return childrenIDs;
    }
}
