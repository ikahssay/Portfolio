package gitlet;

/** Driver class for Gitlet, a subset of the Git version-control system.
 *  @author TODO
 */
public class Main {

    /** Usage: java gitlet.Main ARGS, where ARGS contains
     *  <COMMAND> <OPERAND1> <OPERAND2> ... 
     */
    public static void main(String[] args) {
        if (args.length == 0) {
            Utils.message("Please enter a command.");
            System.exit(0);
        }

        Repository.setupPersistence();

        String firstArg = args[0];
        switch (firstArg) {
            case "init":
                validateNumArgs("init", args, 1);
                Repository.setUpInitialCommit();
                break;
            case "add":
                validateNumArgs("add", args, 2);
                Repository.addToStagingArea(args[1]);
                break;
            case "commit":
                validateNumArgs("commit", args, 2);

                /**Failure case: If commit has a blank message, print error message:
                 * "Please enter a commit message."*/
                if (args[1].equals("")) {
                    Utils.message("Please enter a commit message.");
                    System.exit(0);
                }

                Repository.makeCommit(args[1]);
                break;
            case "rm":
                validateNumArgs("rm", args, 2);
                Repository.unstageFile(args[1]);
                break;
            case "log":
                validateNumArgs("log", args, 1);
                Repository.goThroughCurrentBranch();
                break;
            case "global-log":
                validateNumArgs("global-log", args, 1);
                Repository.printAllCommits();
                break;
            case "find":
                validateNumArgs("find", args, 2);
                Repository.findCommitsWithSameMessage(args[1]);
                break;
            case "status":
                validateNumArgs("status", args, 1);
                Repository.showStatus();
                break;
            case "checkout":
                /**If length is 3, then it's checkout -- [file name]
                 * If length is 4, then it's checkout [commit ID] -- [file name]
                 * If length is 2, then it's checkout [branch name]*/
                if (args.length == 3) {
                    Repository.checkoutFileName(args[2], Repository.headPointer);
                } else if (args.length == 4) {
                    Repository.checkoutCommitIDAndFileName(args[1], args[3]);
                } else if (args.length == 2) {
                    Repository.checkoutBranch(args[1]);
                } else {
                    Utils.message("Incorrect operands.");
                    System.exit(0);
                }
                break;
            case "branch":
                validateNumArgs("branch", args, 2);
                Repository.createNewBranch(args[1]);
                break;
            default:
                Utils.message("No command with that name exists.");
                System.exit(0);

        }
        return;
    }

    public static void validateNumArgs(String cmd, String[] args, int n) {
        if (args.length != n) {
            Utils.message("Incorrect operands.");
            System.exit(0);
        }
    }
}
