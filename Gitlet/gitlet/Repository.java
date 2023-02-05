package gitlet;

import java.io.File;
import java.io.IOException;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Hashtable;
import java.util.Set;
import java.util.List;


import static gitlet.Utils.*;

/** Represents a gitlet repository.
 *  This class is like the 'hood' behind main. It implements
 *  all of the commands in main, and interacts with all the other classes.
 *  does at a high level.
 *
 *  @author Iman
 */
public class Repository implements Serializable {
    /**
     * List all instance variables of the Repository class here with a useful
     * comment above them describing what that variable represents and how that
     * variable is used. We've provided two examples for you.
     */

    /** The current working directory. */
    public static final File CWD = new File(System.getProperty("user.dir"));
    /** The .gitlet directory. */
    public static final File GITLET_REPOSITORY = join(CWD, ".gitlet");

    /**Saves the hashmap of all the pointers/branches (including the master pointer)*/
    public static final File BRANCHES_FOLDER = join(CWD, ".branches");

    /** 'Points' to the most current commit (the 'leaf') by saving the
     *      SHA-1 Hash of that commit (i.e. commitID) in a text file*/
    static File headPointer = join(Repository.GITLET_REPOSITORY, "*");

    /** 'Points' to the most current commit (the 'leaf') by saving the
     *      SHA-1 Hash of that commit (i.e. commitID) in a text file*/
    static File masterPointer = join(BRANCHES_FOLDER, "master");

    /**File that stores the object info of the hashmap of all the
     * pointers/branches (including master pointer/branch)*/
    static File tableOfBranchesFile = join(BRANCHES_FOLDER, "tableOfBranches.txt.");

    /**Saves all of the pointers/branches*/
    static Hashtable<String, File> tableOfBranches;



