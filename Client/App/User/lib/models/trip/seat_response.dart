// lib/models/seat_data_response.dart

import 'trip.dart';

class SeatDataResponse {
  final List<TripData> seats;

  SeatDataResponse({required this.seats});

  factory SeatDataResponse.fromJson(Map<String, dynamic> json) {
    return SeatDataResponse(
      seats:
          (json['seats'] as List)
              .map((seatJson) => TripData.fromJson(seatJson))
              .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
    'seats': seats.map((seat) => seat.toJson()).toList(),
  };
}
