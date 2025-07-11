// lib/ticket_detail_page.dart

import 'dart:io';
import 'dart:typed_data';
import 'package:caoky/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';
import 'package:path_provider/path_provider.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:saver_gallery/saver_gallery.dart';

// --- HELPER FUNCTIONS & MODELS (Dành cho TicketDetailPage) ---

// Helper function to safely parse nullable JSON objects
T? _parseNullable<T>(Map<String, dynamic>? json, String key) {
  if (json != null && json[key] != null && json[key]['Valid'] == true) {
    var valueMap = json[key];
    return valueMap.values.firstWhere((v) => v != null, orElse: () => null);
  }
  return null;
}

// Model for seat information in ticket details
class SeatTicket {
  final int id;
  final int seatId;
  final String ticketId;
  final String tripId;
  final String? seatName;

  SeatTicket({
    required this.id,
    required this.seatId,
    required this.ticketId,
    required this.tripId,
    this.seatName,
  });

  factory SeatTicket.fromJson(Map<String, dynamic> json) {
    return SeatTicket(
      id: json['id'] ?? -1,
      seatId: json['seat_id'] ?? -1,
      ticketId: json['ticket_id'] ?? '',
      tripId: json['trip_id'] ?? '',
      seatName: _parseNullable(json, 'seat_name'),
    );
  }
}

// Model for journey details (pickup/dropoff points)
class JourneyDetail {
  final int detailId;
  final int? pickupLocationBegin;
  final int? dropoffLocationBegin;
  final int? pickupLocationEnd;
  final int? dropoffLocationEnd;

  JourneyDetail({
    required this.detailId,
    this.pickupLocationBegin,
    this.dropoffLocationBegin,
    this.pickupLocationEnd,
    this.dropoffLocationEnd,
  });

  factory JourneyDetail.fromJson(Map<String, dynamic> json) {
    return JourneyDetail(
      detailId: json['detail_id'] ?? -1,
      pickupLocationBegin: _parseNullable(json, 'pickup_location_begin'),
      dropoffLocationBegin: _parseNullable(json, 'dropoff_location_begin'),
      pickupLocationEnd: _parseNullable(json, 'pickup_location_end'),
      dropoffLocationEnd: _parseNullable(json, 'dropoff_location_end'),
    );
  }
}

// Main model for detailed ticket information
class TicketDetails {
  final String ticketId;
  final String? name;
  final String? phone;
  final String? email;
  final double price;
  final int status;
  final int paymentStatus;
  final int bookingChannel;
  final DateTime bookingTime;
  final bool isRoundTrip;
  final List<JourneyDetail> details;
  final List<SeatTicket> seatTicketsBegin;
  final List<SeatTicket> seatTicketsEnd;

  TicketDetails({
    required this.ticketId,
    this.name,
    this.phone,
    this.email,
    required this.price,
    required this.status,
    required this.paymentStatus,
    required this.bookingChannel,
    required this.bookingTime,
    required this.isRoundTrip,
    required this.details,
    required this.seatTicketsBegin,
    required this.seatTicketsEnd,
  });

  factory TicketDetails.fromJson(Map<String, dynamic> json) {
    final ticketData = json['ticket'];
    return TicketDetails(
      ticketId: ticketData['ticket_id'] ?? 'N/A',
      name: _parseNullable(ticketData, 'name'),
      phone: _parseNullable(ticketData, 'phone'),
      email: _parseNullable(ticketData, 'email'),
      price: (ticketData['price'] as num? ?? 0).toDouble(),
      status: ticketData['status'] ?? -1,
      paymentStatus: ticketData['payment_status'] ?? -1,
      bookingChannel: ticketData['booking_channel'] ?? 0,
      bookingTime: DateTime.parse(ticketData['booking_time']),
      isRoundTrip: _parseNullable<String>(ticketData, 'trip_id_end') != null,
      details:
          (ticketData['Details'] as List? ?? [])
              .map((d) => JourneyDetail.fromJson(d))
              .toList(),
      seatTicketsBegin:
          (ticketData['SeatTicketsBegin'] as List? ?? [])
              .map((s) => SeatTicket.fromJson(s))
              .toList(),
      seatTicketsEnd:
          (ticketData['SeatTicketsEnd'] as List? ?? [])
              .map((s) => SeatTicket.fromJson(s))
              .toList(),
    );
  }
}

