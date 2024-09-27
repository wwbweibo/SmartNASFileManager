import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:localstorage/localstorage.dart';

class Env {
  static String baseUrl = 'http://192.168.163.65:8080';

  static void save() {
    if (kIsWeb) {
      localStorage.setItem("baseUrl", baseUrl);
    } else {
      // 写入文件
      File file = File("env.json");
      file.writeAsStringSync('{"baseUrl": "$baseUrl"}');
    }
  }

  static void load() {
    if (kIsWeb) {
      final String? url = localStorage.getItem("baseUrl");
      if (url != null) {
        baseUrl = url;
      }
    } else {
      // 读取文件
      File file = File("env.json");
      if (file.existsSync()) {
        final String content = file.readAsStringSync();
        final Map<String, dynamic> env = json.decode(content);
        baseUrl = env['baseUrl'];
      }
    }
  }
}