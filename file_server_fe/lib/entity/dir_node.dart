class DirNode {
  DirNode({
    required this.name,
    required this.path,
  });
  final String name;
  final String path;
  List<DirNode>? children;
}