// Model for station/location details
class Station {
  final int id;
  final String name;
  Station({required this.id, required this.name});

  factory Station.fromJson(Map<String, dynamic> json) {
    return Station(id: json['id'], name: json['name']);
  }
}

// --- API SERVICE (Dành cho TicketDetailPage) ---
class ApiService {
  final String _baseUrl = "http://57.155.76.74";
  final AuthRepository _authRepository = AuthRepository();

  Future<String?> _getBearerToken() async {
    return await _authRepository.getValidAccessToken();
  }

  Future<TicketDetails> getTicketDetails(String ticketId) async {
    final token = await _getBearerToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/tickets/$ticketId'),
      headers: {'Authorization': 'Bearer $token'},
    );
    if (response.statusCode == 200) {
      return TicketDetails.fromJson(json.decode(response.body)['data']);
    } else {
      throw Exception('Failed to load ticket details: ${response.body}');
    }
  }

  Future<Station> getStationDetails(int stationId) async {
    final token = await _getBearerToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/stations/$stationId'),
      headers: {'Authorization': 'Bearer $token'},
    );
    if (response.statusCode == 200) {
      return Station.fromJson(json.decode(response.body)['data']);
    } else {
      throw Exception('Failed to load station details for ID: $stationId');
    }
  }

  Future<Uint8List> getQrCodeImage(String content) async {
    final token = await _getBearerToken();
    final response = await http.get(
      Uri.parse('$_baseUrl/api/v1/qr/image?content=$content'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );
    if (response.statusCode == 200) {
      return response.bodyBytes;
    } else {
      throw Exception('Failed to generate QR code');
    }
  }
}

// --- VIEWMODEL ---
class TicketPageViewModel {
  final TicketDetails ticket;
  final Map<int, String> locationNames;
  final Map<int, Uint8List> qrCodeImages; // Keyed by seat_id
  final Map<int, bool> qrCodeErrors;

  TicketPageViewModel({
    required this.ticket,
    required this.locationNames,
    required this.qrCodeImages,
    required this.qrCodeErrors,
  });
}

// --- MAIN WIDGET ---
class TicketDetailPage extends StatefulWidget {
  final String ticketId;
  const TicketDetailPage({super.key, required this.ticketId});

  @override
  State<TicketDetailPage> createState() => _TicketDetailPageState();
}

class _TicketDetailPageState extends State<TicketDetailPage> {
  final ApiService _apiService = ApiService();
  late Future<TicketPageViewModel> _pageDataFuture;

  @override
  void initState() {
    super.initState();
    _pageDataFuture = _loadAllPageData();
  }

  Future<TicketPageViewModel> _loadAllPageData() async {
    try {
      // 1. Fetch ticket details
      final ticket = await _apiService.getTicketDetails(widget.ticketId);

      // 2. Collect all unique location IDs from all journey details
      final locationIds = <int>{};
      for (var detail in ticket.details) {
        if (detail.pickupLocationBegin != null)
          locationIds.add(detail.pickupLocationBegin!);
        if (detail.dropoffLocationBegin != null)
          locationIds.add(detail.dropoffLocationBegin!);
        if (detail.pickupLocationEnd != null)
          locationIds.add(detail.pickupLocationEnd!);
        if (detail.dropoffLocationEnd != null)
          locationIds.add(detail.dropoffLocationEnd!);
      }

      // Combine all seats from both trips
      final allSeats = [...ticket.seatTicketsBegin, ...ticket.seatTicketsEnd];

      // 3. Concurrently fetch location names and QR codes
      final results = await Future.wait([
        _fetchLocationNames(locationIds),
        _fetchQrCodes(ticket.ticketId, allSeats),
      ]);

      final locationNames = results[0] as Map<int, String>;
      final qrData = results[1] as Map<String, dynamic>;

      return TicketPageViewModel(
        ticket: ticket,
        locationNames: locationNames,
        qrCodeImages: qrData['images'] as Map<int, Uint8List>,
        qrCodeErrors: qrData['errors'] as Map<int, bool>,
      );
    } catch (e, stackTrace) {
      debugPrint(
        '[TicketDetailPage] Lỗi khi tải dữ liệu trang: $e\n$stackTrace',
      );
      throw Exception('Lỗi khi tải dữ liệu vé: $e');
    }
  }

