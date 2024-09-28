import 'package:file_server_fe/common/env.dart';
import 'package:flutter/material.dart';
import 'package:file_server_fe/entity/file.dart';
import 'package:syncfusion_flutter_sliders/sliders.dart';
import 'package:video_player/video_player.dart';

class VideoPlayerPage extends StatefulWidget {
  const VideoPlayerPage({required this.file, Key? key}) : super(key: key);
  final File file;
  @override
  _VideoPlayerState createState() => _VideoPlayerState(file: file);
}

class _VideoPlayerState extends State<VideoPlayerPage> {
  _VideoPlayerState({required this.file});
  final File file;
  late VideoPlayerController _controller;
  late Future<void> _initializeVideoPlayerFuture;
  late double volume = 100;
  late double currentPosition = 0;
  late double videoLength = 1;
  bool isPlaying = false;

  @override
  void initState() {
    _controller = VideoPlayerController.networkUrl(
        Uri.parse("${Env.baseUrl}/static${file.path}"));
    _initializeVideoPlayerFuture = _controller.initialize();
    _controller.setLooping(true);
    // 启动一个定时器，每秒更新一次进度条
    _controller.addListener(() {
      if (_controller.value.isPlaying) {
        setState(() {
          currentPosition = _controller.value.position.inSeconds.toDouble();
        });
      }
    });
    super.initState();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [ 
        Align(
          alignment: Alignment.center,
          child: SizedBox(
          width: MediaQuery.of(context).size.width,
          height: MediaQuery.of(context).size.height,
          child: FutureBuilder(
            future: _initializeVideoPlayerFuture,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.done) {
                return AspectRatio(
                    aspectRatio: _controller.value.aspectRatio,
                    child: VideoPlayer(_controller));
              } else {
                return const Center(child: CircularProgressIndicator());
              }
            }),
      )),
      Align(
        alignment: Alignment.bottomCenter,
        child:
      SizedBox(
          width: MediaQuery.of(context).size.width - 100,
          height: 100,
          child: _renderPlayControl()),
          )
      ]);
  }

  String _formatSeconds(int seconds) {
    final int hour = seconds ~/ 3600;
    final int minute = seconds % 3600 ~/ 60;
    final int second = seconds % 60;
    return "$hour:$minute:$second";
  }

  Widget _renderPlayControl() {
    return Container(
      color: Color.fromARGB(50, 228, 239, 247),
      child:  Column(children: [
      Row(children: [
        SizedBox(
            width: MediaQuery.of(context).size.width - 200,
            child: SfSlider(
              value: _controller.value.position.inSeconds.toDouble(),
              min: 0,
              max: videoLength,
              onChanged: (dynamic value) {
                setState(() {
                  currentPosition = value;
                });
                _controller.seekTo(Duration(seconds: value.toInt()));
              },
            )),
        Text(
          "${_formatSeconds(_controller.value.position.inSeconds.toInt())} / ${_formatSeconds(videoLength.toInt())}",
          style: const TextStyle(fontSize: 12, color: Colors.black),
        ),
      ]),
      Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          IconButton(
            icon: const Icon(Icons.exit_to_app),
            onPressed: () {
              Navigator.pop(context);
            },
          ),
          IconButton(
            icon: Icon(isPlaying ? Icons.pause : Icons.play_arrow),
            onPressed: () {
              if (isPlaying) {
                _controller.pause();
              } else {
                _controller.play();
              }
              setState(() {
                videoLength = _controller.value.duration.inSeconds.toDouble();
                isPlaying = !isPlaying;
              });
            },
          ),
          SfSlider(
              value: volume,
              min: 0,
              max: 100,
              onChanged: (dynamic value) {
                setState(() {
                  volume = value;
                  _controller.setVolume(volume / 100);
                });
              }
          ),
          IconButton(
            icon: const Icon(Icons.fullscreen),
            onPressed: () {
              _controller.pause();

            },
          )
        ],
      ),
    ]));
  }
}
