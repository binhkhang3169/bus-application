// lib/models/province.dart

class Province {
  final int id;
  final String name;
  // Add other fields like code, region, etc., if your API provides them

  Province({
    required this.id,
    required this.name,
  });

  factory Province.fromJson(Map<String, dynamic> json) {
    return Province(
      id: json['id'] as int, // Assuming 'id' is the key from API
      name: json['name'] as String, // Assuming 'name' is the key from API
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
      };

  @override
  String toString() {
    return 'Province{id: $id, name: $name}';
  }
}