import 'dart:developer';

import 'package:flutter/material.dart';
import 'constants.dart';
import 'coord.dart';
import 'ffi_funcs.dart';
import 'preferences_manager.dart';
import 'square.dart';

//contains whether the button is visible in the ui, and if it is enabled. 
class ButtonState {
  bool isVisible = false;
  bool isEnabled = false;
  ButtonState(this.isVisible, this.isEnabled);
  @override
  bool operator ==(Object other) =>
    other is ButtonState &&
    other.runtimeType == runtimeType &&
    other.isVisible == isVisible && other.isEnabled == isEnabled;
  @override
  int get hashCode => Object.hash(isVisible, isEnabled);
  // @override
  // String toString() {
  //   return "$isVisible,$isEnabled";
  // }
  // void loadFromString(String stateStr) {
  //   List<String> buttonState = stateStr.split(",");
  //   if (buttonState.length != 2){
  //     log("Attempted to load button state from a bad string: $stateStr");
  //   } else {
  //     isVisible = buttonState[0].toLowerCase() == "true";
  //     isEnabled = buttonState[1].toLowerCase() == "true";
  //   }
  // }
}

class BoardState {
  // Objects and information related to the board positions and colors, but  
  // not including whether a piece is selected. 
  // not including theme data. 
  List<List<String>> boardModel = parseBoardString(startingBoard);
  List<Square> boardView = getInitialBoardView(parseBoardString(startingBoard));
  String indicatedCoords = '';
  String gameStatus = statusWhiteMove;

  String getBoardString() {
    String boardString = '';
    for (int i=0; i<boardModel.length; i++) {
      for (int j=0; j<boardModel[i].length; j++) {
        boardString += boardModel[i][j];
      }
    }
    return boardString;
  }
  @override
  String toString() {
    return "${getBoardString()}-$indicatedCoords-$gameStatus";
  }
  void loadFromString(String stateStr) {
    List<String> boardState = stateStr.split("-");
    if (boardState.length != 3) {
      log("Attempted to load board state from a bad string: $stateStr");
    } else {
      boardModel = parseBoardString(boardState[0]);
      boardView = getInitialBoardView(parseBoardString(boardState[0]));
      indicatedCoords = boardState[1];
      gameStatus = boardState[2];
    }
  }
  void resetGame() {
    boardModel = parseBoardString(startingBoard);
    //boardView = getInitialBoardView(parseBoardString(startingBoard));
    indicatedCoords = '';
    log("Reseting game: ${getBoardString()}");
    gameStatus = statusWhiteMove;
  }
}

const humanName = "Human";
const aiName = "Ai";
class PlayerState {
  // Info about who the players are. 
  bool isBlackAi = true;
  bool isWhiteAi = false;
  int aiDropdownDepth = 1;
  @override
  String toString() {
    return "$isBlackAi,$isWhiteAi,$aiDropdownDepth";
  }
  void loadFromString(String stateStr){
    List<String> playerState = stateStr.split(",");
    if (playerState.length != 3){
      log("Tried to load player state from a bad string: $stateStr");
    } else {
      try {
        isBlackAi = playerState[0].toLowerCase() == "true";
        isWhiteAi = playerState[1].toLowerCase() == "true";
        aiDropdownDepth = int.parse(playerState[2]);
      } catch (e) {
        log("Failed to parse player state: $stateStr");
      }
    }
  }
  String getPlayerName(bool isWhite) {
    if (isWhite){
      return isWhiteAi ? aiName : humanName;
    } else {
      return isBlackAi ? aiName : humanName;
    }
  }
}

class MyAppState extends ChangeNotifier {
  Coord? selectedCoord;
  String moveDestinations = '';
  bool isWhiteTurn = true;
  PlayerState players = PlayerState();
  bool isGameOver = false;
  ButtonState undoButtonModel = ButtonState(false, false);
  BoardState board = BoardState();
  PreferencesManager savedData = PreferencesManager();
  @override 
  String toString() {
    return "$board#placeholder";
  }
  void loadBoardStateFromString(String stateStr) {
    List<String> appState = stateStr.split("#");
    if (appState.length != 2) {
      log("Tried to parse bad board state: $stateStr");
    } else {
      board.loadFromString(appState[0]);
      //undoButtonModel.loadFromString(appState[1]);
    }
  }
  void loadPlayerStateFromString(String stateStr) {
    players.loadFromString(stateStr);
  }
  Future<void> resetGame() async {
    await savedData.clearBoardStates();
    moveDestinations = '';
    isWhiteTurn = true;
    isGameOver = false;
    board.resetGame();
    setIsUndoEnabled();
    clearSelection();
    notifyListeners();
  }
  
