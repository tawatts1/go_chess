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
  late bool isUndoPossible;
  bool isLoaded = false;
  final int maxBoards = 2;
  PreferencesManager() {
    prefs = SharedPreferencesAsync();
  }
  loadBoards() async {
    List<String>? loaded = await prefs.getStringList('b');
    boardStrings = loaded ?? [];
    if (boardStrings.isNotEmpty) {
      isUndoPossible = true;
    } else {
      isUndoPossible = false;
    }
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
    saveBoards();
  }
  String popBoard() {
    if (boardStrings.isNotEmpty){
      String out = boardStrings.removeLast();
      saveBoards();
      return out;
    } else {
      log('tried to pop board history when there was none. ');
      return '';
    }
  }
  Future<String> getLastBoard() async {
    if (!isLoaded) {
      await loadBoards();
    }
    return boardStrings.last;
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
  void resetGame() {
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
      selectButton(c);
    }
  }
  void selectButton(Coord c){
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
        selectButton(Coord(i1, j1));
        Coord click2 = await Future.delayed(const Duration(milliseconds: 700),  () => Coord(i2,j2));
        selectButton(click2);
      } catch(ex) {
        log("failed to parse ai move");
      }
    }
  }
 
  void printBoard() {
    log(getBoardString());
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
    if (!appState.savedData.isLoaded){
      appState.savedData.loadBoards();
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
                child: const Text('Print'),
              ),
              ElevatedButton(
                onPressed: () {
                  appState.resetGame();
                },
                child: const Text('Reset\nGame'),
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
          Text(appState.gameStatus, style: const TextStyle(fontSize:24, fontWeight: FontWeight.bold)),
          GridView.count(
            shrinkWrap: true,
            crossAxisCount: BoardWidth,
            children: [...appState.boardView],),
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
