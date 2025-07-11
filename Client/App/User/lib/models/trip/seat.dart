// lib/models/seat.dart

class Seat {
  final int id;
  final String name;

  Seat({required this.id, required this.name});

  factory Seat.fromJson(Map<String, dynamic> json) {
    return Seat(id: json['id'], name: json['name']);
  }
}
