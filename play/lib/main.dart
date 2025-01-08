import 'package:english_words/english_words.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

const piecesList = [
  '0',
  'p',
  'n',
  'b',
  'r',
  'q',
  'k',
  'o',
  'a',
  'P',
  'N',
  'B',
  'R',
  'Q',
  'K',
  'O',
  'A'
];

const Space = '0';

const BlackPawn = 'p';
const BlackKnight = 'n';
const BlackBishop = 'b';
const BlackRookNC = 'r';
const BlackQueen = 'q';
const BlackKing = 'k';
const BlackRookC = 'o';
const BlackPawnEP = 'a';

const WhitePawn =  'P';
const WhiteKnight =  'N';
const WhiteBishop =  'B';
const WhiteRookNC =  'R';
const WhiteQueen =  'Q';
const WhiteKing =  'K';
const WhiteRookC =  'O';
const WhitePawnEP =  'A';

const imageMap = {
  BlackPawn: 'bp.png',
  BlackKnight: 'bn.png',
  BlackBishop: 'bb.png',
  BlackRookNC: 'br.png',
  BlackQueen: 'bq.png',
  BlackKing: 'bk.png',
  BlackRookC: 'br.png',
  BlackPawnEP: 'bp.png',
  WhitePawn: 'wp.png',
  WhiteKnight: 'wn.png',
  WhiteBishop: 'wb.png',
  WhiteRookNC: 'wr.png',
  WhiteQueen: 'wq.png',
  WhiteKing: 'wk.png',
  WhiteRookC: 'wr.png',
  WhitePawnEP: 'wp.png'
};

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

class MyAppState extends ChangeNotifier {
  var current = WordPair.random();
  int? piece_i;
  int? piece_j;
  String? piece_code;
  var white = Color.fromARGB(255, 241, 189, 129);
  var black = const Color.fromARGB(255, 99, 46, 11);
  List<List<String>> board = [
    ['o', 'n', 'b', 'q', 'k', 'b', 'n', 'o',],
    ['p', 'p', 'p', 'p', '0', 'p', 'p', 'p',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['0', '0', '0', '0', 'p', '0', '0', '0',],
    ['0', '0', '0', '0', 'P', '0', '0', '0',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['P', 'P', 'P', 'P', '0', 'P', 'P', 'P',],
    ['O', 'N', 'B', 'Q', 'K', 'B', 'N', 'O',],
  ];
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
        Color color;
        switch ((i+j)%2){
          case 0:
            color=appState.white;
            break;
          default: 
            color=appState.black;
        }
        var pieceCode = row[j];
        rowView.add(
          Square(pieceCode: pieceCode, color: color, i: i, j: j)
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
    required this.i,
    required this.j
  });

  final String pieceCode;
  final Color color;
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
        borderRadius: BorderRadius.circular(1),
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
              ;
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
            ;
          },
          
        )
      );
  }
}


