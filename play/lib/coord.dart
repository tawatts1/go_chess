class Coord{
  final int i;
  final int j;
  const Coord(this.i, this.j);
  @override 
  String toString() {
    return '$i,$j';
  }
  @override
  bool operator ==(Object other) => 
    other is Coord &&
    other.runtimeType == runtimeType &&
    other.i == i && other.j == j;

  @override
  int get hashCode => Object.hash(i, j);
}