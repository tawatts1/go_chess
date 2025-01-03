import 'package:flutter/material.dart';
import 'dart:async';

import 'package:go_chess/go_chess.dart' as go_chess;

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  late int sumResult;
  late Future<int> longSumResult;
  final myController1 = TextEditingController();
  final myController2 = TextEditingController();
  @override
  void dispose() {
    myController1.dispose();
    myController2.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    sumResult = go_chess.sum(0, 0);
    longSumResult = go_chess.longSum(0, 0);
  }



  void _doCalc() {
    setState(() {
      var n1 = int.parse(myController1.text);
      var n2 = int.parse(myController2.text);
      longSumResult = Future<int>(() { return go_chess.longSum(n1, n2);});
      sumResult = go_chess.sum(n1,n2);
      //longSumResult = 
    });
  }

  @override
  Widget build(BuildContext context) {
    const textStyle = TextStyle(fontSize: 25);
    const spacerSmall = SizedBox(height: 10);
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Native Packages'),
        ),
        body: SingleChildScrollView(
          child: Container(
            padding: const EdgeInsets.all(10),
            child: Column(
              children: [
                const Text(
                  'This calls a native function through FFI that is shipped as source in the package. '
                  'The native code is built as part of the Flutter Runner build.',
                  style: textStyle,
                  textAlign: TextAlign.center,
                ),
                spacerSmall,
                Text('First Number:'),
                TextField(
                  controller: myController1,
                ),
                Text('Second number: '),
                TextField(
                  controller: myController2,
                ),
                ElevatedButton(
                  onPressed: () {
                    _doCalc();
                    print('test $sumResult');
                  },
                  child: Text('do calculation'),
                ),
                Text(
                  'The sum is $sumResult',
                  style: textStyle,
                  textAlign: TextAlign.center,
                ),
                spacerSmall,
                FutureBuilder<int>(
                  future: longSumResult,
                  builder: (BuildContext context, AsyncSnapshot<int> value) {
                    final displayValue =
                        (value.hasData) ? value.data : 'loading';
                    return Text(
                      'The longSum is $displayValue',
                      style: textStyle,
                      textAlign: TextAlign.center,
                    );
                  },
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
