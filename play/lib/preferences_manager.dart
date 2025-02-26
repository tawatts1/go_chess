import 'package:mutex/mutex.dart';

import 'package:shared_preferences/shared_preferences.dart';
import 'dart:developer';

const String boardSnapshotsKey = "b";
const String playersKey = "p";

const maxStatesSaved = 5;
class PreferencesManager {
  late final SharedPreferencesAsync prefs;
  List<String> boardSnapshots = [];
  bool isUndoPossible = false;
  bool isLoaded = false;
  final int maxBoardStates = maxStatesSaved;
  String players = "";
  final mtx = ReadWriteMutex(); // Used to access/modify the following: isLoaded, prefs.get* and prefs.set*

  PreferencesManager() {
    prefs = SharedPreferencesAsync();
  }
  @override 
  String toString() {
    return 'PreferencesManager: $boardSnapshots, $isUndoPossible, $isLoaded, $maxBoardStates\n$players';
  }
  load() async {
    await mtx.acquireWrite(); // accessing isLoaded and doing the load operation. 
    try {
      if (!isLoaded) {
        isLoaded = true;
        Future<List<String>?> boardSnapshotsFuture = prefs.getStringList(boardSnapshotsKey);
        Future<String?> playersFuture = prefs.getString(playersKey);
        List<String>? loadedBoardSnapshots = await boardSnapshotsFuture;
        String? loadedPlayers = await playersFuture;
        boardSnapshots = loadedBoardSnapshots ?? [];
        players = loadedPlayers ?? "";
        isUndoPossible = boardSnapshots.length > 2;
      }
    } finally {
      mtx.release();
    }
  }
  savePlayers() async {
    await mtx.acquireWrite();
    if (!isLoaded){
      log("Trying to save players when preferences aren't even loaded. ");
    }
    try {
      prefs.setString(playersKey, players);
    } finally {
      mtx.release();
    } 
  }
  Future<String> getPlayers() async {
    await mtx.acquireRead();
    String out = "";
    try {
      if (!isLoaded){
        log("Trying to get players when preferences aren't even loaded.");
      }
      out = players;
    } finally {
      mtx.release();
    }
    return out;
  }
  addBoardSnapshot(String boardStr) async {
    await mtx.acquireWrite();
    try {
      if (!isLoaded){
        log("Trying to add a board snapshot when preferences aren't even loaded.");
      }
      boardSnapshots.add(boardStr);
      if (boardSnapshots.length > maxBoardStates) {
        boardSnapshots = boardSnapshots.sublist(1);
      }
      isUndoPossible = boardSnapshots.length > 2;
      await prefs.setStringList(boardSnapshotsKey, boardSnapshots);
    } finally {
      mtx.release();
    }
  }
  Future<String> popBoard(int N) async {
    String out = "";
    await mtx.acquireWrite();
    try{
      if (isUndoPossible){
        boardSnapshots.removeLast();
        boardSnapshots.removeLast();
        out = boardSnapshots.last;
        isUndoPossible = boardSnapshots.length > 2;
        await prefs.setStringList(boardSnapshotsKey, boardSnapshots);
      } else {
        log('tried to pop board history when there was none. ');
      }
    } finally {
      mtx.release();
    }
    return out;
  }
  Future<String> getLastSavedState() async {
    String out = "";
    await mtx.acquireRead();
    try {
      if (!isLoaded){
        log("Trying to get last board state when preferences aren't even loaded.");
      }
      if (boardSnapshots.isNotEmpty){
        out =  boardSnapshots.last;
      } else {
        log("No saved board states found");
      }
    } finally {
      mtx.release();
    }
    return out;
  }
  clearBoardStates() async {
    await mtx.acquireWrite();
    try {
      if (!isLoaded){
        log("Trying to clear board state when preferences aren't even loaded.");
      }
      boardSnapshots.clear();
      boardSnapshots = [];
      isUndoPossible = false;
      await prefs.setStringList(boardSnapshotsKey, boardSnapshots);
    } finally {
      mtx.release();
    }
  }
  Future<List<String>?> getInfoForStartup() async {
    List<String>? out = [];
    if (!isLoaded) {
      await load();
    } else {
      log("Getting startup info after the preferences have already been loaded.");
    }
    out.add(await getLastSavedState());
    out.add(await getPlayers());
    return out;
  }
}