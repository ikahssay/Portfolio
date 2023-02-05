package gitlet;

import java.io.File;
import java.io.IOException;
import java.io.Serializable;

public class Blob implements Serializable {

    /**The File object that corresponds to the directory containing all the serialized Blob objects.
     * This is static since all Blob objects are stored within the same directory.*/
    static final File BLOBS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY, ".blobs");

    /** The SHA1 value of this Blob (it's ID).
     * Will also need to put a "B" in front of the ID to distinguish Blobs from Commits
     * OR SINCE THEY'RE IN DIFF FOLDERS, I CAN LOOK FOR HASHES IN JUST THE BLOB FOLDER RIGHT??*/
    private String blobID;

    /** Name of the file Blob will be saving the contents of.*/
    private String fileName;


    /** Saves the target files version/contents into a text file.
     * Thought about saving it as a String, but the contents may be too large.*/
    private String fileContentsName;
    private File fileContents;

    Blob(String file) {
        byte[] serializedBlob = Utils.serialize(this);
        blobID = Utils.sha1(serializedBlob);
        fileName = file;
        fileContentsName = "fileContents for Blob " + this.blobID;
        fileContents = Utils.join(BLOBS_FOLDER, fileContentsName);
        createBlob(file);
    }

    /**Helper method to Constructor*/
    private void createBlob(String file) {

        /**Creating a file to save Blob Object*/
        File newBlob = Utils.join(BLOBS_FOLDER, this.blobID);

        /** Creates a new Blob*/
        try {
            // Makes the Blob if it does not yet exist.
            newBlob.createNewFile();
            // Makes the .blobs/fileContents.txt file if it does not yet exist.
            fileContents.createNewFile();
        } catch (IOException excp) {
            throw new IllegalArgumentException(excp.getMessage());
        }

        /**Blobs gets the contents in a target file and saves it in the
         * fileContents File Object (a text file)*/
        saveFileContents(file);

        /**Everything within Blob object (including File Object fileContents)
         * is now written and saved in a file*/
        Utils.writeObject(newBlob, this);   //Save object's info persistently in a file
    }

    /**Helper method to createBlob method*/
    private void saveFileContents(String file) {
        /**Finding the file Blob needs to read (IN CWD) from and save contents. */
        File readingFile = Utils.join(Repository.CWD, file);
        String contents = Utils.readContentsAsString(readingFile);
        Utils.writeContents(fileContents, contents);

        /**Checking to see that these are all the same
        Utils.readContentsAsString(readingFile));
        System.out.println(Utils.readContentsAsString(fileContents));
        System.out.println(contents);
        */
    }

    /**
     * Saves a Blob to a file for future use.
     */
    public void saveBlob() {
        File newBlob = Utils.join(BLOBS_FOLDER, this.blobID);

        //If the file with the dog (object)'s info doesn't exist.
        if (!newBlob.exists()) {
            try {
                // Makes the .gitlet/commit object file if it does not yet exist.
                newBlob.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        }

        Utils.writeObject(newBlob, this);   //Save object's info persistently in a file
    }

    /**
     * Reads in and deserializes a Blob from a file with name NAME in DOG_FOLDER.
     *
     * @param id ID (or 'name') of Commit to load
     * @return Commit read from file
     */
    public static Blob fromFile(String id) {
        File blobFile = Utils.join(BLOBS_FOLDER, id);
        return Utils.readObject(blobFile, Blob.class);
    }

    public File getFileContents() {
        return this.fileContents;
    }

    /** Gets the Blob's SHA1 Hashcode ID*/
    public String getBlobID() {
        return this.blobID;
    }

    @Override
    public String toString() {
        return String.format(
                "Blob Info: \nID: %s \nVersion of File: %s \nFile Contents: %s \n",
                blobID, fileName, Utils.readContentsAsString(fileContents));
    }
}
