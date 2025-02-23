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
}

class MyAppState extends ChangeNotifier {
  Coord? selectedCoord;
  String moveDestinations = '';
  bool isWhiteTurn = true;
  bool isBlackAi = true;
  bool isWhiteAi = false;
  String gameStatus = statusWhiteMove;
  bool isGameOver = false;
  String indicatedCoords = '';
  ButtonState undoButtonModel = ButtonState(false, false);
  
  int aiDropdownDepth = 4;
  List<List<String>> boardModel = parseBoardString(startingBoard);
  List<Square> boardView = getInitialBoardView(parseBoardString(startingBoard));
  PreferencesManager savedData = PreferencesManager();
  Future<void> resetGame() async {
    await savedData.clearBoards();
    boardModel = parseBoardString(startingBoard);
    moveDestinations = '';
    isWhiteTurn = true;
    isBlackAi = true;
    isWhiteAi = false;
    gameStatus = statusWhiteMove;
    isGameOver = false;
    indicatedCoords = '';
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
    } else if ((isWhiteTurn && isWhiteAi) || (!isWhiteTurn && isBlackAi)){
      log('It is an AIs turn');
    } else {
      selectButton(c, false);
    }
  } 
  void selectButton(Coord c, bool saveBoardOnMove){
    String piece = boardModel[c.i][c.j];
    bool isNotifyAi = false;
    String boardString = getBoardString();
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
          gameStatus = resultList[1];
          if (gameStatus == statusCheckMate || gameStatus == statusStaleMate){
            isGameOver = true;
          }
          boardModel = parseBoardString(newBoardStr);
          isWhiteTurn = !isWhiteTurn;
          isNotifyAi = true;
          indicatedCoords = '$selectedCoord|$c';
          if (saveBoardOnMove && boardString != startingBoard){
            // do not save the starting board. The user can reset the board if they want. 
            savedData.addBoard(newBoardStr);
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
    if (!isGameOver && ((isWhiteTurn && isWhiteAi) || (!isWhiteTurn && isBlackAi))) {
      //it is the ai's turn
      //setBoardString(); // todo: get rid of this function. either implement getBoardString or store it whenever it changes. 
      String aiMove = await getAiChosenMove(getBoardString(), isWhiteTurn, 'simple', aiDropdownDepth);
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
    log(getBoardString());
  }
  void printSavedData() {
    log(savedData.toString());
  }
  void saveBoard() {
    String boardStr = getBoardString();
    savedData.addBoard(boardStr);
  }
  void setIsUndoVisible() {
    // do not show undo if two humans are playing or two ai's are playing
    bool newVal = (isWhiteAi || isBlackAi) && (!isWhiteAi || !isBlackAi);
    if (newVal != undoButtonModel.isVisible){
      undoButtonModel.isVisible = newVal;
      //notifyListeners();
    }
  }
  void setIsUndoEnabled() {
    bool newEnableUndo = ((isWhiteTurn && !isWhiteAi) || (!isWhiteTurn && !isBlackAi)) && savedData.isUndoPossible;
    if (newEnableUndo != undoButtonModel.isEnabled){
      undoButtonModel.isEnabled = newEnableUndo;
      //notifyListeners();
    }
  }
  void undo() {
    String lastBoard = savedData.popBoard();
    boardModel = parseBoardString(lastBoard);
    setIsUndoEnabled();
    notifyListeners();
  }
  Color getColor(Coord c){
    bool isLightSquare = (c.i+c.j)%2==0;
    if (c==selectedCoord) {
      return selectedColor;
    } else if (indicatedCoords.contains(c.toString())){
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
  String getBoardString() {
    String boardString = '';
    for (int i=0; i<boardModel.length; i++) {
      for (int j=0; j<boardModel[i].length; j++) {
        boardString += boardModel[i][j];
      }
    }
    return boardString;
  }
  void setAiDepth(int d) {
    aiDropdownDepth = d;
    notifyListeners();
  }
}
