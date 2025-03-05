
import 'dart:developer';

import 'package:play/constants.dart';
import 'package:play/square.dart';

class BoardState {
  // Objects and information related to the board positions and colors, but  
  // not including whether a piece is selected. 
  // not including theme data. 
  List<List<String>> boardModel = parseBoardString(startingBoard);
  List<Square> boardView = getInitialBoardView(parseBoardString(startingBoard));
  String indicatedCoords = '';
  String gameStatus = statusWhiteMove;
  bool isWhiteTurn = true;

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
      if (gameStatus == statusWhiteMove){
        isWhiteTurn = true;
      } else if (gameStatus == statusBlackMove) {
        isWhiteTurn = false;
      } else {
        log("Loading end of previous game.");
      }
    }
  }
  void resetGame() {
    boardModel = parseBoardString(startingBoard);
    indicatedCoords = '';
    log("Reseting game: ${getBoardString()}");
    gameStatus = statusWhiteMove;
    isWhiteTurn = true;
  }
  bool isGameOver() {
    return gameStatus == statusCheckMate || gameStatus == statusStaleMate;
  }
}