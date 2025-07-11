import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:table_calendar/table_calendar.dart';
import 'package:taixe/models/trip.dart'; // Import model Trip
import 'package:taixe/services/trip_service.dart'; // Import service Trip

class CalendarPage extends StatefulWidget {
  const CalendarPage({super.key});

  @override
  _CalendarPageState createState() => _CalendarPageState();
}

class _CalendarPageState extends State<CalendarPage> {
  // Trạng thái của UI
  bool _isLoading = true;
  String? _errorMessage;

  // Dữ liệu lịch
  CalendarFormat _calendarFormat = CalendarFormat.month;
  DateTime _focusedDay = DateTime.now();
  DateTime? _selectedDay;
  
  // Service để lấy dữ liệu
  final TripService _tripService = TripService();

  // Dữ liệu sự kiện được lấy từ API, nhóm theo ngày
  Map<DateTime, List<Trip>> _events = {};

  @override
  void initState() {
    super.initState();
    _selectedDay = _focusedDay; // Mặc định chọn ngày hôm nay
    _fetchTrips(); // Bắt đầu lấy dữ liệu từ API khi trang được tạo
  }

  /// Lấy dữ liệu chuyến đi từ API và cập nhật UI
  Future<void> _fetchTrips() async {
    try {
      final trips = await _tripService.getDriverTrips();
      final Map<DateTime, List<Trip>> eventsMap = {};

      for (var trip in trips) {
        // Chuyển đổi chuỗi ngày "yyyy-MM-dd" thành đối tượng DateTime
        final departureDate = DateTime.parse(trip.departureDate);
        // Chuẩn hóa ngày về 0 giờ 0 phút để so sánh chính xác
        final normalizedDay = DateTime(departureDate.year, departureDate.month, departureDate.day);
        
        if (eventsMap[normalizedDay] == null) {
          eventsMap[normalizedDay] = [];
        }
        eventsMap[normalizedDay]!.add(trip);
      }

      if (mounted) {
        setState(() {
          _events = eventsMap;
          _isLoading = false;
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

  /// Trả về danh sách các chuyến đi cho một ngày cụ thể
  List<Trip> _getEventsForDay(DateTime day) {
    final normalizedDay = DateTime(day.year, day.month, day.day);
    return _events[normalizedDay] ?? [];
  }

  /// Chuyển đổi mã trạng thái từ API thành chuỗi văn bản dễ hiểu
  String _getTripStatusText(int status) {
    // Dựa trên đặc tả API, 1 là trạng thái hoạt động.
    // Bạn có thể mở rộng logic này nếu có nhiều mã trạng thái hơn.
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
      return DateFormat('hh:mm a').format(parsedTime);
    } catch (e) {
      return time; // Trả về thời gian gốc nếu không thể định dạng
    }
  }

  /// Widget chính để xây dựng giao diện
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.blue,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white, size: 28),
          onPressed: () => Navigator.pop(context),
        ),
        title: const Text(
          'Lịch Chuyến',
          style: TextStyle(
            fontFamily: 'Pacifico',
            fontWeight: FontWeight.normal,
            fontSize: 18,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [Colors.blue[50]!, Colors.white70],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: _buildBody(), // Gọi hàm xây dựng phần thân
      ),
    );
  }

  /// Xây dựng phần thân của trang dựa trên trạng thái (loading, error, success)
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
                onPressed: () {
                  setState(() {
                    _isLoading = true;
                    _errorMessage = null;
                  });
                  _fetchTrips();
                },
                child: const Text('Thử lại'),
              )
            ],
          ),
        ),
      );
    }

    return Column(
      children: [
        TableCalendar(
          locale: 'vi_VN', // Hỗ trợ tiếng Việt
          firstDay: DateTime.utc(2025, 1, 1),
          lastDay: DateTime.utc(2030, 12, 31),
          focusedDay: _focusedDay,
          calendarFormat: _calendarFormat,
          selectedDayPredicate: (day) => isSameDay(_selectedDay, day),
          onDaySelected: (selectedDay, focusedDay) {
            setState(() {
              _selectedDay = selectedDay;
              _focusedDay = focusedDay;
            });
          },
          onFormatChanged: (format) {
            if (_calendarFormat != format) {
              setState(() {
                _calendarFormat = format;
              });
            }
          },
          onPageChanged: (focusedDay) {
            _focusedDay = focusedDay;
          },
          eventLoader: _getEventsForDay,
          calendarStyle: CalendarStyle(
            todayDecoration: BoxDecoration(color: Colors.blue[300], shape: BoxShape.circle),
            selectedDecoration: BoxDecoration(color: Colors.blue[700], shape: BoxShape.circle),
            markerDecoration: const BoxDecoration(color: Colors.red, shape: BoxShape.circle),
          ),
          headerStyle: HeaderStyle(
            formatButtonVisible: true,
            titleCentered: true,
            formatButtonShowsNext: false,
            titleTextStyle: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            formatButtonDecoration: BoxDecoration(
              color: Colors.blue[700],
              borderRadius: BorderRadius.circular(10),
            ),
            formatButtonTextStyle: const TextStyle(color: Colors.white),
          ),
        ),
        const SizedBox(height: 8.0),
        Expanded(
          child: ListView.builder(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.fromLTRB(16.0, 8.0, 16.0, 80.0),
            itemCount: _getEventsForDay(_selectedDay!).length,
            itemBuilder: (context, index) {
              final trip = _getEventsForDay(_selectedDay!)[index];
              return _buildTripCard(trip); // Xây dựng card cho mỗi chuyến đi
            },
          ),
        ),
      ],
    );
  }

  /// Xây dựng card hiển thị thông tin chi tiết của một chuyến đi
  Widget _buildTripCard(Trip trip) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8.0),
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
          contentPadding: const EdgeInsets.all(12.0),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15)),
          title: Text(
            '${trip.route.start.name} - ${trip.route.end.name}',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 18,
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
                'Trạng thái: ${_getTripStatusText(trip.status)}',
                style: TextStyle(
                  fontSize: 12,
                  color: _getStatusColor(trip.status),
                  fontWeight: FontWeight.w500,
                ),
              ),
              const SizedBox(height: 2),
              Text(
                'Số hành khách: ${trip.passengers}',
                style: TextStyle(fontSize: 12, color: Colors.grey[700]),
              ),
              const SizedBox(height: 2),
              Text(
                'Xuất phát: ${trip.route.start.name}',
                style: TextStyle(fontSize: 12, color: Colors.grey[700]),
              ),
              const SizedBox(height: 2),
              Text(
                'Đích đến: ${trip.route.end.name}',
                style: TextStyle(fontSize: 12, color: Colors.grey[700]),
              ),
            ],
          ),
          trailing: Icon(Icons.directions_bus, color: Colors.blue[700], size: 32),
          onTap: () {
            // TODO: Thêm hành động khi nhấn vào một chuyến đi, ví dụ: xem chi tiết
            print('Đã nhấn vào chuyến đi ID: ${trip.id}');
          },
        ),
      ),
    );
  }
}