import 'dart:convert';

import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/entity/dir_node.dart';
import 'package:flutter/material.dart';
import 'package:flutter_fancy_tree_view/flutter_fancy_tree_view.dart';
import 'package:http/http.dart' as http;

class DirTree extends StatefulWidget {
  final Function(String) onDirChange;
  DirTree({Key? key, required this.onDirChange}) : super(key: key);
  @override
  _DirTreeState createState() => _DirTreeState();
}

class _DirTreeState extends State<DirTree> {
  _DirTreeState() {
    super.initState();
    fetchDirTree();
  }

  final treeController = TreeController<DirNode>(
    roots: [DirNode(name: "/", path: "/")],
    childrenProvider: (DirNode node) => node.children ?? [],
  );

  @override
  Widget build(BuildContext context) {
    return Container( 
      margin: const EdgeInsets.all(8),
      width: 100,
      child:  AnimatedTreeView<DirNode>(
        treeController: treeController,
        nodeBuilder: (BuildContext context, TreeEntry<DirNode> entry) {
          return InkWell(
            onTap: () =>  {
              // trigger dir change event
              widget.onDirChange(entry.node.path),
              treeController.toggleExpansion(entry.node)
            },
            child: TreeIndentation(
              entry: entry,
              child: Text(entry.node.name),
            ),
          );
        })
    );
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
      node.children = List.from((resp['children'] as List<dynamic>).map((e) => _parseDirTree(e)).toList());
    }
    return node;
  }
}
