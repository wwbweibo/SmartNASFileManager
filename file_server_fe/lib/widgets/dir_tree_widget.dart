import 'dart:convert';

import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/entity/dir_node.dart';
import 'package:flutter/material.dart';
import 'package:flutter_fancy_tree_view/flutter_fancy_tree_view.dart';
import 'package:http/http.dart' as http;

class DirTreeWidget extends StatefulWidget {
  final String selectedPath;
  final Function(String) onDirChange;
  DirTreeWidget(
      {Key? key, required this.selectedPath, required this.onDirChange})
      : super(key: key);
  @override
  _DirTreeState createState() => _DirTreeState();
}

class _DirTreeState extends State<DirTreeWidget> {
  _DirTreeState() {
    super.initState();
    fetchDirTree();
  }

  final treeController = TreeController<DirNode>(
    roots: [DirNode(name: "/", path: "/")],
    childrenProvider: (DirNode node) => node.children ?? [],
    defaultExpansionState: true,
  );

  @override
  Widget build(BuildContext context) {
    return Container(
        margin: const EdgeInsets.all(8),
        width: 100,
        child: AnimatedTreeView<DirNode>(
            treeController: treeController,
            nodeBuilder: (BuildContext context, TreeEntry<DirNode> entry) {
              return InkWell(
                onTap: () => {treeController.toggleExpansion(entry.node)},
                onDoubleTap: () => {
                  treeController.expand(entry.node),
                  widget.onDirChange(entry.node.path),
                },
                child: TreeIndentation(
                  entry: entry,
                  child: Text(entry.node.name),
                ),
              );
            }));
  }

  Future fetchDirTree() async {
    final response = await http.get(Uri.parse("${Env.baseUrl}/api/v1/dir"));
    if (response.statusCode == 200) {
      final Map<String, dynamic> resp = json.decode(response.body);
      treeController.roots = [_parseDirTree(resp)];
    } else {
      throw Exception('Failed to load dir tree');
    }
  }

  DirNode _parseDirTree(Map<String, dynamic> resp) {
    final DirNode node = DirNode(name: resp['name'], path: resp['path']);
    if (resp['children'] != null) {
      node.children = List.from((resp['children'] as List<dynamic>)
          .map((e) => _parseDirTree(e))
          .toList());
    }
    return node;
  }
}
