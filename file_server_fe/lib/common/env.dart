import 'package:flutter/foundation.dart';
import 'package:localstorage/localstorage.dart';

class Env {
  static String baseUrl = 'http://192.168.163.65:8080';

  static void save() {
    if (kIsWeb) {
      localStorage.setItem("baseUrl", baseUrl);
    }
  }

  static void load() {
    if (kIsWeb) {
      final String? url = localStorage.getItem("baseUrl");
      if (url != null) {
        baseUrl = url;
      }
    }
  }
}