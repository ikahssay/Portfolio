package byow.lab12;
import org.junit.Test;
import static org.junit.Assert.*;

import byow.TileEngine.TERenderer;
import byow.TileEngine.TETile;
import byow.TileEngine.Tileset;

import java.util.Random;

/**
 * Draws a world consisting of hexagonal regions.
 */
public class HexWorld {
    //Width and length of the world
    private static final int WIDTH = 50;
    private static final int HEIGHT = 50;

    //Random color and size of hexagon
    private static final long SEED = 287313232;
    private static final Random RANDOM = new Random(SEED);

    public static void initializeTiles(TETile[][] tiles) {
        //Initialize world by filling it with "nothing" tiles
        for (int x = 0; x < WIDTH; x += 1) {
            for (int y = 0; y < HEIGHT; y += 1) {
                tiles[x][y] = Tileset.NOTHING;
            }
        }
    }

    public static void addHexagon(TETile[][] tiles, int size, int widthStart, int heightStart,
                                  TETile randomTerrain) {
        //Creates the lower-portion of a Hexagon
        lowerHexagon(widthStart, heightStart, size, size, tiles, randomTerrain);

        //Creates the upper-portion of a Hexagon
        upperHexagon(widthStart - size + 1, heightStart + size,
                        size + (2 * (size-1)), size, tiles, randomTerrain);
    }

    private static void lowerHexagon(int widthStart, int heightStart, int size, int counter, TETile[][] tiles, TETile randomTerrain) {
        if (counter == 0) {
            ;
        } else {

            for (int x = widthStart; x < widthStart + size; x++) {
                tiles[x][heightStart] = randomTerrain;
            }

            lowerHexagon(widthStart - 1, heightStart + 1, size + 2, counter - 1, tiles, randomTerrain);
        }
    }

    private static void upperHexagon(int widthStart, int heightStart, int size, int counter, TETile[][] tiles, TETile randomTerrain) {
        if (counter == 0) {
            ;
        } else {

            for (int x = widthStart; x < widthStart + size; x++) {
                tiles[x][heightStart] = randomTerrain;
            }

            upperHexagon(widthStart + 1, heightStart + 1, size - 2, counter - 1, tiles, randomTerrain);
        }
    }

    private static TETile randomTile() {
        int tileNum = RANDOM.nextInt(6);
        switch (tileNum) {
            case 0: return Tileset.GRASS;
            case 1: return Tileset.FLOWER;
            case 2: return Tileset.SAND;
            case 3: return Tileset.TREE;
            case 4: return Tileset.MOUNTAIN;
            case 5: return Tileset.WATER;
            default: return Tileset.NOTHING;
        }
    }

    public static int buildHexagonsUpward(TETile[][] tiles, int size, int maxHeight, int capOff,
                                          int startingWidth, int startingHeight, TETile randomTerrain) {
        int y = 0;
        for (y = startingHeight; y < HEIGHT - (size * maxHeight) - capOff; y = y + (size* 2)) {
            //Gets rid of out of bound errors
            if (y + (size* 2) <= HEIGHT) {
                addHexagon(tiles, size, startingWidth,
                        y, randomTerrain);
            }
            randomTerrain = randomTile();
        }

        return y - HEIGHT;
    }

    public static void main(String[] args) {
        // initialize the tile rendering engine with a window of size WIDTH x HEIGHT
        TERenderer ter = new TERenderer();
        ter.initialize(WIDTH, HEIGHT);

        // initialize the 'world' with empty tiles
        TETile[][] world = new TETile[WIDTH][HEIGHT];
        initializeTiles(world);

        //Make the hexagon a random size and terrain
        int randomSize =  RANDOM.nextInt(7);
        System.out.println(randomSize);
        TETile randomTerrain = randomTile();

        //Add a hexagon that begins in the bottom middle of the world
        int startingWidth = (WIDTH / 2) - (randomSize/2);
        int startingHeight = 0;


        /** Builds hexagons all around the world*/

        //Builds the center part
        int maxHeight = 0;
        int capOff = buildHexagonsUpward(world, randomSize, maxHeight, 0, startingWidth,
                startingHeight, randomTerrain);

        //Builds the left-half of world
        int y = randomSize;
        for (int x = startingWidth - (2 * randomSize) + 1; x > 0 + randomSize; x = x - (2 * randomSize) + 1) {
            maxHeight = maxHeight + 1;

            buildHexagonsUpward(world, randomSize, maxHeight, capOff, x,
                    y, randomTerrain);

            y = y + randomSize;
            randomTerrain = randomTile();
        }

        //Builds the right-half of the world
        y = randomSize;
        maxHeight = 0;
        for (int x = startingWidth + (2 * randomSize) - 1; x < WIDTH - (2 * randomSize) + 1; x = x + (2 * randomSize) - 1) {
            maxHeight = maxHeight + 1;

            buildHexagonsUpward(world, randomSize, maxHeight, capOff, x,
                    y, randomTerrain);

            y = y + randomSize;
            randomTerrain = randomTile();
        }


        ter.renderFrame(world);
    }
}
