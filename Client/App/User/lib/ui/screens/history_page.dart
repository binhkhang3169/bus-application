// lib/history_page.dart

import 'package:caoky/global/toast.dart'; // Giả sử bạn có tệp này cho Toast, nếu không hãy thay thế
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';
import 'package:dotted_line/dotted_line.dart';

// Import trang chi tiết vé đã được cập nhật
import 'ticket_detail_page.dart';

// Helper function to safely parse nullable JSON objects like {"String": "value", "Valid": true}
T? _parseNullable<T>(Map<String, dynamic>? json, String key) {
  if (json != null && json[key] != null && json[key]['Valid'] == true) {
    var valueMap = json[key];
    // Lấy giá trị đầu tiên không phải là key 'Valid' (ví dụ: 'String', 'Int32')
    return valueMap.values.firstWhere((v) => v != null, orElse: () => null);
  }
  return null;
}

// --- DATA MODELS (Cập nhật cho HistoryPage) ---

// Model for seat information
class SeatTicketInfo {
  final int id;
  final int seatId;
  final String? seatName;

  SeatTicketInfo({required this.id, required this.seatId, this.seatName});

  factory SeatTicketInfo.fromJson(Map<String, dynamic> json) {
    return SeatTicketInfo(
      id: json['id'] ?? -1,
      seatId: json['seat_id'] ?? -1,
      seatName: _parseNullable(json, 'seat_name'),
    );
  }
}

// Model for the main ticket data in the list
class Ticket {
  final String ticketId;
  final String tripIdBegin;
  final String? tripIdEnd;
  final double price;
  final int status;
  final DateTime bookingTime;
  final List<SeatTicketInfo> seatTicketsBegin;
  final List<SeatTicketInfo> seatTicketsEnd;
  final Trip? trip; // Populated after fetching trip details

  Ticket({
    required this.ticketId,
    required this.tripIdBegin,
    this.tripIdEnd,
    required this.price,
    required this.status,
    required this.bookingTime,
    required this.seatTicketsBegin,
    required this.seatTicketsEnd,
    this.trip,
  });

  factory Ticket.fromJson(Map<String, dynamic> json) {
    var seatTicketsBeginList = (json['SeatTicketsBegin'] as List? ?? [])
        .map((i) => SeatTicketInfo.fromJson(i))
        .toList();
    var seatTicketsEndList = (json['SeatTicketsEnd'] as List? ?? [])
        .map((i) => SeatTicketInfo.fromJson(i))
        .toList();

    return Ticket(
      ticketId: json['ticket_id'] ?? 'N/A',
      tripIdBegin: json['trip_id_begin'] ?? '',
      tripIdEnd: _parseNullable(json, 'trip_id_end'),
      price: (json['price'] as num? ?? 0).toDouble(),
      status: json['status'] ?? -1,
      bookingTime: json['booking_time'] != null
          ? DateTime.parse(json['booking_time'])
          : DateTime.now(),
      seatTicketsBegin: seatTicketsBeginList,
      seatTicketsEnd: seatTicketsEndList,
    );
  }

  // Helper to create a new instance with updated trip info
  Ticket withTrip(Trip newTrip) {
    return Ticket(
      ticketId: this.ticketId,
      tripIdBegin: this.tripIdBegin,
      tripIdEnd: this.tripIdEnd,
      price: this.price,
      status: this.status,
      bookingTime: this.bookingTime,
      seatTicketsBegin: this.seatTicketsBegin,
      seatTicketsEnd: this.seatTicketsEnd,
      trip: newTrip, // Set the fetched trip data
    );
  }

  // Get total number of seats for display
  int get totalSeats => seatTicketsBegin.length + seatTicketsEnd.length;

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
}

// Model for the trip data (Verified to match new structure)
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
}

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
}

class Location {
  final String name;
  Location({required this.name});

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(name: json['name'] ?? 'Không rõ');
  }
}


// --- API SERVICE (Cập nhật cho HistoryPage) ---
class HistoryApiService {
  final String _baseUrl = "http://57.155.76.74"; // Use 10.0.2.2 for Android emulator

  Future<String> _getBearerToken() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    return prefs.getString('accessToken') ?? "";
  }

  Future<List<Ticket>> getTickets() async {
    final _bearerToken = await _getBearerToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/tickets'),
      headers: {
        'Authorization': 'Bearer $_bearerToken',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode == 200) {
      // Adjusted parsing to match the new JSON structure
      final data = json.decode(response.body)['data'] as List;
      return data.map((ticketJson) => Ticket.fromJson(ticketJson)).toList();
    } else {
      print('Failed to load tickets: ${response.statusCode} ${response.body}');
      throw Exception('Failed to load tickets');
    }
  }

  Future<Trip> getTripDetails(String tripId) async {
    final _bearerToken = await _getBearerToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/trips/$tripId'),
      headers: {
        'Authorization': 'Bearer $_bearerToken',
        'Content-Type': 'application/json',
      },
    );

    if (response.statusCode == 200) {
      // Adjusted parsing to match the new JSON structure
      final data = json.decode(response.body)['data'];
      return Trip.fromJson(data);
    } else {
      print('Failed to load trip details for $tripId: ${response.statusCode} ${response.body}');
      throw Exception('Failed to load trip details');
    }
  }
}


