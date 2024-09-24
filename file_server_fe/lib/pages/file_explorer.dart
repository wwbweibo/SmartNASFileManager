import 'dart:developer';

import 'package:easy_image_viewer/easy_image_viewer.dart';
import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/widgets/dir_path_widget.dart';
import 'package:file_server_fe/widgets/dir_tree_widget.dart';
import 'package:file_server_fe/widgets/file_widget.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:file_server_fe/entity/file.dart';

class FileExplorer extends StatefulWidget {
  const FileExplorer({Key? key});

  @override
  _FileExplorerState createState() => _FileExplorerState();
}

class _FileExplorerState extends State<FileExplorer> {
  final String fileListUrl = "/api/v1/file";
  List<File> files = [];
  String currentSelectedPath = "/";
  bool showImageViewer = false;
  late MultiImageProvider imageProvider;

  Widget _buildGrid() => GridView.extent(
      maxCrossAxisExtent: 150,
      // mainAxisSpacing: 4,
      // crossAxisSpacing: 4,
      children: _buildGridTileList()
    );

  List<Widget> _buildGridTileList() =>
      files.map(
        (file) => FileWidget(
          file: file,
          onFileClick: onFileClick,
        ),
      ).toList();

  DirTreeWidget _buildDirTree() => DirTreeWidget(key: Key(currentSelectedPath), selectedPath: currentSelectedPath, onDirChange: onPathChanged);
  
  DirPathWidget _buildDirPath() => DirPathWidget(key: Key(currentSelectedPath), path: currentSelectedPath, onPathChange: onPathChanged);

  @override
  void initState() {
    super.initState();
    onPathChanged("/");
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        SizedBox(
            width: 200,
            child: _buildDirTree(),
        ),
        Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(
                height: 50,
                child: _buildDirPath(),
              ),
              Expanded(child: _buildGrid())
            ],
        )
        )
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
      return resp.map((e) => File(name: e['name'], path: e['path'], size: e['size'], type: e['type'], group: e['group'])).toList();
    } else {
      throw Exception('Failed to load file list');
    }
  }

  onPathChanged(String path) {
    log("onPathChanged: $path");
    listFile(path).then((value) {
      setState(() {
        this.files = value;
        this.currentSelectedPath = path;
        // this.treeWidget = _buildDirTree();
      });
    });
  }

  onFileClick(File file) {
    if (file.group == "dir") {
      onPathChanged(file.path);
    }
    if (file.group == "image") {
      // 弹出图片浏览
      var initialIndex = 0;
      var index = 0;
      List<NetworkImage> images = [];
      files
        .where((item) => item.group == "image")
        .forEach((item) {
          if (item.path == file.path) {
            initialIndex = index;
          }
          index = index + 1;
          images.add(NetworkImage("${Env.baseUrl}/static${item.path}"));
        });
      MultiImageProvider multiImageProvider = MultiImageProvider(images.toList(), initialIndex: initialIndex);
      showImageViewerPager(context, multiImageProvider);
      // setState(() {
      //   this.imageProvider = multiImageProvider;
      //   this.showImageViewer = true;
      // });
    }
  }
}
