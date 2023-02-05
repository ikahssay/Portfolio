package gitlet;

import java.io.IOException;
import java.io.Serializable;
import java.util.Hashtable;
import java.util.Set;
import java.io.File;
import java.util.Date;
import java.text.SimpleDateFormat;


import static gitlet.Utils.readContentsAsString;

/** Represents a gitlet commit object.
 *  This class represents a Commit that will be stored in a file (i.e. it will persist).
 *  Because each Commit will have a unique ID, we may simply use
 *  that as the name of the file that the object is serialized to.
 *
 * All Commit objects are serialized within the COMMITS_FOLDER which
 *      is within the GITLET_REPOSITORY_FOLDER.
 * The Commit class has helpful methods that will return the Commit object
 *      corresponding to some String ID given to it, as well as write that
 *      Commit to a file to persist its changes.
 *
 *  @author: Iman Kahssay
 */
public class Commit implements Serializable {
    /**
     * List all instance variables of the Commit class here with a useful
     * comment above them describing what that variable represents and how that
     * variable is used. We've provided one example for `message`.
     */

    /**The File object that corresponds to the directory containing
     *  all the serialized Commit objects.
     * This is static since all Commit objects are stored within the same directory.*/
    static final File COMMITS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY, ".commits");

    /** The message of this Commit. */
    private String message;

    /** The time this Commit was made.*/
    private String timestamp;

    /** The SHA1 value of this Commit (it's ID).
     * Will also need to put a "C" in front of the ID to distinguish Commits from Blobs
     * OR SINCE THEY'RE IN DIFF FOLDERS, I CAN LOOK FOR HASHES IN JUST THE COMMIT FOLDER RIGHT??*/
    private String commitID;

    /** Hashtable of Blob references (the The SHA1 value of Blobs (it's ID).)
     *     -Key = fileName
     *     -Value = SHA-1 Hash */
    Hashtable<String, String> tableOfBlobs;

    /** Commit's parent (the commit that was made before this one) ID*/
    private String commitParentID;

    /** Format of timestamp always stays the same.*/
    final SimpleDateFormat timeFormat = new SimpleDateFormat("EEE MMM d HH:mm:ss yyyy Z");

    /**Setting up the empty constructor to be the initial commit.*/
    public Commit() {
        message = "initial commit";
        timestamp = getBeginningOfTime();

        byte[] serializedCommit = Utils.serialize(this);
        commitID = Utils.sha1(Utils.serialize(serializedCommit));

        commitParentID = null;
        tableOfBlobs = new Hashtable<String, String>();
    }

    /**Setting up the loaded constructor to follow the initial commit.
     * First, 'clone' the parent commit, and then change the metadata.
     * Then, (in Repo) change the hashtable of Blobs*/
    public Commit(String logMessage, String parentID) {
        timestamp = getTime();
        message = logMessage;

        byte[] serializedCommit = Utils.serialize(this);
        commitID = Utils.sha1(Utils.serialize(serializedCommit));

        commitParentID = parentID;
        Commit parent = fromFile(parentID);

        /**Making a deep copy of the parent's tableOfBlobs*/
        Hashtable<String, String> copy = new Hashtable<String, String>(parent.tableOfBlobs.size());
        Set<String> originalSet = parent.tableOfBlobs.keySet();

        for (String fileNames: originalSet) {
            copy.put(fileNames, parent.tableOfBlobs.get(fileNames));
        }

        tableOfBlobs = copy;
    }

    /**Helper Function -> Helps initialize initial commit's timestamp*/
    private String getBeginningOfTime() {
        Date beginningOfTime = new Date();
        beginningOfTime.setTime(0);
        return timeFormat.format(beginningOfTime.getTime());
    }

    /**Helper Function -> Helps initialize new commit's (AFTER initial commit) timestamp*/
    private String getTime() {
        Date date = new Date();
        return timeFormat.format(date.getTime());
    }

    /** Gets the Commits's SHA1 Hashcode ID*/
    public String getCommitID() {
        return this.commitID;
    }

    public String getMessage() {
        return this.message;
    }

    /**
     * Reads in and deserializes a Commit from a file with name NAME in DOG_FOLDER.
     *
     * @param id ID (or 'name') of Commit to load
     * @return Commit read from file
     */
    public static Commit fromFile(String id) {
        File commitFile = Utils.join(COMMITS_FOLDER, id);
        return Utils.readObject(commitFile, Commit.class);
    }

    /**
     * Saves a Commit to a file for future use.
     */
    public void saveCommit() {
        File newCommit = Utils.join(COMMITS_FOLDER, this.commitID);

        //If the file with the Commit (object)'s info doesn't exist.
        if (!newCommit.exists()) {
            try {
                // Makes the .gitlet/commit object file if it does not yet exist.
                newCommit.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        }

        Utils.writeObject(newCommit, this);   //Save object's info persistently in a file
    }

    /**Gets parent Commit object*/
    public static Commit getParentCommit(File parentID) {
        /** HEAD points to parent commit since it hasn't
         *  been updated to point to most recent commit.*/
        String hp = readContentsAsString(parentID);
        Commit parentCommit = Commit.fromFile(hp);
        return parentCommit;
    }

    /** Helper function to goThroughCurrentBranch()
     *      -Follows parent pointers (recursively) and prints all of the commits. */
    public static void recursivelyIterate(Commit currentCommit) {
        Commit tracker = currentCommit;
        System.out.println(tracker);

        if (tracker.commitParentID == null) {
            return;
        } else {
            //commit is now the parent commit
            tracker = fromFile(tracker.commitParentID);
            recursivelyIterate(tracker);
        }
    }


    @Override
    public String toString() {

        return String.format(
                "=== \ncommit %s \nDate: %s \n%s \n",
                this.commitID, this.timestamp, this.message);

        /**
        return String.format(
                "=== \ncommit %s \nDate: %s \n%s \nParentID: %s \n",
                commitID, timestamp, message, commitParentID);

         */
    }


}

