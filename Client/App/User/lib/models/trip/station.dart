import 'Province.dart';

class Station {
  final int id;
  final String name;
  final String address;
  final Province province;

  Station({
    required this.id,
    required this.name,
    required this.address,
    required this.province,
  });

  factory Station.fromJson(Map<String, dynamic> json) {
    return Station(
      id: json['id'],
      name: json['name'],
      address: json['address'],
      province: Province.fromJson(json['province']),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'address': address,
    'province': province.toJson(),
  };
}
