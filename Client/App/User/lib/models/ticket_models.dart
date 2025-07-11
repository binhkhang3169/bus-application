// lib/models/ticket_models.dart

import 'package:flutter/material.dart';

// --- DATA MODELS ---

// Model for the main ticket data
class Ticket {
  final String ticketId;
  final int? customerId;
  final String? phone;
  final String? email;
  final String? name;
  final double price;
  final int status;
  final DateTime bookingTime;
  final int paymentStatus;
  final int bookingChannel;
  final DateTime createdAt;
  final DateTime updatedAt;
  final String? bookedBy;
  final int policyId;
  final String tripId;
  final List<SeatTicket> seatTickets;
  final List<TicketDetail> details; // Corresponds to 'Details' in React
  final Trip? trip; // Populated after fetching trip details

  Ticket({
    required this.ticketId,
    this.customerId,
    this.phone,
    this.email,
    this.name,
    required this.price,
    required this.status,
    required this.bookingTime,
    required this.paymentStatus,
    required this.bookingChannel,
    required this.createdAt,
    required this.updatedAt,
    required this.bookedBy,
    required this.policyId,
    required this.tripId,
    required this.seatTickets,
    required this.details,
    this.trip,
  });

  factory Ticket.fromJson(Map<String, dynamic> json) {
    return Ticket(
      ticketId: json['TicketID'] ?? 'N/A',
      customerId: json['CustomerID']?['Int32'],
      phone: json['phone']?['String'],
      email: json['email']?['String'],
      name: json['name']?['String'],
      price: (json['Price'] as num? ?? 0).toDouble(),
      status: json['Status'] ?? -1,
      bookingTime:
          json['BookingTime'] != null
              ? DateTime.parse(json['BookingTime'])
              : DateTime.now(),
      paymentStatus: json['PaymentStatus'] ?? -1,
      bookingChannel: json['BookingChannel'] ?? -1,
      bookedBy: json['BookedBy'] ?? '',
      createdAt:
          json['CreatedAt'] != null
              ? DateTime.parse(json['CreatedAt'])
              : DateTime.now(),
      updatedAt:
          json['UpdatedAt'] != null
              ? DateTime.parse(json['UpdatedAt'])
              : DateTime.now(),
      policyId: json['PolicyID'] ?? -1,
      tripId: json['trip_id'] ?? '',
      seatTickets:
          (json['SeatTickets'] as List? ?? [])
              .map((seatJson) => SeatTicket.fromJson(seatJson))
              .toList(),
      details:
          (json['Details'] as List? ?? [])
              .map((detailJson) => TicketDetail.fromJson(detailJson))
              .toList(),
    );
  }

  // Helper to determine the status text and color
  String get statusText {
    switch (status) {
      case 0:
        return 'Chờ thanh toán';
      case 1:
        return 'Thành công';
      case 2:
        return 'Đã huỷ';
      default:
        return 'Không xác định';
    }
  }

