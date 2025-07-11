// lib/ui/ticket/seat_selection_page.dart
import 'dart:developer';

import 'package:caoky/global/convert_money.dart';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/api_response.dart';
import 'package:caoky/models/trip/seat.dart';
import 'package:caoky/models/trip/seat_data.dart';
import 'package:caoky/models/trip/trip_info.dart';
import 'package:caoky/services/api_trip_service.dart';
import 'package:caoky/services/dio_client.dart';
import 'package:caoky/ui/payment_page.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

// NEW: Enum để quản lý việc chọn chuyến đi
enum TripSelection { departure, returnTrip }

class SeatSelectionScreen extends StatefulWidget {
  // Dữ liệu chuyến đi
  final TripInfo initialTripInfo;
  final DateTime selectedDisplayDate;

  // NEW: Dữ liệu cho chuyến về (có thể null)
  final TripInfo? returnTripInfo;
  final DateTime? returnDisplayDate;

  SeatSelectionScreen({
    super.key,
    required this.initialTripInfo,
    required this.selectedDisplayDate,
    this.returnTripInfo, // Chấp nhận chuyến về
    this.returnDisplayDate, // Chấp nhận ngày về
  });

  @override
  _SeatSelectionScreenState createState() => _SeatSelectionScreenState();
}

class _SeatSelectionScreenState extends State<SeatSelectionScreen> {
  late ApiTripService _apiService;
  bool _isLoading = true;

  // NEW: Biến để theo dõi chuyến đi đang được chọn (đi hoặc về)
  TripSelection _currentTripSelection = TripSelection.departure;
  bool get _isReturnTripSelected =>
      _currentTripSelection == TripSelection.returnTrip;
  bool get _isRoundTrip => widget.returnTripInfo != null;

  // MODIFIED: Tách biệt trạng thái cho chuyến đi và chuyến về
  TripInfo? _departureTripDetails;
  Map<String, int> _departureAvailableSeatMap = {};
  List<String> _departureSelectedSeatNames = [];
  List<int> _departureSelectedSeatIds = [];

  TripInfo? _returnTripDetails;
  Map<String, int> _returnAvailableSeatMap = {};
  List<String> _returnSelectedSeatNames = [];
  List<int> _returnSelectedSeatIds = [];

  // Trạng thái chung
  bool _isUpperDeck = false;
  final List<String> _defaultLowerDeckSeatNames = List.generate(
    15,
    (i) => "A${i + 1}",
  );
  final List<String> _defaultUpperDeckSeatNames = List.generate(
    15,
    (i) => "B${i + 1}",
  );

  @override
  void initState() {
    super.initState();
    _apiService = ApiTripService(DioClient.createDio());
    _fetchTripAndSeatData();
  }

