import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:taixe/calendar_page.dart';
import 'package:taixe/checkin_page.dart';
import 'package:taixe/models/trip.dart'; // Import model Trip
import 'package:taixe/services/trip_service.dart'; // Import TripService

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  // Trạng thái của UI
  bool _isLoading = true;
  String? _errorMessage;

  // Dữ liệu
  String _fullName = "";
  List<Trip> _todaysTrips = [];
  final TripService _tripService = TripService();

  @override
  void initState() {
    super.initState();
    // Bắt đầu lấy tất cả dữ liệu cần thiết khi trang được tải
    _fetchData();
  }

  /// Lấy cả dữ liệu người dùng và danh sách chuyến đi
  Future<void> _fetchData() async {
    // Nếu không phải đang làm mới, hiển thị vòng quay loading chính
    if (!_isLoading) {
      setState(() {
        _isLoading = true;
      });
    }

    try {
      // Lấy tên người dùng từ SharedPreferences
      final prefs = await SharedPreferences.getInstance();
      final fullName = prefs.getString('fullName') ?? "Tài xế";

      // Lấy danh sách chuyến đi từ API
      final allTrips = await _tripService.getDriverTrips();

      // Lọc ra các chuyến đi của ngày hôm nay
      final now = DateTime.now();
      final today = DateTime(now.year, now.month, now.day);

      final todaysTrips =
          allTrips.where((trip) {
            final departureDate = DateTime.parse(trip.departureDate);
            return isSameDay(departureDate, today);
          }).toList();

      // Sắp xếp các chuyến đi trong ngày theo thời gian khởi hành
      todaysTrips.sort((a, b) => a.departureTime.compareTo(b.departureTime));

      if (mounted) {
        setState(() {
          _fullName = fullName;
          _todaysTrips = todaysTrips;
          _isLoading = false;
          _errorMessage = null;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = e.toString();
          _isLoading = false;
        });
      }
    }
  }

  bool isSameDay(DateTime date1, DateTime date2) {
    return date1.year == date2.year &&
        date1.month == date2.month &&
        date1.day == date2.day;
  }

  /// Chuyển đổi mã trạng thái từ API thành chuỗi văn bản dễ hiểu
  String _getTripStatusText(int status) {
    switch (status) {
      case 0:
        return 'Đã hủy';
      case 1:
        return 'Sắp khởi hành';
      case 2:
        return 'Đang chạy';
      case 3:
        return 'Hoàn thành';
      default:
        return 'Không xác định';
    }
  }

  /// Lấy màu sắc tương ứng với trạng thái
  Color _getStatusColor(int status) {
    switch (status) {
      case 0:
        return Colors.red;
      case 1:
        return Colors.orange;
      case 2:
        return Colors.green;
      case 3:
        return Colors.grey;
      default:
        return Colors.black;
    }
  }

  /// Chuyển đổi định dạng thời gian "HH:mm:ss" thành "hh:mm a"
  String _formatTime(String time) {
    try {
      final parsedTime = DateFormat('HH:mm:ss').parse(time);
      return DateFormat('hh:mm a', 'vi_VN').format(parsedTime);
    } catch (e) {
      return time;
    }
  }

  @override
  Widget build(BuildContext context) {
    SystemChrome.setSystemUIOverlayStyle(
      const SystemUiOverlayStyle(
        statusBarColor: Colors.blue,
        statusBarIconBrightness: Brightness.light,
      ),
    );

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [Colors.blue[50]!, Colors.white70],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: Column(
          children: [
            _buildHeader(),
            Expanded(
              // RefreshIndicator cho phép "kéo để làm mới"
              child: RefreshIndicator(
                onRefresh: _fetchData,
                child: _buildBody(),
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// Widget xây dựng phần Header
  Widget _buildHeader() {
    return Container(
      color: Colors.blue,
      padding: const EdgeInsets.only(
        top: 40.0,
        left: 16.0,
        right: 16.0,
        bottom: 16.0,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Expanded(
            child: Text(
              'Tài xế: $_fullName',
              style: const TextStyle(
                fontFamily: 'Pacifico',
                fontWeight: FontWeight.normal,
                fontSize: 20,
                color: Colors.white,
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.calendar_today,
              color: Colors.white,
              size: 28,
            ),
            onPressed:
                () => Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const CalendarPage()),
                ),
          ),
        ],
      ),
    );
  }

  /// Widget xây dựng phần thân chính, xử lý các trạng thái UI
  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_errorMessage != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                _errorMessage!,
                textAlign: TextAlign.center,
                style: const TextStyle(color: Colors.red, fontSize: 16),
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: _fetchData,
                child: const Text('Thử lại'),
              ),
            ],
          ),
        ),
      );
    }

    if (_todaysTrips.isEmpty) {
      return const Center(
        child: Text(
          'Không có chuyến đi nào được giao cho hôm nay.',
          style: TextStyle(fontSize: 16, color: Colors.grey),
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.only(top: 8.0, bottom: 35.0),
      itemCount: _todaysTrips.length,
      itemBuilder: (context, index) {
        final trip = _todaysTrips[index];
        return _buildTripCard(trip);
      },
    );
  }

  /// Widget xây dựng Card cho mỗi chuyến đi
  Widget _buildTripCard(Trip trip) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [Colors.blue[100]!, Colors.white],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(15),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.2),
            spreadRadius: 2,
            blurRadius: 8,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: ListTile(
          leading: CircleAvatar(
            radius: 30,
            backgroundColor: Colors.blue.shade200,
            child: Icon(Icons.directions_bus, size: 30, color: Colors.white),
          ),
          title: Text(
            '${trip.route.start.name} - ${trip.route.end.name}',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 16,
              color: Colors.blue[900],
            ),
          ),
          subtitle: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 4),
              Text(
                'Giờ chạy: ${_formatTime(trip.departureTime)}',
                style: TextStyle(fontSize: 14, color: Colors.grey[700]),
              ),
              const SizedBox(height: 2),
              Text(
                _getTripStatusText(trip.status),
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: _getStatusColor(trip.status),
                ),
              ),
            ],
          ),
          trailing: Icon(
            Icons.qr_code_scanner,
            color: Colors.blue[700],
            size: 32,
          ),
          onTap: () {
            print('Navigating to CheckInPage with tripId: ${trip.id}');
            Navigator.push(
              context,
              MaterialPageRoute(
                // Truyền tripId vào CheckInPage
                builder: (context) => CheckInPage(trip: trip),
              ),
            );
          },
          contentPadding: const EdgeInsets.all(12.0),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
        ),
      ),
    );
  }
}