  Future<Map<int, String>> _fetchLocationNames(Set<int> locationIds) async {
    if (locationIds.isEmpty) return {};
    final futures = locationIds.map(
      (id) => _apiService
          .getStationDetails(id)
          .then((s) => {s.id: s.name})
          .catchError((_) => <int, String>{}),
    );
    final results = await Future.wait(futures);
    final names = <int, String>{};
    for (var map in results) {
      names.addAll(map);
    }
    return names;
  }

  Future<Map<String, dynamic>> _fetchQrCodes(
    String ticketId,
    List<SeatTicket> seats,
  ) async {
    final images = <int, Uint8List>{};
    final errors = <int, bool>{};

    final qrFutures =
        seats.map((seat) {
          final content = 'TICKET:$ticketId-SEAT:${seat.seatId}';
          return _apiService
              .getQrCodeImage(content)
              .then((bytes) {
                images[seat.seatId] = bytes;
              })
              .catchError((e) {
                errors[seat.seatId] = true;
                debugPrint(
                  '[TicketDetailPage] Lỗi tải QR cho ghế ${seat.seatId}: $e',
                );
              });
        }).toList();

    await Future.wait(qrFutures);
    return {'images': images, 'errors': errors};
  }

  // --- Utility Formatters ---
  String formatDate(DateTime dt) =>
      DateFormat('HH:mm dd/MM/yyyy', 'vi_VN').format(dt);
  String formatPrice(double price) =>
      NumberFormat.currency(locale: 'vi_VN', symbol: '₫').format(price);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[100],
      appBar: AppBar(
        title: Text('Chi Tiết Vé'),
        backgroundColor: Colors.blueAccent,
        foregroundColor: Colors.white,
        centerTitle: true,
      ),
      body: FutureBuilder<TicketPageViewModel>(
        future: _pageDataFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(
              child: Text(
                'Lỗi tải chi tiết vé: ${snapshot.error}',
                style: TextStyle(color: Colors.red),
              ),
            );
          }
          if (!snapshot.hasData) {
            return const Center(child: Text('Không tìm thấy dữ liệu vé.'));
          }

