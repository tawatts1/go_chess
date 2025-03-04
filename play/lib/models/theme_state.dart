import 'dart:developer';
import 'dart:ui';

import 'package:play/constants.dart';

class ThemeState {
  bool isDarkTheme = true;

  @override
  String toString() {
    log("Theme toString is $isDarkTheme");
    return "$isDarkTheme";
  }
  void loadFromString(String stateStr) {
    isDarkTheme = stateStr.toLowerCase() == "true";
  }
  Color getSquareColor(bool isLightSquare){
    if (isDarkTheme){
      if (isLightSquare){
        return darkWhite;
      } else {
        return darkBlack;
      }
    } else {
      if (isLightSquare) {
        return lightWhite;
      } else {
        return lightBlack;
      }
    }
  }
}