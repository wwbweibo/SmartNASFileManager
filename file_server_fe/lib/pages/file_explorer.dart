import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/widgets/dir_tree.dart';
import 'package:file_server_fe/widgets/file.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class FileExplorer extends StatefulWidget {
  const FileExplorer({Key? key});

  @override
  _FileExplorerState createState() => _FileExplorerState();
}

class _FileExplorerState extends State<FileExplorer> {
  final String fileListUrl = "/api/v1/file";
  List<File> files = [];

  @override
  void initState() {
    super.initState();
  }

  Widget _buildGrid() => GridView.extent(
      maxCrossAxisExtent: 150,
      padding: const EdgeInsets.all(4),
      mainAxisSpacing: 4,
      crossAxisSpacing: 4,
      children: _buildGridTileList());

  List<Widget> _buildGridTileList() =>
      files.map((file) => File(fileName:  file.fileName)).toList();

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        SizedBox(
            width: 200,
            child: DirTree(
              onDirChange: (p0) => {
                listFile(p0).then((value) {
                  setState(() {
                    this.files = value;
                  });
                })
              },
            )),
        Expanded(child: _buildGrid())
      ],
    );
  }

  Future<List<File>> listFile(String path) async {
    final response =
        await http.get(Uri.parse("${Env.baseUrl}$fileListUrl?path=$path"));
    if (response.statusCode == 200) {
      if (response.body == "null") {
        return List<File>.empty();
      }
      final List<dynamic> resp = json.decode(response.body);
      return resp.map((e) => File(fileName:  e['name'])).toList();
    } else {
      throw Exception('Failed to load file list');
    }
  }
}
