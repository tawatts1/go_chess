
import 'dart:ui';

const statusCheckMate = "Check Mate";
const statusStaleMate = "Stale Mate";
const statusBlackMove = "Black's Move";
const statusWhiteMove = "White's Move";

const aiDropdownList = <int>[0,1,2,3,4,5,6];

const startingBoard = "onbqkbnopppppppp00000000000000000000000000000000PPPPPPPPONBQKBNO";
const boardBeforeWhiteCheckmate = "o00qkb0op0pp0pppnp000p000000000Q00BP000000000000PPP00PbPONB0K0NO";

var white = const Color.fromARGB(255, 223, 150, 82);
var greyedWhite = const Color.fromARGB(255, 180,170,170);
var black = const Color.fromARGB(255, 116, 59, 6);
var greyedBlack = const Color.fromARGB(255, 90,75,75);
var selectedColor = const Color.fromARGB(255, 120, 0, 100);

const piecesList = [
  '0',
  'p',
  'n',
  'b',
  'r',
  'q',
  'k',
  'o',
  'a',
  'P',
  'N',
  'B',
  'R',
  'Q',
  'K',
  'O',
  'A'
];

const BoardHeight = 8;
const BoardWidth = 8;

const Space = '0';

const BlackPawn = 'p';
const BlackKnight = 'n';
const BlackBishop = 'b';
const BlackRookNC = 'r';
const BlackQueen = 'q';
const BlackKing = 'k';
const BlackRookC = 'o';
const BlackPawnEP = 'a';

const WhitePawn =  'P';
const WhiteKnight =  'N';
const WhiteBishop =  'B';
const WhiteRookNC =  'R';
const WhiteQueen =  'Q';
const WhiteKing =  'K';
const WhiteRookC =  'O';
const WhitePawnEP =  'A';


const imageMap = {
  BlackPawn: 'bp.png',
  BlackKnight: 'bn.png',
  BlackBishop: 'bb.png',
  BlackRookNC: 'br.png',
  BlackQueen: 'bq.png',
  BlackKing: 'bk.png',
  BlackRookC: 'br.png',
  BlackPawnEP: 'bp.png',
  WhitePawn: 'wp.png',
  WhiteKnight: 'wn.png',
  WhiteBishop: 'wb.png',
  WhiteRookNC: 'wr.png',
  WhiteQueen: 'wq.png',
  WhiteKing: 'wk.png',
  WhiteRookC: 'wr.png',
  WhitePawnEP: 'wp.png'
};

const Map<String, bool> whiteMap = {
  BlackPawn: false,
  BlackKnight: false,
  BlackBishop: false,
  BlackRookNC: false,
  BlackQueen: false,
  BlackKing: false,
  BlackRookC: false,
  BlackPawnEP: false,
  WhitePawn: true,
  WhiteKnight: true,
  WhiteBishop: true,
  WhiteRookNC: true,
  WhiteQueen: true,
  WhiteKing: true,
  WhiteRookC: true,
  WhitePawnEP: true,
};
