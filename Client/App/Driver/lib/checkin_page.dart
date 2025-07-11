import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_speed_dial/flutter_speed_dial.dart';
import 'package:image_picker/image_picker.dart';
import 'package:intl/intl.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:qr_code_tools/qr_code_tools.dart';
import 'package:taixe/global/toast.dart';
import 'package:taixe/models/trip.dart';
import 'package:taixe/services/checkin_service.dart';

class CheckInPage extends StatefulWidget {
  final Trip trip;

  const CheckInPage({super.key, required this.trip});

  @override
  State<CheckInPage> createState() => _CheckInPageState();
}

class _CheckInPageState extends State<CheckInPage> {
  // Trạng thái UI
  bool _isLoading = true;
  String? _errorMessage;
  bool _isProcessingCheckin = false;

  // Dữ liệu
  final CheckinService _checkinService = CheckinService();
  late List<Map<String, String>> _allSeats;
  Set<String> _checkedInSeatNames = {};

  @override
  void initState() {
    super.initState();
    _loadInitialData();
  }

  Future<void> _loadInitialData() async {
    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    try {
      _allSeats = _generateSeatLayout(widget.trip.totalSeats);
      final checkedInList = await _checkinService.getCheckedInSeats(widget.trip.id);
      
      if (mounted) {
        setState(() {
          _checkedInSeatNames = checkedInList
              .where((c) => c.seatName != null)
              .map((c) => c.seatName!)
              .toSet();
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

  /// Tự động tạo sơ đồ ghế dựa trên tổng số ghế
  List<Map<String, String>> _generateSeatLayout(int totalSeats) {
    List<Map<String, String>> seats = [];
    int seatsPerDeck = (totalSeats / 2).ceil();
    // Tầng dưới: A01, A02, ...
    for (int i = 1; i <= seatsPerDeck; i++) {
      seats.add({"seat": "A${i.toString().padLeft(2, '0')}", "deck": "Tầng dưới"});
    }
    // Tầng trên: B01, B02, ...
    for (int i = 1; i <= totalSeats - seatsPerDeck; i++) {
      seats.add({"seat": "B${i.toString().padLeft(2, '0')}", "deck": "Tầng trên"});
    }
    return seats;
  }

  /// Xử lý nội dung QR sau khi quét
  Future<void> _handleQrScan(String? qrContent) async {
    if (qrContent == null || qrContent.isEmpty) {
      ToastUtils.show("Không tìm thấy mã QR.");
      return;
    }

    setState(() { _isProcessingCheckin = true; });
    _showProcessingDialog("Đang xử lý check-in...");

    try {
      final result = await _checkinService.performCheckin(
        qrContent: qrContent,
        tripId: widget.trip.id,
      );
      
      Navigator.of(context).pop(); // Đóng dialog xử lý
      
      if (result.seatName != null) {
        setState(() {
          _checkedInSeatNames.add(result.seatName!);
        });
        _showResultDialog(
          isSuccess: true,
          title: "Thành Công!",
          content: "Đã check-in thành công cho ghế ${result.seatName}.",
        );
      }
    } catch (e) {
       Navigator.of(context).pop(); // Đóng dialog xử lý
       _showResultDialog(
          isSuccess: false,
          title: "Check-in Thất Bại",
          content: e.toString(),
        );
    } finally {
        if(mounted){
             setState(() { _isProcessingCheckin = false; });
        }
    }
  }
  
  /// Mở màn hình quét QR bằng camera
  void _startCameraScan() {
    Navigator.push(context, MaterialPageRoute(
      builder: (context) => Scaffold(
        appBar: AppBar(title: const Text("Quét mã QR")),
        body: MobileScanner(
          onDetect: (capture) {
            final List<Barcode> barcodes = capture.barcodes;
            if (barcodes.isNotEmpty && barcodes.first.rawValue != null) {
              Navigator.pop(context, barcodes.first.rawValue);
            }
          },
        ),
      ),
    )).then((qrContent) => _handleQrScan(qrContent));
  }
  
  /// Mở thư viện để chọn và quét ảnh QR
  Future<void> _pickAndScanImage() async {
     final picker = ImagePicker();
     final pickedFile = await picker.pickImage(source: ImageSource.gallery);
     if (pickedFile != null) {
        final data = await QrCodeToolsPlugin.decodeFrom(pickedFile.path);
        _handleQrScan(data);
     }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [Colors.blue[50]!, Colors.white70],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: _buildBody(),
      ),
      floatingActionButton: _buildFloatingActionButton(),
    );
  }

  AppBar _buildAppBar() {
    return AppBar(
      backgroundColor: Colors.blue,
      elevation: 0,
      leading: IconButton(
        icon: const Icon(Icons.arrow_back, color: Colors.white, size: 28),
        onPressed: () => Navigator.pop(context),
      ),
      title: const Text(
        'Điểm Danh Chuyến Đi',
        style: TextStyle(fontFamily: 'Pacifico', fontSize: 18, color: Colors.white),
      ),
      centerTitle: true,
    );
  }
  
  Widget _buildBody() {
    if (_isLoading) return const Center(child: CircularProgressIndicator());
    if (_errorMessage != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(_errorMessage!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _loadInitialData, child: const Text("Tải lại"))
          ],
        ),
      );
    }

    final lowerDeckSeats = _allSeats.where((s) => s["deck"] == "Tầng dưới").toList();
    final upperDeckSeats = _allSeats.where((s) => s["deck"] == "Tầng trên").toList();

    return Column(
      children: [
        _buildTripInfo(),
        const Divider(height: 1, thickness: 1),
        Expanded(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                _buildDeckView(lowerDeckSeats, "Tầng dưới"),
                const SizedBox(height: 24),
                _buildDeckView(upperDeckSeats, "Tầng trên"),
                const SizedBox(height: 24),
                _buildLegend(),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildTripInfo() {
    String formattedDate = DateFormat('dd/MM/yyyy').format(DateTime.parse(widget.trip.departureDate));
    String formattedTime = DateFormat('hh:mm a').format(DateFormat('HH:mm:ss').parse(widget.trip.departureTime));

    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text(
            "${widget.trip.route.start.name} → ${widget.trip.route.end.name}",
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 5),
          Text(
            "$formattedTime, $formattedDate",
            style: const TextStyle(fontSize: 14, color: Colors.grey),
          ),
        ],
      ),
    );
  }

  Widget _buildDeckView(List<Map<String, String>> seats, String title) {
    return Column(
      children: [
        Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
        const SizedBox(height: 10),
        GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 5, // Hiển thị 5 ghế mỗi hàng
            crossAxisSpacing: 8,
            mainAxisSpacing: 8,
            childAspectRatio: 1.5, // Điều chỉnh tỉ lệ cho ghế
          ),
          itemCount: seats.length,
          itemBuilder: (context, index) => _buildSeat(seats[index]),
        ),
      ],
    );
  }

  Widget _buildSeat(Map<String, String> seatData) {
    String seatName = seatData["seat"]!;
    bool isCheckedIn = _checkedInSeatNames.contains(seatName);

    return Container(
      decoration: BoxDecoration(
        color: isCheckedIn ? Colors.green : Colors.white,
        border: Border.all(color: Colors.blue.shade300),
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.2),
            blurRadius: 2,
            offset: const Offset(1,1)
          )
        ]
      ),
      child: Center(
        child: Text(
          seatName,
          style: TextStyle(
            color: isCheckedIn ? Colors.white : Colors.black,
            fontSize: 12,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
    );
  }

  Widget _buildLegend() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        _legendItem(Colors.white, "Chưa Check-in"),
        _legendItem(Colors.green, "Đã Check-in"),
      ],
    );
  }
  
  Widget _legendItem(Color color, String text) {
     return Row(
      children: [
        Container(
          width: 20, height: 20,
          decoration: BoxDecoration(
             color: color,
             border: Border.all(color: Colors.blue.shade300),
             borderRadius: BorderRadius.circular(4)
          ),
        ),
        const SizedBox(width: 8),
        Text(text),
      ],
    );
  }

  Widget? _buildFloatingActionButton() {
    return SpeedDial(
      icon: Icons.qr_code_scanner,
      activeIcon: Icons.close,
      backgroundColor: Colors.blue,
      foregroundColor: Colors.white,
      overlayColor: Colors.black,
      overlayOpacity: 0.5,
      spacing: 12,
      spaceBetweenChildren: 12,
      children: [
        SpeedDialChild(
          child: const Icon(Icons.camera_alt),
          label: 'Quét bằng Camera',
          onTap: _isProcessingCheckin ? null : _startCameraScan,
        ),
        SpeedDialChild(
          child: const Icon(Icons.image),
          label: 'Quét từ ảnh',
          onTap: _isProcessingCheckin ? null : _pickAndScanImage,
        ),
      ],
    );
  }

  void _showProcessingDialog(String text) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        content: Row(
          children: [
            const CircularProgressIndicator(),
            const SizedBox(width: 20),
            Text(text),
          ],
        ),
      ),
    );
  }
  
  void _showResultDialog({required bool isSuccess, required String title, required String content}) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(title, style: TextStyle(color: isSuccess ? Colors.green : Colors.red)),
        content: Text(content),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text("OK"),
          ),
        ],
      ),
    );
  }
}