package gitlet;
import java.io.File;
//import java.io.IOException;
import java.io.Serializable;
import java.util.Hashtable;

public class StagingArea implements Serializable {

    static final File STAGING_AREA_FOLDER = Utils.join(Repository.GITLET_REPOSITORY,
            ".stagingArea");
    static final File STATE = Utils.join(STAGING_AREA_FOLDER, "state.txt");

    /**
     * This is static since all Blob objects are stored within the same directory.
     * Hashtable of Blob references (the The SHA1 value of
     * Blobs (it's ID)) that are staged for addition.
             * -Key = fileName
             * -Value = Blob's SHA-1 Hash
     */
    Hashtable<String, String> blobsForAdd;

    /**
     * This is static since all Blob objects are stored within the same directory.
     * Hashtable of Blob references (the The SHA1 value
     * of Blobs (it's ID)) that are staged for removal.
             * -Key = fileName
             * -Value = Blob's SHA-1 Hash
     */
    Hashtable<String, String> blobsForRemove;

    public StagingArea() {
        blobsForAdd = new Hashtable<String, String>();
        blobsForRemove = new Hashtable<String, String>();
    }


    /**
     * Saves state of StagingArea to a file for future use.
     */
    public void saveState() {
        Utils.writeObject(STATE, this);   //Save object's info persistently in a file
    }

    /** Gets the object (serialized in a file) from persisting 'state.txt' file*/
    public static StagingArea fromFile() {
        File stagingAreaFile = Utils.join(STAGING_AREA_FOLDER, "state.txt");
        return Utils.readObject(stagingAreaFile, StagingArea.class);
    }
}
