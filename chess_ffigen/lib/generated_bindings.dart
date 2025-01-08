// ignore_for_file: always_specify_types
// ignore_for_file: camel_case_types
// ignore_for_file: non_constant_identifier_names
// ignore_for_file: unused_field
// ignore_for_file: unused_element

// AUTO GENERATED FILE, DO NOT EDIT.
//
// Generated by `package:ffigen`.
// ignore_for_file: type=lint
import 'dart:ffi' as ffi;

/// Bindings to `src/sum.h`.
class NativeLibrary {
  /// Holds the symbol lookup function.
  final ffi.Pointer<T> Function<T extends ffi.NativeType>(String symbolName)
      _lookup;

  /// The symbols are looked up in [dynamicLibrary].
  NativeLibrary(ffi.DynamicLibrary dynamicLibrary)
      : _lookup = dynamicLibrary.lookup;

  /// The symbols are looked up with [lookup].
  NativeLibrary.fromLookup(
      ffi.Pointer<T> Function<T extends ffi.NativeType>(String symbolName)
          lookup)
      : _lookup = lookup;

  ffi.Pointer<ffi.Char> GetNextMoves(
    ffi.Pointer<ffi.Char> boardStr,
    int y,
    int x,
  ) {
    return _GetNextMoves(
      boardStr,
      y,
      x,
    );
  }

  late final _GetNextMovesPtr = _lookup<
      ffi.NativeFunction<
          ffi.Pointer<ffi.Char> Function(
              ffi.Pointer<ffi.Char>, ffi.Int, ffi.Int)>>('GetNextMoves');
  late final _GetNextMoves = _GetNextMovesPtr.asFunction<
      ffi.Pointer<ffi.Char> Function(ffi.Pointer<ffi.Char>, int, int)>();

  int sum(
    int a,
    int b,
  ) {
    return _sum(
      a,
      b,
    );
  }

  late final _sumPtr =
      _lookup<ffi.NativeFunction<ffi.Int Function(ffi.Int, ffi.Int)>>('sum');
  late final _sum = _sumPtr.asFunction<int Function(int, int)>();

  int longSum(
    int a,
    int b,
  ) {
    return _longSum(
      a,
      b,
    );
  }

  late final _longSumPtr =
      _lookup<ffi.NativeFunction<ffi.Int Function(ffi.Int, ffi.Int)>>(
          'longSum');
  late final _longSum = _longSumPtr.asFunction<int Function(int, int)>();

  void enforce_binding() {
    return _enforce_binding();
  }

  late final _enforce_bindingPtr =
      _lookup<ffi.NativeFunction<ffi.Void Function()>>('enforce_binding');
  late final _enforce_binding =
      _enforce_bindingPtr.asFunction<void Function()>();
}

final class max_align_t extends ffi.Opaque {}

final class _GoString_ extends ffi.Struct {
  external ffi.Pointer<ffi.Char> p;

  @ptrdiff_t()
  external int n;
}

typedef ptrdiff_t = ffi.Long;
typedef Dartptrdiff_t = int;

final class GoInterface extends ffi.Struct {
  external ffi.Pointer<ffi.Void> t;

  external ffi.Pointer<ffi.Void> v;
}

final class GoSlice extends ffi.Struct {
  external ffi.Pointer<ffi.Void> data;

  @GoInt()
  external int len;

  @GoInt()
  external int cap;
}

typedef GoInt = GoInt64;
typedef GoInt64 = ffi.LongLong;
typedef DartGoInt64 = int;

const int NULL = 0;
