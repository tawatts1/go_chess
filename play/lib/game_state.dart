import 'dart:developer';

import 'package:flutter/material.dart';
import 'constants.dart';
import 'coord.dart';
import 'ffi_funcs.dart';
import 'preferences_manager.dart';
import 'square.dart';
import 'models/button_state.dart';
import 'models/board_state.dart';
import 'models/player_state.dart';
import 'models/theme_state.dart';

enum PlayStatus {play, pause, undefined}

class MyAppState extends ChangeNotifier {
  Coord? selectedCoord;
  String moveDestinations = '';
  PlayerState players = PlayerState();
  bool isGameOver = false;
  ButtonState undoButtonModel = ButtonState(false, false);
  ButtonState playButtonModel = ButtonState(false, true);
  ButtonState pauseButtonModel = ButtonState(false, true);
  PlayStatus playPauseStatus = PlayStatus.undefined;
  BoardState board = BoardState();
  ThemeState theme = ThemeState();
  PreferencesManager savedData = PreferencesManager();
  @override 
  String toString() {
    return "$board";
  }
  void loadBoardStateFromString(String stateStr) {
    board.loadFromString(stateStr);
  }
  void loadPlayerStateFromString(String stateStr) {
    players.loadFromString(stateStr);
  }
  void loadThemeStateFromString(String stateStr) {
    theme.loadFromString(stateStr);
  }
  Future<void> resetGame() async {
    await savedData.clearBoardStates();
    moveDestinations = '';
    isGameOver = false;
    board.resetGame();
    setUndoState();
    clearSelection();
    notifyListeners();
    if (players.isWhiteAi && !players.isBlackAi){
      notifyAi();
    }
  }
  
  void clearSelection() {
    selectedCoord = null;
    moveDestinations = '';
  }
  void humanSelectButton(Coord c){
    //functions that humans have to use to select the buttons
    if (isGameOver) {
      log('game is over');
    } else if ((board.isWhiteTurn && players.isWhiteAi) || (!board.isWhiteTurn && players.isBlackAi)){
      log('It is an AIs turn');
    } else {
      selectButton(c);
    }
  } 
  void selectButton(Coord c) async {
    String piece = board.boardModel[c.i][c.j];
    bool isNotifyAi = false;
    String boardString = board.getBoardString();
    Future<void>? dataSaved;
    if (selectedCoord == null) {
      // no selection has been made
      if (piece == Space) {
        log('try clicking on a piece!');
      } else if ((board.isWhiteTurn && (whiteMap[piece] ?? false)) ||
                 (!board.isWhiteTurn && !(whiteMap[piece] ?? true))) {
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
          board.isWhiteTurn = !board.isWhiteTurn;
          isNotifyAi = true;
          board.indicatedCoords = '$selectedCoord|$c';
          if (boardString != startingBoard){
            // do not save the starting board. The user can reset the board if they want.
            dataSaved = savedData.addBoardSnapshot(toString());
          }
        } 
      } else {
        log('not one of the legal moves. Clearing selection');
      }
      clearSelection();
    }
    if (dataSaved != null) {
      await dataSaved;
    }
    setUndoState(); 
    notifyListeners();
    if (isNotifyAi) {
      notifyAi();
    }
  }
  Future<void> notifyAi() async {
    if (!isGameOver && 
        ((board.isWhiteTurn && players.isWhiteAi) || (!board.isWhiteTurn && players.isBlackAi)) &&
        playPauseStatus != PlayStatus.pause) {
      //it is the ai's turn
      String aiMove = await getAiChosenMove(board.getBoardString(), board.isWhiteTurn, 'simple', players.aiDropdownDepth);
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
        selectButton(Coord(i1, j1));
        Coord click2 = await Future.delayed(const Duration(milliseconds: 700),  () => Coord(i2,j2));
        selectButton(click2);
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
  void setUndoState() {
    // do not show undo if two humans are playing or two ai's are playing
    bool newVisibleUndo = !players.isBothAi() && !players.isNeitherAi();
    if (newVisibleUndo != undoButtonModel.isVisible){
      undoButtonModel.isVisible = newVisibleUndo;
    }
    bool newEnableUndo = ((board.isWhiteTurn && !players.isWhiteAi) || (!board.isWhiteTurn && !players.isBlackAi)) && savedData.isUndoPossible;
    if (newEnableUndo != undoButtonModel.isEnabled){
      undoButtonModel.isEnabled = newEnableUndo;
    }
  }
  void undo() async {
    String lastAppState = await savedData.popBoard(2);
    loadBoardStateFromString(lastAppState);
    setUndoState();
    notifyListeners();
  }
  void setPlayPauseButtonState(){
    if (players.isBothAi()) {
      if (playPauseStatus == PlayStatus.play){
        playButtonModel.isVisible = false;
        pauseButtonModel.isVisible = true;
      } else {
        if (playPauseStatus == PlayStatus.undefined){
          playPauseStatus = PlayStatus.pause;
        }
        playButtonModel.isVisible = true;
        pauseButtonModel.isVisible = false;
      }
    } else {
      playButtonModel.isVisible = false;
      pauseButtonModel.isVisible = false;
      playPauseStatus = PlayStatus.undefined;
    }
  }
  void playPause() {
    if (playPauseStatus == PlayStatus.play){
      playPauseStatus = PlayStatus.pause;
      setPlayPauseButtonState();
      notifyListeners();
      log("Pause was just pushed. ");
    } else if (playPauseStatus == PlayStatus.pause){
      playPauseStatus = PlayStatus.play;
      setPlayPauseButtonState();
      notifyAi();
      notifyListeners();
      log("Play was just pushed");
    } else {
      log("Error: Play/Pause was pushed when play pause status was not defined. ");
    }
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
      return theme.getSquareColor(true);
    } else {
      return theme.getSquareColor(false);
    }
  }
  double getRadius(Coord c){
    if (moveDestinations.contains(c.toString())) {
      return 15;
    } else if (selectedCoord != null && selectedCoord == c) {
      return 15;
    } else {
      return 0;
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
    setUndoState();
    setPlayPauseButtonState();
    notifyListeners(); // necessary because the undo button may change visibility. 
    if (players.isAtLeastOneAi() && !players.isBothAi()){
      notifyAi();
    }
    
  }
  bool shouldBoardBeFlipped() {
    //Determines if board should show black at the bottom or white. 
    if (players.isBothAi()){
      //two ai playing
      return false;
    } else if (players.isNeitherAi()){
      //two humans playing
      if (board.isWhiteTurn){
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
  doDuringStartup() {
    setUndoState();
    setPlayPauseButtonState();
    if (!players.isBothAi()){
      if ((board.isWhiteTurn && players.isWhiteAi) || (!board.isWhiteTurn && players.isBlackAi)){
        notifyAi();
      }
    }
  }
  void setTheme(bool newVal){
    if (theme.isDarkTheme != newVal){
      log("Setting theme...");
      theme.isDarkTheme = newVal;
      savedData.theme = theme.toString();
      savedData.saveTheme();
      notifyListeners();
    } else {
      log("Error: changing theme to the same theme");
    }
  }
}
