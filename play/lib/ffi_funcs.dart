import 'dart:developer';
import 'dart:ffi' as ffi;
import 'package:ffi/ffi.dart';
import 'package:go_chess/go_chess.dart' as go_chess;
import 'coord.dart';

/// This file contains functions that call functions from the foreign function interface (FFI) 

Future<String> getAiChosenMove(String boardStr, bool isWhite, String aiName, int N) async {
  log('$aiName,$boardStr,$N,${isWhite?"w":"b"}');
  final ffi.Pointer<ffi.Char> cBoardStr = boardStr.toNativeUtf8().cast<ffi.Char>();
  int isWhiteInt = isWhite ? 1 : 0;
  final ffi.Pointer<ffi.Char> cAiName = aiName.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> movePtr = (await go_chess.getAiChosenMove(cBoardStr, isWhiteInt, cAiName, N)).cast<Utf8>();
  final moveString = movePtr.toDartString();
  calloc.free(cBoardStr);
  calloc.free(cAiName);
  calloc.free(movePtr);
  return moveString;
}

String getBoardAfterMove(String boardStr, Coord c1, Coord c2){
  final ffi.Pointer<ffi.Char> cBoardStr1 = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> boardPtr = go_chess.getBoardAfterMove(cBoardStr1, c1.i, c1.j, c2.i, c2.j).cast<Utf8>();
  final newBoardStr = boardPtr.toDartString();
  calloc.free(boardPtr);
  calloc.free(cBoardStr1);
  return newBoardStr;
}

String getMoves(String boardStr, Coord c){
  final cBoardStr = boardStr.toNativeUtf8().cast<ffi.Char>();
  final ffi.Pointer<Utf8> movesPtr = go_chess.getNextMoves(cBoardStr, c.i, c.j).cast<Utf8>();
  final movesStr = movesPtr.toDartString();
  calloc.free(movesPtr);
  calloc.free(cBoardStr);
  return movesStr;
}