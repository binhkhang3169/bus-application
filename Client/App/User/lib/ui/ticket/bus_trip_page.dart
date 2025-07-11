// lib/ui/ticket/bus_trip_page.dart
import 'package:caoky/global/convert_money.dart';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/api_response.dart';
import 'package:caoky/models/api_response1.dart';
import 'package:caoky/models/trip/trip.dart';
import 'package:caoky/models/trip/trip_info.dart';
import 'package:caoky/services/api_trip_service.dart';
import 'package:caoky/services/dio_client.dart';
import 'package:caoky/ui/ticket/seat_selection_page.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class BusTripScreen extends StatefulWidget {
  final String startStationName;
  final int startStationId;
  final String endStationName;
  final int endStationId;
  final String departureDateString; // dd/MM/yyyy (for display)
  final String departureApiDateString; // yyyy-MM-dd (for API)
  final String departureDayOfWeek;
  final int ticketCount;
  final bool isRoundTrip;
  final String? endDateString; // dd/MM/yyyy (for display)
  final String? endDateApiDateString; // yyyy-MM-dd (for API)

  const BusTripScreen({
    super.key,
    required this.startStationName,
    required this.startStationId,
    required this.endStationName,
    required this.endStationId,
    required this.departureDateString,
    required this.departureApiDateString,
    required this.departureDayOfWeek,
    required this.ticketCount,
    required this.isRoundTrip,
    this.endDateString,
    this.endDateApiDateString,
  });

  @override
  _BusTripScreenState createState() => _BusTripScreenState();
}

class _BusTripScreenState extends State<BusTripScreen> {
  late ApiTripService _apiService;
  late String _selectedRouteDisplay;
  late DateTime _currentDisplayDate; // For the DateSelectionWidget
  late String _currentApiDate; // For API calls, format yyyy-MM-dd

  // Round trip state management
  bool _isSelectingDepartureTrip = true;
  TripInfo? _selectedDepartureTrip;

  List<TripInfo> _foundTripsData = [];
  bool _isLoading = true;
  String? _sortPrice;
  String? _seatType;
  String? _timeRange;

  final ScrollController _scrollController = ScrollController();
  bool _showHeader = true;

  final List<Map<String, List<String>>> _filterOptions = [
    {
      "Giá": ["Tăng dần", "Giảm dần"],
    },
    {
      "Loại ghế": ["Limousine", "Giường"],
    },
    {
      "Giờ": [
        "Buổi sáng (00:00 - 11:59)",
        "Buổi chiều (12:00 - 17:59)",
        "Buổi tối (18:00 - 23:59)",
      ],
    },
  ];

  @override
  void initState() {
    super.initState();
    _apiService = ApiTripService(DioClient.createDio());
    _updateUIForSelectionMode(); // Set initial UI based on selection mode
    _scrollController.addListener(_onScroll);
    _fetchTrips();
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (mounted) {
      final show = _scrollController.offset <= 50;
      if (show != _showHeader) {
        setState(() {
          _showHeader = show;
        });
      }
    }
  }

  // Update UI text and dates based on whether we are selecting departure or return
  void _updateUIForSelectionMode() {
    if (_isSelectingDepartureTrip) {
      _selectedRouteDisplay = "${widget.startStationName} → ${widget.endStationName}";
      _currentApiDate = widget.departureApiDateString;
      try {
        _currentDisplayDate = DateFormat('dd/MM/yyyy', 'vi_VN').parse(widget.departureDateString);
      } catch (e) {
        _currentDisplayDate = DateTime.now();
        print("Error parsing departure date: $e");
      }
    } else {
      // Swapping for return trip
      _selectedRouteDisplay = "${widget.endStationName} → ${widget.startStationName}";
      _currentApiDate = widget.endDateApiDateString!;
      try {
        _currentDisplayDate = DateFormat('dd/MM/yyyy', 'vi_VN').parse(widget.endDateString!);
      } catch (e) {
        _currentDisplayDate = DateTime.now();
        print("Error parsing return date: $e");
      }
    }
  }

