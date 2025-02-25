import 'dart:developer';

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
    //check if saved data is initialized, and initialize it if not. 
    Future<List<String>?>? loadedStartupInfo;
    bool skipStartup = false;
    String currentBoardString = appState.board.getBoardString();
    String currentBoardStateString = appState.toString();
    String currentPlayersString = appState.players.toString();
    if (currentBoardString == startingBoard && currentPlayersString == defaultPlayers) {
      log("Starting to load startup info...");
      loadedStartupInfo = appState.savedData.getInfoForStartup();
    } else {
      loadedStartupInfo = null;
      skipStartup = true;
    }
    
    return Scaffold(
      body: FutureBuilder<List<String>?>(
        future: loadedStartupInfo,
        builder: (BuildContext context, AsyncSnapshot<List<String>?>? snapshot) {
          if (snapshot != null && snapshot.hasData && snapshot.data != null){
            if (!skipStartup) {
              if (currentBoardString == startingBoard){
                // The user is currently on the starting board, but there was a history that hasn't been deleted. 
                // This means the user was just playing a game and the app may have gotten closed, but the 
                // history wasn't deleted by a user action, such as resetting the board. 
                String lastSavedStateVal = snapshot.data![0];
                if (lastSavedStateVal != currentBoardStateString && lastSavedStateVal != "") {
                  appState.loadBoardStateFromString(lastSavedStateVal);
                }
              } else {
                log("warning: detected starting board that was later not starting board...");
              }
              if (currentPlayersString == defaultPlayers){
                String savedPlayerString = snapshot.data![1];
                if (savedPlayerString != currentPlayersString){
                  appState.loadPlayerStateFromString(savedPlayerString);
                }
              } else {
                log("warning: detected default player state that later changed...");
              }
              appState.setIsUndoEnabled();
              appState.setIsUndoVisible();
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
                  if (appState.undoButtonModel.isVisible) 
                  IconButton(
                    icon: const Icon(Icons.undo),
                    onPressed: appState.undoButtonModel.isEnabled ? () => appState.undo() : null
                  )
                ],
              ),
              Row(
                children: [
                  Expanded(
                    child: DropdownMenu<String>(
                      initialSelection: appState.players.getPlayerName(true), //appState.players.aiDropdownDepth,
                      label: const Text("White Player"),
                      requestFocusOnTap: true,
                      onSelected: (String? value) {
                        appState.setPlayer(value!, true); 
                      },
                      dropdownMenuEntries: ["Human", "Ai"].map<DropdownMenuEntry<String>>((String value) {
                        return DropdownMenuEntry<String>(value:value, label:value,
                        );
                      }).toList(),
                    ),
                  ),
                  Expanded(
                    child: DropdownMenu<String>(
                      initialSelection: appState.players.getPlayerName(false), //appState.players.aiDropdownDepth,
                      label: const Text("Black Player"),
                      requestFocusOnTap: true,
                      onSelected: (String? value) {
                        appState.setPlayer(value!, false); 
                      },
                      dropdownMenuEntries: ["Human", "Ai"].map<DropdownMenuEntry<String>>((String value) {
                        return DropdownMenuEntry<String>(value:value, label:value,
                        );
                      }).toList(),
                    ),
                  ),
                  Expanded(
                    child: DropdownMenu<int>(
                      initialSelection: appState.players.aiDropdownDepth,
                      label: const Text("Ai depth"),
                      requestFocusOnTap: true,
                      onSelected: (int? value) {
                        appState.setAiDepth(value!); 
                      },
                      dropdownMenuEntries: aiDropdownList.map<DropdownMenuEntry<int>>((int value) {
                        return DropdownMenuEntry<int>(value:value, label:'$value',
                        );
                      }).toList(),
                    ),
                  ),
                ]
              ),
              Text(appState.board.gameStatus, style: const TextStyle(fontSize:24, fontWeight: FontWeight.bold)),
              GridView.count(
                    shrinkWrap: true,
                    crossAxisCount: BoardWidth,
                    children: List.from(
                      appState.shouldBoardBeFlipped() ? appState.board.boardView.reversed : appState.board.boardView
                          ),
              ),
            ]  
          );
        }
      ),
    );
  }
}