  Color get statusColor {
    switch (status) {
      case 0:
        return Colors.blue;
      case 1:
        return Colors.green;
      case 2:
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  Ticket copyWith({
    String? ticketId,
    int? customerId,
    String? phone,
    String? email,
    String? name,
    double? price,
    int? status,
    DateTime? bookingTime,
    int? paymentStatus,
    int? bookingChannel,
    DateTime? createdAt,
    DateTime? updatedAt,
    int? policyId,
    String? tripId,
    List<SeatTicket>? seatTickets,
    List<TicketDetail>? details,
    Trip? trip,
  }) {
    return Ticket(
      ticketId: ticketId ?? this.ticketId,
      customerId: customerId ?? this.customerId,
      phone: phone ?? this.phone,
      email: email ?? this.email,
      name: name ?? this.name,
      price: price ?? this.price,
      status: status ?? this.status,
      bookingTime: bookingTime ?? this.bookingTime,
      paymentStatus: paymentStatus ?? this.paymentStatus,
      bookingChannel: bookingChannel ?? this.bookingChannel,
      bookedBy: bookedBy ?? this.bookedBy,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      policyId: policyId ?? this.policyId,
      tripId: tripId ?? this.tripId,
      seatTickets: seatTickets ?? this.seatTickets,
      details: details ?? this.details,
      trip: trip ?? this.trip,
    );
  }

  Map<String, dynamic> toJson() => {
    'ticketId': ticketId,
    'tripId': tripId,
    'name': name,
    'phone': phone,
    'email': email,
    'price': price,
    'bookingTime': bookingTime.toIso8601String(),
    'bookingChannel': bookingChannel,
    'policyId': policyId,
    'bookedBy': bookedBy,
    'createdAt': createdAt.toIso8601String(),
    'updatedAt': updatedAt.toIso8601String(),
    'status': status,
    'paymentStatus': paymentStatus,
    'details': details.map((d) => d.toJson()).toList(),
    'seatTickets': seatTickets.map((s) => s.toJson()).toList(),
  };
}

// Model for seat information
class SeatTicket {
  final String seatId;
  final int status;

  SeatTicket({required this.seatId, required this.status});

  factory SeatTicket.fromJson(Map<String, dynamic> json) {
    return SeatTicket(
      seatId: json['seat_id'] ?? 'N/A',
      status: json['status'] ?? -1,
    );
  }
  Map<String, dynamic> toJson() => {'seat_id': seatId, 'status': status};
}

// Model for journey details (pickup/dropoff)
class TicketDetail {
  final String detailId;
  final int? pickupLocationId;
  final int? dropoffLocationId;

  TicketDetail({
    required this.detailId,
    this.pickupLocationId,
    this.dropoffLocationId,
  });

  factory TicketDetail.fromJson(Map<String, dynamic> json) {
    return TicketDetail(
      detailId: json['detail_id'] ?? 'N/A',
      pickupLocationId: json['pickup_location']?['Int32'],
      dropoffLocationId: json['dropoff_location']?['Int32'],
    );
  }
  Map<String, dynamic> toJson() => {
    'detail_id': detailId,
    'pickup_location': {'Int32': pickupLocationId},
    'dropoff_location': {'Int32': dropoffLocationId},
  };
}

// Model for station/location details fetched separately
class Station {
  final int id;
  final String name;

  Station({required this.id, required this.name});

  factory Station.fromJson(Map<String, dynamic> json) {
    final data = json['data'] ?? {};
    return Station(id: data['id'] ?? -1, name: data['name'] ?? 'Không rõ');
  }
  Map<String, dynamic> toJson() => {
    'data': {'id': id, 'name': name},
  };
}

// Model for the trip data
class Trip {
  final int id;
  final DateTime departureDateTime;
  final DateTime arrivalDateTime;
  final RouteInfo route;

  Trip({
    required this.id,
    required this.departureDateTime,
    required this.arrivalDateTime,
    required this.route,
  });

  factory Trip.fromJson(Map<String, dynamic> json) {
    String departureDate = json['departureDate'] ?? '1970-01-01';
    String departureTime = json['departureTime'] ?? '00:00:00';
    String arrivalDate = json['arrivalDate'] ?? '1970-01-01';
    String arrivalTime = json['arrivalTime'] ?? '00:00:00';

    return Trip(
      id: json['id'] ?? -1,
      departureDateTime: DateTime.parse('$departureDate $departureTime'),
      arrivalDateTime: DateTime.parse('$arrivalDate $arrivalTime'),
      route: RouteInfo.fromJson(json['route'] ?? {}),
    );
  }
  Map<String, dynamic> toJson() => {
    'id': id,
    'departureDate': departureDateTime.toIso8601String().split('T').first,
    'departureTime': departureDateTime.toIso8601String().split('T').last,
    'arrivalDate': arrivalDateTime.toIso8601String().split('T').first,
    'arrivalTime': arrivalDateTime.toIso8601String().split('T').last,
    'route': route.toJson(),
  };
}

// Model for the route information within a trip
class RouteInfo {
  final Location start;
  final Location end;

  RouteInfo({required this.start, required this.end});

  factory RouteInfo.fromJson(Map<String, dynamic> json) {
    return RouteInfo(
      start: Location.fromJson(json['start'] ?? {}),
      end: Location.fromJson(json['end'] ?? {}),
    );
  }
  Map<String, dynamic> toJson() => {
    'start': start.toJson(),
    'end': end.toJson(),
  };
}

// Model for location data
class Location {
  final String name;

  Location({required this.name});

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(name: json['name'] ?? 'Không rõ');
  }
  Map<String, dynamic> toJson() => {'name': name};
}
