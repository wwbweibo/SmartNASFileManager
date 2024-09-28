import 'package:flutter/material.dart';
import 'package:file_server_fe/entity/file.dart';

class VideoPlayer extends StatefulWidget {
  const VideoPlayer({required this.file, Key? key}) : super(key: key);
  final File file;
  @override
  _VideoPlayerState createState() =>  _VideoPlayerState(file: file );
}

class _VideoPlayerState extends State<VideoPlayer> {
  _VideoPlayerState({required this.file});
  final File file;
  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Text('Video Player'),
    );
  }
}