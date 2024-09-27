import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:file_server_fe/common/env.dart';

class Setting extends StatefulWidget {
  const Setting({Key? key}) : super(key: key);

  @override
  _SettingState createState() => _SettingState();
}

class _SettingState extends State<Setting> {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          margin: const EdgeInsets.all(8),
          child: 
        SizedBox(
            width: 260,
            child: TextField(
              key: const Key( "serverAddress"),
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: '请输入服务器地址',
              ),
              controller: TextEditingController(text: Env.baseUrl),
              onChanged: (text) {
                Env.baseUrl = text;
              },
            )
          )
        ),
        Container(
          margin: const EdgeInsets.all(8),
          child: SizedBox(
          width: 100,
          child: ElevatedButton(
            onPressed: () {
              Env.save();
              // 提示保存成功
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('保存成功')),
              );
            },
            child: const Text('保存'),
          ),
        )
        )
      ],
    );
  }
}
