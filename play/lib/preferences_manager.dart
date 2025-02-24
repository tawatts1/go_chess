import 'package:shared_preferences/shared_preferences.dart';
import 'dart:developer';

final String boardSnapshotsKey = "b";

class PreferencesManager {
  late final SharedPreferencesAsync prefs;
  late List<String> boardSnapshots;
  bool isUndoPossible = false;
  bool isLoaded = false;
  final int maxBoards = 3;


  PreferencesManager() {
    prefs = SharedPreferencesAsync();
  }
  @override 
  String toString() {
    return 'PreferencesManager: $boardSnapshots, $isUndoPossible, $isLoaded, $maxBoards';
  }
  loadBoards() async {
    isLoaded = true;
    List<String>? loaded = await prefs.getStringList(boardSnapshotsKey);
    boardSnapshots = loaded ?? [];
    isUndoPossible = boardSnapshots.length > 1;
  }
  saveBoards() async {
    prefs.setStringList(boardSnapshotsKey, boardSnapshots);
  }
  addBoardSnapshot(String boardStr) async {
    boardSnapshots.add(boardStr);
    if (boardSnapshots.length > maxBoards) {
      boardSnapshots = boardSnapshots.sublist(1);
    }
    isUndoPossible = boardSnapshots.length > 1;
    saveBoards();
  }
  String popBoard() {
    if (isUndoPossible){
      boardSnapshots.removeLast();
      String out = boardSnapshots.last;
      isUndoPossible = boardSnapshots.length > 1;
      saveBoards();
      return out;
    } else {
      log('tried to pop board history when there was none. ');
      return '';
    }
  }
  Future<String?> getLastSavedState() async {
    String? out;
    if (!isLoaded) {
      await loadBoards();
    }
    if (boardSnapshots.isNotEmpty){
      out =  boardSnapshots.last;
    } else {
      out =  null;
      log("No saved board states found");
    }
    return out;
  }
  Future<void> clearBoardStates() async {
    boardSnapshots.clear();
    boardSnapshots = [];
    isUndoPossible = false;
    await saveBoards();
  }
}