  Future<void> _fetchTrips({String? apiDateOverride}) async {
    if (!mounted) return;
    setState(() {
      _isLoading = true;
      _foundTripsData.clear();
    });

    final String startStation = _isSelectingDepartureTrip ? widget.startStationName : widget.endStationName;
    final int startStationId = _isSelectingDepartureTrip ? widget.startStationId : widget.endStationId;
    final String endStation = _isSelectingDepartureTrip ? widget.endStationName : widget.startStationName;
    final int endStationId = _isSelectingDepartureTrip ? widget.endStationId : widget.startStationId;
    final String dateToFetch = apiDateOverride ?? _currentApiDate;

    try {
      final ApiResponse1<List<TripInfo>> response = await _apiService.getTrips(
        startStation,
        startStationId,
        endStation,
        endStationId,
        dateToFetch,
        widget.ticketCount,
      );

      if (mounted) {
        if (response.code == 200 && response.data != null) {
          setState(() {
            _foundTripsData = response.data!;
            if (_foundTripsData.isEmpty) {
              final formattedDate = DateFormat('dd/MM/yyyy').format(DateFormat('yyyy-MM-dd').parse(dateToFetch));
              ToastUtils.show("Không tìm thấy chuyến xe nào cho ngày $formattedDate.");
            }
          });
        } else {
          ToastUtils.show("Lỗi ${response.code}: ${response.message}");
        }
      }
    } on DioException catch (e) {
      print("API DioException: $e");
      if (mounted) {
        ToastUtils.show("Lỗi kết nối: ${e.message ?? 'Không thể tải dữ liệu'}");
      }
    } catch (e) {
      print("Lỗi không xác định khi gọi API: $e");
      if (mounted) {
        ToastUtils.show("Đã xảy ra lỗi không mong muốn.");
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  void _onTripSelected(TripInfo selectedTrip) {
    if (!widget.isRoundTrip) {
      // One-way trip: navigate directly
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => SeatSelectionScreen(
            initialTripInfo: selectedTrip,
            // returnTripInfo: null, // No return trip
            selectedDisplayDate: _currentDisplayDate,
          ),
        ),
      );
    } else {
      // Round trip logic
      if (_isSelectingDepartureTrip) {
        // Step 1: Departure trip selected, now select return trip
        setState(() {
          _selectedDepartureTrip = selectedTrip;
          _isSelectingDepartureTrip = false; // Switch to selecting return trip
          _updateUIForSelectionMode(); // Update AppBar and dates for return trip
          _resetFilters(); // Reset filters for the new list
        });
        _fetchTrips(); // Fetch return trips
      } else {
        // Step 2: Return trip selected, navigate to seat selection
        final TripInfo returnTrip = selectedTrip;
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => SeatSelectionScreen(
              initialTripInfo: _selectedDepartureTrip!,
              returnTripInfo: returnTrip,
              selectedDisplayDate: DateFormat('dd/MM/yyyy', 'vi_VN').parse(widget.departureDateString),
              returnDisplayDate: DateFormat('dd/MM/yyyy', 'vi_VN').parse(widget.endDateString!),
            ),
          ),
        );
      }
    }
  }

  void _onDateSelectedFromWidget(DateTime newSelectedDate) {
    if (!mounted) return;
    final newApiDate = DateFormat('yyyy-MM-dd').format(newSelectedDate);

    setState(() {
      _currentDisplayDate = newSelectedDate;
      _currentApiDate = newApiDate;
      // Update the correct date string for the current selection phase
      if (_isSelectingDepartureTrip) {
        // This is complex as it deviates from initial widget data.
        // For simplicity, we just refetch. A more robust solution might update widget data.
      } else {
        // Same complexity for return date.
      }
    });
    _fetchTrips(apiDateOverride: newApiDate);
  }

  void _resetFilters() {
    setState(() {
      _sortPrice = null;
      _seatType = null;
      _timeRange = null;
    });
  }

  List<TripInfo> _getFilteredAndSortedTrips() {
    List<TripInfo> filteredList = List.from(_foundTripsData);

    if (_seatType != null) {
      filteredList = filteredList.where((trip) => trip.vehicleType.toLowerCase().contains(_seatType!.toLowerCase())).toList();
    }

    if (_timeRange != null) {
      filteredList = filteredList.where((trip) {
        try {
          final departureHour = int.parse(trip.departureTime.split(':')[0]);
          if (_timeRange == "Buổi sáng (00:00 - 11:59)" && (departureHour >= 0 && departureHour < 12)) return true;
          if (_timeRange == "Buổi chiều (12:00 - 17:59)" && (departureHour >= 12 && departureHour < 18)) return true;
          if (_timeRange == "Buổi tối (18:00 - 23:59)" && (departureHour >= 18 && departureHour < 24)) return true;
          return false;
        } catch (e) {
          return false;
        }
      }).toList();
    }

    if (_sortPrice != null) {
      filteredList.sort((a, b) {
        if (_sortPrice == "Tăng dần") return a.price.compareTo(b.price);
        if (_sortPrice == "Giảm dần") return b.price.compareTo(a.price);
        return 0;
      });
    }
    return filteredList;
  }

  @override
  Widget build(BuildContext context) {
    final displayTrips = _getFilteredAndSortedTrips();
    String displayDateForAppBar = DateFormat('EEEE, dd/MM/yyyy', 'vi_VN').format(_currentDisplayDate);

    return Scaffold(
      backgroundColor: Colors.grey[100],
      appBar: AppBar(
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () {
            // If selecting return trip, go back to selecting departure trip
            if (widget.isRoundTrip && !_isSelectingDepartureTrip) {
              setState(() {
                _isSelectingDepartureTrip = true;
                _selectedDepartureTrip = null;
                _updateUIForSelectionMode();
                _resetFilters();
              });
              _fetchTrips();
            } else {
              Navigator.pop(context);
            }
          },
        ),
        title: Column(
          children: [
            Text(
              _selectedRouteDisplay,
              style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold),
              overflow: TextOverflow.ellipsis,
            ),
            Text(
              displayDateForAppBar,
              style: TextStyle(color: Colors.white, fontSize: 13),
            ),
          ],
        ),
        centerTitle: true,
        backgroundColor: Colors.blueAccent,
        elevation: _showHeader ? 4.0 : 0.0,
      ),
      body: Column(
        children: [
          if (widget.isRoundTrip)
            RoundTripProgressIndicator(isSelectingDeparture: _isSelectingDepartureTrip),
          AnimatedContainer(
            duration: Duration(milliseconds: 200),
            height: _showHeader ? 80 : 0,
            child: AnimatedOpacity(
              opacity: _showHeader ? 1.0 : 0.0,
              duration: Duration(milliseconds: 200),
              child: DateSelectionWidget(
                currentDate: DateTime.now(),
                selectedDate: _currentDisplayDate,
                onDateSelected: _onDateSelectedFromWidget,
                initialScrollDate: _currentDisplayDate,
              ),
            ),
          ),
          AnimatedContainer(
            duration: Duration(milliseconds: 200),
            height: _showHeader ? 50 : 0,
            child: AnimatedOpacity(
              opacity: _showHeader ? 1.0 : 0.0,
              duration: Duration(milliseconds: 200),
              child: FilterWidget(
                filterOptions: _filterOptions,
                selectedPriceSort: _sortPrice,
                selectedSeatType: _seatType,
                selectedTimeRange: _timeRange,
                onPriceSort: _showPriceSortBottomSheet,
                onSeatType: _showSeatTypeBottomSheet,
                onTimeRange: _showTimeRangeBottomSheet,
              ),
            ),
          ),
          Expanded(
            child: _isLoading
                ? Center(child: CircularProgressIndicator())
                : displayTrips.isEmpty
                    ? Center(
                        child: Text(
                          "Không có chuyến xe nào phù hợp.",
                          style: TextStyle(fontSize: 16, color: Colors.grey[700]),
                        ),
                      )
                    : ListView.builder(
                        controller: _scrollController,
                        itemCount: displayTrips.length,
                        padding: EdgeInsets.only(top: _showHeader ? 0 : 8, bottom: 8),
                        itemBuilder: (context, index) {
                          final tripInfo = displayTrips[index];
                          return BusTripCardWidget(
                            tripInfo1: tripInfo,
                            onTripSelected: _onTripSelected,
                          );
                        },
                      ),
          ),
        ],
      ),
    );
  }

  // --- Bottom Sheet Methods ---
  // (No changes needed for these methods, they are kept as they are)
  void _showPriceSortBottomSheet() {
    showModalBottomSheet(
      backgroundColor: Colors.white,
      context: context,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _filterOptions[0].keys.first,
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: Icon(Icons.close),
                        onPressed: () => Navigator.pop(context),
                      ),
                    ],
                  ),
                  Divider(),
                  ..._filterOptions[0]["Giá"]!.map((option) {
                    return ListTile(
                      title: Text(option),
                      trailing: _sortPrice == option
                          ? Icon(
                              Icons.check_circle,
                              color: Colors.blueAccent,
                            )
                          : Icon(Icons.circle_outlined),
                      onTap: () {
                        setState(() {
                          _sortPrice = (_sortPrice == option) ? null : option;
                        }); // Allow unselecting
                        Navigator.pop(context);
                      },
                    );
                  }).toList(),
                  if (_sortPrice != null) SizedBox(height: 10),
                  if (_sortPrice != null)
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.redAccent,
                        ),
                        onPressed: () {
                          setState(() {
                            _sortPrice = null;
                          });
                          Navigator.pop(context);
                        },
                        child: Text(
                          "Bỏ lọc giá",
                          style: TextStyle(color: Colors.white),
                        ),
                      ),
                    ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  void _showSeatTypeBottomSheet() {
    showModalBottomSheet(
      backgroundColor: Colors.white,
      context: context,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _filterOptions[1].keys.first,
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: Icon(Icons.close),
                        onPressed: () => Navigator.pop(context),
                      ),
                    ],
                  ),
                  Divider(),
                  ..._filterOptions[1]["Loại ghế"]!.map((option) {
                    return ListTile(
                      title: Text(option),
                      trailing: _seatType == option
                          ? Icon(
                              Icons.check_circle,
                              color: Colors.blueAccent,
                            )
                          : Icon(Icons.circle_outlined),
                      onTap: () {
                        setState(() {
                          _seatType = (_seatType == option) ? null : option;
                        });
                        Navigator.pop(context);
                      },
                    );
                  }).toList(),
                  if (_seatType != null) SizedBox(height: 10),
                  if (_seatType != null)
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.redAccent,
                        ),
                        onPressed: () {
                          setState(() {
                            _seatType = null;
                          });
                          Navigator.pop(context);
                        },
                        child: Text(
                          "Bỏ lọc loại ghế",
                          style: TextStyle(color: Colors.white),
                        ),
                      ),
                    ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  void _showTimeRangeBottomSheet() {
    showModalBottomSheet(
      backgroundColor: Colors.white,
      context: context,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _filterOptions[2].keys.first,
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: Icon(Icons.close),
                        onPressed: () => Navigator.pop(context),
                      ),
                    ],
                  ),
                  Divider(),
                  ..._filterOptions[2]["Giờ"]!.map((option) {
                    return ListTile(
                      title: Text(option),
                      trailing: _timeRange == option
                          ? Icon(
                              Icons.check_circle,
                              color: Colors.blueAccent,
                            )
                          : Icon(Icons.circle_outlined),
                      onTap: () {
                        setState(() {
                          _timeRange = (_timeRange == option) ? null : option;
                        });
                        Navigator.pop(context);
                      },
                    );
                  }).toList(),
                  if (_timeRange != null) SizedBox(height: 10),
                  if (_timeRange != null)
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.redAccent,
                        ),
                        onPressed: () {
                          setState(() {
                            _timeRange = null;
                          });
                          Navigator.pop(context);
                        },
                        child: Text(
                          "Bỏ lọc giờ",
                          style: TextStyle(color: Colors.white),
                        ),
                      ),
                    ),
                ],
              ),
            );
          },
        );
      },
    );
  }
}

