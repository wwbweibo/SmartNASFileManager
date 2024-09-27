import 'package:file_server_fe/common/env.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

class ImageViewer extends StatelessWidget {
  ImageViewer(List<String> images, {int? selectedIndex,  Key? key}) : super(key: key);
  List<String> images = [];
  int selectedIndex = 0;
  double scale = 1.0;

  Image _renderImage(String url) {
    return Image.network(
      url,
      scale: scale,
    );
  }

  Widget _renderImageList() {
    return ListView.builder(
      itemCount: images.length,
      itemBuilder: (context, index) {
        return _renderImage(_formatImageUrl(images[index], useThumbnail: true));
      },
    );
  }

  Widget _renderViewer() {
    return Column(
      children: [
        _renderImage(_formatImageUrl(images[selectedIndex])),
        _renderImageList(),
      ]
    );
  }

  @override
  Widget build(BuildContext context) {
    return const Center(
      
    );
  }

  String _formatImageUrl(String url, {bool useThumbnail = false}) {
    if (useThumbnail) {
      return "${Env.baseUrl}/cache$url";
    } else {
      return "${Env.baseUrl}/static/$url";
    }
  }
}