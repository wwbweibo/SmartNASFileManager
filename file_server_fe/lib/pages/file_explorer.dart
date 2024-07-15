import 'package:file_server_fe/widgets/file.dart';
import 'package:flutter/material.dart';

class FileExplorer extends StatefulWidget {
  const FileExplorer({Key? key});

  @override
  _FileExplorerState createState() => _FileExplorerState();
}

class _FileExplorerState extends State<FileExplorer> {
  Widget _buildGrid() => GridView.extent(
      maxCrossAxisExtent: 150,
      padding: const EdgeInsets.all(4),
      mainAxisSpacing: 4,
      crossAxisSpacing: 4,
      children: _buildGridTileList(30)
  );

  List<Widget> _buildGridTileList(int count) => List.generate(
      count, (i) => Container(child: File( "File $i"))
  );

  @override
  Widget build(BuildContext context) {
    return _buildGrid();
  }
}