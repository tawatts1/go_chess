import 'dart:developer';

const humanName = "Human";
const aiName = "Ai";


enum PlayStatus {play, pause, undefined}
class PlayerState {
  // Info about who the players are. 
  bool isBlackAi = true;
  bool isWhiteAi = false;
  int aiDropdownDepth = 1;
  PlayStatus playPauseStatus = PlayStatus.undefined;
  @override
  String toString() {
    return "$isBlackAi,$isWhiteAi,$aiDropdownDepth";
  }
  void loadFromString(String stateStr){
    List<String> playerState = stateStr.split(",");
    if (playerState.length != 3){
      log("Tried to load player state from a bad string: $stateStr");
    } else {
      try {
        isBlackAi = playerState[0].toLowerCase() == "true";
        isWhiteAi = playerState[1].toLowerCase() == "true";
        aiDropdownDepth = int.parse(playerState[2]);
      } catch (e) {
        log("Failed to parse player state: $stateStr");
      }
    }
  }
  String getPlayerName(bool isWhite) {
    String out = "";
    if (isWhite){
      out = isWhiteAi ? aiName : humanName;
    } else {
      out = isBlackAi ? aiName : humanName;
    }
    return out;
  }
  bool isBothAi() {
    return isWhiteAi && isBlackAi;
  }
  bool isAtLeastOneAi() {
    return isWhiteAi || isBlackAi;
  }
  bool isNeitherAi() {
    return !isWhiteAi && !isBlackAi;
  }
}