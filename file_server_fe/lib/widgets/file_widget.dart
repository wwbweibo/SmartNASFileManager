import 'package:file_server_fe/common/env.dart';
import 'package:flutter/material.dart';
import 'package:file_server_fe/entity/file.dart';

class FileWidget extends StatefulWidget {
  final File file;
  final Function(File) onFileClick;
  const FileWidget({Key? key, required this.file, required this.onFileClick})
      : super(key: key);
  @override
  _FileState createState() => _FileState();
}

class _FileState extends State<FileWidget> {
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onDoubleTap: () => widget.onFileClick(widget.file),
        child: Container(
          color: const Color.fromARGB(255, 228, 239, 247),
          padding: const EdgeInsets.all(8),
          margin: const EdgeInsets.fromLTRB(8,8,8,0),
          child: Column (children: [
            Container(
              margin: const EdgeInsets.all(2),
              child: () {
                if (widget.file.group == "dir") {
                  return const Icon(Icons.folder, size: 85,);
                } else if (widget.file.group == "image") {
                  return Image.network("${Env.baseUrl}/cache${widget.file.path}", height: 85,); 
                } else {
                  return const Icon(Icons.file_open, size: 85);
                }
              }(),
            ),
            Text(widget.file.name, 
              style: const TextStyle(
                fontSize: 12,
              ),
              softWrap: false ,
              overflow: TextOverflow.clip,
            ),
          ],
        ),
      )
    );
  }
}
