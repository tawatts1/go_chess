import 'dart:ffi' as ffi;

import 'package:english_words/english_words.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:ffi/ffi.dart';
import 'constants.dart';

import 'package:go_chess/go_chess.dart' as go_chess;

var boardString = '';

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

String getBoardAfterMove(String boardStr, int i1, int j1, int i2, int j2){
  final cBoardStr1 = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> boardPtr = go_chess.getBoardAfterMove(cBoardStr1, i1, j1, i2, j2).cast<Utf8>();
  final newBoardStr = boardPtr.toDartString();
  calloc.free(boardPtr);
  return newBoardStr;
}

String getMoves(String boardStr, int i, int j){
  final cBoardStr = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> movesPtr = go_chess.getNextMoves(cBoardStr, i, j).cast<Utf8>();
  final movesStr = movesPtr.toDartString();
  calloc.free(movesPtr);
  return movesStr;
}

const colorChange = 15;

class MyAppState extends ChangeNotifier {
  var current = WordPair.random();
  int? selectedI;
  int? selectedJ;
  String moveDestinations = '';
  bool isWhiteTurn = true;
  var white = const Color.fromARGB(255, 223, 150, 82);
  var greyedWhite = const Color.fromARGB(255, 223-colorChange, 150-colorChange, 82-colorChange);
  var black = const Color.fromARGB(255, 116, 59, 6);
  var greyedBlack = const Color.fromARGB(255, 116-colorChange, 59-colorChange, 0);
  var selectedColor = const Color.fromARGB(255, 100, 20, 200);
  List<List<String>> board = [
    [BlackRookC, BlackKnight, BlackBishop, BlackQueen, BlackKing, BlackBishop, BlackKnight, BlackRookC,],
    [BlackPawn,BlackPawn,BlackPawn,BlackPawn,BlackPawn,BlackPawn,BlackPawn,BlackPawn,],
    [Space,Space,Space,Space,Space,Space,Space,Space,],
    [Space,Space,Space,Space,Space,Space,Space,Space,],
    [Space,Space,Space,Space,Space,Space,Space,Space,],
    [Space,Space,Space,Space,Space,Space,Space,Space,],
    [WhitePawn,WhitePawn,WhitePawn,WhitePawn,WhitePawn,WhitePawn,WhitePawn,WhitePawn,],
    [WhiteRookC, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, WhiteBishop, WhiteKnight, WhiteRookC,],
  ];
  void clearSelection() {
    selectedI = null;
    selectedJ = null;
    moveDestinations = '';
  }
  void selectButton(int i, int j){
    String piece = board[i][j];
    if (selectedI == null || selectedJ == null) {
      // no selection has been made
      if (piece == Space) {
        print('try clicking on a piece!');
      } else if ((isWhiteTurn && (whiteMap[piece] ?? false)) ||
                 (!isWhiteTurn && !(whiteMap[piece] ?? true))) {
        //mark down the selection and populate the move destinations
        selectedI = i;
        selectedJ = j;
        moveDestinations = getMoves(boardString, i, j);
        print(moveDestinations);
      } else {
        print('not that colors turn');
      }
    } else if (selectedI != null && selectedJ != null) {
      // a move has already been selected. 
      var moveStr = '$i,$j';
      if (moveDestinations.contains(moveStr)){
        print('legal move');
        String newBoardStr = getBoardAfterMove(boardString, selectedI!, selectedJ!, i, j);
        if (newBoardStr.length == BoardHeight * BoardWidth) {
          for (int k=0; k<newBoardStr.length; k++){
            board[k ~/BoardWidth][k%BoardWidth] = newBoardStr[k];
          }
          isWhiteTurn = !isWhiteTurn;
        }
      } else {
        print('not one of the legal moves. Clearing selection');
      }
      clearSelection();
    }
    notifyListeners();
  }
  void printBoard() {
    print(boardString);
  }
  void clearPieces() {
    for (int i= 0; i<board.length; i++) {
      for (int j=0; j<board[0].length; j++) {
        String p = board[i][j];
        if (p != BlackKing && p != WhiteKing) {
          board[i][j] = Space;
        }
      }
    }
    notifyListeners();
  }
  Color getColor(int i, int j){
    bool isLightSquare = (i+j)%2==0;
    if (i==selectedI && j==selectedJ) {
      return selectedColor;
    } else if (moveDestinations.contains('$i,$j')){
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
          Text(appState.current.asLowerCase),
          Row(children: [
            ElevatedButton(
              onPressed: () {
                appState.printBoard();
              },
              child: Text('Print'),
            ),
            ElevatedButton(
              onPressed: () {
                appState.clearPieces();
              },
              child: Text('Delete all\nexcept kings'),
            ),
          ],)
          ] 
          + myBoard(appState),
      ),
    );
  }

  List<Widget> myBoard(MyAppState appState) {
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
    return out;
  }
}

class Square extends StatelessWidget {
  Square({
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
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    Size screenSize = MediaQuery.of(context).size;
    double screenWidth = screenSize.width;
    double screenHeight = screenSize.height;
    double squareW = (screenWidth-20)/8;

    ButtonStyle style = ElevatedButton.styleFrom(
      backgroundColor: color,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(radius),
      ),
      );
    String iconLoc = '';
    if (pieceCode != Space) {
      iconLoc = 'images/${imageMap[pieceCode] ?? ''}';
    }

    if (iconLoc=='') {
      return SizedBox(
          width: squareW,
          height: squareW,
          child: ElevatedButton(
            style:style,
            onPressed: () {
              appState.selectButton(i,j);
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
            appState.selectButton(i,j);
          },
          
        )
      );
  }
}


