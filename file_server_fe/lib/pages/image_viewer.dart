import 'dart:core';
import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/entity/file.dart';
import 'package:flutter/material.dart';
import 'package:syncfusion_flutter_sliders/sliders.dart';

class ImageViewer extends StatefulWidget {
  const ImageViewer({required this.images, this.selectedIndex = 0, Key? key})
      : super(key: key);
  final List<File> images;
  final int selectedIndex;
  @override
  State<ImageViewer> createState() =>
      _ImageViewerState(images, selectedIndex: selectedIndex);
}

class _ImageViewerState extends State<ImageViewer> {
  _ImageViewerState(this.images, {int? selectedIndex}) {
    this.selectedIndex = selectedIndex ?? 0;
  }

  List<File> images = [];
  int selectedIndex = 0;
  double scale = 1.0;
  late double positionXBeforeDrag;
  late double positionYBeforeDrag;
  double imagePositionX = 0;
  double imagePositionY = 0;

  @override
  Widget build(BuildContext context) {
    return Container(
      color: const Color.fromARGB(255, 10, 10, 10),
      child: _renderViewer(),
    );
  }

  Widget _renderImage(String url) {
    return Image.network(
      url,
    );
  }

  double calcScrollOffset() {
    var itemWidth = 116;
    var screenWidth = MediaQuery.of(context).size.width;
    var offset =  (itemWidth * selectedIndex) - (screenWidth / 2) + (itemWidth / 2);
    // 如果 offset 小于 0，说明当前选中的图片在屏幕左侧，需要滚动到最左侧
    if (offset < 0) {
      offset = 0;
    }
    // 如果 offset 大于最大值，说明当前选中的图片在屏幕右侧，需要滚动到最右侧
    if (offset > (images.length * itemWidth - screenWidth)) {
      offset = images.length * itemWidth - screenWidth;
    }
    return offset;
  }

  Widget _renderImageList() {
    final ScrollController rowController = ScrollController(initialScrollOffset: calcScrollOffset());
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
            children: images.map((item) {
              return GestureDetector(
                onTap: () {
                  setState(() {
                    selectedIndex = images.indexOf(item);
                  });
                  rowController.animateTo(
                    calcScrollOffset(),
                    duration: const Duration(milliseconds: 300),
                    curve: Curves.easeInOut,
                  );
                },
                child: Container(
                  margin: const EdgeInsets.all(8),
                  child: Container(
                      padding: const EdgeInsets.all(5),
                      width: 100,
                      height: 100,
                      color: selectedIndex == images.indexOf(item)
                          ? Colors.blue
                          : Colors.grey,
                      child: Image.network(
                        _formatImageUrl(item.path, useThumbnail: true),
                        width: 100,
                        height: 100,
                      )),
                ),
              );
            }).toList(),
          ),
        ));
  }

  Widget _renderViewer() {
    return Stack(children: [
      GestureDetector(
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
                            _formatImageUrl(images[selectedIndex].path)),
                      ),
                    ),
                  ),
                ],
              ))),
      Align(alignment: Alignment.bottomCenter, child: _renderImageList()),
      Align(
          alignment: const Alignment(0.9, 0.5),
          child: SizedBox(
            height: 200,
            width: 10,
            child: SfSlider.vertical(
                value: scale,
                min: 1,
                max: 5,
                onChanged: (dynamic value) {
                  setState(() {
                    scale = value;
                  });
                }),
          )),
    ]);
  }

  String _formatImageUrl(String url, {bool useThumbnail = false}) {
    if (useThumbnail) {
      return "${Env.baseUrl}/cache$url";
    } else {
      return "${Env.baseUrl}/static/$url";
    }
  }

  void _updateSelectedIndex(int index) {
    var _index = selectedIndex + index;
    if (_index < 0) {
      _index = images.length - 1;
    }
    if (_index >= images.length) {
      _index = 0;
    }
    setState(() {
      selectedIndex = _index;
    });
  }
}
