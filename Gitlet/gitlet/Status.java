package gitlet;
//import jdk.javadoc.doclet.Reporter;

import static gitlet.Utils.*;

import java.io.File;
import java.util.ArrayList;
import java.util.Hashtable;
import java.util.Set;
import java.util.List;

public class Status {

    /** Displays what branches currently exist, and marks the current branch with a *.*/
    public static String showBranches() {
        String allTheBranches = "";

        String headPointerAddress = Utils.readContentsAsString(Repository.headPointer);
        Repository.tableOfBranches = Utils.readObject(Repository.tableOfBranchesFile,
                Hashtable.class);
        Set<String> setOfBranchNames = Repository.tableOfBranches.keySet();
        for (String branchName : setOfBranchNames) {
            File findPointer = Utils.join(Repository.BRANCHES_FOLDER, branchName);
            String findPointerAdress = Utils.readContentsAsString(findPointer);

            if (findPointerAdress.equals(headPointerAddress)) {
                allTheBranches = allTheBranches + '\n' + "*" + branchName;
            } else {
                allTheBranches = allTheBranches + '\n' + allTheBranches;
            }
        }

        String showStagedFiles = String.format(
                "=== Branches === %s \n",
                allTheBranches);

        return showStagedFiles;
    }

    public static String showStagedFilesBranch() {
        StagingArea currentState = StagingArea.fromFile();

        /** Get the textfiles saved in the Stage For Addition and save them in a String*/
        String filesInAddStage = "";
        Set<String> setOfFiles = currentState.blobsForAdd.keySet();

        //Gets all of the file names Dtaged For Addition
        for (String file : setOfFiles) {
            filesInAddStage = filesInAddStage + '\n' + file;
        }

        /**Make the "Staged Files" category in status
         *  -Have all of the filenames saved in the add stage in String filesInAddStage.
         *   Now put that in a String with the correct 'format'
         *  */
        String showStagedFiles = String.format(
                "=== Staged Files === %s \n",
                filesInAddStage);

        return showStagedFiles;
    }

    public static String showRemovedFiles() {
        StagingArea currentState = StagingArea.fromFile();

        /** Get the textfiles saved in the Stage For Removal and save them in a String*/
        String filesInRemovalStage = "";
        Set<String> setOfFiles = currentState.blobsForRemove.keySet();

        //If there are files in
        for (String file : setOfFiles) {
            filesInRemovalStage = filesInRemovalStage + '\n' + file;
        }

        /**Make the "Removed Files" category in status
         *  -Have all of the filenames saved in the removal stage in String filesInRemovalStage.
         *   Now put that in a String with the correct 'format'
         *  */
        String showRemovedFiles = String.format(
                "=== Removed Files === %s \n",
                filesInRemovalStage);

        return showRemovedFiles;
    }

    /** File must be:
     *      - Tracked in current commit, changed in CWD, but NOT staged, OR
     *      - Staged for Addition, but with different contents than in the working directory, OR
     *      - Staged for Addition, but deleted in the working Directory, OR
     *      - Not staged for removal, but tracked in the current commit
     *              and deleted in the working directory.*/
    public static String
        showModsNotStagedForCommit() {

        //Got a list of text all the files in CWD
        List<String> filesInCWD = new ArrayList<String>();
        filesInCWD = plainFilenamesIn(Repository.CWD);

        //Got current commit and current state of Staging Area
        Commit currentCommit = Commit.getParentCommit(Repository.headPointer);
        StagingArea currentState = StagingArea.fromFile();
        String filesNotStagedForCommit = "";

        /**Looks at the files in working directory, and checks that it's not being tracked
         nor staged for addition. Saves it in String filesUntracked*/
        for (String file : filesInCWD) {
            if (file.contains(".txt")) {
                File fileAsAnObject = Utils.join(Repository.CWD, file);

                /**If file was tracked in current commit, changed in CWD, but NOT staged*/
                if (!Repository.doCommitAndFileHaveSameContents(currentCommit,
                        fileAsAnObject, file)
                        && !currentState.blobsForAdd.containsKey(file)) {

                    filesNotStagedForCommit = filesNotStagedForCommit + '\n' + file;
                }

                /**If file is Staged For Addition, but with different contents
                 * than in the working directory**/
                if (currentState.blobsForAdd.containsKey(file)) {
                    //Get contents of the file in Staged For Addition
                    String stagedBlobID = currentState.blobsForAdd.get(file);

                    //Get contents of file in working Directory
                    Blob fileBlob = Repository.makeBlob(file);
                    String fileBlobID = fileBlob.getBlobID();

                    //If they don't equal, add it to this cateogry
                    if (!stagedBlobID.equals(fileBlobID)) {
                        filesNotStagedForCommit = filesNotStagedForCommit + '\n' + file;
                    }
                }
            }
        }

        /**If there's a file tagged for Addition, but deleted in the working Directory,*/
        Set<String> setOfFilesTaggedForAddition = currentState.blobsForAdd.keySet();
        for (String fileTagged : setOfFilesTaggedForAddition) {
            if (!filesInCWD.contains(fileTagged)) {
                filesNotStagedForCommit = filesNotStagedForCommit + '\n' + fileTagged;
            }
        }

        /**If file is not staged for removal, but tracked in the
         * current commit and deleted in the working directory.*/
        Set<String> filesTrackedInCommit = currentCommit.tableOfBlobs.keySet();
        for (String fileTracked : filesTrackedInCommit) {
            if (!currentState.blobsForRemove.containsKey(fileTracked)
                    && !filesInCWD.contains(fileTracked)) {
                filesNotStagedForCommit = filesNotStagedForCommit + '\n' + fileTracked;
            }
        }
    /** KEPT GETTING ERRORS!!!
        String showNotStagedForCommitFiles = String.format(
                "=== Modifications Not Staged For Commit === %s \n",
                filesNotStagedForCommit);
     */

        String showNotStagedForCommitFiles = String.format(
                "=== Modifications Not Staged For Commit === \n");
        return showNotStagedForCommitFiles;
    }

    /**Files present in the working directory (CWD) but neither
     * STAGED FOR ADDITION nor TRACKED.*/
    public static String showUntrackedFiles() {
        List<String> filesInCWD = new ArrayList<String>();
        filesInCWD = plainFilenamesIn(Repository.CWD);

        Commit currentCommit = Commit.getParentCommit(Repository.headPointer);
        StagingArea currentState = StagingArea.fromFile();
        String filesUntracked = "";

        /**Looks at the files in working directory, and checks that it's not being tracked
            nor staged for addition. Saves it in String filesUntracked*/
        for (String file : filesInCWD) {
            if (file.contains(".txt") && !currentCommit.tableOfBlobs.containsKey(file)
                    && !currentState.blobsForAdd.containsKey(file)) {
                filesUntracked = filesUntracked + '\n' + file;
            }
        }

        /**Make the "Untracked Files" category in status
         *  -Have all of the filenames saved in the removal stage in String filesUntracked.
         *   Now put that in a String with the correct 'format'
         */
        String showUntrackedFiles = String.format(
                "=== Untracked Files === %s \n",
                filesUntracked);

        /**^- KEPT GETTING ERRORS!!!


        String showUntrackedFiles = String.format(
                "=== Untracked Files === \n");
         */

        return showUntrackedFiles;
    }

    public static void main(String[] args) {
        System.out.println(showBranches());
        System.out.println(showStagedFilesBranch());
        System.out.println(showRemovedFiles());
        System.out.println(showModsNotStagedForCommit());
        System.out.println(showUntrackedFiles());
    }
}
