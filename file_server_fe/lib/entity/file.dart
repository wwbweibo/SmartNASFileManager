class File {
  const File({
    required this.name,
    required this.path,
    required this.size,
    required this.type,
    required this.group,
  });
  final String name;
  final String path;
  final int size;
  final String type;
  final String group;
}
