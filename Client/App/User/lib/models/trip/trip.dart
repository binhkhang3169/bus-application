// lib/models/trip/trip_data.dart

class TripData {
  final int id;
  final String name;

  TripData({required this.id, required this.name});

  factory TripData.fromJson(Map<String, dynamic> json) {
    return TripData(id: json['id'] as int, name: json['name'] as String);
  }

  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name};
  }
}