          final viewModel = snapshot.data!;
          final ticket = viewModel.ticket;

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildInfoSection(ticket),
                const SizedBox(height: 20),
                if (ticket.details.isNotEmpty)
                  _buildJourneySection(ticket, viewModel.locationNames),
                const SizedBox(height: 20),
                if (ticket.seatTicketsBegin.isNotEmpty)
                  _buildQrSection(
                    'Mã QR Chiều Đi',
                    ticket.seatTicketsBegin,
                    viewModel.qrCodeImages,
                    viewModel.qrCodeErrors,
                  ),
                if (ticket.isRoundTrip && ticket.seatTicketsEnd.isNotEmpty) ...[
                  const SizedBox(height: 20),
                  _buildQrSection(
                    'Mã QR Chiều Về',
                    ticket.seatTicketsEnd,
                    viewModel.qrCodeImages,
                    viewModel.qrCodeErrors,
                  ),
                ],
              ],
            ),
          );
        },
      ),
    );
  }

  // --- UI Building Widgets ---

  Widget _buildInfoSection(TicketDetails ticket) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Thông Tin Chung',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            const Divider(height: 24),
            _buildInfoItem('Mã vé:', ticket.ticketId),
            _buildInfoItem('Tên người đặt:', ticket.name ?? 'Không có'),
            _buildInfoItem('Số điện thoại:', ticket.phone ?? 'Không có'),
            _buildInfoItem('Email:', ticket.email ?? 'Không có'),
            _buildInfoItem('Tổng giá:', formatPrice(ticket.price)),
            _buildInfoItem('Thời gian đặt:', formatDate(ticket.bookingTime)),
            _buildInfoItem(
              'Kênh đặt:',
              _getBookingChannelText(ticket.bookingChannel),
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _buildStatusItem(
                  'Thanh toán:',
                  _getPaymentStatus(ticket.paymentStatus),
                ),
                _buildStatusItem('Trạng thái:', _getTicketStatus(ticket)),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildJourneySection(
    TicketDetails ticket,
    Map<int, String> locationNames,
  ) {
    // Assuming the first detail has all the info needed for display
    if (ticket.details.isEmpty) return const SizedBox.shrink();

    final detail = ticket.details.first;
    final pickupBegin =
        locationNames[detail.pickupLocationBegin] ??
        'ID: ${detail.pickupLocationBegin}';
    final dropoffBegin =
        locationNames[detail.dropoffLocationBegin] ??
        'ID: ${detail.dropoffLocationBegin}';
    final pickupEnd =
        locationNames[detail.pickupLocationEnd] ??
        'ID: ${detail.pickupLocationEnd}';
    final dropoffEnd =
        locationNames[detail.dropoffLocationEnd] ??
        'ID: ${detail.dropoffLocationEnd}';

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Chi Tiết Hành Trình',
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            const Divider(height: 24),
            Text(
              'CHIỀU ĐI',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: Colors.blueAccent,
              ),
            ),
            const SizedBox(height: 8),
            _buildJourneyRow(
              Icons.my_location,
              'Điểm đón:',
              pickupBegin,
              Colors.blue,
            ),
            _buildJourneyRow(Icons.flag, 'Điểm trả:', dropoffBegin, Colors.red),
            if (ticket.isRoundTrip) ...[
              const Divider(height: 24),
              Text(
                'CHIỀU VỀ',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  color: Colors.green,
                ),
              ),
              const SizedBox(height: 8),
              _buildJourneyRow(
                Icons.my_location,
                'Điểm đón:',
                pickupEnd,
                Colors.blue,
              ),
              _buildJourneyRow(Icons.flag, 'Điểm trả:', dropoffEnd, Colors.red),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildJourneyRow(
    IconData icon,
    String label,
    String value,
    Color color,
  ) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: color, size: 20),
          const SizedBox(width: 8),
          Text(
            label,
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.black87,
            ),
          ),
          const SizedBox(width: 4),
          Expanded(child: Text(value)),
        ],
      ),
    );
  }

  Widget _buildQrSection(
    String title,
    List<SeatTicket> seats,
    Map<int, Uint8List> qrImages,
    Map<int, bool> qrErrors,
  ) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: Theme.of(
                context,
              ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
            ),
            const Divider(height: 24),
            GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                maxCrossAxisExtent: 200,
                childAspectRatio: 3 / 4.2,
                crossAxisSpacing: 16,
                mainAxisSpacing: 16,
              ),
              itemCount: seats.length,
              itemBuilder: (context, index) {
                final seat = seats[index];
                final qrImage = qrImages[seat.seatId];
                final hasError = qrErrors[seat.seatId] ?? false;
                final seatName = seat.seatName ?? 'N/A';

                return Card(
                  elevation: 1,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'Ghế $seatName',
                          style: TextStyle(fontWeight: FontWeight.bold),
                        ),
                        const SizedBox(height: 8),
                        Expanded(
                          child:
                              hasError
                                  ? const Center(
                                    child: Text(
                                      "Lỗi QR",
                                      style: TextStyle(color: Colors.red),
                                    ),
                                  )
                                  : qrImage != null
                                  ? GestureDetector(
                                    onTap:
                                        () => _showQrDialog(
                                          context,
                                          qrImage,
                                          seatName,
                                        ),
                                    child: Image.memory(
                                      qrImage,
                                      fit: BoxFit.contain,
                                    ),
                                  )
                                  : const Center(
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                    ),
                                  ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  // --- Dialogs and Helpers ---

  Future<void> _showQrDialog(
    BuildContext context,
    Uint8List qrImage,
    String seatName,
  ) async {
    showDialog(
      context: context,
      builder:
          (_) => Dialog(
            backgroundColor: Colors.white,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
            child: Padding(
              padding: const EdgeInsets.all(20.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    "Mã QR cho ghế $seatName",
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                  const SizedBox(height: 16),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: Image.memory(qrImage),
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton.icon(
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(
                        horizontal: 24,
                        vertical: 12,
                      ),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(30),
                      ),
                    ),
                    icon: const Icon(Icons.download_rounded),
                    label: const Text("Tải về máy"),
                    onPressed: () async {
                      final status = await Permission.photos.request();
                      print(status.isGranted);
                      if (status.isGranted) {
                        try {
                          final tempDir = await getTemporaryDirectory();
                          final fileName =
                              'qr_code_ve_${widget.ticketId}_ghe_$seatName.png';
                          final file =
                              await File('${tempDir.path}/$fileName').create();
                          await file.writeAsBytes(qrImage);

                          // await SaverGallery.saveFile(
                          //   filePath: file.path,
                          //   fileName: fileName,
                          //   androidRelativePath: "Pictures/MaQRVeXe",
                          // );

                          await SaverGallery.saveFile(
                            filePath: file.path,
                            fileName: fileName,
                            androidRelativePath: "Pictures/MaQRVeXe",
                            skipIfExists: true,
                          );

                          if (mounted) {
                            Navigator.of(context).pop();
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(
                                content: Text("Đã lưu mã QR vào thư viện ảnh!"),
                              ),
                            );
                          }
                        } catch (e) {
                          if (mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text("Lỗi khi lưu ảnh: $e")),
                            );
                          }
                        }
                      } else if (status.isPermanentlyDenied) {
                        // TRƯỜNG HỢP 2: NGƯỜI DÙNG TỪ CHỐI VĨNH VIỄN
                        // Ứng dụng không thể xin quyền lại.
                        // Cách tốt nhất là hướng dẫn họ mở cài đặt của ứng dụng.
                        if (mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Text(
                                "Bạn đã từ chối quyền truy cập ảnh. Vui lòng vào cài đặt để cấp quyền.",
                              ),
                              action: SnackBarAction(
                                label: "Mở Cài đặt",
                                onPressed: () {
                                  // Mở trang cài đặt của ứng dụng này
                                  openAppSettings();
                                },
                              ),
                            ),
                          );
                        }
                      } else {
                        if (mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text("Không được cấp quyền để lưu ảnh."),
                            ),
                          );
                        }
                      }
                    },
                  ),
                ],
              ),
            ),
          ),
    );
  }

  Widget _buildInfoItem(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '$label ',
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.grey[700],
            ),
          ),
          Expanded(child: Text(value)),
        ],
      ),
    );
  }

  Widget _buildStatusItem(String label, Widget statusBadge) {
    return Row(
      children: [
        Text(
          '$label ',
          style: TextStyle(
            fontWeight: FontWeight.w600,
            color: Colors.grey[700],
          ),
        ),
        statusBadge,
      ],
    );
  }

  Widget _buildStatusBadge(String text, Color bgColor, Color textColor) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        text,
        style: TextStyle(
          color: textColor,
          fontSize: 12,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  String _getBookingChannelText(int channel) => switch (channel) {
    0 => "Online",
    1 => "Tại quầy",
    _ => "N/A",
  };

  Widget _getPaymentStatus(int status) => switch (status) {
    1 => _buildStatusBadge(
      'Đã thanh toán',
      Colors.green.shade100,
      Colors.green.shade800,
    ),
    0 => _buildStatusBadge(
      'Chưa thanh toán',
      Colors.orange.shade100,
      Colors.orange.shade800,
    ),
    _ => _buildStatusBadge(
      'Không xác định',
      Colors.grey.shade200,
      Colors.grey.shade800,
    ),
  };

  Widget _getTicketStatus(TicketDetails ticket) {
    if (ticket.status == 2) {
      return _buildStatusBadge(
        'Đã hủy',
        Colors.red.shade100,
        Colors.red.shade800,
      );
    }
    return _buildStatusBadge(
      'Hoạt động',
      Colors.blue.shade100,
      Colors.blue.shade800,
    );
  }
}
