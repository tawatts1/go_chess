import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'constants.dart';
import 'coord.dart';
import 'ffi_funcs.dart';

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
    String piece = board[c.i][c.j];
    bool isNotifyAi = false;
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

          board = parseBoardString(newBoardStr);
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
      setBoardString(); // todo: get rid of this function. either implement getBoardString or store it whenever it changes. 
      String aiMove = await getAiChosenMove(boardString, isWhiteTurn, 'simple', aiDropdownDepth);
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
    log(boardString);
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
        var c = Coord(i,j);
        Color color = appState.getColor(c);
        var radius = appState.getRadius(c);      
        var pieceCode = row[j];
        rowView.add(
          Square(pieceCode: pieceCode, color: color, radius: radius, c: c)
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
    required this.c
  });

  final String pieceCode;
  final Color color;
  final double radius;
  final Coord c;

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
              appState.humanSelectButton(widget.c);
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
            appState.humanSelectButton(widget.c);
          },
          
        )
      );
  }
}