// --- MAIN WIDGET ---
class HistoryPage extends StatefulWidget {
  @override
  State<HistoryPage> createState() => _HistoryPageState();
}

class _HistoryPageState extends State<HistoryPage> {
  late Future<List<Ticket>> _ticketsFuture;
  final HistoryApiService _apiService = HistoryApiService();

  @override
  void initState() {
    super.initState();
    _ticketsFuture = _loadAndCombineTicketData();
  }

  Future<List<Ticket>> _loadAndCombineTicketData() async {
    try {
      final baseTickets = await _apiService.getTickets();
      if (baseTickets.isEmpty) return [];

      final ticketWithTripFutures = baseTickets.map((ticket) async {
        try {
          // Fetch trip details for the beginning trip
          final trip = await _apiService.getTripDetails(ticket.tripIdBegin);
          // Return a new Ticket object with combined data
          return ticket.withTrip(trip);
        } catch (e) {
          print("Error fetching trip details for ${ticket.tripIdBegin}: $e");
          return ticket; // Return original ticket if trip details fail
        }
      }).toList();

      return await Future.wait(ticketWithTripFutures);
    } catch (e) {
      print("Error loading tickets: $e");
      if (mounted) {
        // Sử dụng ToastUtils hoặc ScaffoldMessenger để hiển thị lỗi
        ToastUtils.show("Không thể tải danh sách vé: $e");
      }
      return [];
    }
  }
  
