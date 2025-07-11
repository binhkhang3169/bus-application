class AddressInfo {
  final int id;
  final String name;

  AddressInfo({required this.id, required this.name});

  factory AddressInfo.fromJson(Map<String, dynamic> json) {
    return AddressInfo(id: json['id'], name: json['name']);
  }
}
