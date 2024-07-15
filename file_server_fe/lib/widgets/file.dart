import 'package:flutter/material.dart';

class File extends StatefulWidget {
  File(String fileName, {Key? key}) {
    this.fileName = fileName;
  }
  var fileName = "file_name";
  @override
  _FileState createState() => _FileState();
}

class _FileState extends State<File> {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          width: 150,
          height: 150,
          color: Colors.blue,
          child: const Icon(Icons.file_open),
        ),
        Text(widget.fileName)
      ],
    );   
  }
}