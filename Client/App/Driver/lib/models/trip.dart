// lib/models/trip.dart

// Lớp chính đại diện cho một chuyến đi (Trip)
class Trip {
  final int id;
  final String departureDate;
  final String departureTime;
  final String arrivalDate;
  final String arrivalTime;
  final RouteInfo route;
  final int totalSeats;
  final int availableSeats; // 'stock' trong API là số ghế còn lại
  final int status;

  Trip({
    required this.id,
    required this.departureDate,
    required this.departureTime,
    required this.arrivalDate,
    required this.arrivalTime,
    required this.route,
    required this.totalSeats,
    required this.availableSeats,
    required this.status,
  });

  factory Trip.fromJson(Map<String, dynamic> json) {
    return Trip(
      id: json['id'],
      departureDate: json['departureDate'],
      departureTime: json['departureTime'],
      arrivalDate: json['arrivalDate'],
      arrivalTime: json['arrivalTime'],
      route: RouteInfo.fromJson(json['route']),
      totalSeats: json['total'],
      availableSeats: json['stock'],
      status: json['status'],
    );
  }

  // Tính số hành khách đã đặt vé
  int get passengers => totalSeats - availableSeats;
}

// Lớp thông tin về tuyến đường
class RouteInfo {
  final int id;
  final Location start;
  final Location end;
  final String distance;
  final int price;

  RouteInfo({
    required this.id,
    required this.start,
    required this.end,
    required this.distance,
    required this.price,
  });

  factory RouteInfo.fromJson(Map<String, dynamic> json) {
    return RouteInfo(
      id: json['id'],
      start: Location.fromJson(json['start']),
      end: Location.fromJson(json['end']),
      distance: json['distance'],
      price: json['price'],
    );
  }
}

// Lớp thông tin về địa điểm (điểm đầu/cuối)
class Location {
  final int id;
  final String name;

  Location({
    required this.id,
    required this.name,
  });

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(
      id: json['id'],
      name: json['name'],
    );
  }
}