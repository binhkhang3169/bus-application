// lib/api_service.dart
import 'dart:convert';
import 'dart:typed_data';
import 'package:caoky/services/auth_service.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

// --- DATA MODELS ---

// Helper functions remain the same...
String _getString(Map<String, dynamic>? json, String key) {
  if (json != null && json[key] != null && json[key]['Valid'] == true) {
    return json[key]['String'] ?? 'N/A';
  }
  return 'N/A';
}

int _getInt(Map<String, dynamic>? json, String key) {
  if (json != null && json[key] != null && json[key]['Valid'] == true) {
    return (json[key]['Int32'] as num? ?? -1).toInt();
  }
  return -1;
}

// Model for a seat ticket
class SeatTicket {
  final String seatId;
  final int status;

  SeatTicket({required this.seatId, required this.status});

  factory SeatTicket.fromJson(Map<String, dynamic> json) {
    return SeatTicket(
      // FIX: Handle seat_id which is an int in JSON but a String in the model.
      seatId: json['seat_id']?.toString() ?? 'N/A',
      status: (json['status'] as num? ?? -1).toInt(),
    );
  }
}

// Model for journey details (pickup/dropoff)
class JourneyDetail {
  final String detailId;
  final int pickupLocationId;
  final int dropoffLocationId;

  JourneyDetail({
    required this.detailId,
    required this.pickupLocationId,
    required this.dropoffLocationId,
  });

  factory JourneyDetail.fromJson(Map<String, dynamic> json) {
    return JourneyDetail(
      detailId:
          json['detail_id']?.toString() ??
          'N/A', // Ensure detail_id is a string
      pickupLocationId: _getInt(json, 'pickup_location'),
      dropoffLocationId: _getInt(json, 'dropoff_location'),
    );
  }
}

// Model for the main ticket details
class TicketDetails {
  final String ticketId;
  final String tripId;
  final String name;
  final String phone;
  final String email;
  final double price;
  final DateTime bookingTime;
  final int status;
  final int paymentStatus;
  final int bookingChannel;
  final int policyId;
  final String bookedBy;
  final DateTime createdAt;
  final DateTime updatedAt;
  final List<JourneyDetail> details;
  final List<SeatTicket> seatTickets;

  TicketDetails({
    required this.ticketId,
    required this.tripId,
    required this.name,
    required this.phone,
    required this.email,
    required this.price,
    required this.bookingTime,
    required this.status,
    required this.paymentStatus,
    required this.bookingChannel,
    required this.policyId,
    required this.bookedBy,
    required this.createdAt,
    required this.updatedAt,
    required this.details,
    required this.seatTickets,
  });

  factory TicketDetails.fromJson(Map<String, dynamic> json) {
    return TicketDetails(
      ticketId: json['TicketID'] ?? 'N/A',
      tripId: json['trip_id'] ?? 'N/A',
      name: _getString(json, 'name'),
      phone: _getString(json, 'phone'),
      email: _getString(json, 'email'),
      price: (json['Price'] as num? ?? 0).toDouble(),
      bookingTime:
          json['BookingTime'] != null
              ? DateTime.parse(json['BookingTime'])
              : DateTime.now(),
      status: (json['Status'] as num? ?? -1).toInt(),
      paymentStatus: (json['PaymentStatus'] as num? ?? -1).toInt(),
      bookingChannel: (json['BookingChannel'] as num? ?? -1).toInt(),
      policyId: (json['PolicyID'] as num? ?? -1).toInt(),
      bookedBy: _getString(json, 'BookedBy'),
      createdAt:
          json['CreatedAt'] != null
              ? DateTime.parse(json['CreatedAt'])
              : DateTime.now(),
      updatedAt:
          json['UpdatedAt'] != null
              ? DateTime.parse(json['UpdatedAt'])
              : DateTime.now(),
      // FIX: Use 'Details' (PascalCase) to match the JSON response.
      details:
          (json['Details'] as List? ?? [])
              .map((d) => JourneyDetail.fromJson(d))
              .toList(),
      // FIX: Use 'SeatTickets' (PascalCase) to match the JSON response.
      seatTickets:
          (json['SeatTickets'] as List? ?? [])
              .map((s) => SeatTicket.fromJson(s))
              .toList(),
    );
  }
}

// Model for station data
class Station {
  final int id;
  final String name;

  Station({required this.id, required this.name});

  factory Station.fromJson(Map<String, dynamic> json) {
    return Station(
      id: (json['id'] as num? ?? -1).toInt(),
      name: json['name'] ?? 'Không rõ',
    );
  }
}

// --- API SERVICE CLASS ---

class ApiService {
  final String _baseUrl = "http://57.155.76.74";
 final AuthRepository  _authRepository = AuthRepository();
  Future<String> _getBearerToken() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    return prefs.getString('accessToken') ?? "";
  }

  // Get details for a single ticket
  Future<TicketDetails> getTicketDetails(String ticketId) async {
    final token = _authRepository.getValidAccessToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/tickets/$ticketId'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode == 200) {
      // FIX: Extract the nested 'ticket' object from the 'data' field.
      final data = json.decode(response.body)['data']['ticket'];
      return TicketDetails.fromJson(data);
    } else {
      throw Exception('Failed to load ticket details: ${response.statusCode}');
    }
  }

  // Other methods (getStationDetails, getQrCodeImage) remain the same...
  Future<Station> getStationDetails(int stationId) async {
    final token = _authRepository.getValidAccessToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/stations/$stationId'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body)['data'];
      return Station.fromJson(data);
    } else {
      throw Exception('Failed to load station details for ID $stationId');
    }
  }

  Future<Uint8List> getQrCodeImage(String content) async {
    final token = _authRepository.getValidAccessToken();
    final url =
        '$_baseUrl/api/v1/qr/image?content=${Uri.encodeComponent(content)}';
    final response = await http.get(
      Uri.parse(url),
      headers: {'Authorization': 'Bearer $token'},
    );

    if (response.statusCode == 200) {
      return response.bodyBytes;
    } else {
      throw Exception('Failed to fetch QR code image');
    }
  }
}