// New Widget to show Round Trip Progress
class RoundTripProgressIndicator extends StatelessWidget {
  final bool isSelectingDeparture;

  const RoundTripProgressIndicator({
    Key? key,
    required this.isSelectingDeparture,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 12.0),
      color: Colors.blue.shade50,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          _buildStep(
            icon: Icons.arrow_upward,
            label: "Chọn chuyến đi",
            isActive: isSelectingDeparture,
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Divider(
                color: isSelectingDeparture
                    ? Colors.grey.shade300
                    : Colors.blueAccent,
                thickness: 2,
              ),
            ),
          ),
          _buildStep(
            icon: Icons.arrow_downward,
            label: "Chọn chuyến về",
            isActive: !isSelectingDeparture,
          ),
        ],
      ),
    );
  }

  Widget _buildStep({
    required IconData icon,
    required String label,
    required bool isActive,
  }) {
    final color = isActive ? Colors.blueAccent : Colors.grey.shade600;
    final fontWeight = isActive ? FontWeight.bold : FontWeight.normal;

    return Column(
      children: [
        Icon(icon, color: color),
        SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: color,
            fontWeight: fontWeight,
            fontSize: 13,
          ),
        ),
      ],
    );
  }
}


// --- Extracted Widgets ---
// (No changes needed for these widgets, they are kept as they are)

