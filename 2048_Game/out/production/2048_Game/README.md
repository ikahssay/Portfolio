# 2048 Game
In this project, I built the core logic of the “2048” game - a popular single-player computer game written by Gabriele Cirulli. This game is compromised on a 4x4 grid of squares, each of which can either be empty or contain a tile bearing an integer - a power of 2 greater than or equal to 2. 

<h> <b> How the game works: </b> </h>
In this game, the player can choose a direction via their arrow keys to move all tiles either north, south, east, or west on the board . All tiles slide in that direction until there is no empty space left in the direction of motion. 
When two adjacent tiles in the direction of motion have the same number, the tiles will <i>merge</i> into one tile containing double that initial number and will earn the player points. Each time two tiles merge to form a larger tile, the player earns the number of points on the new tile. The points will be tracked in the field called “Score” at the bottom of the screen and will be updated with each move. 

Here are the full rules for when merges occur:
    1. Two tiles of the same value merge into one tile containing double the initial number.
    2. A tile that is the result of a merge will not merge again on that tilt. For example, if we have [X, 2, 2, 4], where X represents an empty space, and we move the tiles to the left, we should end up with [4, 4, X, X], not [8, X, X, X]. This is because the leftmost 4 was already part of a merge so should not merge again.
    3. When three adjacent tiles in the direction of motion have the same number, then the leading two tiles in the direction of motion merge, and the trailing tile does not. For example, if we have [X, 2, 2, 2] and move tiles left, we should end up with [4, 2, X, X] not [2, 4, X, X].

    As a corollary of these rules, if there are four adjacent tiles with the same number in the direction of motion, they form two merged tiles. For example, if we have [4, 4, 4, 4], then if we move to the left, we end up with [8, 8, X, X]. This is because the leading two tiles will be merged as a result of rule 3, then the trailing two tiles will be merged, but because of rule 2 these merged tiles (8 in our example) will not merge themselves on that tilt.
    If the tilt did not change the board state, then no new tiles will be randomly generated. Otherwise, a single randomly generated tile will be added to the board on an empty square.
    
The game ends when the player has no available moves (when all tiles fill the board and no direction will merge tiles), or a move forms a tile containing the 2048 value. 

<h> <b> How The Game is Modeled: </h> </b>
The game exhibits two design patterns in common use: the <i> Model-View-Controller Pattern (MVC)</i>, and the <i> Observer Pattern </i>.

The MVC pattern divides the problem into three parts:

1. The <b> model </b> represents the subject matter being represented and acted upon – in this case incorporating the state of a board game and the rules by which it may be modified. The model resides in the <i> Model, Side, Board, </i> and <i> Tile </i> classes. The instance variables of Model fully determine what the state of the game is. 
2. A <b> view </b> of the model, which displays the game state to the user. The view resides in the GUI and BoardWidget classes.
3. A <b>controller</b> for the game, which translates user actions into operations on the model. The controller resides mainly in the Game class, although it also uses the GUI class to read keystrokes.

The second pattern utilized is the “Observer pattern”. Basically this means that the model doesn’t actually report changes to the view. Instead, the view registers itself as an observer of the Model object. 

