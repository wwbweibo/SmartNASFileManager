import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:file_server_fe/entity/file.dart';
import 'package:file_server_fe/common/env.dart';
import 'package:flutter/widgets.dart';
import 'package:syncfusion_flutter_sliders/sliders.dart';
import 'package:video_player/video_player.dart';

class MediaPlayerPage extends StatefulWidget {
  const MediaPlayerPage({required this.files, this.selectedIndex = 0, Key? key})
      : super(key: key);
  final List<File> files;
  final int selectedIndex;

  @override
  _MediaPlayerPageState createState() =>
      _MediaPlayerPageState(files: files, selectedIndex: selectedIndex);
}

class _MediaPlayerPageState extends State<MediaPlayerPage> {
  _MediaPlayerPageState({required this.files, this.selectedIndex = 0});
  final List<File> files;
  int selectedIndex = 0;

  // 图片展示相关的控制属性
  double scale = 1.0;
  double positionXBeforeDrag = 0;
  double positionYBeforeDrag = 0;
  double imagePositionX = 0;
  double imagePositionY = 0;

  // 视频播放相关的控制属性
  late VideoPlayerController _controller;
  late Future<void> _initializeVideoPlayerFuture;
  late double volume = 100;
  late double currentPosition = 0;
  late double videoLength = 1;
  bool isPlaying = false;

  Widget _renderImage(String url) {
    return Image.network(
      url,
    );
  }

  Widget _renderMediaList() {
    final ScrollController rowController =
        ScrollController(initialScrollOffset: calcScrollOffset());
    return Container(
        margin: const EdgeInsets.all(8),
        height: 110,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          color: const Color.fromARGB(100, 240, 240, 240),
        ),
        child: Scrollbar(
          controller: rowController,
          child: ListView(
            scrollDirection: Axis.horizontal,
            controller: rowController,
            children: files.map((item) {
              return GestureDetector(
                onTap: () {
                  _onSelectedMediaChanged(item, rowController);
                },
                child: Container(
                  margin: const EdgeInsets.all(8),
                  child: Container(
                      padding: const EdgeInsets.all(5),
                      width: 100,
                      height: 100,
                      color: selectedIndex == files.indexOf(item)
                          ? Colors.blue
                          : Colors.grey,
                      child: item.group == "image"
                          ? Image.network(
                              _formatImageUrl(item.path, useThumbnail: true),
                              width: 100,
                              height: 100,
                            )
                          : const Icon(Icons.video_file)),
                ),
              );
            }).toList(),
          ),
        ));
  }

  Widget _renderImageViewer() {
    return GestureDetector(
        onScaleStart: (details) => {
              setState(() {
                positionXBeforeDrag = details.focalPoint.dx;
                positionYBeforeDrag = details.focalPoint.dy;
              })
            },
        onScaleUpdate: (details) => {
              setState(() {
                // scale = details.scale;
                if (details.scale != 1) {
                  scale = details.scale;
                }
                imagePositionX = imagePositionX -
                    (positionXBeforeDrag - details.focalPoint.dx) / scale;
                imagePositionY = imagePositionY -
                    (positionYBeforeDrag - details.focalPoint.dy) / scale;
                positionXBeforeDrag = details.focalPoint.dx;
                positionYBeforeDrag = details.focalPoint.dy;
              })
            },
        onDoubleTap: () => {
              setState(() {
                scale = 1;
                imagePositionX = 0;
                imagePositionY = 0;
              })
            },
        onTap: () => {Navigator.of(context).pop()},
        child: SizedBox(
            height: MediaQuery.of(context).size.height - 150,
            width: MediaQuery.of(context).size.width,
            child: Stack(
              children: [
                Align(
                  alignment: Alignment.center,
                  child: Transform.scale(
                    scale: scale,
                    child: Transform.translate(
                      offset: Offset(imagePositionX, imagePositionY),
                      child: _renderImage(
                          _formatImageUrl(files[selectedIndex].path)),
                    ),
                  ),
                ),
              ],
            )));
  }

  Widget _renderVideoViewer() {
    return Stack(children: [
      Align(
          alignment: Alignment.topCenter,
          child: SizedBox(
              height: MediaQuery.of(context).size.height - 150,
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
                  }))),
      Align(
          alignment: Alignment.bottomCenter,
          child: SizedBox(
              width: MediaQuery.of(context).size.width - 100,
              height: 100,
              child: _renderVideoPlayControl()))
    ]);
  }

  Widget _renderVideoPlayControl() {
    return Container(
        color: Color.fromARGB(50, 228, 239, 247),
        child: Column(children: [
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
                    videoLength =
                        _controller.value.duration.inSeconds.toDouble();
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
                  }),
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

  Widget _renderViewer() {
    return SizedBox(
      height: MediaQuery.of(context).size.height - 150,
      child: files[selectedIndex].group == "image"
          ? _renderImageViewer()
          : _renderVideoViewer(),
    );
  }

  Widget _renderPage() {
    return Stack(
      children: [
        Align(
          alignment: Alignment.topCenter,
          child: _renderViewer(),
        ),
        Align(
          alignment: Alignment.bottomCenter,
          child: _renderMediaList(),
        ),
      ],
    );
  }
  
  @override 
  void initState() {
    super.initState();
    if (files[selectedIndex].group == "video") {
      _initVideoController();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: const Color.fromARGB(255, 10, 10, 10),
      child: _renderPage(),
    );
  }

  double calcScrollOffset() {
    var itemWidth = 116;
    var screenWidth = MediaQuery.of(context).size.width;
    var offset =
        (itemWidth * selectedIndex) - (screenWidth / 2) + (itemWidth / 2);
    // 如果 offset 小于 0，说明当前选中的图片在屏幕左侧，需要滚动到最左侧
    if (offset < 0) {
      offset = 0;
    }
    // 如果 offset 大于最大值，说明当前选中的图片在屏幕右侧，需要滚动到最右侧
    if (offset > (files.length * itemWidth - screenWidth)) {
      offset = files.length * itemWidth - screenWidth;
    }
    return offset;
  }

  String _formatImageUrl(String url, {bool useThumbnail = false}) {
    if (useThumbnail) {
      return "${Env.baseUrl}/cache$url";
    } else {
      return "${Env.baseUrl}/static/$url";
    }
  }

  void _initVideoController() {
    _controller = VideoPlayerController.networkUrl(
        Uri.parse("${Env.baseUrl}/static${files[selectedIndex].path}"));
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
  }

  void _onSelectedMediaChanged(File file, ScrollController rowController) {
    var currentFile = files[selectedIndex];
    if (currentFile.group == "video") {
      // 当前选中的是视频文件，需要清除播放状态，释放资源
      _controller.dispose();
    }
    setState(() {
      selectedIndex = files.indexOf(file);
      // 视频控制状态清除
      isPlaying = false;
      currentPosition = 0;
      videoLength = 1;
    });
    if (file.group == "video") {
      // 新选中的是视频文件，需要初始化视频播放器
      _initVideoController();
    }
    rowController.animateTo(
      calcScrollOffset(),
      duration: const Duration(milliseconds: 300),
      curve: Curves.easeInOut,
    );
  }

  String _formatSeconds(int seconds) {
    final int hour = seconds ~/ 3600;
    final int minute = seconds % 3600 ~/ 60;
    final int second = seconds % 60;
    return "$hour:$minute:$second";
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}
