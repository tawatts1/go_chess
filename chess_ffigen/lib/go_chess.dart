// lib/go_chess.dart
import 'dart:async';
import 'dart:ffi';
import 'dart:io';
import 'dart:isolate';
import 'package:ffi/ffi.dart';

import 'generated_bindings.dart';

int sum(int a, int b) => _bindings.sum(a, b);

Future<Pointer<Char>> getAiChosenMove(Pointer<Char> boardStr, int isWhite, Pointer<Char> aiName, int N) async {
  final SendPort helperIsolateSendPort = await _helperIsolateSendPort;
  final int requestId = _nextMoveRequestId++;
  final _MoveRequest request = _MoveRequest(requestId, boardStr, isWhite, aiName, N);
  final Completer<Pointer<Char>> completer = Completer<Pointer<Char>>();
  _moveRequests[requestId] = completer;
  helperIsolateSendPort.send(request);
  return completer.future;
}
//=> _bindings.GetAiChosenMove(boardStr, isWhite, aiName, N);

Pointer<Char> getBoardAfterMove(Pointer<Char> boardStr, int y1, int x1, int y2, int x2, Pointer<Char> specialChar) => _bindings.GetBoardAfterMove(boardStr, y1, x1, y2, x2, specialChar);

Pointer<Char> getNextMoves(Pointer<Char> boardStr, int y, int x) => _bindings.GetNextMoves(boardStr, y, x);

Future<int> longSum(int a, int b) async {
  final SendPort helperIsolateSendPort = await _helperIsolateSendPort;
  final int requestId = _nextSumRequestId++;
  final _SumRequest request = _SumRequest(requestId, a, b);
  final Completer<int> completer = Completer<int>();
  _sumRequests[requestId] = completer;
  helperIsolateSendPort.send(request);
  return completer.future;
}

const String _libName = 'go_chess';

/// The dynamic library in which the symbols for [NativeAddBindings] can be found.
final DynamicLibrary _dylib = () {
  if (Platform.isMacOS || Platform.isIOS) {
    return DynamicLibrary.open('$_libName.framework/$_libName');
  }
  if (Platform.isAndroid || Platform.isLinux) {
    return DynamicLibrary.open('libsum.so');
  }
  if (Platform.isWindows) {
    return DynamicLibrary.open('$_libName.dll');
  }
  throw UnsupportedError('Unknown platform: ${Platform.operatingSystem}');
}();

/// The bindings to the native functions in [_dylib].
final NativeLibrary _bindings = NativeLibrary(_dylib);


/// A request to compute `sum`.
///
/// Typically sent from one isolate to another.
class _SumRequest {
  final int id;
  final int a;
  final int b;

  const _SumRequest(this.id, this.a, this.b);
}

// request to compute a move
class _MoveRequest {
  final int id;
  final Pointer<Char> boardStr;
  final int isWhite;
  final Pointer<Char> aiName;
  final int N;

  const _MoveRequest(this.id, this.boardStr, this.isWhite, this.aiName, this.N);
}

/// A response with the result of `sum`.
///
/// Typically sent from one isolate to another.
class _SumResponse {
  final int id;
  final int result;

  const _SumResponse(this.id, this.result);
}

//response with move request result
class _MoveResponse {
  final int id;
  final Pointer<Char> result;
  
  const _MoveResponse(this.id, this.result);
}

/// Counter to identify [_SumRequest]s and [_SumResponse]s.
int _nextSumRequestId = 0;
int _nextMoveRequestId = 0;

/// Mapping from [_SumRequest] `id`s to the completers corresponding to the correct future of the pending request.
final Map<int, Completer<int>> _sumRequests = <int, Completer<int>>{};
final Map<int, Completer<Pointer<Char>>> _moveRequests = <int, Completer<Pointer<Char>>>{};

/// The SendPort belonging to the helper isolate.
Future<SendPort> _helperIsolateSendPort = () async {
  // The helper isolate is going to send us back a SendPort, which we want to
  // wait for.
  final Completer<SendPort> completer = Completer<SendPort>();

  // Receive port on the main isolate to receive messages from the helper.
  // We receive two types of messages:
  // 1. A port to send messages on.
  // 2. Responses to requests we sent.
  final ReceivePort receivePort = ReceivePort()
    ..listen((dynamic data) {
      if (data is SendPort) {
        // The helper isolate sent us the port on which we can sent it requests.
        completer.complete(data);
        return;
      }
      if (data is _SumResponse) {
        // The helper isolate sent us a response to a request we sent.
        final Completer<int> completer = _sumRequests[data.id]!;
        _sumRequests.remove(data.id);
        completer.complete(data.result);
        return;
      } else if (data is _MoveResponse) {
        final Completer<Pointer<Char>> completer = _moveRequests[data.id]!;
        _moveRequests.remove(data.id);
        completer.complete(data.result);
        return;
      }
      throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
    });

  // Start the helper isolate.
  await Isolate.spawn((SendPort sendPort) async {
    final ReceivePort helperReceivePort = ReceivePort()
      ..listen((dynamic data) {
        // On the helper isolate listen to requests and respond to them.
        if (data is _SumRequest) {
          final int result = _bindings.longSum(data.a, data.b);
          final _SumResponse response = _SumResponse(data.id, result);
          sendPort.send(response);
          return;
        } else if (data is _MoveRequest) {
          final Pointer<Char> result = _bindings.GetAiChosenMove(data.boardStr, data.isWhite, data.aiName, data.N);
          final _MoveResponse response = _MoveResponse(data.id, result);
          sendPort.send(response);
          return;
        }
        throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
      });

    // Send the port to the main isolate on which we can receive requests.
    sendPort.send(helperReceivePort.sendPort);
  }, receivePort.sendPort);

  // Wait until the helper isolate has sent us back the SendPort on which we
  // can start sending requests.
  return completer.future;
}();