  void clearSelection() {
    selectedCoord = null;
    moveDestinations = '';
  }
  void humanSelectButton(Coord c){
    //functions that humans have to use to select the buttons
    if (isGameOver) {
      log('game is over');
    } else if ((isWhiteTurn && players.isWhiteAi) || (!isWhiteTurn && players.isBlackAi)){
      log('It is an AIs turn');
    } else {
      selectButton(c, false);
    }
  } 
  void selectButton(Coord c, bool saveBoardOnMove){
    String piece = board.boardModel[c.i][c.j];
    bool isNotifyAi = false;
    String boardString = board.getBoardString();
    if (selectedCoord == null) {
      // no selection has been made
      if (piece == Space) {
        log('try clicking on a piece!');
      } else if ((isWhiteTurn && (whiteMap[piece] ?? false)) ||
                 (!isWhiteTurn && !(whiteMap[piece] ?? true))) {
        //mark down the selection and populate the move destinations
        selectedCoord = c;
        moveDestinations = getMoves(boardString, c);
        log(moveDestinations);
      } else {
        log('not that colors turn');
      }
    } else if (selectedCoord != null) {
      // a move has already been selected. 
      if (moveDestinations.contains(c.toString())){
        log('legal move');
        String boardResult = getBoardAfterMove(boardString, selectedCoord!, c);
        List<String> resultList = boardResult.split(',');
        if (resultList.length == 2) {
          String newBoardStr = resultList[0];
          board.gameStatus = resultList[1];
          if (board.gameStatus == statusCheckMate || board.gameStatus == statusStaleMate){
            isGameOver = true;
          }
          board.boardModel = parseBoardString(newBoardStr);
          isWhiteTurn = !isWhiteTurn;
          isNotifyAi = true;
          board.indicatedCoords = '$selectedCoord|$c';
          if (saveBoardOnMove && boardString != startingBoard){
            // do not save the starting board. The user can reset the board if they want.
            savedData.addBoardSnapshot(toString());
          }
        } 
      } else {
        log('not one of the legal moves. Clearing selection');
      }
      clearSelection();
    }
    setIsUndoEnabled(); 
    setIsUndoVisible();
    notifyListeners();
    if (isNotifyAi) {
      notifyAi();
    }
  }
  Future<void> notifyAi() async {
    if (!isGameOver && ((isWhiteTurn && players.isWhiteAi) || (!isWhiteTurn && players.isBlackAi))) {
      //it is the ai's turn
      String aiMove = await getAiChosenMove(board.getBoardString(), isWhiteTurn, 'simple', players.aiDropdownDepth);
      parseAndDoAiMove(aiMove);
    }
  }
  void parseAndDoAiMove(String moveStr) async {
    List<String> indexList = moveStr.split(',');
    if (indexList.length == 4) {
      try {
        int i1 = int.parse(indexList[0]);
        int j1 = int.parse(indexList[1]);
        int i2 = int.parse(indexList[2]);
        int j2 = int.parse(indexList[3]);
        selectButton(Coord(i1, j1), false);
        Coord click2 = await Future.delayed(const Duration(milliseconds: 700),  () => Coord(i2,j2));
        // Save the board after the ai makes a move, and the undo button is visible
        bool saveNewBoard = undoButtonModel.isVisible;
        selectButton(click2, saveNewBoard);
      } catch(ex) {
        log("failed to parse ai move");
      }
    }
  }
 
  void printBoard() {
    log(board.getBoardString());
  }
  void printSavedData() {
    log(savedData.toString());
  }
  void saveBoard() {
    savedData.addBoardSnapshot(toString());
  }
  void setIsUndoVisible() {
    // do not show undo if two humans are playing or two ai's are playing
    bool newVal = (players.isWhiteAi || players.isBlackAi) && (!players.isWhiteAi || !players.isBlackAi);
    if (newVal != undoButtonModel.isVisible){
      undoButtonModel.isVisible = newVal;
      //notifyListeners();
    }
  }
  void setIsUndoEnabled() {
    bool newEnableUndo = ((isWhiteTurn && !players.isWhiteAi) || (!isWhiteTurn && !players.isBlackAi)) && savedData.isUndoPossible;
    if (newEnableUndo != undoButtonModel.isEnabled){
      undoButtonModel.isEnabled = newEnableUndo;
      //notifyListeners();
    }
  }
  void undo() async {
    String lastAppState = await savedData.popBoard();
    loadBoardStateFromString(lastAppState);
    setIsUndoEnabled();
    notifyListeners();
  }
  Color getColor(Coord c){
    bool isLightSquare = (c.i+c.j)%2==0;
    if (c==selectedCoord) {
      return selectedColor;
    } else if (board.indicatedCoords.contains(c.toString())){
      if (isLightSquare){
        return greyedWhite;
      } else {
        return greyedBlack;
      }
    } else if (isLightSquare){
      return white;
    } else {
      return black;
    }
  }
  double getRadius(Coord c){
    if (moveDestinations.contains(c.toString())) {
      return 15;
    } else if (selectedCoord != null && selectedCoord == c) {
      return 15;
    } else {
      return 1;
    }
  }
  void setAiDepth(int d) {
    players.aiDropdownDepth = d;
    savedData.players = players.toString();
    savedData.savePlayers();
    notifyListeners();
  }
  void setPlayer(String playerName, bool isWhite) {
    if (playerName == "Human") {
      if (isWhite) {
        players.isWhiteAi = false;
      } else {
        players.isBlackAi = false;
      }
    } else {
      if (isWhite) {
        players.isWhiteAi = true;
      } else {
        players.isBlackAi = true;
      }
    }
    savedData.players = players.toString();
    savedData.savePlayers();
    setIsUndoEnabled();
    setIsUndoVisible();
    notifyListeners(); // necessary because the undo button may change visibility. 
    notifyAi();
  }
  bool shouldBoardBeFlipped() {
    //Determines if board should show black at the bottom or white. 
    if (players.isWhiteAi && players.isBlackAi){
      //two ai playing
      return false;
    } else if (!players.isWhiteAi && !players.isBlackAi){
      //two humans playing
      if (isWhiteTurn){
        return false;
      } else {
        return true;
      }
    } else if (players.isWhiteAi) {
      return true;
    } else {
      return false;
    }
  }
}
