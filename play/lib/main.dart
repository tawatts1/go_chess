import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'constants.dart';
import 'coord.dart';
import 'ffi_funcs.dart';

List<List<String>> parseBoardString(String boardStr) {
  final List<List<String>> board = [[],[],[],[],[],[],[],[],];
  if (boardStr.length == BoardHeight * BoardWidth) {
    for (int k=0; k<boardStr.length; k++){
      board[k ~/BoardWidth].add(boardStr[k]);
    }
  }
  return board;
}

List<Square> getInitialBoardView(List<List<String>> boardModel) {
  final List<Square> out = [];
  for (int i=0; i<boardModel.length; i++) {
    for (int j=0; j<boardModel[i].length; j++) {
      Color color = black;
      if ((i+j)%2==0) {
        color = white;
      }
      final Coord c = Coord(i,j);
      final SquareModel sq = SquareModel(boardModel[i][j], color, 1);
      final Square newSq = Square(c: c, sq: sq, key: ValueKey(Object.hash(c, sq)),);
      out.add(newSq);
    }
  }
  return out;
}

class SquareModel {
  final String pieceCode;
  final Color color;
  final double radius;
  SquareModel(this.pieceCode, this.color, this.radius);
  @override
  bool operator ==(Object other) =>
    other is SquareModel &&
    other.runtimeType == runtimeType &&
    other.pieceCode == pieceCode && other.color == color && other.radius == radius;
  @override
  int get hashCode => Object.hash(pieceCode, color, radius);
}

class PreferencesManager {
  late final SharedPreferencesAsync prefs;
  late List<String> boardStrings;
  bool isUndoPossible = false;
  bool isLoaded = false;
  final int maxBoards = 2;

