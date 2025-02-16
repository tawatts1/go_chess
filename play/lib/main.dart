import 'dart:developer';
import 'dart:ffi' as ffi;

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:ffi/ffi.dart';
import 'constants.dart';

import 'package:go_chess/go_chess.dart' as go_chess;

var boardString = '';
List<List<String>> parseBoardString(String boardStr) {
  List<List<String>> board = [[],[],[],[],[],[],[],[],];
  if (boardStr.length == BoardHeight * BoardWidth) {
    for (int k=0; k<boardStr.length; k++){
      board[k ~/BoardWidth].add(boardStr[k]);
    }
  }
  return board;
}

class Coord{
  final int i;
  final int j;
  const Coord(this.i, this.j);
}

Future<Coord> returnAfterDelay(Coord out, Duration t) async {
    return Future<Coord>.delayed(t, () {
      return out;
    });
    //sleep(t);
    //return out;
}


void main() {
  runApp(MyApp());
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
        home: MyHomePage(),
      ),
    );
  }
}

Future<String> getAiChosenMove(String boardStr, bool isWhite, String aiName, int N) async {
  log('$aiName,$boardStr,$N,${isWhite?"w":"b"}');
  final ffi.Pointer<ffi.Char> cBoardStr = boardStr.toNativeUtf8().cast<ffi.Char>();
  int isWhiteInt = isWhite ? 1 : 0;
  final ffi.Pointer<ffi.Char> cAiName = aiName.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> movePtr = (await go_chess.getAiChosenMove(cBoardStr, isWhiteInt, cAiName, N)).cast<Utf8>();
  final moveString = movePtr.toDartString();
  calloc.free(cBoardStr);
  calloc.free(cAiName);
  calloc.free(movePtr);
  return moveString;
}

String getBoardAfterMove(String boardStr, int i1, int j1, int i2, int j2){
  final ffi.Pointer<ffi.Char> cBoardStr1 = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> boardPtr = go_chess.getBoardAfterMove(cBoardStr1, i1, j1, i2, j2).cast<Utf8>();
  final newBoardStr = boardPtr.toDartString();
  calloc.free(boardPtr);
  calloc.free(cBoardStr1);
  return newBoardStr;
}

String getMoves(String boardStr, int i, int j){
  final cBoardStr = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> movesPtr = go_chess.getNextMoves(cBoardStr, i, j).cast<Utf8>();
  final movesStr = movesPtr.toDartString();
  calloc.free(movesPtr);
  calloc.free(cBoardStr);
  return movesStr;
}

const colorChange = 30;

