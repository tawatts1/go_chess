import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'constants.dart';
import 'coord.dart';
import 'square.dart';
import 'game_state.dart';



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

class MyHomePage extends StatelessWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    //check if saved data is initialized, ana initialize it if not. 
    Future<String?>? lastSavedBoard;
    String currentBoardString = appState.board.getBoardString();
    if (currentBoardString == startingBoard) {
      lastSavedBoard = appState.savedData.getLastSavedBoard();
    }
    return Scaffold(
      body: FutureBuilder<String?>(
        future: lastSavedBoard,
        builder: (BuildContext context, AsyncSnapshot<String?>? snapshot) {
          if (snapshot != null && snapshot.hasData && snapshot.data != null){
            String lastBoardString = snapshot.data!;
            appState.setIsUndoEnabled();
            appState.setIsUndoVisible();
            if (lastBoardString != currentBoardString && currentBoardString == startingBoard){
              // The user is currently on the starting board, but there was a history that hasn't been deleted. 
              // This means the user was just playing a game and the app may have gotten closed, but the 
              // history wasn't deleted by a user action, such as resetting the board. 
              appState.board.boardModel = parseBoardString(snapshot.data!);
            }
          }
          for (int i=0; i<appState.board.boardModel.length; i++){
            for (int j=0; j<appState.board.boardModel[i].length; j++) {
              var c = Coord(i,j);
              var k = i*BoardWidth + j;
              Color color = appState.getColor(c);
              var radius = appState.getRadius(c);      
              var pieceCode = appState.board.boardModel[i][j];
              final SquareModel sq = SquareModel(pieceCode, color, radius);
              if (sq != appState.board.boardView[k].sq) {
                final Square newSq = Square(c: c, sq: sq, key: ValueKey(Object.hash(c, sq)),);
                appState.board.boardView[k] = newSq;
              }
            } 
          }
          return Column( 
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
                    if (appState.undoButtonModel.isVisible) 
                    IconButton(
                      icon: const Icon(Icons.undo),
                      onPressed: appState.undoButtonModel.isEnabled ? () => appState.undo() : null
                    )
                  ],),
              Text(appState.board.gameStatus, style: const TextStyle(fontSize:24, fontWeight: FontWeight.bold)),
              GridView.count(
                    shrinkWrap: true,
                    crossAxisCount: BoardWidth,
                    children: [...appState.board.boardView],),
            ]  
          );
        }
      ),
    );
  }
}

