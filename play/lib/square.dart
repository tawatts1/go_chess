import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'constants.dart';
import 'coord.dart';
import 'game_state.dart';



List<List<String>> parseBoardString(String boardStr) {
  final List<List<String>> board = [[],[],[],[],[],[],[],[],];
  if (boardStr.length == BoardHeight * BoardWidth) {
    for (int k=0; k<boardStr.length; k++){
      board[k ~/BoardWidth].add(boardStr[k]);
    }
  }
  return board;
}

List<Square> getInitialBoardView(List<List<String>> boardModel) {
  final List<Square> out = [];
  for (int i=0; i<boardModel.length; i++) {
    for (int j=0; j<boardModel[i].length; j++) {
      Color color = lightBlack;
      if ((i+j)%2==0) {
        color = lightWhite;
      }
      final Coord c = Coord(i,j);
      final SquareModel sq = SquareModel(boardModel[i][j], color, 1);
      final Square newSq = Square(c: c, sq: sq, key: ValueKey(Object.hash(c, sq)),);
      out.add(newSq);
    }
  }
  return out;
}

class SquareModel {
  final String pieceCode;
  final Color color;
  final double radius;
  SquareModel(this.pieceCode, this.color, this.radius);
  @override
  bool operator ==(Object other) =>
    other is SquareModel &&
    other.runtimeType == runtimeType &&
    other.pieceCode == pieceCode && other.color == color && other.radius == radius;
  @override
  int get hashCode => Object.hash(pieceCode, color, radius);
}


class Square extends StatelessWidget {
  const Square({
    super.key,
    required this.c,
    required this.sq
  });
  final Coord c;
  final SquareModel sq;

  @override
  Widget build(BuildContext context) {
    var appState = context.watch<MyAppState>();
    Size screenSize = MediaQuery.of(context).size;
    double screenWidth = screenSize.width;
    double squareW = (screenWidth-20)/8;

    ButtonStyle style = ElevatedButton.styleFrom(
      backgroundColor: sq.color,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(sq.radius),
      ),
      );
    String iconLoc = '';
    if (sq.pieceCode != Space) {
      iconLoc = 'images/${imageMap[sq.pieceCode] ?? ''}';
    }
    if (iconLoc=='') {
      return SizedBox(
          width: squareW,
          height: squareW,
          child: ElevatedButton(
            style:style,
            onPressed: () {
              appState.humanSelectButton(c);
            }, 
            child: const Text('')
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
            appState.humanSelectButton(c);
          },  
        )
      );
  }
}