  void _refreshData() {
    setState(() {
      _ticketsFuture = _loadAndCombineTicketData();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: Center(
          child: Text(
            "Lịch sử vé",
            style: TextStyle(
                fontWeight: FontWeight.bold, fontSize: 18, color: Colors.white),
          ),
        ),
        backgroundColor: Colors.blueAccent,
        automaticallyImplyLeading: false, // Bỏ nút back nếu đây là tab chính
      ),
      body: FutureBuilder<List<Ticket>>(
        future: _ticketsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(child: Text('Đã xảy ra lỗi khi tải dữ liệu.'));
          }
          if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return Center(child: Text('Không tìm thấy vé nào.'));
          }

          final allTickets = snapshot.data!;
          final now = DateTime.now();

          allTickets.sort((a, b) => b.bookingTime.compareTo(a.bookingTime));

          // Chuyến đi được xác định dựa trên thời gian khởi hành của chiều đi
          final upcomingTickets = allTickets.where((t) =>
                  t.trip != null && t.trip!.departureDateTime.isAfter(now))
              .toList();
          final historicalTickets = allTickets.where((t) =>
                  t.trip == null || t.trip!.departureDateTime.isBefore(now))
              .toList();

          return DefaultTabController(
            length: 2,
            child: Column(
              children: [
                TabBar(
                  tabs: [Tab(text: 'SẮP ĐI'), Tab(text: 'ĐÃ ĐI')],
                  indicatorColor: Colors.blue,
                  labelColor: Colors.blue,
                  unselectedLabelColor: Colors.black,
                  indicatorWeight: 4.0,
                ),
                Expanded(
                  child: TabBarView(
                    children: [
                      RefreshIndicator(
                        onRefresh: () async => _refreshData(),
                        child: ListView.builder(
                          padding: EdgeInsets.all(12),
                          itemCount: upcomingTickets.length,
                          itemBuilder: (context, index) =>
                              buildTicketCard(upcomingTickets[index]),
                        ),
                      ),
                      RefreshIndicator(
                        onRefresh: () async => _refreshData(),
                        child: ListView.builder(
                          padding: EdgeInsets.all(12),
                          itemCount: historicalTickets.length,
                          itemBuilder: (context, index) =>
                              buildTicketCard(historicalTickets[index]),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
  
  // Widget `buildTicketCard` để sử dụng model mới
  Widget buildTicketCard(Ticket ticket) {
    return GestureDetector(
      onTap: () async {
        // Điều hướng và đợi kết quả trả về (ví dụ: nếu vé được hủy)
        final result = await Navigator.push(
          context,
          MaterialPageRoute(
            // Truyền ticketId vào trang chi tiết
            builder: (context) => TicketDetailPage(ticketId: ticket.ticketId),
          ),
        );
        // Nếu có kết quả trả về (ví dụ: "refresh"), tải lại dữ liệu
        if (result == 'refresh') {
          _refreshData();
        }
      },
      child: ClipPath(
        clipper: TicketClipper(),
        child: Container(
          margin: EdgeInsets.only(bottom: 16),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(10),
            color: ticket.statusColor == Colors.green
                ? Color(0xFFE6F7EC)
                : (ticket.statusColor == Colors.red
                    ? Color(0xFFFEEFEF)
                    : Color(0xFFEBF5FF)),
            border: Border.all(color: ticket.statusColor.withOpacity(0.5)),
          ),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Flexible(
                    child: RichText(
                      text: TextSpan(
                        style: TextStyle(fontSize: 16, color: Colors.black),
                        children: [
                          TextSpan(text: 'Mã vé/Code: '),
                          TextSpan(
                            text: ticket.ticketId,
                            style: TextStyle(fontWeight: FontWeight.bold),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Icon(Icons.arrow_forward_ios, color: Colors.grey[600], size: 16),
                ],
              ),
              SizedBox(height: 4),
              Text(
                ticket.statusText,
                style: TextStyle(
                    color: ticket.statusColor, fontWeight: FontWeight.bold),
              ),
              SizedBox(height: 10),
              DottedLine(dashColor: Colors.grey.withOpacity(0.7)),
              SizedBox(height: 10),
              Row(
                children: [
                  Text(
                    'Thông tin chiều đi',
                    style: TextStyle(
                        color: Colors.red.shade700,
                        fontWeight: FontWeight.w500),
                  ),
                  SizedBox(width: 4),
                  Icon(Icons.info_outline, size: 16, color: Colors.red.shade700),
                ],
              ),
              SizedBox(height: 8),
              Row(
                children: [
                  Icon(Icons.location_on, color: Colors.green),
                  SizedBox(width: 3),
                  Expanded(child: Text(ticket.trip?.route.start.name ?? 'Đang tải...', overflow: TextOverflow.ellipsis,)),
                  SizedBox(width: 5),
                  Icon(Icons.arrow_right_alt, color: Colors.grey),
                  SizedBox(width: 5),
                  Icon(Icons.flag, color: Colors.red),
                  SizedBox(width: 3),
                  Expanded(child: Text(ticket.trip?.route.end.name ?? '...', overflow: TextOverflow.ellipsis,)),
                ],
              ),
              SizedBox(height: 6),
              RichText(
                text: TextSpan(
                  style: TextStyle(fontSize: 15, color: Colors.black87),
                  children: [
                    TextSpan(text: 'Số ghế: '),
                    TextSpan(
                      // Sử dụng totalSeats getter
                      text: ticket.totalSeats.toString(),
                      style: TextStyle(fontWeight: FontWeight.bold),
                    ),
                  ],
                ),
              ),
              SizedBox(height: 4),
              RichText(
                text: TextSpan(
                  style: TextStyle(fontSize: 15, color: Colors.black87),
                  children: [
                    TextSpan(text: 'Giờ xuất bến: '),
                    TextSpan(
                      text: ticket.trip != null
                          ? DateFormat('HH:mm dd/MM/yyyy', 'vi_VN')
                              .format(ticket.trip!.departureDateTime)
                          : 'Đang tải...',
                      style: TextStyle(
                          color: Colors.red, fontWeight: FontWeight.w600),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// --- CUSTOM CLIPPER (Không thay đổi) ---
class TicketClipper extends CustomClipper<Path> {
  @override
  Path getClip(Size size) {
    double radius = 12;
    double curveCenterY = size.height * 0.35;
    Path path = Path();
    path.moveTo(0, radius);
    path.arcToPoint(Offset(radius, 0), radius: Radius.circular(radius));
    path.lineTo(size.width - radius, 0);
    path.arcToPoint(Offset(size.width, radius), radius: Radius.circular(radius));
    path.lineTo(size.width, curveCenterY - radius);
    path.arcToPoint(Offset(size.width, curveCenterY + radius), radius: Radius.circular(radius), clockwise: false);
    path.lineTo(size.width, size.height - radius);
    path.arcToPoint(Offset(size.width - radius, size.height), radius: Radius.circular(radius));
    path.lineTo(radius, size.height);
    path.arcToPoint(Offset(0, size.height - radius), radius: Radius.circular(radius));
    path.lineTo(0, curveCenterY + radius);
    path.arcToPoint(Offset(0, curveCenterY - radius), radius: Radius.circular(radius), clockwise: false);
    path.close();
    return path;
  }
  @override
  bool shouldReclip(covariant CustomClipper<Path> oldClipper) => false;
}