    /** Setting up persistence for the repository, the Head and Master Pointer,
     * Commit Objects, and Blob Objects.*/
    public static void setupPersistence() {
        /**Makes the .gitlet/ folder if it does not yet exist.*/
        if (!GITLET_REPOSITORY.exists()) {
            GITLET_REPOSITORY.mkdir();
        }

        /**Makes the .gitlet/.commit/ folder if it does not yet exist.*/
        if (!Commit.COMMITS_FOLDER.exists()) {
            Commit.COMMITS_FOLDER.mkdir();
        }

        /**Makes the .gitlet/blobs/ folder if it does not yet exist.*/
        if (!Blob.BLOBS_FOLDER.exists()) {
            Blob.BLOBS_FOLDER.mkdir();
        }

        /**Makes the .gitlet/stagingArea/ folder if it does not yet exist.*/
        if (!StagingArea.STAGING_AREA_FOLDER.exists()) {
            StagingArea.STAGING_AREA_FOLDER.mkdir();
        }

        /**Makes the .stagingArea/state.txt file if it does not yet exist.*/
        if (!StagingArea.STATE.exists()) {
            try {
                StagingArea.STATE.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }

            initializeStagingArea();
        }

        /**Makes the .gitlet/.branches/ folder if it does not yet exist.*/
        if (!BRANCHES_FOLDER.exists()) {
            BRANCHES_FOLDER.mkdir();
        }

        if (!tableOfBranchesFile.exists()) {
            try {
                tableOfBranchesFile.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        }

        /**Makes the .capers/MASTER.txt file if it does not yet exist.
         * Saves the reference/address of (usually) the most recent commit in the text file*/
        if (!masterPointer.exists()) {
            try {
                masterPointer.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        }

        /**Makes the .capers/HEAD.txt file if it does not yet exist.
         * Saves the reference/address of (usually) the most recent commit in the text file*/
        if (!headPointer.exists()) {
            try {
                headPointer.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        }

    }

    /**Command words will be set as methods in this class (just lke Dog class in Lab 6)*/
    public static void setUpInitialCommit() {
        /** If there is already an initial commit,
         *      then there is an initialized head pointer pointing to it.
         *  If there is already an initial commit,
         *      don't make a new initial commit.*/
        String headPointerContents = readContentsAsString(headPointer);
        if (headPointerContents.equals("")) {
            Commit initialCommit = new Commit();
            initialCommit.saveCommit();
            //System.out.println(initial_commit);

            String pointers = initialCommit.getCommitID();

            /**POINTER AND MASTER MUST BE EQUAL TO EACH OTHER AT THIS POINT!*/
            Utils.writeContents(headPointer, pointers);
            Utils.writeContents(masterPointer, pointers);

            /**There is only one branch, and that is the master branch*/
            tableOfBranches = new Hashtable<>();
            tableOfBranches.put("master", masterPointer);
            Utils.writeObject(tableOfBranchesFile, tableOfBranches);
        } else {
            /**Failure case*/
            Utils.message("A Gitlet version-control system "
                    + "already exists in the current directory.");
        }
    }


    public static void makeCommit(String message) {
        /** HEAD is pointing at parent commit rn, didn't change it yet.*/
        String parentPointerID = readContentsAsString(headPointer);

        /** Making a commit that follows the parent (and copies its tracking of Blobs).*/
        Commit newCommit = new Commit(message, parentPointerID);

        /** Looks at Staging Area for Addition, checking if there's any Blobs we can add/update
         *      - Gets all the keys in blobsForAdd Hashtable so that I can get the values/Blobs */
        StagingArea currentState = StagingArea.fromFile();
        Set<String> setOfFileNames = currentState.blobsForAdd.keySet(); //Set of key of blobsForAdd

        /**FAILURE CASE: If no files have been stored, abort.
         *      Print the message "No changes added to the commit."*/
        if (currentState.blobsForAdd.isEmpty()
                && currentState.blobsForRemove.isEmpty()) {

            Utils.message("No changes added to the commit.");
            System.exit(0);
        } else {

            /**Looks at all the keys-values in the staging area add hashmap
             * And replaces/adds the new fileVersions/Blobs into new commit's
             *    hashmap (i.e. tracking modified files.)*/
            for (String key : setOfFileNames) {
                //Gets value/BlobID mapped from every key in StagingArea's blobsForAdd HashTable
                String value = currentState.blobsForAdd.get(key);

                //Adds/replaces/Updates those key-values to Commit's Hashtable of Blobs
                if (newCommit.tableOfBlobs.containsKey(key)) {
                    //If there was already a key/file/old_blob in the new/copied commit's hashmap,
                    // replace/update it with hashmaps updated filecontent (new blob).
                    newCommit.tableOfBlobs.replace(key, value);
                } else {
                    //Otherwise, just add it to new commit's hashmap
                    newCommit.tableOfBlobs.put(key, value);
                }
            }

            /**Clear staging area after moving contents to commit's hashmap
             * COMMIT CLEARS THE STAGING AREA - INCLUDING REMOVE STAGING AREA*/
            currentState.blobsForAdd.clear();
            currentState.blobsForRemove.clear();

            /** At the end of our a new commit, the HEAD and ACTIVE BRANCH/POINTER moves*/
            Utils.writeContents(headPointer, newCommit.getCommitID());

            //Finding active branch and 'moving' it also
            tableOfBranches = Utils.readObject(Repository.tableOfBranchesFile, Hashtable.class);
            Set<String> setOfBranchNames = tableOfBranches.keySet();
            for (String branchName : setOfBranchNames) {
                File findPointer = Utils.join(BRANCHES_FOLDER, branchName);
                String findPointerAdress = Utils.readContentsAsString(findPointer);

                if (findPointerAdress.equals(parentPointerID)) {
                    Utils.writeContents(findPointer, newCommit.getCommitID());
                }
            }

            /**Saving this new commit, and all of the changes made*/
            currentState.saveState();
            newCommit.saveCommit();
            //System.out.println(newCommit);
        }
    }

    public static Blob makeBlob(String fileName) {
        Blob newBlob = new Blob(fileName);
        newBlob.saveBlob();
        //System.out.println(newBlob);
        return newBlob;
    }

    public static void initializeStagingArea() {
        StagingArea newStagingArea = new StagingArea(); //Creates an empty Staging Area.
        newStagingArea.saveState();
        //System.out.println(newStagingArea);
    }

    public static void addToStagingArea(String fileVersion) {
        /**Check if repo is initialized first! (Did user call init?)*/
        String headP = Utils.readContentsAsString(headPointer);
        if (headP.equals("")) {
            Utils.message("Not in an initialized Gitlet directory.");
            System.exit(0);
        }


        /**Finding the file needed to add to Staging Area*/
        File findFile = Utils.join(CWD, fileVersion);
        if (findFile.exists()) {
            /** If the current working version of the file is identical to the version
             *      in the current commit, do not stage it to be added, and remove it
             *      from the staging area if it is already there */
            Commit parentCommit = Commit.getParentCommit(headPointer);

            //If they are the same, STOP.
            if (doCommitAndFileHaveSameContents(parentCommit, findFile, fileVersion)) {
                //Deleting any changes made (and objects created) and getting out.
                headP = null;
                findFile = null;
                parentCommit = null;
                //Utils.message("No changes were made to this file!!");
                System.exit(0);
            }

            //Create a Blob Object that contains the saved contents in a file, and save it.
            Blob addingBlob = makeBlob(fileVersion);
            //Gets the persistent Staging Area state
            StagingArea changingState = StagingArea.fromFile();
            //Add Blob to staging area. Key = fileName and Value = Blob reference/ID
            changingState.blobsForAdd.put(fileVersion, addingBlob.getBlobID());

            //Save changes in Staging Area
            changingState.saveState();  //Save changes made to Staging Area state
        } else {
            Utils.message("File does not exist.");
            System.exit(0);
        }

    }

    /**If commit has the same blob as file, then they both have
     * the same contents/version of the file.*/
    public static boolean doCommitAndFileHaveSameContents(Commit commitObject,
                                                          File fileObject, String fileName) {

        //Get file content from commit
        String commitBlobID = commitObject.tableOfBlobs.get(fileName);

        /** If fileBlobID == null, then the only commit there is is the initial commit
            Therefore, there wouldn't be any tracked versions of this file
            (file wouldn't be saved anywhere) and no commits would have the same contents
            as target file.*/
        if (commitBlobID != null) {
            Blob commitBlob = Blob.fromFile(commitBlobID);
            File commitBlobFileObject = commitBlob.getFileContents();
            String commitBlobFileContent = Utils.readContentsAsString(commitBlobFileObject);

            //Get file content from target file
            String targetFileConten = Utils.readContentsAsString(fileObject);

            //See if the file contents are the same
            if (commitBlobFileContent.equals(targetFileConten)) {
                return true;
            }
        }
        return false;
    }

    /** Unstages the file if it is currently staged for addition,
     *  IF FILE IS TRACKED IN CURRENT COMMIT, stage it for removal (in Staging Area), and
     *  remove the file from  working directory (CWD) IF user hasn't already done so*/
    public static void unstageFile(String fileName) {
        //Gets the persistent Staging Area state
        StagingArea changingState = StagingArea.fromFile();
        Commit currentCommit = Commit.getParentCommit(headPointer);

        /**Failure Case: If the file is NEITHER STAGED NOR TRACKED BY HEAD COMMIT,
         * print the error message "No reason to remove the file."*/
        if (!changingState.blobsForAdd.containsKey(fileName)
                && !currentCommit.tableOfBlobs.containsKey(fileName)) {
            Utils.message("No reason to remove the file.");
            System.exit(0);
        }

        /**If the file/blob is currently staged for addition, take it out*/
        if (changingState.blobsForAdd.containsKey(fileName)) {
            //Remove key-value from the staging area for addition
            changingState.blobsForAdd.remove(fileName);
        }

        /**If the file/blob is tracked/in the current commit, stage it for removal*/
        if (currentCommit.tableOfBlobs.containsKey(fileName)) {
            //Gets the Blobid/Blobreference/value
            String value = currentCommit.tableOfBlobs.get(fileName);
            changingState.blobsForRemove.put(fileName, value);

            /** Remove the file from  working directory (CWD) IF user hasn't already done so*/
            //Finding the file needed to remove from CWD
            File findFile = Utils.join(CWD, fileName);
            if (findFile.exists()) {
                Utils.restrictedDelete(findFile);
            }
        }

        /**SAVE CHANGES FOR PERSISTENCE!
         * CURRENTCOMMIT SHOULD NOT HAVE BEEN CHANGED!!*/
        changingState.saveState();
    }

    /** STARTING AT THE HEAD COMMIT, display to information about each commit backwards
     * (display should start from current commit UP TOP to root/initial commit AT BOTTOM).
     *      -Follow the parent commit links/pointers/references,
     *        IGNORING ANY SECOND PARENTS FOUND IN MERGE COMMITS
     *      -FOR MERGECOMMITS (if there are 2 parent commits)......
     * */
    public static void goThroughCurrentBranch() {
        //Follow parent links RECURSIVELY and print commits as you go
        Commit currentCommit = Commit.getParentCommit(headPointer);
        Commit.recursivelyIterate(currentCommit);
        //FILL IN WITH YOU KNOW WHAT!: Fill out the merge case!
    }

    /**Displays information about ALL COMMITS EVER MADE.
     *   -ORDER OF COMMITS DOESN'T MATTER*/
    public static void printAllCommits() {
        //Grabs all the files of commits in Commit folder, and saves it in a list
        List<String> listOfAllCommits = new ArrayList<String>();
        listOfAllCommits = Utils.plainFilenamesIn(Commit.COMMITS_FOLDER);

        //Gets the commitID of each file, converts it into a commit, and prints the commit's info
        //Does this for every commit/file in the CommitFolder
        for (String commitIDs : listOfAllCommits) {
            Commit currentCommit = Commit.fromFile(commitIDs);
            System.out.println(currentCommit);
        }
    }

    /** Prints out ALL of the IDS of all commits that have the same given log message
     *      -If there are multiple such commits, it prints the IDs out in SEPERATE LINES.
     *      -Hint: Same as the one for global-log.*/
    public static void findCommitsWithSameMessage(String doYouHaveThisMessage) {
        //Grabs all the files of commits in Commit folder, and saves it in a list
        List<String> listOfAllCommits = new ArrayList<String>();
        listOfAllCommits = Utils.plainFilenamesIn(Commit.COMMITS_FOLDER);
        boolean haveMessage = false;

        //Gets the commitID of each file, converts it into a commit, and prints the commit's info
        //Does this for every commit/file in the CommitFolder
        for (String commitIDs : listOfAllCommits) {
            Commit currentCommit = Commit.fromFile(commitIDs);
            if (currentCommit.getMessage().equals(doYouHaveThisMessage)) {
                haveMessage = true;
                System.out.println(currentCommit.getCommitID());
            }
        }

        /** Failure Case: If there are no commits that have the message,
         * then print error message: "Found no commit with that message."**/
        if (!haveMessage) {
            Utils.message("Found no commit with that message.");
        }
    }

    public static void showStatus() {
        System.out.println(Status.showBranches());
        System.out.println(Status.showStagedFilesBranch());
        System.out.println(Status.showRemovedFiles());
        System.out.println(Status.showModsNotStagedForCommit());
        System.out.println(Status.showUntrackedFiles());
    }

    /**Takes the version of the file as it exists in HEAD COMMIT and
     *  puts it in the working directory (CWD), OVERWRITING the version
     *  of the file that's already there if this is one.
     *  The new version of the file is not staged.*/
    public static void checkoutFileName(String fileName, File pointer) {
        Commit headCommit = Commit.getParentCommit(pointer);
        /**If file does not exist in the head commit, abort,
         *  printing error message: "File does not exist in that commit."
         *  DO NOT CHANGE THE CWD*/
        if (!headCommit.tableOfBlobs.containsKey(fileName)) {
            Utils.message("File does not exist in that commit.");
            System.exit(0);
        }

        //Find head commit's filecontents/Blob
        String headCommitsBlobID = headCommit.tableOfBlobs.get(fileName);
        Blob headCommitsBlob = Blob.fromFile(headCommitsBlobID);
        File headCommitBlobFileObject = headCommitsBlob.getFileContents();
        byte[] headCommitBlobFileContents = readContents(headCommitBlobFileObject);
        //System.out.println(headCommitsBlob);

        //Finds the file in CWD and overwrites (if needed)
        // the file in CWD with the fileContents in commit
        File findFile = Utils.join(CWD, fileName);
        if (findFile.exists()) {
            Utils.writeContents(findFile, headCommitBlobFileContents);
        }

        //System.out.println(Utils.readContentsAsString(findFile));

    }

    /**Takes the version of the file as it EXISTS in the COMMIT  WITH THE GIVEN ID
     * and PUTS IT IN THE CWD, OVERWRITING the version of the file that's already
     * there if there is one.
     * The new version of the file is not staged.*/
    public static void checkoutCommitIDAndFileName(String commitID, String fileName) {
        /**Failure case: If no commit with the given ID exists,
         * print error message: No commit with that id exists."
         * and abort.
         */
        File noCommitID = Utils.join(Commit.COMMITS_FOLDER, commitID);
        if (!noCommitID.exists()) {
            Utils.message("No commit with that id exists.");
            System.exit(0);
        }

        Commit targetCommit = Commit.fromFile(commitID);
        /**If file does not exist in the given commit, abort,
         *  printing error message: "File does not exist in that commit."
         *  DO NOT CHANGE THE CWD*/
        if (!targetCommit.tableOfBlobs.containsKey(fileName)) {
            Utils.message("File does not exist in that commit.");
            System.exit(0);
        }

        String targetCommitsBlobID = targetCommit.tableOfBlobs.get(fileName);
        Blob targetCommitsBlob = Blob.fromFile(targetCommitsBlobID);
        File targetCommitBlobFileObject = targetCommitsBlob.getFileContents();
        byte[] targetCommitBlobFileContents = readContents(targetCommitBlobFileObject);
        System.out.println(targetCommit);
        System.out.println(targetCommitsBlob);

        //Finds the file in CWD and overwrites (if needed)
        // the file in CWD with the fileContents in commit
        File findFile = Utils.join(CWD, fileName);
        System.out.println("Before: " + Utils.readContentsAsString(findFile));

        if (findFile.exists()) {
            Utils.writeContents(findFile, targetCommitBlobFileContents);
        }
        System.out.println("After: " + Utils.readContentsAsString(findFile));
    }

    /** - Takes all files in the commit at the head of THIS GIVEN BRANCH, and
     *      puts them in the CWD, OVERWRITING the version of the files that are
     *      already there if they exist.
     *      -> All the files in CWD is overwritten by fileVersions in OTHER BRANCHE'S HEAD COMMIT
     *  - Any files that are tracked in the CURRENT BRANCH (HEAD) but are NOT PRESENT in the
     *        CHECKED OUT BRANCH (other branch/pointer) are DELETED.
     *
     *  - The STAGING AREA IS CLEARED, UNLESS the checked out branch is the current branch.
     *  - At the end of this command, the GIVEN BRANCH will now be considered the CURRENT BRANCH
     *         -> HEAD POINTER WILL BE POINTING TO THIS BRANCH/POINTER!
     */
    public static void checkoutBranch(String branchName) {
        //Gets the commit head of THIS GIVEN BRANCH
        tableOfBranches = (Hashtable<String, File>) Utils.readObject(tableOfBranchesFile,
                Hashtable.class);

        File branchPointerToCheckOut = tableOfBranches.get(branchName);
        Commit commitOfBranchToCheckOut = Commit.getParentCommit(branchPointerToCheckOut);
        //Gets all of the files in the commit
        Set<String> setOfFileNamesInCommit = commitOfBranchToCheckOut.tableOfBlobs.keySet();

        //Gets the commit head of THE CURRENT/ACTIVE BRANCH
        Commit currentBranch = Commit.getParentCommit(headPointer);
        Set<String> setOfFileNamesCurrentBranch = currentBranch.tableOfBlobs.keySet();

        //Gets the Staging Area's current state.
        StagingArea currentState = StagingArea.fromFile();

        /**FAILURE CASES:*/
        failureCasesForCheckoutBranch(branchName, setOfFileNamesInCommit,
                currentBranch, currentState);

        /**Takes all files in the commit at the head of THIS GIVEN BRANCH, and
         * puts them in the CWD, OVERWRITING the version of the files that are
         * already there if they exist.*/
        //For each file
        for (String fileName : setOfFileNamesInCommit) {
            File isFileInCWD = Utils.join(CWD, fileName);

            if (isFileInCWD.exists()) {
                //takes the version of the file and puts them in the CWD,
                // OVERWRITING the version of the files that are already there if they exist.
                checkoutFileName(fileName, branchPointerToCheckOut);
            } else {
                //Creates the file in the CWD and puts contents into it
                try {
                    isFileInCWD.createNewFile();
                } catch (IOException excp) {
                    throw new IllegalArgumentException(excp.getMessage());
                }

                String blobID = commitOfBranchToCheckOut.tableOfBlobs.get(fileName);
                Blob savedBlob = Blob.fromFile(blobID);
                File fileSavingContents = savedBlob.getFileContents();

                Utils.writeContents(isFileInCWD, Utils.readContents(fileSavingContents));
            }
        }

        /**Any files that are tracked in the CURRENT BRANCH (HEAD) but are NOT PRESENT in the
         * CHECKED OUT BRANCH (other branch/pointer) are DELETED.*/

        for (String fileTrackedInCurrentBranch: setOfFileNamesCurrentBranch) {
            if (!commitOfBranchToCheckOut.tableOfBlobs.containsKey(fileTrackedInCurrentBranch)) {
                File deleteFile = Utils.join(CWD, fileTrackedInCurrentBranch);
                Utils.restrictedDelete(deleteFile);
            }
        }

        /**Staging Area is cleared*/
        currentState.blobsForAdd.clear();
        currentState.blobsForRemove.clear();
        currentState.saveState();

        /**At the end of this command, the GIVEN BRANCH will now be considered the CURRENT BRANCH
         * -> HEAD POINTER WILL BE POINTING TO THIS BRANCH/POINTER!*/
        Utils.writeContents(headPointer, Utils.readContentsAsString(branchPointerToCheckOut));

        /**Saving the tableofBranches*/
        Utils.writeObject(tableOfBranchesFile, tableOfBranches);
    }

    static void failureCasesForCheckoutBranch(String branchName, Set<String> setOfFileNamesInCommit,
                                              Commit currentBranch, StagingArea currentState) {
        /**FAILURE CASES:*/
        /**If no branch with that name exists, print No such branch exists. */
        if (!tableOfBranches.containsKey(branchName)) {
            Utils.message("No such branch exists.");
            System.exit(0);
        }
        /** If that branch is the current branch, print No need to checkout the current branch.*/
        String currentBranchPointer = Utils.readContentsAsString(headPointer);
        String thisBranchPointer = Utils.readContentsAsString(tableOfBranches.get(branchName));
        if (currentBranchPointer.equals(thisBranchPointer)) {
            Utils.message("No need to checkout the current branch.");
            System.exit(0);
        }
        /**If a working file is untracked in the current branch and would be overwritten by the
         * checkout, print "There is an untracked file in the way; delete it, or add and commit
         * it first."
         * and exit; perform this check before doing anything else. Do not change the CWD.*/
        //File must be in the given branch's head commit in order for anything to be overwritten
        //UNTRACKED DEFINITION: File is in the CWD but not tracked in current head commit and
        // not staged for addition (right?)
        for (String fileName : setOfFileNamesInCommit) {
            File isFileInCWD = Utils.join(CWD, fileName);

            if (isFileInCWD.exists()
                    && !currentBranch.tableOfBlobs.containsKey(fileName)
                    && !currentState.blobsForAdd.containsKey(fileName)) {
                Utils.message("There is an untracked file in the way; "
                        + "delete it, or add and commit it first.");
                System.exit(0);
            }
        }
    }

    public static void createNewBranch(String branchName) {
        //Creates a branch/pointer that is PERSISTENT -> MUST BE SAVED IN A FILE AND SERIALIZED!!
        File otherBranch = join(BRANCHES_FOLDER, branchName);

        if (!otherBranch.exists()) {
            try {
                otherBranch.createNewFile();
            } catch (IOException excp) {
                throw new IllegalArgumentException(excp.getMessage());
            }
        } else {
            /** Failure case: If a branch with the given name already exists,
             *      print the error message: "A branch with that name already exists."*/
            Utils.message("A branch with that name already exists.");
            System.exit(0);
        }

        /** Initially, a new branch/pointer must point to the most current (i.e. HEAD commit).
         *  Then, branch must be saved/serialized so that branch persists.
         *  -> HEAD, MASTER, AND OTHER BRANCH SHOULD ALL BE POINTING TO THE SAME COMMIT!!*/
        Utils.writeContents(otherBranch, Utils.readContents(headPointer));

        /** Adds that branch into the table ofbranches
         * (in case there are multiple branches), and saves the changed table*/
        tableOfBranches = (Hashtable<String, File>) Utils.readObject(tableOfBranchesFile,
                Hashtable.class);
        tableOfBranches.put(branchName, otherBranch);
        Utils.writeObject(tableOfBranchesFile, tableOfBranches);
    }


    public static void main(String[] args) {
        File testerFile = Utils.join(CWD, "TesterForMe.txt");
        Commit testerCommit = Commit.getParentCommit(headPointer);

        //ERROR: Once checkout change is made to CW folder,
        // all of the previous commit's Blobs also all change!
        //checkoutCommitIDAndFileName("e30605afac70977e4cd412e0db7cc9814ee93369",
        // "TesterForMe.txt");
        //makeCommit("AnotherTester.txt with Q");

        Commit commitWithQ = Commit.fromFile("a7a9ad1d95151700297ffd620e939f8ef3dfa80e");
        Commit commitWithE = Commit.fromFile("7023d388764c5a4bf63d0b001e5277efc7d3ac77");
        System.out.println(commitWithQ);
        System.out.println(commitWithE);

        String qBlobID = commitWithQ.tableOfBlobs.get("AnotherTesterForMe.txt");
        String eBlobID = commitWithQ.tableOfBlobs.get("AnotherTesterForMe.txt");

        Blob qBlob = Blob.fromFile(qBlobID);
        Blob eBlob = Blob.fromFile(eBlobID);

        File qBFile = qBlob.getFileContents();
        File eBFile = eBlob.getFileContents();

        System.out.println(Utils.readContentsAsString(qBFile));
        System.out.println(Utils.readContentsAsString(eBFile));

    }
}
