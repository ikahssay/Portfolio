package byow.lab13;

import byow.Core.RandomUtils;
import edu.princeton.cs.introcs.StdDraw;

import java.awt.Color;
import java.awt.Font;
import java.util.Random;

public class MemoryGame {
    /** The width of the window of this game. */
    private int width;
    /** The height of the window of this game. */
    private int height;
    /** The current round the user is on. */
    private int round;
    /** The Random object used to randomly generate Strings. */
    private Random rand;
    /** Whether or not the game is over. */
    private boolean gameOver;
    /** Whether or not it is the player's turn. Used in the last section of the
     * spec, 'Helpful UI'. */
    private boolean playerTurn;
    /** The characters we generate random Strings from. */
    private static final char[] CHARACTERS = "abcdefghijklmnopqrstuvwxyz".toCharArray();
    /** Encouraging phrases. Used in the last section of the spec, 'Helpful UI'. */
    private static final String[] ENCOURAGEMENT = {"You can do this!", "I believe in you!",
            "You got this!", "You're a star!", "Go Bears!",
            "Too easy for you!", "Wow, so impressive!"};

    public static void main(String[] args) {
        if (args.length < 1) {
            System.out.println("Please enter a seed");
            return;
        }

        long seed = Long.parseLong(args[0]);
        MemoryGame game = new MemoryGame(40, 40, seed);

        game.startGame();
    }

    public MemoryGame(int width, int height, long seed) {
        /* Sets up StdDraw so that it has a width by height grid of 16 by 16 squares as its canvas
         * Also sets up the scale so the top left is (0,0) and the bottom right is (width, height)
         */
        this.width = width;
        this.height = height;
        StdDraw.setCanvasSize(this.width * 16, this.height * 16);
        Font font = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(font);
        StdDraw.setXscale(0, this.width);
        StdDraw.setYscale(0, this.height);
        StdDraw.clear(Color.BLACK);
        StdDraw.enableDoubleBuffering();

        //TODO: Initialize random number generator
        rand = new Random(seed);
    }

    public String generateRandomString(int n) {
        //TODO: Generate random string of letters of length n
        String targetString = "";

        for (int i = 0; i < n; i++) {
            int randomIndex = rand.nextInt(CHARACTERS.length);
            char randomChar = CHARACTERS[randomIndex];
            targetString = targetString + randomChar;
        }

        return targetString;
    }

    public void drawFrame(String s) {
        //TODO: Take the string and display it in the center of the screen ONE LETTER AT A TIME

        /**Clear entire screen*/
        StdDraw.clear(StdDraw.BLACK);   //Clears canvas to white by default. Want it black.
        StdDraw.setPenColor(StdDraw.WHITE); //Pen color is black by default. Want it to be white.

        /**Sets the font to be large and Bold*/
        Font f = new Font("Monaco", Font.BOLD, 30);
        StdDraw.setFont(f);

        /**Draws the input String so that it is CENTERED on the canvas*/
        StdDraw.text(width/2, height/2, s);

        //TODO: If game is not over, display relevant game information at the top of the screen
        if (!gameOver) {
            Font smallerF = new Font("Monaco", Font.BOLD, 20);
            StdDraw.setFont(smallerF);

            /**Display the current round on the top left corner*/
            StdDraw.textLeft(0, height - 1, "Round: " + round);

            /**Display the current task on the top Center*/
            if (playerTurn) {
                StdDraw.text(width / 2, height - 1, "Type!");
            } else {
                StdDraw.text(width / 2, height - 1, "Watch!");
            }

            /**Display an encouraging phrase on the top right corner*/
            String randomEncouragement = ENCOURAGEMENT[RandomUtils.uniform(rand, ENCOURAGEMENT.length)];
            StdDraw.textRight(width, height - 1, randomEncouragement);

            /**Draw a line/"bar" to separate game info from everything else
             * -Start line from beginning point (width = 0, height = gameInfo height - 1)
             * -And draw that straight horizontal line all the way till the end point
             *      (width = width, height = gameInfo height - 1)*/
            StdDraw.line(0, this.height - 2, this.width, this.height - 2);
        }

        StdDraw.show();
    }

    public void flashSequence(String letters) {
        //TODO: Display each character in letters, making sure to blank the screen between letters
        for (int i = 0; i < letters.length(); i++) {
            /**Blanks the screen for 0.5 seconds -> 500 ms*/
            drawFrame(""); //Pauses for t milliseconds.
            StdDraw.pause(500);

            /**Draws the input String so that it is CENTERED on the canvas for 1 second*/
            /**NOTE: toString ONLY WORKS ON Object Character, and not primitive char */
            Character c = letters.charAt(i);
            String letter = c.toString();
            drawFrame(letter);
            StdDraw.pause(1000); //Pauses for 1 sec.
        }
    }

    public String solicitNCharsInput(int n) {
        //TODO: Reads N keystrokes using STDDraw and returns the String
        // corresponding to those keystrokes.
        String inputString = "";

        /**NOTE TO SELF: IF YOU PUT StdDraw.hasNextKeyTyped() IN THE WHILE CONDITION
         * GAME WILL NEVER WORK!!! (loop never starts) */
        while (n > 0) {
            if (StdDraw.hasNextKeyTyped()) {
                /**Gets the letter user typed in and saves it in String inputString*/
                char letterTyped = StdDraw.nextKeyTyped();
                inputString = inputString + letterTyped;

                /**String built up so far should appear centered on the screen as keys
                 * are being typed by the user*/
                drawFrame(inputString);

                n--;
            }
        }
        return inputString;
    }

    public void startGame() {
        /**Once you start the game at round 1, check to see if the player got the targetString correct
         *  -> If they got it correct, REPEAT beginning from displaying the "Round: "
         *      message, generate random string, etc.
         *  -> If they got it wrong, end the game and display the message "Game Over!
         *      You made it to round: " + round number*/

        //TODO: Set any relevant variables before the game starts
        /**Start the game at round 1*/
        round = 0;
        String inputString = "";
        String targetString = "";
        playerTurn = false;
        gameOver = false;

        //TODO: Establish Engine loop
        while(inputString.equals(targetString)) {
            //Increments the game if user passed the round.
            round++;

            /**Display the message "Round: " + round number in the center of the screen*/
            drawFrame("Round: " + round);
            StdDraw.pause(1000);

            /**Generate a random string of length equal to the current round number*/
            targetString = generateRandomString(round);
            /**Display the random string one letter at a time*/
            flashSequence(targetString);
            drawFrame("");

            /**Wait for player to type in a string the same length as the target string*/
            playerTurn = true;
            inputString = solicitNCharsInput(round);
            StdDraw.pause(1000);
        }

        playerTurn = false;
        gameOver = true;
        drawFrame("Game Over! You made it to round: " + round);
    }
}
