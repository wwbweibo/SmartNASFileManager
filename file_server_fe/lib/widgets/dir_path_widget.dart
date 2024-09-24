import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class DirPathWidget extends StatefulWidget{
  final String path;
  final Function(String) onPathChange;
  const DirPathWidget({Key? key, required this.path, required this.onPathChange}) : super(key: key);
  @override
  _DirPathState createState() => _DirPathState(path);
}

class _DirPathState extends State<DirPathWidget> {
  _DirPathState(String path) {
    super.initState();
    if (path.endsWith("/")) {
      path = path.substring(0, path.length - 1);
    }
    path.split("/").skip(1).forEach((element) {
      paths.add(element);
    });
    selectedPath = path;
    if (selectedPath.isEmpty) {
      selectedPath = "/";
    }
  }

  final List<String> paths = ["/"];
  String selectedPath = "/";

  @override
  Widget build(BuildContext context) {
    List<ButtonSegment<String>> segments = List<ButtonSegment<String>>.empty(growable: true);
    String path = "";
    for (int i = 0; i < paths.length; i++) {
      segments.add(ButtonSegment<String>(
        value: path + paths[i],
        label: Text(paths[i]),
      ));
      if (paths[i] == "/") {
        path = "/";
      } else {
        path += paths[i] + "/";
      }
    }
    return SegmentedButton<String>(
      style: SegmentedButton.styleFrom(
        backgroundColor: Colors.green,
        foregroundColor: Colors.white,
        selectedForegroundColor: Colors.white,
        selectedBackgroundColor: Colors.green,
        side: const BorderSide(width: 0, style: BorderStyle.none),
        shape: const LinearBorder()
      ),
      segments: segments,
      selected: {selectedPath},
      onSelectionChanged: (p0) => {
        if(!p0.first.endsWith("/"))
          {widget.onPathChange("${p0.first}/")}
        else
          widget.onPathChange(p0.first),
      },
    );
  }
}