class MyAppState extends ChangeNotifier {
  int? selectedI;
  int? selectedJ;
  String moveDestinations = '';
  bool isWhiteTurn = true;
  bool isBlackAi = true;
  bool isWhiteAi = false;
  String gameStatus = statusWhiteMove;
  bool isGameOver = false;
  String indicatedCoords = '';
  var white = const Color.fromARGB(255, 223, 150, 82);
  var greyedWhite = const Color.fromARGB(255, 180,170,170);
  var black = const Color.fromARGB(255, 116, 59, 6);
  var greyedBlack = const Color.fromARGB(255, 90,75,75);
  var selectedColor = const Color.fromARGB(255, 120, 0, 100);
  int aiDropdownDepth = 4;
  List<List<String>> board = parseBoardString(startingBoard);
  void resetGame() {
    board = parseBoardString(startingBoard);
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
    selectedI = null;
    selectedJ = null;
    moveDestinations = '';
  }
  void humanSelectButton(int i, int j){
    //functions that humans have to use to select the buttons
    if (isGameOver) {
      log('game is over');
    } else if ((isWhiteTurn && isWhiteAi) || (!isWhiteTurn && isBlackAi)){
      log('It is an AIs turn');
    } else {
      selectButton(i,j);
    }
  }
  void selectButton(int i, int j){
    String piece = board[i][j];
    bool isNotifyAi = false;
    if (selectedI == null || selectedJ == null) {
      // no selection has been made
      if (piece == Space) {
        log('try clicking on a piece!');
      } else if ((isWhiteTurn && (whiteMap[piece] ?? false)) ||
                 (!isWhiteTurn && !(whiteMap[piece] ?? true))) {
        //mark down the selection and populate the move destinations
        selectedI = i;
        selectedJ = j;
        moveDestinations = getMoves(boardString, i, j);
        log(moveDestinations);
      } else {
        log('not that colors turn');
      }
    } else if (selectedI != null && selectedJ != null) {
      // a move has already been selected. 
      var moveStr = '$i,$j';
      if (moveDestinations.contains(moveStr)){
        log('legal move');
        String boardResult = getBoardAfterMove(boardString, selectedI!, selectedJ!, i, j);
        List<String> resultList = boardResult.split(',');
        if (resultList.length == 2) {
          String newBoardStr = resultList[0];
          gameStatus = resultList[1];
          if (gameStatus == statusCheckMate || gameStatus == statusStaleMate){
            isGameOver = true;
          }

          board = parseBoardString(newBoardStr);
          isWhiteTurn = !isWhiteTurn;
          isNotifyAi = true;
          indicatedCoords = '$selectedI,$selectedJ|$i,$j';
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
  void notifyAi() {
    if (!isGameOver && ((isWhiteTurn && isWhiteAi) || (!isWhiteTurn && isBlackAi))) {
      //it is the ai's turn
      setBoardString(); // todo: get rid of this function. either implement getBoardString or store it whenever it changes. 
      Future<String> aiMove = getAiChosenMove(boardString, isWhiteTurn, 'simple', aiDropdownDepth);
      aiMove.then((value) => parseAndDoAiMove(value))
      .catchError((error) => log(error));
    }
  }
  void parseAndDoAiMove(String moveStr) {
    List<String> indexList = moveStr.split(',');
    if (indexList.length == 4) {
      try {
        int i1 = int.parse(indexList[0]);
        int j1 = int.parse(indexList[1]);
        int i2 = int.parse(indexList[2]);
        int j2 = int.parse(indexList[3]);
        simulateClickBoard(Coord(i1,j1), Coord(i2, j2));
      } catch(ex) {
        log("failed to parse ai move");
      }
    }
  }
  void simulateClickBoard(Coord c1, Coord c2) async {
    if (!isGameOver && ((isWhiteTurn && isWhiteAi) || (!isWhiteTurn && isBlackAi))) {
      setBoardString();
      //Coord click1 = await returnAfterDelay(c1, Duration(milliseconds: 50));
      //click.then((value) => selectButton(value.i, value.j))
      //.catchError((error) => log(error));
      //secondClick.then((click) => )
      //selectButton(click1.i, click1.j);
      selectButton(c1.i, c1.j);
      
      Coord click2 = await returnAfterDelay(c2, Duration(milliseconds: 700));
      selectButton(click2.i, click2.j);
      
    }
  }
  void printBoard() {
    log(boardString);
  }
  Color getColor(int i, int j){
    bool isLightSquare = (i+j)%2==0;
    if (i==selectedI && j==selectedJ) {
      return selectedColor;
    } else if (indicatedCoords.contains('$i,$j')){
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
  double getRadius(int i, int j){
    if (moveDestinations.contains('$i,$j')) {
      return 15;
    } else if (selectedI!= null && selectedJ != null && selectedI! == i && selectedJ! == j) {
      return 15;
    } else {
      return 1;
    }
  }
  void setBoardString() {
    boardString = '';
    for (int i=0; i<board.length; i++) {
      for (int j=0; j<board[i].length; j++) {
        boardString += board[i][j];
      }
    }
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
                child: Text('Print'),
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
          myBoard(appState)
          ] 
          //+ myBoard(appState),
      ),
    );
  }

  Column myBoard(MyAppState appState) {
    // out is the list of row widgets that make up the board. 
    List<Widget> out = [];
    boardString = '';
    for (int i=0; i<appState.board.length; i++){
      // row is a list of strings
      var row = appState.board[i];
      List<Widget> rowView = [];
      for (int j=0; j<row.length; j++) {
        Color color = appState.getColor(i, j);
        var radius = appState.getRadius(i,j);      
        var pieceCode = row[j];
        rowView.add(
          Square(pieceCode: pieceCode, color: color, radius: radius, i: i, j: j)
        );
        boardString += pieceCode;
      } 
      out.add(Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children:rowView,));
    }
    return Column(children: out,);
  }
}

class Square extends StatefulWidget {
  const Square({
    super.key,
    required this.pieceCode,
    required this.color,
    required this.radius,
    required this.i,
    required this.j
  });

  final String pieceCode;
  final Color color;
  final double radius;
  final int i;
  final int j;

  @override
  State<Square> createState() => _SquareState();
}

class _SquareState extends State<Square> {
  @override
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    Size screenSize = MediaQuery.of(context).size;
    double screenWidth = screenSize.width;
    double screenHeight = screenSize.height;
    double squareW = (screenWidth-20)/8;

    ButtonStyle style = ElevatedButton.styleFrom(
      backgroundColor: widget.color,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(widget.radius),
      ),
      );
    String iconLoc = '';
    if (widget.pieceCode != Space) {
      iconLoc = 'images/${imageMap[widget.pieceCode] ?? ''}';
    }

    if (iconLoc=='') {
      return SizedBox(
          width: squareW,
          height: squareW,
          child: ElevatedButton(
            style:style,
            onPressed: () {
              appState.humanSelectButton(widget.i,widget.j);
            }, 
            child: Text('')
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
            appState.humanSelectButton(widget.i,widget.j);
          },
          
        )
      );
  }
}
