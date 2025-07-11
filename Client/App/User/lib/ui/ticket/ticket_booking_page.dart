import 'package:cached_network_image/cached_network_image.dart';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/trip/address_info.dart';
import 'package:caoky/ui/ticket/bus_trip_page.dart';
import 'package:caoky/ui/ticket/localation_selected_page.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:intl/date_symbol_data_local.dart';

class TicketBookingPage extends StatefulWidget {
  @override
  _TicketBookingPageState createState() => _TicketBookingPageState();
}

class _TicketBookingPageState extends State<TicketBookingPage> {
  int _currentIndex = 0;
  bool _isRoundTrip = false;
  final List<String> recentSearches = [
    'TP. Hồ Chí Minh - Bến Tre',
    'TP. Hồ Chí Minh - Đồng Tháp',
  ];
  final CarouselController _carouselController = CarouselController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: TicketBookingContent(
        isRoundTrip: _isRoundTrip,
        recentSearches: recentSearches,
        currentIndex: _currentIndex,
        carouselController: _carouselController,
        onRoundTripChanged: (value) {
          setState(() {
            _isRoundTrip = value;
          });
        },
        onClearSearches: () {
          setState(() {
            recentSearches.clear();
          });
        },
        onPageChanged: (index, reason) {
          setState(() {
            _currentIndex = index;
          });
        },
      ),
    );
  }
}

class TicketBookingContent extends StatefulWidget {
  final bool isRoundTrip;
  final List<String> recentSearches;
  final int currentIndex;
  final ValueChanged<bool> onRoundTripChanged;
  final VoidCallback onClearSearches;
  final Function(int index, CarouselPageChangedReason reason) onPageChanged;
  final CarouselController carouselController;

  const TicketBookingContent({
    super.key,
    required this.isRoundTrip,
    required this.recentSearches,
    required this.currentIndex,
    required this.onRoundTripChanged,
    required this.onClearSearches,
    required this.onPageChanged,
    required this.carouselController,
  });

  static const List<String> imgBanner = [
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/Thiet_ke_chua_co_ten_1ca8aaade3/Thiet_ke_chua_co_ten_1ca8aaade3.png',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/dat_ve_xe_khach_giam_300_K_Futa_599x337_15730c90c4_7a6d3012ea/dat_ve_xe_khach_giam_300_K_Futa_599x337_15730c90c4_7a6d3012ea.jpg',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/599_x_337_px_266806d0c2/599_x_337_px_266806d0c2.jpg',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/599_x_337_px_897c640899/599_x_337_px_897c640899.png',
  ];

  static const List<String> imgBanner1 = [
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/2_343_x_184_px_f365e0f9c8/2_343_x_184_px_f365e0f9c8.png',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/599_c4c05c0b3a/599_c4c05c0b3a.png',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/599_x_337_px_897c640899/599_x_337_px_897c640899.png',
    'https://cdn.futabus.vn/futa-busline-web-cms-prod/599_x_337_px_57901c4bd1/599_x_337_px_57901c4bd1.png',
  ];

  @override
  State<TicketBookingContent> createState() => _TicketBookingContentState();
}

class _TicketBookingContentState extends State<TicketBookingContent> {
  AddressInfo? startLocation;
  AddressInfo? endLocation;
  DateTime? selectedStartDate;
  DateTime? selectedEndDate;
  bool _localeInitialized = false;
  int? selectedTicketQuantity = 1;

  @override
  void initState() {
    super.initState();
    initializeDateFormatting('vi_VN', null).then((_) {
      if (mounted) {
        setState(() {
          _localeInitialized = true;
        });
      }
    });
  }

