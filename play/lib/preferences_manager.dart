import 'package:shared_preferences/shared_preferences.dart';
import 'dart:developer';

class PreferencesManager {
  late final SharedPreferencesAsync prefs;
  late List<String> boardStrings;
  bool isUndoPossible = false;
  bool isLoaded = false;
  final int maxBoards = 3;


  PreferencesManager() {
    prefs = SharedPreferencesAsync();
  }
  @override 
  String toString() {
    return 'PreferencesManager: $boardStrings, $isUndoPossible, $isLoaded, $maxBoards';
  }
  loadBoards() async {
    isLoaded = true;
    List<String>? loaded = await prefs.getStringList('b');
    boardStrings = loaded ?? [];
    isUndoPossible = boardStrings.length > 1;
  }
  saveBoards() async {
    prefs.setStringList('b',boardStrings);
  }
  addBoard(String boardStr) async {
    boardStrings.add(boardStr);
    if (boardStrings.length > maxBoards) {
      boardStrings = boardStrings.sublist(1);
    }
    isUndoPossible = boardStrings.length > 1;
    saveBoards();
  }
  String popBoard() {
    if (isUndoPossible){
      boardStrings.removeLast();
      String out = boardStrings.last;
      isUndoPossible = boardStrings.length > 1;
      saveBoards();
      return out;
    } else {
      log('tried to pop board history when there was none. ');
      return '';
    }
  }
  Future<String?> getLastSavedBoard() async {
    String? out;
    if (!isLoaded) {
      await loadBoards();
    }
    if (boardStrings.isNotEmpty){
      out =  boardStrings.last;
    } else {
      out =  null;
    }
    return out;
  }
  Future<void> clearBoards() async {
    boardStrings.clear();
    boardStrings = [];
    isUndoPossible = false;
    await saveBoards();
  }
}