  // MODIFIED: Tải dữ liệu cho cả hai chuyến nếu có
  Future<void> _fetchTripAndSeatData() async {
    if (!mounted) return;
    setState(() => _isLoading = true);

    try {
      // Xây dựng danh sách các yêu cầu API cần thực hiện
      final apiCalls = <Future<ApiResponse>>[
        _apiService.getTripDetails(widget.initialTripInfo.tripId),
        _apiService.getAvailableSeats(widget.initialTripInfo.tripId),
      ];

      // Nếu là chuyến khứ hồi, thêm API call cho chuyến về
      if (_isRoundTrip) {
        apiCalls.add(_apiService.getTripDetails(widget.returnTripInfo!.tripId));
        apiCalls.add(
          _apiService.getAvailableSeats(widget.returnTripInfo!.tripId),
        );
      }

      final responses = await Future.wait(apiCalls);
      if (!mounted) return;

      // Xử lý kết quả cho chuyến đi
      _handleTripDetailsResponse(
        responses[0] as ApiResponse<TripInfo>,
        isReturn: false,
      );
      _handleSeatsResponse(
        responses[1] as ApiResponse<SeatsData>,
        isReturn: false,
      );

      // Xử lý kết quả cho chuyến về (nếu có)
      if (_isRoundTrip) {
        _handleTripDetailsResponse(
          responses[2] as ApiResponse<TripInfo>,
          isReturn: true,
        );
        _handleSeatsResponse(
          responses[3] as ApiResponse<SeatsData>,
          isReturn: true,
        );
      }
    } on DioException catch (e) {
      log("API DioException in SeatSelection: $e");
      if (mounted)
        ToastUtils.show("Lỗi kết nối: ${e.message ?? 'Không thể tải dữ liệu'}");
    } catch (e) {
      log("Lỗi không xác định: $e");
      if (mounted) ToastUtils.show("Đã xảy ra lỗi không mong muốn.");
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  void _handleTripDetailsResponse(
    ApiResponse<TripInfo> response, {
    required bool isReturn,
  }) {
    if (response.code == 200 && response.data != null) {
      setState(() {
        if (isReturn) {
          _returnTripDetails = response.data!;
        } else {
          _departureTripDetails = response.data!;
        }
      });
    } else {
      ToastUtils.show(
        "Lỗi tải chuyến ${isReturn ? 'về' : 'đi'}: ${response.message}",
      );
    }
  }

  void _handleSeatsResponse(
    ApiResponse<SeatsData> response, {
    required bool isReturn,
  }) {
    if (response.code == 200 && response.data != null) {
      final seatMap = {
        for (var seat in response.data!.seats) seat.name: seat.id,
      };
      setState(() {
        if (isReturn) {
          _returnAvailableSeatMap = seatMap;
        } else {
          _departureAvailableSeatMap = seatMap;
        }
      });
    } else {
      ToastUtils.show(
        "Lỗi tải ghế chuyến ${isReturn ? 'về' : 'đi'}: ${response.message}",
      );
    }
  }

  // MODIFIED: Logic chọn ghế dựa trên chuyến đi hiện tại
  void _toggleSeatSelection(String seatName) {
    if (!mounted) return;

    final currentAvailableMap =
        _isReturnTripSelected
            ? _returnAvailableSeatMap
            : _departureAvailableSeatMap;
    final currentSelectedNames =
        _isReturnTripSelected
            ? _returnSelectedSeatNames
            : _departureSelectedSeatNames;
    final currentSelectedIds =
        _isReturnTripSelected
            ? _returnSelectedSeatIds
            : _departureSelectedSeatIds;

    final int? seatId = currentAvailableMap[seatName];

    setState(() {
      if (currentSelectedNames.contains(seatName)) {
        currentSelectedNames.remove(seatName);
        if (seatId != null) currentSelectedIds.remove(seatId);
      } else {
        if (seatId != null) {
          currentSelectedNames.add(seatName);
          currentSelectedIds.add(seatId);
        } else {
          log("Warning: Could not find ID for selectable seat $seatName.");
        }
      }
    });
  }

  // MODIFIED: Hiển thị ngày/giờ cho chuyến đi hoặc về
  String _formatTripDateForDisplay() {
    final isReturn = _isReturnTripSelected;
    final tripDetails = isReturn ? _returnTripDetails : _departureTripDetails;
    final initialTripInfo =
        isReturn ? widget.returnTripInfo! : widget.initialTripInfo;
    final displayDate =
        isReturn ? widget.returnDisplayDate! : widget.selectedDisplayDate;

    if (tripDetails == null) {
      return DateFormat('EEEE, dd/MM/yyyy', 'vi_VN').format(displayDate) +
          " ${initialTripInfo.departureTime.substring(0, 5)}";
    }
    try {
      DateTime departureDateTime = DateFormat("yyyy-MM-dd HH:mm:ss").parse(
        "${tripDetails.departureDate} ${tripDetails.departureTime.split('.')[0]}",
      );
      return DateFormat(
        'EEEE, dd/MM/yyyy HH:mm',
        'vi_VN',
      ).format(departureDateTime);
    } catch (e) {
      log("Error formatting trip date: $e");
      return DateFormat('EEEE, dd/MM/yyyy', 'vi_VN').format(displayDate);
    }
  }

  @override
  Widget build(BuildContext context) {
    // MODIFIED: Chọn đúng dữ liệu để hiển thị dựa trên chuyến đi đang xem
    final TripInfo displayTripInfo =
        (_isReturnTripSelected ? _returnTripDetails : _departureTripDetails) ??
        (_isReturnTripSelected
            ? widget.returnTripInfo!
            : widget.initialTripInfo);

    // MODIFIED: Tính tổng giá vé cho cả 2 chuyến
    final int departurePrice =
        (_departureTripDetails ?? widget.initialTripInfo).price;
    final int returnPrice =
        (_returnTripDetails ?? widget.returnTripInfo)?.price ?? 0;
    final int totalPrice =
        (_departureSelectedSeatNames.length * departurePrice) +
        (_returnSelectedSeatNames.length * returnPrice);

    return Scaffold(
      backgroundColor: Colors.grey[100],
      appBar: AppBar(
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
        title: Column(
          children: [
            Text(
              "${displayTripInfo.departureStation} - ${displayTripInfo.arrivalStation}",
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
              overflow: TextOverflow.ellipsis,
            ),
            Text(
              _formatTripDateForDisplay(),
              style: TextStyle(color: Colors.white, fontSize: 13),
            ),
          ],
        ),
        centerTitle: true,
        backgroundColor: Colors.blueAccent,
        elevation: 3.0,
      ),
      body:
          _isLoading
              ? Center(child: CircularProgressIndicator())
              : Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: [
                        _buildLegendItem(
                          Colors.green.shade100,
                          Colors.green.shade500,
                          "Còn trống",
                        ),
                        _buildLegendItem(
                          Colors.orangeAccent,
                          Colors.orange.shade700,
                          "Đang chọn",
                        ),
                        _buildLegendItem(
                          Colors.grey.shade400,
                          Colors.grey.shade500,
                          "Đã bán",
                        ),
                      ],
                    ),
                  ),
                  // NEW: Bộ chuyển đổi Chuyến đi / Chuyến về
                  if (_isRoundTrip) ...[
                    ToggleButtons(
                      isSelected: [
                        !_isReturnTripSelected,
                        _isReturnTripSelected,
                      ],
                      onPressed: (index) {
                        setState(() {
                          _currentTripSelection =
                              index == 0
                                  ? TripSelection.departure
                                  : TripSelection.returnTrip;
                          _isUpperDeck =
                              false; // Reset về tầng dưới khi chuyển chuyến
                        });
                      },
                      borderRadius: BorderRadius.circular(10),
                      selectedColor: Colors.white,
                      fillColor: Colors.blueAccent.shade200,
                      color: Colors.black87,
                      textStyle: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                      ),
                      children: const [
                        Padding(
                          padding: EdgeInsets.symmetric(
                            horizontal: 24,
                            vertical: 8,
                          ),
                          child: Text("Chuyến đi"),
                        ),
                        Padding(
                          padding: EdgeInsets.symmetric(
                            horizontal: 24,
                            vertical: 8,
                          ),
                          child: Text("Chuyến về"),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                  ],
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.all(8),
                      child: _buildSeatGrid(
                        _isUpperDeck
                            ? _defaultUpperDeckSeatNames
                            : _defaultLowerDeckSeatNames,
                      ),
                    ),
                  ),
                ],
              ),
      bottomNavigationBar: _isLoading ? null : _buildBottomBar(totalPrice),
    );
  }

  Widget _buildSeatGrid(List<String> seatNames) {
    // MODIFIED: Di chuyển bộ chọn tầng xuống đây
    return Container(
      padding: const EdgeInsets.all(12.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 3,
              crossAxisSpacing: 10,
              mainAxisSpacing: 10,
              childAspectRatio: 2,
            ),
            itemCount: seatNames.length,
            itemBuilder: (context, index) => _buildSeatWidget(seatNames[index]),
          ),
          const SizedBox(height: 20),
          // MOVED: Bộ chọn tầng trên/dưới
          ToggleButtons(
            isSelected: [!_isUpperDeck, _isUpperDeck],
            onPressed: (index) {
              setState(() {
                _isUpperDeck = index == 1;
              });
            },
            borderRadius: BorderRadius.circular(10),
            selectedColor: Colors.white,
            fillColor: Colors.grey.shade400,
            color: Colors.black87,
            children: const [
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Text("Tầng dưới"),
              ),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Text("Tầng trên"),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSeatWidget(String seatName) {
    // MODIFIED: Logic hiển thị ghế dựa trên chuyến đi đang được chọn
    final currentAvailableMap =
        _isReturnTripSelected
            ? _returnAvailableSeatMap
            : _departureAvailableSeatMap;
    final currentSelectedNames =
        _isReturnTripSelected
            ? _returnSelectedSeatNames
            : _departureSelectedSeatNames;

    final isSelectable = currentAvailableMap.containsKey(seatName);
    final isSelectedByUser = currentSelectedNames.contains(seatName);

    Color backgroundColor;
    Color textColor;
    Color borderColor;
    bool isTapEnabled;

    if (isSelectedByUser) {
      backgroundColor = Colors.orange.shade400;
      textColor = Colors.white;
      borderColor = Colors.orange.shade700;
      isTapEnabled = true;
    } else if (isSelectable) {
      backgroundColor = Colors.teal.shade50;
      textColor = Colors.teal.shade800;
      borderColor = Colors.green.shade100;
      isTapEnabled = true;
    } else {
      backgroundColor = Colors.grey.shade200;
      textColor = Colors.grey.shade600;
      borderColor = Colors.grey.shade300;
      isTapEnabled = false;
    }

    return GestureDetector(
      onTap: isTapEnabled ? () => _toggleSeatSelection(seatName) : null,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        curve: Curves.easeInOut,
        decoration: BoxDecoration(
          color: backgroundColor,
          border: Border.all(color: borderColor, width: 1.5),
          borderRadius: BorderRadius.circular(8),
          boxShadow:
              isSelectedByUser
                  ? [
                    BoxShadow(
                      color: Colors.blueAccent.withOpacity(0.3),
                      blurRadius: 6,
                      offset: const Offset(0, 2),
                    ),
                  ]
                  : [
                    BoxShadow(
                      color: Colors.grey.withOpacity(0.1),
                      blurRadius: 4,
                      offset: const Offset(0, 1),
                    ),
                  ],
        ),
        child: Center(
          child: Text(
            seatName,
            style: TextStyle(
              color: textColor,
              fontSize: 12,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildBottomBar(int totalPrice) {
    final hasSelectedAnySeat =
        _departureSelectedSeatNames.isNotEmpty ||
        _returnSelectedSeatNames.isNotEmpty;

    return Material(
      elevation: 10.0,
      child: Container(
        padding: EdgeInsets.fromLTRB(
          20,
          15,
          20,
          15 + MediaQuery.of(context).padding.bottom / 2,
        ),
        decoration: BoxDecoration(color: Colors.white),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // MODIFIED: Hiển thị thông tin vé cho cả 2 chiều
            if (_departureSelectedSeatNames.isNotEmpty) ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "Vé chuyến đi:",
                    style: TextStyle(fontSize: 15, color: Colors.grey.shade700),
                  ),
                  Text(
                    "${_departureSelectedSeatNames.length} vé - Ghế: ${_departureSelectedSeatNames.join(", ")}",
                    style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold),
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
              SizedBox(height: 4),
            ],
            if (_isRoundTrip && _returnSelectedSeatNames.isNotEmpty) ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "Vé chuyến về:",
                    style: TextStyle(fontSize: 15, color: Colors.grey.shade700),
                  ),
                  Text(
                    "${_returnSelectedSeatNames.length} vé - Ghế: ${_returnSelectedSeatNames.join(", ")}",
                    style: TextStyle(fontSize: 14, fontWeight: FontWeight.bold),
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ],
            SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text(
                  "Tổng tiền:",
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
                Text(
                  "${ConvertMoney.currencyFormatter.format(totalPrice)} ₫",
                  style: TextStyle(
                    fontSize: 22,
                    fontWeight: FontWeight.bold,
                    color: Colors.orange.shade800,
                  ),
                ),
              ],
            ),
            SizedBox(height: 15),
            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor:
                      hasSelectedAnySeat ? Colors.orange : Colors.grey.shade400,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
                onPressed:
                    !hasSelectedAnySeat
                        ? null
                        : () {
                          // MODIFIED: Truyền dữ liệu của cả 2 chuyến đến PaymentScreen
                          // Lưu ý: Bạn cần cập nhật PaymentScreen để nhận các tham số mới này
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder:
                                  (context) => PaymentScreen(
                                    // Dữ liệu chuyến đi
                                    selectedRoute:
                                        "${widget.initialTripInfo.departureStation} → ${widget.initialTripInfo.arrivalStation}",
                                    rawTotalPrice: totalPrice,
                                    seatType:
                                        widget.initialTripInfo.vehicleType,
                                    time: DateFormat(
                                      'dd/MM/yyyy HH:mm',
                                      'vi_VN',
                                    ).format(
                                      widget.selectedDisplayDate,
                                    ), // Ví dụ, bạn có thể định dạng lại
                                    initialPickupLocation:
                                        widget.initialTripInfo.departureStation,
                                    initialDropoffLocation:
                                        widget.initialTripInfo.arrivalStation,
                                    selectedSeatsNames:
                                        _departureSelectedSeatNames,
                                    selectedSeatIds: _departureSelectedSeatIds,
                                    tripId: widget.initialTripInfo.tripId,

                                    // // NEW: Dữ liệu chuyến về
                                    returnTripId: widget.returnTripInfo?.tripId,
                                    selectedReturnSeatsNames:
                                        _returnSelectedSeatNames,
                                    selectedReturnSeatIds:
                                        _returnSelectedSeatIds,
                                    initialPickupLocationEnd:
                                        widget
                                            .returnTripInfo
                                            ?.departureStation ??
                                        "",
                                    initialDropoffLocationEnd:
                                        widget.returnTripInfo?.arrivalStation ??
                                        "",

                                    // Các tham số khác
                                    fullRoute: widget.initialTripInfo.fullRoute,
                                    estimatedDistance:
                                        widget
                                            .initialTripInfo
                                            .estimatedDistance,
                                  ),
                            ),
                          );
                        },
                child: Text(
                  "Tiếp tục",
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLegendItem(Color bgColor, Color borderColor, String text) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 16,
          height: 16,
          decoration: BoxDecoration(
            color: bgColor,
            border: Border.all(color: borderColor),
            borderRadius: BorderRadius.circular(4),
          ),
        ),
        SizedBox(width: 6),
        Text(text, style: TextStyle(fontSize: 12)),
      ],
    );
  }
}