  void _findRoute() {
    if (startLocation == null) {
      ToastUtils.show("Vui lòng chọn điểm đi");
      return;
    }
    if (endLocation == null) {
      ToastUtils.show("Vui lòng chọn điểm đến");
      return;
    }
    if (selectedStartDate == null) {
      ToastUtils.show("Vui lòng chọn ngày đi");
      return;
    }
    if (widget.isRoundTrip && selectedEndDate == null) {
      ToastUtils.show("Vui lòng chọn ngày về cho chuyến khứ hồi");
      return;
    }
    if (widget.isRoundTrip && selectedEndDate != null && selectedEndDate!.isBefore(selectedStartDate!)) {
      ToastUtils.show("Ngày về không thể trước ngày đi");
      return;
    }
    if (selectedTicketQuantity == null || selectedTicketQuantity! <= 0) {
      ToastUtils.show("Vui lòng chọn số lượng vé hợp lệ");
      return;
    }

    String apiDepartureDate = DateFormat('yyyy-MM-dd').format(selectedStartDate!);
    String displayDepartureDate = DateFormat('dd/MM/yyyy', 'vi_VN').format(selectedStartDate!);
    String displayDepartureDayOfWeek = DateFormat('EEEE', 'vi_VN').format(selectedStartDate!);

    String searchKey = "${startLocation!.name} - ${endLocation!.name}${widget.isRoundTrip ? ' (Khứ hồi)' : ''}";
    if (!widget.recentSearches.contains(searchKey)) {
      setState(() {
        widget.recentSearches.add(searchKey);
      });
    }

    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => BusTripScreen(
          startStationName: startLocation!.name,
          startStationId: startLocation!.id,
          endStationName: endLocation!.name,
          endStationId: endLocation!.id,
          departureDateString: displayDepartureDate,
          departureApiDateString: apiDepartureDate,
          departureDayOfWeek: displayDepartureDayOfWeek,
          ticketCount: selectedTicketQuantity!,
          isRoundTrip: widget.isRoundTrip,
          endDateString: widget.isRoundTrip ? DateFormat('dd/MM/yyyy', 'vi_VN').format(selectedEndDate!) : null,
          endDateApiDateString: widget.isRoundTrip ? DateFormat('yyyy-MM-dd', 'vi_VN').format(selectedEndDate!) : null,
        ),
      ),
    );
  }

  Future<void> _pickStartDate() async {
    if (!_localeInitialized) return;
    final now = DateTime.now();
    final DateTime firstSelectableDate = DateTime(now.year, now.month, now.day);

    final picked = await showDatePicker(
      context: context,
      locale: const Locale('vi', 'VN'),
      initialDate: selectedStartDate ?? firstSelectableDate,
      firstDate: firstSelectableDate,
      lastDate: DateTime(now.year + 1, now.month, now.day),
      builder: (context, child) => Theme(
        data: ThemeData.light().copyWith(
          colorScheme: ColorScheme.light(primary: Colors.blueAccent),
          textButtonTheme: TextButtonThemeData(
            style: TextButton.styleFrom(foregroundColor: Colors.blueAccent),
          ),
        ),
        child: child!,
      ),
    );
    if (picked != null && picked != selectedStartDate) {
      setState(() {
        selectedStartDate = picked;
        if (widget.isRoundTrip && selectedEndDate != null && selectedEndDate!.isBefore(picked)) {
          selectedEndDate = null;
        }
      });
    }
  }

  Future<void> _pickEndDate() async {
    if (!_localeInitialized) return;
    if (!widget.isRoundTrip) {
      ToastUtils.show("Chỉ chọn ngày về cho chuyến khứ hồi.");
      return;
    }
    if (selectedStartDate == null) {
      ToastUtils.show("Vui lòng chọn ngày đi trước.");
      return;
    }

    final picked = await showDatePicker(
      context: context,
      locale: const Locale('vi', 'VN'),
      initialDate: selectedEndDate ?? selectedStartDate!.add(Duration(days: 1)),
      firstDate: selectedStartDate!,
      lastDate: DateTime(selectedStartDate!.year + 1, selectedStartDate!.month, selectedStartDate!.day),
      builder: (context, child) => Theme(
        data: ThemeData.light().copyWith(
          colorScheme: ColorScheme.light(primary: Colors.blueAccent),
          textButtonTheme: TextButtonThemeData(
            style: TextButton.styleFrom(foregroundColor: Colors.blueAccent),
          ),
        ),
        child: child!,
      ),
    );

    if (picked != null && picked != selectedEndDate) {
      setState(() {
        selectedEndDate = picked;
      });
    }
  }

  void _showTicketQuantityBottomSheet() {
    showModalBottomSheet(
      backgroundColor: Colors.white,
      context: context,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20.0)),
      ),
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter modalSetState) {
            return Container(
              padding: EdgeInsets.symmetric(vertical: 20, horizontal: 16),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    'Chọn số lượng vé',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  SizedBox(height: 20),
                  ListView.builder(
                    shrinkWrap: true,
                    itemCount: 5,
                    itemBuilder: (context, index) {
                      final quantity = index + 1;
                      return ListTile(
                        title: Text('$quantity vé'),
                        onTap: () {
                          setState(() {
                            selectedTicketQuantity = quantity;
                          });
                          Navigator.pop(context);
                        },
                        trailing: selectedTicketQuantity == quantity
                            ? Icon(Icons.check_circle, color: Colors.blueAccent)
                            : Icon(Icons.radio_button_unchecked),
                      );
                    },
                  ),
                  SizedBox(height: 10),
                ],
              ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Column(
          children: [
            Container(
              height: 200,
              color: Colors.blueAccent,
              child: AppBar(
                backgroundColor: Colors.transparent,
                elevation: 0,
                leading: IconButton(
                  icon: Icon(Icons.arrow_back, color: Colors.white),
                  onPressed: () {
                    if (Navigator.canPop(context)) {
                      Navigator.pop(context);
                    }
                  },
                ),
                title: Text(
                  'Mua vé xe',
                  style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
                ),
                centerTitle: true,
              ),
            ),
            Expanded(
              child: Container(
                color: Colors.white,
                child: SingleChildScrollView(
                  physics: BouncingScrollPhysics(),
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        SizedBox(
                          height: MediaQuery.of(context).padding.top + kToolbarHeight - 20 + 100,
                        ),
                        if (widget.recentSearches.isNotEmpty) ...[
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'Tìm kiếm gần đây',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold,
                                  color: Colors.black87,
                                ),
                              ),
                              TextButton(
                                onPressed: widget.onClearSearches,
                                child: Text(
                                  'Xóa lịch sử',
                                  style: TextStyle(fontSize: 14, color: Colors.redAccent),
                                ),
                              ),
                            ],
                          ),
                          SizedBox(
                            height: 50,
                            child: ListView.builder(
                              scrollDirection: Axis.horizontal,
                              itemCount: widget.recentSearches.length,
                              itemBuilder: (context, index) {
                                return Padding(
                                  padding: const EdgeInsets.only(right: 8.0),
                                  child: Card(
                                    elevation: 2,
                                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                                    child: InkWell(
                                      onTap: () {
                                        ToastUtils.show("Nạp lại tìm kiếm: ${widget.recentSearches[index]}");
                                      },
                                      borderRadius: BorderRadius.circular(8),
                                      child: Container(
                                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                                        decoration: BoxDecoration(
                                          color: Colors.white,
                                          borderRadius: BorderRadius.circular(8),
                                        ),
                                        child: Center(
                                          child: Text(
                                            widget.recentSearches[index],
                                            style: TextStyle(fontSize: 14, color: Colors.black87),
                                          ),
                                        ),
                                      ),
                                    ),
                                  ),
                                );
                              },
                            ),
                          ),
                          SizedBox(height: 20),
                        ],
                        _buildBanner(
                          TicketBookingContent.imgBanner,
                          widget.carouselController,
                          widget.onPageChanged,
                          widget.currentIndex,
                        ),
                        SizedBox(height: 20),
                        _buildBanner(
                          TicketBookingContent.imgBanner1,
                          widget.carouselController,
                          widget.onPageChanged,
                          widget.currentIndex,
                        ),
                        SizedBox(height: 20),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
        Positioned(
          top: kToolbarHeight + MediaQuery.of(context).padding.top - 10,
          left: 10,
          right: 10,
          child: Container(
            height: 270,
            child: Card(
              color: Colors.white,
              elevation: 8,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                          child: _locationColumn(
                            title: 'Điểm đi',
                            station: startLocation,
                            icon: Icons.my_location,
                            onTap: () async {
                              final result = await Navigator.push(
                                context,
                                MaterialPageRoute(builder: (context) => LocationSelectPage()),
                              );
                              if (result != null && result is AddressInfo) {
                                setState(() {
                                  startLocation = result;
                                });
                              }
                            },
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 0),
                          child: IconButton(
                            icon: Icon(Icons.swap_horiz, color: Colors.blueAccent, size: 25),
                            onPressed: () {
                              if (startLocation != null || endLocation != null) {
                                setState(() {
                                  final temp = startLocation;
                                  startLocation = endLocation;
                                  endLocation = temp;
                                });
                              }
                            },
                          ),
                        ),
                        Expanded(
                          child: _locationColumn(
                            title: 'Điểm đến',
                            station: endLocation,
                            icon: Icons.location_on_outlined,
                            alignRight: true,
                            onTap: () async {
                              final result = await Navigator.push(
                                context,
                                MaterialPageRoute(builder: (context) => LocationSelectPage()),
                              );
                              if (result != null && result is AddressInfo) {
                                setState(() {
                                  endLocation = result;
                                });
                              }
                            },
                          ),
                        ),
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Icon(Icons.find_replace_outlined, color: Colors.grey[700], size: 20),
                        SizedBox(width: 8),
                        Text(
                          'Khứ hồi',
                          style: TextStyle(fontSize: 15, color: Colors.grey[700]),
                        ),
                        Spacer(),
                        Transform.scale(
                          scale: 0.8,
                          child: Switch(
                            value: widget.isRoundTrip,
                            onChanged: widget.onRoundTripChanged,
                            activeColor: Colors.blueAccent,
                            inactiveThumbColor: Colors.grey.shade400,
                            inactiveTrackColor: Colors.grey.shade300,
                            activeTrackColor: Colors.blueAccent.withOpacity(0.5),
                          ),
                        ),
                      ],
                    ),
                    Container(
                      padding: EdgeInsets.symmetric(horizontal: 10, vertical: 10),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        border: Border.all(color: Colors.grey.shade300, width: 1),
                        borderRadius: BorderRadius.circular(12),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.grey.withOpacity(0.08),
                            spreadRadius: 1,
                            blurRadius: 3,
                            offset: Offset(0, 1),
                          ),
                        ],
                      ),
                      child: IntrinsicHeight(
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Expanded(
                              child: InkWell(
                                onTap: _pickStartDate,
                                borderRadius: BorderRadius.circular(8),
                                child: _infoColumn(
                                  'Ngày đi',
                                  selectedStartDate != null
                                      ? DateFormat('dd/MM', 'vi_VN').format(selectedStartDate!)
                                      : 'Chọn ngày',
                                  selectedStartDate != null
                                      ? DateFormat('EEEE', 'vi_VN').format(selectedStartDate!)
                                      : null,
                                  iconData: Icons.calendar_today_outlined,
                                ),
                              ),
                            ),
                            if (widget.isRoundTrip)
                              VerticalDivider(
                                width: 10,
                                thickness: 1,
                                color: Colors.grey.shade300,
                                indent: 8,
                                endIndent: 8,
                              ),
                            if (widget.isRoundTrip)
                              Expanded(
                                child: InkWell(
                                  onTap: _pickEndDate,
                                  borderRadius: BorderRadius.circular(8),
                                  child: _infoColumn(
                                    'Ngày về',
                                    selectedEndDate != null
                                        ? DateFormat('dd/MM', 'vi_VN').format(selectedEndDate!)
                                        : 'Chọn ngày',
                                    selectedEndDate != null
                                        ? DateFormat('EEEE', 'vi_VN').format(selectedEndDate!)
                                        : null,
                                    iconData: Icons.calendar_month_outlined,
                                    isEnabled: widget.isRoundTrip,
                                  ),
                                ),
                              ),
                            VerticalDivider(
                              width: 10,
                              thickness: 1,
                              color: Colors.grey.shade300,
                              indent: 8,
                              endIndent: 8,
                            ),
                            Expanded(
                              child: InkWell(
                                onTap: _showTicketQuantityBottomSheet,
                                borderRadius: BorderRadius.circular(8),
                                child: _infoColumn(
                                  'Số vé',
                                  selectedTicketQuantity != null ? '$selectedTicketQuantity vé' : 'Chọn vé',
                                  null,
                                  iconData: Icons.people_alt_outlined,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    SizedBox(height: 20),
                  ],
                ),
              ),
            ),
          ),
        ),
        Positioned(
          top: kToolbarHeight + MediaQuery.of(context).padding.top + 230,
          left: 10,
          right: 10,
          child: Center(
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [Colors.blue, Colors.deepPurpleAccent],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(24),
                boxShadow: [
                  BoxShadow(
                    color: Colors.blue.withOpacity(0.4),
                    blurRadius: 8,
                    offset: Offset(0, 4),
                  ),
                ],
              ),
              child: ElevatedButton(
                onPressed: _findRoute,
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.transparent,
                  shadowColor: Colors.transparent,
                  elevation: 0,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
                  padding: EdgeInsets.fromLTRB(30, 10, 30, 10),
                ),
                child: Text(
                  'Tìm Tuyến Xe',
                  style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildBanner(
    List<String> imgList,
    CarouselController controller,
    Function(int, CarouselPageChangedReason) onPageChanged,
    int currentIndex,
  ) {
    return Column(
      children: [
        CarouselSlider.builder(
          itemCount: imgList.length,
          itemBuilder: (context, index, realIndex) {
            return Container(
              margin: EdgeInsets.symmetric(horizontal: 5.0),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(10),
                child: CachedNetworkImage(
                  imageUrl: imgList[index],
                  fit: BoxFit.cover,
                  width: double.infinity,
                  placeholder: (context, url) => Center(child: CircularProgressIndicator(strokeWidth: 2.0)),
                  errorWidget: (context, url, error) => Icon(Icons.error_outline, color: Colors.red),
                ),
              ),
            );
          },
          options: CarouselOptions(
            height: 150,
            autoPlay: true,
            enlargeCenterPage: true,
            viewportFraction: 0.85,
            autoPlayInterval: Duration(seconds: 4),
            autoPlayAnimationDuration: Duration(milliseconds: 800),
            autoPlayCurve: Curves.fastOutSlowIn,
            onPageChanged: onPageChanged,
          ),
        ),
        SizedBox(height: 8),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(
            imgList.length,
            (idx) => Container(
              width: currentIndex == idx ? 12.0 : 8.0,
              height: 8.0,
              margin: EdgeInsets.symmetric(vertical: 0.0, horizontal: 3.0),
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: (Theme.of(context).brightness == Brightness.dark ? Colors.white : Colors.blueAccent)
                    .withOpacity(currentIndex == idx ? 0.9 : 0.4),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _locationColumn({
    required String title,
    required AddressInfo? station,
    required VoidCallback onTap,
    bool alignRight = false,
    IconData? icon,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        height: 63,
        padding: EdgeInsets.all(8),
        decoration: BoxDecoration(
          color: Colors.grey.shade50,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: Colors.grey.shade300, width: 0.8),
        ),
        child: Column(
          crossAxisAlignment: alignRight ? CrossAxisAlignment.end : CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Row(
              mainAxisAlignment: alignRight ? MainAxisAlignment.end : MainAxisAlignment.start,
              children: [
                if (!alignRight && icon != null) Icon(icon, color: Colors.blueAccent, size: 18),
                if (!alignRight && icon != null) SizedBox(width: 5),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey.shade600,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                if (alignRight && icon != null) SizedBox(width: 5),
                if (alignRight && icon != null) Icon(icon, color: Colors.blueAccent, size: 18),
              ],
            ),
            SizedBox(height: 5),
            Expanded(
              child: station != null
                  ? Column(
                      crossAxisAlignment: alignRight ? CrossAxisAlignment.end : CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Text(
                          station.name,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                            color: Colors.black87,
                          ),
                        ),
                      ],
                    )
                  : Text(
                      'Chọn điểm',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                        color: Colors.grey.shade500,
                      ),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _infoColumn(
    String label,
    String value,
    String? sub, {
    IconData? iconData,
    bool isEnabled = true,
  }) {
    Color valueColor = isEnabled
        ? (value == "Chọn ngày" || value == "Chọn vé" ? Colors.grey.shade600 : Colors.black87)
        : Colors.grey.shade400;
    Color labelColor = isEnabled ? Colors.blueAccent : Colors.grey.shade500;
    Color iconColor = isEnabled ? Colors.blueAccent : Colors.grey.shade400;
    FontWeight valueFontWeight = (value == "Chọn ngày" || value == "Chọn vé") ? FontWeight.normal : FontWeight.w600;

    return Container(
      constraints: BoxConstraints(minHeight: 70),
      padding: const EdgeInsets.symmetric(vertical: 6.0, horizontal: 4.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (iconData != null) Icon(iconData, size: 15, color: iconColor),
              if (iconData != null) SizedBox(width: 4),
              Text(
                label,
                style: TextStyle(
                  color: labelColor,
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
          SizedBox(height: 4),
          Text(
            value,
            style: TextStyle(
              fontSize: 14,
              color: valueColor,
              fontWeight: valueFontWeight,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          if (sub != null && sub.isNotEmpty) SizedBox(height: 2),
          if (sub != null && sub.isNotEmpty)
            Text(
              sub,
              style: TextStyle(
                color: valueColor.withOpacity(0.9),
                fontSize: 11,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
        ],
      ),
    );
  }
}