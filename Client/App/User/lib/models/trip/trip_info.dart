// lib/models/trip/trip_info.dart
// Ensure you have this class defined or update it.
// This class is used by BusTripScreen and SeatSelectionScreen.
class TripInfo {
  final String status;
  final String departureDate;
  final String departureTime;
  final String arrivalDate;
  final String arrivalTime;
  final int tripId;
  final int price;
  final String estimatedTime;
  final int stock;
  final String license;
  final String vehicleId;
  final String vehicleType;
  final String estimatedDistance;
  final String departureStation;
  final String arrivalStation;
  final String fullRoute;

  TripInfo({
    required this.status,
    required this.departureDate,
    required this.departureTime,
    required this.arrivalDate,
    required this.arrivalTime,
    required this.tripId,
    required this.price,
    required this.estimatedTime,
    required this.stock,
    required this.license,
    required this.vehicleId,
    required this.vehicleType,
    required this.estimatedDistance,
    required this.departureStation,
    required this.arrivalStation,
    required this.fullRoute,
  });

  factory TripInfo.fromJson(Map<String, dynamic> json) {
    return TripInfo(
      status: json['status'] as String,
      departureDate: json['departureDate'] as String,
      departureTime: json['departureTime'] as String,
      arrivalDate: json['arrivalDate'] as String,
      arrivalTime: json['arrivalTime'] as String,
      tripId: json['tripId'] as int,
      price: json['price'] as int,
      estimatedTime: json['estimatedTime'] as String,
      stock: json['stock'] as int,
      license: json['license'] as String,
      vehicleId: json['vehicleId'] as String,
      vehicleType: json['vehicleType'] as String,
      estimatedDistance: json['estimatedDistance'] as String,
      departureStation: json['departureStation'] as String,
      arrivalStation: json['arrivalStation'] as String,
      fullRoute: json['fullRoute'] as String,
    );
  }
}
