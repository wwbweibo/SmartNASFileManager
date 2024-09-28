import 'package:file_server_fe/common/env.dart';
import 'package:file_server_fe/entity/file.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'dart:developer';

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
      child: Container(color: const Color.fromARGB(255, 10,10,10), child:  _renderViewer())
    );
  }

  Widget _renderImage(String url) {
    return Image.network(
      url,
    );
  }

  Widget _renderImageList() {
    // 仅展示当前选中图片的前后3张，供7张图片
    List<File> showShowImages = [];
    if (images.length <= 7) {
      showShowImages = images;
    } else if (selectedIndex <= 3) {
      showShowImages = images.sublist(0, 6);
    } else if (selectedIndex > 3 && selectedIndex < images.length - 3) {
      showShowImages = images.sublist(selectedIndex - 3, selectedIndex + 3);
    } else {
      showShowImages = images.sublist(images.length - 6, images.length);
    }
    return Container(
        margin: const EdgeInsets.all(8),
        color: Color.fromARGB(100, 240,240,240),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: showShowImages.map((item) {
            return GestureDetector(
              onTap: () {
                setState(() {
                  selectedIndex = images.indexOf(item);
                });
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
        ));
  }

  Widget _renderViewer() {
    return KeyboardListener(focusNode: FocusNode(
      onKeyEvent: (node, event) {
        log("event: ${event.logicalKey}");
        if (event.logicalKey == LogicalKeyboardKey.arrowLeft) {
            _updateSelectedIndex(-1);
        } else if (event.logicalKey == LogicalKeyboardKey.arrowRight) {
            _updateSelectedIndex(1);
        }
        return KeyEventResult.handled;
      },
    ), 
    child: Stack(children: [
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
            imagePositionX = imagePositionX - (positionXBeforeDrag - details.focalPoint.dx) / scale;
            imagePositionY = imagePositionY - (positionYBeforeDrag - details.focalPoint.dy) / scale;
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
        onTap: () => {
          Navigator.of(context).pop()
        },
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
          )
      )),
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
    ]));
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
