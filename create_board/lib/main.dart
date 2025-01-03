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
  var white = Color.fromARGB(255, 241, 189, 129);
  var black = const Color.fromARGB(255, 99, 46, 11);
  List<List<String>> board = [
    ['o', 'n', 'b', 'k', 'q', '0', '0', 'o',],
    ['p', 'p', 'p', 'p', 'p', '0', 'p', 'p',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['0', '0', '0', '0', '0', '0', '0', '0',],
    ['P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',],
    ['O', 'N', 'B', 'K', 'Q', '0', '0', 'O',],
  ];
  void getNext() {
    current = WordPair.random();
    print(boardString);
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
          ElevatedButton(
            onPressed: () {
              appState.getNext();
            },
            child: Text('Next'),
          ),
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
          Square(pieceCode: pieceCode, color: color)
        );
        boardString += pieceCode;
      } 
      out.add(Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children:rowView,));
      boardString += '\n';
    }

    out.add(Expanded(
      child: Row(children: [
        getButtonColumn(0, 1, appState), 
        getButtonColumn(1, 9, appState),
        getButtonColumn(7, 9, appState),
        getButtonColumn(9, 15, appState),
        getButtonColumn(15, piecesList.length, appState)
      ]),
    ));
      
    // for (int i=0; i<piecesList.length; i++){
      
    //   }
    // }
    return out;
  }
}

Column getButtonColumn(int start, int end, var appState) {
  List<Widget> extraPieces = [];
  for (int i=start; i<end; i++){
      var pieceCode = piecesList[i];
      var iconLoc = '';
      if (pieceCode != Space) {
        iconLoc = 'images/${imageMap[pieceCode] ?? ''}';
      }
      if (iconLoc == '') {
        extraPieces.add(Expanded(
          child: Row(children: [Text(piecesList[i]), ElevatedButton(
            onPressed: () {
              appState.getNext();
            }, 
            child: Text(''))]),
        )
      );
      } else {
        extraPieces.add(Expanded(
          child: Row(children: [
              Text(pieceCode),
              IconButton(
                icon: Image.asset(iconLoc),
                onPressed: () {
                  appState.getNext();
                }
              )
            ]
          ),
        ));
      }
  }
  return Column(children: extraPieces);
}

class Square extends StatelessWidget {
  Square({
    super.key,
    required this.pieceCode,
    required this.color,
  });

  final String pieceCode;
  final Color color;
  
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
              appState.getNext();
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
            appState.getNext();
          },
          
        )
      );
  }
}


