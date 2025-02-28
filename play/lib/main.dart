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
      child: const MaterialApp(
        title: 'Namer App',
        home: MyHomePage(),
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
    
    return FutureBuilder<List<String>?>(
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
              String savedThemeString = snapshot.data![2];
              appState.loadThemeStateFromString(savedThemeString);
              appState.doDuringStartup();
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
          ThemeData calculatedTheme = ThemeData(
              colorScheme: ColorScheme.fromSeed(
                seedColor: lightBlack, 
                brightness: appState.theme.isDarkTheme ? Brightness.dark : Brightness.light),
            );
          Color primaryColor = calculatedTheme.colorScheme.primary;
          Color secondaryColor = calculatedTheme.colorScheme.secondary;
          return Theme(
            data: calculatedTheme,
            child: Scaffold(
              body: Column( 
                //mainAxisAlignment: MainAxisAlignment.center,  
                children: [
                  Padding(
                    padding: const EdgeInsets.only(top:20.0),
                    child: Row(
                      children: [
                        ElevatedButton(
                          onPressed: () {
                            appState.printBoard();
                          },
                          child: Text('Print Board',
                          style: TextStyle(color: primaryColor),),
                        ),
                        ElevatedButton(
                          onPressed: () {
                            appState.printSavedData();
                          },
                          child: const Text('Print\nSaved Data'),
                        ),
                        Switch(
                          value: appState.theme.isDarkTheme,
                          onChanged: (bool val) {appState.setTheme(val);}
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
                      if (appState.undoButtonModel.isVisible) 
                      IconButton(
                        icon: Icon(Icons.undo, 
                          color: appState.undoButtonModel.isEnabled ? primaryColor : secondaryColor,
                        ),
                        onPressed: appState.undoButtonModel.isEnabled ? () => appState.undo() : null
                      ),
                      if (appState.playButtonModel.isVisible)
                      IconButton(
                        icon: Icon(Icons.play_arrow,
                        color: appState.playButtonModel.isEnabled ? primaryColor : secondaryColor,
                        ),
                        onPressed: appState.playButtonModel.isEnabled ? () => appState.playPause() : null
                      ),
                      if (appState.pauseButtonModel.isVisible)
                      IconButton(
                        icon: Icon(Icons.pause,
                        color: appState.pauseButtonModel.isEnabled ? primaryColor : secondaryColor,
                        ),
                        onPressed: appState.pauseButtonModel.isEnabled ? () => appState.playPause() : null
                      ),
                    ],
                  ),
                  Row(
                    children: [
                      Expanded(
                        child: CustomDropdownMenu(
                          appState: appState, 
                          textColor: primaryColor, 
                          dropdownLabelText: "White Player",
                          getter: () {return appState.players.getPlayerName(true);},
                          setter: (String val) {appState.setPlayer(val, true);},
                          entries: const ["Human", "Ai"],
                          ),
                      ),
                      Expanded(
                        child: CustomDropdownMenu<String>(
                          appState: appState, 
                          textColor: primaryColor, 
                          dropdownLabelText: "Black Player",
                          getter: () {return appState.players.getPlayerName(false);},
                          setter: (String val) {appState.setPlayer(val, false);},
                          entries: const ["Human", "Ai"],
                        )

                      ),
                      Expanded(
                        child: CustomDropdownMenu<int>(
                          appState: appState, 
                          textColor: primaryColor, 
                          dropdownLabelText: "Ai depth",
                          getter: () {return appState.players.aiDropdownDepth;},
                          setter: (int val) {appState.setAiDepth(val);},
                          entries: aiDropdownList,
                        ),
                      ),
                    ]
                  ),
                  Text(appState.board.gameStatus, style: TextStyle(fontSize:24, fontWeight: FontWeight.bold, color: primaryColor)),
                  GridView.count(
                        shrinkWrap: true,
                        crossAxisCount: BoardWidth,
                        children: List.from(
                          appState.shouldBoardBeFlipped() ? appState.board.boardView.reversed : appState.board.boardView
                              ),
                  ),
                ]  
              ),
            ),
          );
        }
      );
  }
}

class CustomDropdownMenu<T> extends StatelessWidget {
  const CustomDropdownMenu({
    super.key,
    required this.appState,
    required this.textColor,
    required this.dropdownLabelText,
    required this.getter,
    required this.setter,
    required this.entries,
  });

  final MyAppState appState;
  final Color textColor;
  final String dropdownLabelText;
  final Function getter;
  final Function setter;
  final List<T> entries;

  @override
  Widget build(BuildContext context) {
    return DropdownMenu<T>(
      initialSelection: getter(), //appState.players.aiDropdownDepth,
      textStyle: TextStyle(color: textColor),
      trailingIcon: Icon(Icons.arrow_drop_down, 
        color:textColor,
      ),
      selectedTrailingIcon: Icon(Icons.arrow_drop_up, 
        color: textColor,
      ),
      inputDecorationTheme: InputDecorationTheme(
        enabledBorder: OutlineInputBorder(
          borderSide: BorderSide(
            color: textColor
          )
        )
      ),
      label: Text(dropdownLabelText, style: TextStyle(color: textColor)),
      requestFocusOnTap: true,
      onSelected: (T? value) {
        setter(value!); 
      },
      dropdownMenuEntries: entries.map<DropdownMenuEntry<T>>((T value) {
        return DropdownMenuEntry<T>(value:value, 
        label:"$value",
        labelWidget: Text("$value", style: TextStyle(color: textColor)),
        );
      }).toList(),
    );
  }
}