// Extracted DateSelectionWidget
class DateSelectionWidget extends StatefulWidget {
  final DateTime currentDate; // The actual current date (e.g., DateTime.now())
  final DateTime selectedDate; // The currently selected date in the UI
  final DateTime initialScrollDate; // The date to scroll to initially
  final Function(DateTime) onDateSelected;

  DateSelectionWidget({
    required this.currentDate,
    required this.selectedDate,
    required this.initialScrollDate,
    required this.onDateSelected,
  });

  @override
  _DateSelectionWidgetState createState() => _DateSelectionWidgetState();
}

class _DateSelectionWidgetState extends State<DateSelectionWidget> {
  late ScrollController _dateScrollController;
  final int _numberOfDaysToShow = 60; // Show dates for the next 60 days

  @override
  void initState() {
    super.initState();
    _dateScrollController = ScrollController();
    WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToInitialDate());
  }

  @override
  void didUpdateWidget(covariant DateSelectionWidget oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.initialScrollDate != oldWidget.initialScrollDate) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _scrollToInitialDate());
    }
  }

  void _scrollToInitialDate() {
    if (mounted && _dateScrollController.hasClients) {
      DateTime startDate = DateTime(widget.currentDate.year, widget.currentDate.month, widget.currentDate.day);
      int initialIndex = widget.initialScrollDate.difference(startDate).inDays;

      if (initialIndex < 0) initialIndex = 0;
      if (initialIndex >= _numberOfDaysToShow) initialIndex = _numberOfDaysToShow -1;

      double scrollOffset = initialIndex * 73.0; // 65 width + 8 horizontal margin

      if (scrollOffset > _dateScrollController.position.maxScrollExtent) {
        scrollOffset = _dateScrollController.position.maxScrollExtent;
      }
      _dateScrollController.animateTo(
        scrollOffset,
        duration: Duration(milliseconds: 400),
        curve: Curves.easeInOut,
      );
    }
  }


  @override
  void dispose() {
    _dateScrollController.dispose();
    super.dispose();
  }

  String _getDayOfWeek(int weekday) {
    // Using intl for localization
    return DateFormat('E', 'vi_VN').format(DateTime(2023, 1, 1 + weekday));
  }

  @override
  Widget build(BuildContext context) {
    DateTime startDate = DateTime(
      widget.currentDate.year,
      widget.currentDate.month,
      widget.currentDate.day,
    );

    return Container(
      height: 70,
      color: Colors.white,
      padding: EdgeInsets.symmetric(vertical: 8),
      child: ListView.builder(
        controller: _dateScrollController,
        scrollDirection: Axis.horizontal,
        itemCount: _numberOfDaysToShow,
        itemBuilder: (context, index) {
          DateTime date = startDate.add(Duration(days: index));
          String dayOfWeek = _getDayOfWeek(date.weekday);
          String formattedDate = DateFormat('dd/MM').format(date);
          bool isSelected =
              date.year == widget.selectedDate.year &&
              date.month == widget.selectedDate.month &&
              date.day == widget.selectedDate.day;

          return GestureDetector(
            onTap: () => widget.onDateSelected(date),
            child: Container(
              width: 65,
              margin: EdgeInsets.symmetric(horizontal: 4),
              decoration: BoxDecoration(
                color: isSelected ? Colors.blueAccent : Colors.grey[200],
                borderRadius: BorderRadius.circular(8),
                border: isSelected
                    ? Border.all(
                        color: Colors.blueAccent.shade700,
                        width: 1.5,
                      )
                    : null,
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    dayOfWeek,
                    style: TextStyle(
                      color: isSelected ? Colors.white : Colors.black87,
                      fontWeight:
                          isSelected ? FontWeight.bold : FontWeight.normal,
                      fontSize: 13,
                    ),
                  ),
                  SizedBox(height: 4),
                  Text(
                    formattedDate,
                    style: TextStyle(
                      color: isSelected ? Colors.white : Colors.black54,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

// Extracted FilterWidget
class FilterWidget extends StatelessWidget {
  final List<Map<String, List<String>>> filterOptions;
  final String? selectedPriceSort;
  final String? selectedSeatType;
  final String? selectedTimeRange;
  final VoidCallback onPriceSort;
  final VoidCallback onSeatType;
  final VoidCallback onTimeRange;

  FilterWidget({
    required this.filterOptions,
    this.selectedPriceSort,
    this.selectedSeatType,
    this.selectedTimeRange,
    required this.onPriceSort,
    required this.onSeatType,
    required this.onTimeRange,
  });

  Widget _buildFilterButton(
    BuildContext context,
    String title,
    String? selectedValue,
    VoidCallback onTap,
  ) {
    bool isSelected = selectedValue != null;
    return Expanded(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 4.0),
        child: OutlinedButton.icon(
          onPressed: onTap,
          icon: Icon(
            Icons.filter_list,
            size: 16,
            color: isSelected ? Colors.blueAccent : Colors.grey[700],
          ),
          label: Text(
            selectedValue ?? title,
            style: TextStyle(
              fontSize: 12,
              color: isSelected ? Colors.blueAccent : Colors.grey[700],
              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
            ),
            overflow: TextOverflow.ellipsis,
          ),
          style: OutlinedButton.styleFrom(
            padding: EdgeInsets.symmetric(horizontal: 8, vertical: 6),
            backgroundColor:
                isSelected ? Colors.blueAccent.withOpacity(0.1) : Colors.white,
            side: BorderSide(
              color: isSelected ? Colors.blueAccent : Colors.grey.shade400,
            ),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(20),
            ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
      color: Colors.white,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          _buildFilterButton(
            context,
            filterOptions[0].keys.first,
            selectedPriceSort,
            onPriceSort,
          ),
          _buildFilterButton(
            context,
            filterOptions[1].keys.first,
            selectedSeatType,
            onSeatType,
          ),
          _buildFilterButton(
            context,
            filterOptions[2].keys.first,
            selectedTimeRange,
            onTimeRange,
          ),
        ],
      ),
    );
  }
}

// Extracted BusTripCardWidget
class BusTripCardWidget extends StatelessWidget {
  final TripInfo tripInfo1;
  final Function(TripInfo) onTripSelected;

  BusTripCardWidget({required this.tripInfo1, required this.onTripSelected});

  @override
  Widget build(BuildContext context) {
    final tripInfo = tripInfo1;
    String formattedDepartureTime = "N/A";
    String formattedArrivalTime = "N/A";
    try {
      if (tripInfo.departureTime.length >= 5) {
        formattedDepartureTime = tripInfo.departureTime.substring(0, 5);
      }
      if (tripInfo.arrivalTime.length >= 5) {
        formattedArrivalTime = tripInfo.arrivalTime.substring(0, 5);
      }
    } catch (e) {
      print("Error formatting time: $e");
    }

    return Card(
      color: Colors.white,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 3,
      margin: EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      child: InkWell(
        onTap: () => onTripSelected(tripInfo),
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              width: double.infinity,
              padding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(12),
                  topRight: Radius.circular(12),
                ),
                gradient: LinearGradient(
                  colors: [
                    Colors.blueAccent.shade700,
                    Colors.blueAccent.shade400,
                  ],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        tripInfo.vehicleType,
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 15,
                          color: Colors.white,
                        ),
                      ),
                      SizedBox(height: 4),
                      Row(
                        children: [
                          Text(
                            formattedDepartureTime,
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: Colors.white,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6.0,
                            ),
                            child: Icon(
                              Icons.arrow_forward_rounded,
                              color: Colors.white,
                              size: 18,
                            ),
                          ),
                          Text(
                            formattedArrivalTime,
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: Colors.white,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  Text(
                    "${ConvertMoney.currencyFormatter.format(tripInfo.price)} ₫",
                    style: TextStyle(
                      color: Colors.yellowAccent,
                      fontWeight: FontWeight.bold,
                      fontSize: 18,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        "Còn " + tripInfo.stock.toString() + " vé",
                        style: TextStyle(
                          color: Colors.green.shade700,
                          fontWeight: FontWeight.w600,
                          fontSize: 13,
                        ),
                      ),
                      Text(
                        tripInfo.license,
                        style: TextStyle(
                          fontSize: 13,
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 12),
                  _buildRouteDetailRow(
                    Icons.adjust,
                    Colors.green,
                    tripInfo.departureStation,
                    "Thời gian dự kiến: ${tripInfo.estimatedTime}",
                  ),
                  _buildDashedLine(context),
                  _buildRouteDetailRow(
                    Icons.location_on_outlined,
                    Colors.red,
                    tripInfo.arrivalStation,
                    "Quãng đường: ${tripInfo.estimatedDistance}",
                  ),
                  Align(
                    alignment: Alignment.centerRight,
                    child: TextButton(
                      onPressed: () {
                        ToastUtils.show(
                          "Xem lịch trình chi tiết cho chuyến ${tripInfo.tripId}",
                        );
                      },
                      style: TextButton.styleFrom(
                        foregroundColor: Colors.blueAccent.shade700,
                      ),
                      child: Text(
                        "Chi tiết",
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRouteDetailRow(
    IconData icon,
    Color iconColor,
    String stationName,
    String subtitle,
  ) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(icon, color: iconColor, size: 20),
        SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                stationName,
                style: TextStyle(fontWeight: FontWeight.w600, fontSize: 14),
              ),
              SizedBox(height: 2),
              Text(
                subtitle,
                style: TextStyle(color: Colors.grey[700], fontSize: 12),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildDashedLine(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(
        left: 9.0,
        top: 4,
        bottom: 4,
      ),
      child: CustomPaint(
        size: Size(
          2,
          25,
        ),
        painter: DashedLinePainter(color: Colors.grey.shade400),
      ),
    );
  }
}

// Custom Painter for Dashed Line
class DashedLinePainter extends CustomPainter {
  final Color color;
  DashedLinePainter({this.color = Colors.grey});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    const double dashHeight = 4;
    const double dashSpace = 3;
    double startY = 0;

    while (startY < size.height) {
      canvas.drawLine(
        Offset(0, startY),
        Offset(0, startY + dashHeight),
        paint,
      );
      startY += dashHeight + dashSpace;
    }
  }

  @override
  bool shouldRepaint(CustomPainter oldDelegate) => false;
}