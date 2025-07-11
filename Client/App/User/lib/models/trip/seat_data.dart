// lib/models/seats_data.dart

import 'package:caoky/models/trip/seat.dart';

class SeatsData {
  final List<Seat> seats;

  SeatsData({required this.seats});

  factory SeatsData.fromJson(Map<String, dynamic> json) {
    // The data from the API is a list under the 'seats' key.
    var seatList = json['seats'] as List;
    List<Seat> seats = seatList.map((i) => Seat.fromJson(i)).toList();
    return SeatsData(seats: seats);
  }
}