  PreferencesManager() {
    prefs = SharedPreferencesAsync();
  }
  @override 
  String toString() {
    return 'PreferencesManager: $boardStrings, $isUndoPossible, $isLoaded, $maxBoards';
  }
  loadBoards() async {
    List<String>? loaded = await prefs.getStringList('b');
    boardStrings = loaded ?? [];
    isUndoPossible = boardStrings.length > 1;
    isLoaded = true;
  }
  saveBoards() async {
    prefs.setStringList('b',boardStrings);
  }
  addBoard(String boardStr) async {
    boardStrings.add(boardStr);
    if (boardStrings.length > maxBoards) {
      boardStrings = boardStrings.sublist(1);
    }
    isUndoPossible = boardStrings.length > 1;
    saveBoards();
  }
  String popBoard() {
    if (isUndoPossible){
      boardStrings.removeLast();
      String out = boardStrings.last;
      isUndoPossible = boardStrings.length > 1;
      saveBoards();
      return out;
    } else {
      log('tried to pop board history when there was none. ');
      return '';
    }
  }
  Future<String?> getLastSavedBoard() async {
    String? out;
    if (!isLoaded) {
      await loadBoards();
    }
    if (boardStrings.isNotEmpty){
      out =  boardStrings.last;
    } else {
      out =  null;
    }
    return out;
  }
  Future<void> clearBoards() async {
    boardStrings.clear();
    boardStrings = [];
    isUndoPossible = false;
    await saveBoards();
  }
}

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (context) => MyAppState(),
      child: MaterialApp(
        title: 'Namer App',
        theme: ThemeData(
          useMaterial3: true,
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepOrange),
        ),
        home: const MyHomePage(),
      ),
    );
  }
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
        bool saveNewBoard = isUndoVisible();
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
  bool isUndoVisible() {
    // do not show undo if two humans are playing or two ai's are playing
    return (isWhiteAi || isBlackAi) && (!isWhiteAi || !isBlackAi);
  }
  bool isUndoEnabled() {
    return savedData.isUndoPossible && 
      ((isWhiteTurn && !isWhiteAi) || (!isWhiteTurn && !isBlackAi));
  }
  void undo() {
    String lastBoard = savedData.popBoard();
    boardModel = parseBoardString(lastBoard);
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

class MyHomePage extends StatelessWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    //check if saved data is initialized, ana initialize it if not. 
    Future<String?>? lastSavedBoard;
    String currentBoardString = appState.getBoardString();
    if (currentBoardString == startingBoard) {
      lastSavedBoard = appState.savedData.getLastSavedBoard();
    }
    return Scaffold(
      body: Column( 
        //mainAxisAlignment: MainAxisAlignment.center,  
        children: [
          Padding(
            padding: const EdgeInsets.only(top:20.0),
            child: Row(children: [
              ElevatedButton(
                onPressed: () {
                  appState.printBoard();
                },
                child: const Text('Print Board'),
              ),
              ElevatedButton(
                onPressed: () {
                  appState.printSavedData();
                },
                child: const Text('Print\nSaved Data'),
              ),
              
              DropdownButton<int>(
                value: appState.aiDropdownDepth,
                icon: const Icon(Icons.arrow_downward),
                elevation: 16,
                underline: Container(height: 2, color:Colors.grey,),
                onChanged: (int? value) {
                  appState.setAiDepth(value!); 
                },
                items: aiDropdownList.map<DropdownMenuItem<int>>((int value) {
                  return DropdownMenuItem<int>(value:value, child:Text('$value'),
                  );
                }).toList(),
              )
            ]
            ),
            
          ),
          Row(
              children: [
                ElevatedButton(
                  onPressed: () {
                    appState.resetGame();
                  },
                  child: const Text('Reset\nGame'),
                ),
                ElevatedButton(
                  onPressed: () {
                    appState.saveBoard();
                  },
                  child: const Text('Save board'),
                ),
                if (appState.isUndoVisible()) IconButton(
                  icon: const Icon(Icons.undo),
                  onPressed: appState.isUndoEnabled() ? () => appState.undo() : null
                )
              ],),
          Text(appState.gameStatus, style: const TextStyle(fontSize:24, fontWeight: FontWeight.bold)),
          FutureBuilder<String?>(
            future: lastSavedBoard,
            builder: (BuildContext context, AsyncSnapshot<String?>? snapshot) {
              if (snapshot != null && snapshot.hasData && snapshot.data != null){
                String lastBoardString = snapshot.data!;
                
                if (lastBoardString != currentBoardString && currentBoardString == startingBoard){
                  // The user is currently on the starting board, but there was a history that hasn't been deleted. 
                  // This means the user was just playing a game and the app may have gotten closed, but the 
                  // history wasn't deleted by a user action, such as resetting the board. 
                  appState.boardModel = parseBoardString(snapshot.data!);
                }
              }
              for (int i=0; i<appState.boardModel.length; i++){
                for (int j=0; j<appState.boardModel[i].length; j++) {
                  var c = Coord(i,j);
                  var k = i*BoardWidth + j;
                  Color color = appState.getColor(c);
                  var radius = appState.getRadius(c);      
                  var pieceCode = appState.boardModel[i][j];
                  final SquareModel sq = SquareModel(pieceCode, color, radius);
                  if (sq != appState.boardView[k].sq) {
                    final Square newSq = Square(c: c, sq: sq, key: ValueKey(Object.hash(c, sq)),);
                    appState.boardView[k] = newSq;
                  }
                } 
              }
              return GridView.count(
                shrinkWrap: true,
                crossAxisCount: BoardWidth,
                children: [...appState.boardView],);
            },
          ),
        ]  
      ),
    );
  }
}

class Square extends StatelessWidget {
  const Square({
    super.key,
    required this.c,
    required this.sq
  });
  final Coord c;
  final SquareModel sq;

  @override
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    Size screenSize = MediaQuery.of(context).size;
    double screenWidth = screenSize.width;
    double squareW = (screenWidth-20)/8;

    ButtonStyle style = ElevatedButton.styleFrom(
      backgroundColor: sq.color,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(sq.radius),
      ),
      );
    String iconLoc = '';
    if (sq.pieceCode != Space) {
      iconLoc = 'images/${imageMap[sq.pieceCode] ?? ''}';
    }
    if (iconLoc=='') {
      return SizedBox(
          width: squareW,
          height: squareW,
          child: ElevatedButton(
            style:style,
            onPressed: () {
              appState.humanSelectButton(c);
            }, 
            child: const Text('')
          ),
        );
    }
    return SizedBox(
        width: squareW,
        height: squareW,
        child: IconButton(
          style:style,
          icon: Image.asset(iconLoc),
          onPressed: () {
            appState.humanSelectButton(c);
          },  
        )
      );
  }
}
