import 'package:cached_network_image/cached_network_image.dart';
import 'package:caoky/global/toast.dart';
import 'package:caoky/models/account_model.dart';
import 'package:caoky/services/account_service.dart';
import 'package:caoky/ui/account/create_account_page.dart';
import 'package:caoky/ui/account/deposit_page.dart';
import 'package:caoky/ui/account/transaction_history_page.dart';
import 'package:caoky/ui/chat/chat_list_page.dart';
import 'package:caoky/ui/ticket/bus_trip_page.dart';
import 'package:caoky/ui/ticket/ticket_booking_page.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:caoky/models/announcement_model.dart';
import 'package:caoky/services/news_service.dart';
import 'package:shared_preferences/shared_preferences.dart';

class HomePage extends StatefulWidget {
  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  Account? _account;
  bool _isLoading = true;
  String? _errorMessage;
  bool _isAccountNotFound = false;
  final AccountService _accountService = AccountService();
  final NewsService _newsService = NewsService();
  late Future<List<Announcement>> _announcementsFuture;

  double _scrollOffset = 0.0;
  final GlobalKey _cardKey = GlobalKey();
  double _cardHeight = 0.0;

  bool _isObscured = true;

  int _currentIndex = 0;

  String fullName = "";

  final List<Category> categories = [
    Category(
      name: "Mua vé xe",
      image:
          "https://cdn4.iconfinder.com/data/icons/transportation-190/1000/double_double_decker_bus_double_decker_london_double_deck_bus_double_decker_bus_london_double_decker_bus-256.png",
    ),
    Category(
      name: "Xe Buýt",
      image:
          "https://cdn0.iconfinder.com/data/icons/back-to-school-284/512/School_buss_Side_view-256.png",
    ),
    Category(
      name: "Dịch vụ SHB",
      image:
          "https://cdn2.iconfinder.com/data/icons/finance-253/24/banking-money-bank-finance-512.png",
    ),

    Category(
      name: "Gửi hàng hóa",
      image:
          "https://cdn4.iconfinder.com/data/icons/logistics-and-shipping-5/85/delivery_man_box_package_courier-256.png",
    ),
    Category(
      name: "Khuyến mãi",
      image:
          "https://cdn2.iconfinder.com/data/icons/aami-web-internet/64/aami7-94-256.png",
    ),
    Category(
      name: "Ưu đãi sinh viên",
      image:
          "https://cdn2.iconfinder.com/data/icons/education-582/64/Reading-study-student-homework-learning-512.png",
    ),
    Category(
      name: "Dịch vụ SHB",
      image:
          "https://cdn2.iconfinder.com/data/icons/finance-253/24/banking-money-bank-finance-512.png",
    ),

    Category(
      name: "Gửi hàng hóa",
      image:
          "https://cdn4.iconfinder.com/data/icons/logistics-and-shipping-5/85/delivery_man_box_package_courier-256.png",
    ),
    Category(
      name: "Khuyến mãi",
      image:
          "https://cdn2.iconfinder.com/data/icons/aami-web-internet/64/aami7-94-256.png",
    ),
    Category(
      name: "Ưu đãi sinh viên",
      image:
          "https://cdn2.iconfinder.com/data/icons/education-582/64/Reading-study-student-homework-learning-512.png",
    ),
    Category(
      name: "Dịch vụ A",
      image:
          "https://cdn2.iconfinder.com/data/icons/finance-253/24/banking-money-bank-finance-512.png",
    ),

    Category(
      name: "Gửi hàng B",
      image:
          "https://cdn4.iconfinder.com/data/icons/logistics-and-shipping-5/85/delivery_man_box_package_courier-256.png",
    ),
    Category(
      name: "Khuyến C",
      image:
          "https://cdn2.iconfinder.com/data/icons/aami-web-internet/64/aami7-94-256.png",
    ),
    Category(
      name: "Ưu đãi sinh D",
      image:
          "https://cdn2.iconfinder.com/data/icons/education-582/64/Reading-study-student-homework-learning-512.png",
    ),
  ];

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

  final List<Map<String, String>> announcements = [
    {
      'image':
          'https://media-cdn-v2.laodong.vn/storage/newsportal/2019/12/27/775022/Tieu-Su-Jack-2.jpg?w=660',
      'title': 'THÔNG BÁO ĐIỀU CHỈNH LỘ TRÌNH TUYẾN...',
      'description':
          'Nhằm nâng cao trải nghiệm di chuyển và tối ưu khả năng kết nối với tuyến Metro số 1, Công ty Phương Trang chính thức điều chỉnh lộ trình và d...',
    },
    {
      'image':
          'https://media-cdn-v2.laodong.vn/storage/newsportal/2019/12/27/775022/Tieu-Su-Jack-2.jpg?w=660',
      'title': 'TUNG BỪNG KHAI TRƯỚNG TUYẾN XE B...',
      'description':
          'Công ty Phương Trang trân trọng thông báo khai trương tuyến xe buýt mới kết nối hai tỉnh Thừa Thiên Huế và Quảng Trị vào ngày 29/12/2024.',
    },
    {
      'image':
          'https://media-cdn-v2.laodong.vn/storage/newsportal/2019/12/27/775022/Tieu-Su-Jack-2.jpg?w=660',
      'title': 'KHAI TRƯỚNG 17 TUYẾN XE BUÝT THUẬN...',
      'description':
          'Sáng ngày 20/12/2024 Công ty Phương Trang chính thức khai trương 17 tuyến xe buýt thuận tiện – EV kết nối tuyến Metro số 1, đánh dấu bước tiế...',
    },
  ];

  final List<Map<String, String>> popularRoutes = [
    {
      'title': 'Sài Gòn - Đà Lạt',
      'description':
          'Ra mắt dịch vụ xe VIP mới 34 giường cho bạn thoải mái di chuyển.',
      'image':
          'https://cdn.futabus.vn/futa-busline-cms-dev/Rectangle_23_2_8bf6ed1d78/Rectangle_23_2_8bf6ed1d78.png', // Placeholder for the single image
    },
    {
      'title': 'Hà Nội - Sapa',
      'description':
          'Xe giường nằm cao cấp, dịch vụ 5 sao, giá chỉ từ 300.000đ.',
      'image':
          'https://cdn.futabus.vn/futa-busline-cms-dev/Rectangle_23_3_2d8ce855bc/Rectangle_23_3_2d8ce855bc.png',
    },
    {
      'title': 'Đà Nẵng - Hội An',
      'description': 'Tuyến xe mới, tiện nghi, khởi hành hàng ngày.',
      'image':
          'https://cdn.futabus.vn/futa-busline-cms-dev/Rectangle_23_4_061f4249f6/Rectangle_23_4_061f4249f6.png',
    },
    {
      'title': 'Cần Thơ - Vũng Tàu',
      'description': 'Xe VIP 34 giường, dịch vụ cao cấp, giá ưu đãi.',
      'image':
          'https://cdn.futabus.vn/futa-busline-cms-dev/Rectangle_23_4_061f4249f6/Rectangle_23_4_061f4249f6.png',
    },
  ];

  Future<void> _loadUserData() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    setState(() {
      fullName = prefs.getString('fullName') ?? "";
    });
  }

  @override
  void initState() {
    super.initState();
    _loadAccountData();
    _loadInitialData();
    _loadUserData();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final RenderBox? renderBox =
          _cardKey.currentContext?.findRenderObject() as RenderBox?;
      if (renderBox != null) {
        setState(() {
          _cardHeight = renderBox.size.height;
        });
      }
    });
  }

  Future<void> _loadInitialData() async {
    setState(() {
      _announcementsFuture = _newsService.getAnnouncements(
        limit: 3,
      ); // Lấy 3 tin mới nhất
    });
  }

  Future<void> _loadAccountData() async {
    setState(() {
      _isLoading = true;
      _errorMessage = null;
      _isAccountNotFound = false;
    });

    try {
      final accountData = await _accountService.getMyAccount();
      setState(() {
        _account = accountData;
      });
    } on AccountNotFoundException {
      setState(() {
        _isAccountNotFound = true;
      });
    } catch (e) {
      setState(() {
        _errorMessage = e.toString().replaceFirst(
          "Exception: ",
          "",
        ); // Bỏ tiền tố "Exception: "
      });
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  // Hàm format tiền tệ
  String _formatCurrency(int amount) {
    final format = NumberFormat.currency(locale: 'vi_VN', symbol: 'đ');
    return format.format(amount);
  }

  @override
  Widget build(BuildContext context) {
    const double appBarHeight = 150;

    return Scaffold(
      body: Container(
        color: Colors.white,
        child: RefreshIndicator(
          // Thêm RefreshIndicator để dễ dàng tải lại
          onRefresh: () async {
            // Tải lại cả hai loại dữ liệu khi kéo để làm mới
            await _loadAccountData();
            await _loadInitialData();
          },
          child: Stack(
            children: [
              NotificationListener<ScrollNotification>(
                onNotification: (scrollNotification) {
                  if (scrollNotification is ScrollUpdateNotification &&
                      scrollNotification.depth == 0) {
                    setState(() {
                      _scrollOffset = scrollNotification.metrics.pixels.clamp(
                        0.0,
                        100.0,
                      );
                    });
                  }
                  return true;
                },
                child: CustomScrollView(
                  slivers: [
                    SliverAppBar(
                      expandedHeight: appBarHeight,
                      floating: false,
                      pinned: true,
                      toolbarHeight: 80,
                      shape:
                          _scrollOffset > 50
                              ? null
                              : RoundedRectangleBorder(
                                borderRadius: BorderRadius.only(
                                  bottomLeft: Radius.circular(15),
                                  bottomRight: Radius.circular(15),
                                ),
                              ),
                      flexibleSpace: ClipRRect(
                        borderRadius:
                            _scrollOffset > 50
                                ? BorderRadius.zero
                                : const BorderRadius.only(
                                  bottomLeft: Radius.circular(15),
                                  bottomRight: Radius.circular(15),
                                ),
                        child: Container(
                          color: Colors.blueAccent,
                          // decoration: BoxDecoration(
                          //   image: DecorationImage(
                          //     image: AssetImage('assets/images/background1.jpg'),
                          //     fit: BoxFit.cover,
                          //   ),
                          // ),
                        ),
                      ),
                      leading: Padding(
                        padding: const EdgeInsets.only(left: 20.0, top: 8.0),
                        child: CircleAvatar(
                          radius: 15,
                          backgroundColor: Colors.white,
                          child: Icon(
                            Icons.person,
                            size: 25,
                            color: Colors.blueAccent,
                          ),
                        ),
                      ),
                      title:
                          _scrollOffset > 50
                              ? Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceAround,
                                children: [
                                  IconButton(
                                    icon: Icon(
                                      Icons.add_circle_outline,
                                      color: Colors.white,
                                      size: 30,
                                    ),
                                    onPressed: () async {
                                      final bool? success =
                                          await Navigator.push<bool>(
                                            context,
                                            MaterialPageRoute(
                                              builder:
                                                  (context) =>
                                                      const DepositPage(),
                                            ),
                                          );
                                      if (success == true) _loadAccountData();
                                    },
                                  ),
                                  IconButton(
                                    icon: Icon(
                                      Icons.refresh,
                                      color: Colors.white,
                                      size: 30,
                                    ),
                                    onPressed: () {},
                                  ),
                                  IconButton(
                                    icon: Icon(
                                      Icons.account_balance_wallet,
                                      color: Colors.white,
                                      size: 30,
                                    ),
                                    onPressed: () {
                                      Navigator.push(
                                        context,
                                        MaterialPageRoute(
                                          builder:
                                              (_) =>
                                                  const TransactionHistoryPage(),
                                        ),
                                      );
                                    },
                                  ),
                                ],
                              )
                              : Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Text(
                                    "Xin chào,",
                                    style: TextStyle(
                                      fontSize: 14,
                                      color: Colors.white,
                                    ),
                                  ),
                                  Text(
                                    fullName,
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: const Color.fromRGBO(
                                        255,
                                        255,
                                        255,
                                        1,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                      actions: [
                        IconButton(
                          icon: Icon(
                            Icons.question_answer,
                            color: Colors.white,
                          ),
                          onPressed: () {
                            ToastUtils.show("Tính năng sẽ sớm phát triển");
                            // Navigator.push(
                            //   context,
                            //   MaterialPageRoute(
                            //     builder: (context) => ChatListPage(),
                            //   ),
                            // );
                          },
                        ),
                      ],
                    ),
                    SliverToBoxAdapter(
                      child: Column(
                        children: [
                          SizedBox(height: _cardHeight * 0.7 + 150),
                          _buildHorizontalList(),
                          _buildBanner(imgBanner),
                          _buildDynamicAnnouncementSection(
                            "DACNTT City Bus",
                            _announcementsFuture,
                          ),
                          _buildBanner(imgBanner1),
                          _buildDynamicAnnouncementSection(
                            "DACNTT Express",
                            _announcementsFuture,
                          ),
                          _buildPopularRoutesList(),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              Positioned(
                left: 16,
                right: 16,
                // Điều chỉnh lại vị trí top một cách an toàn
                top: appBarHeight - 70,
                child: AnimatedOpacity(
                  opacity: (1.0 - (_scrollOffset / 100.0)).clamp(0.0, 1.0),
                  duration: const Duration(milliseconds: 200),
                  child: Card(
                    elevation: 4,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                      side: BorderSide(color: Colors.grey.shade200, width: 1),
                    ),

                    color: Colors.white,
                    shadowColor: Colors.grey.withOpacity(0.5),

                    child: Container(
                      constraints: const BoxConstraints(
                        minHeight: 160,
                      ), // Đảm bảo thẻ có chiều cao tối thiểu
                      padding: const EdgeInsets.all(16.0),

                      // Gọi hàm để build nội dung động cho thẻ
                      child: _buildCardContent(),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  /// PHẦN LOGIC QUAN TRỌNG: QUYẾT ĐỊNH HIỂN THỊ GÌ TRONG THẺ
  Widget _buildCardContent() {
    if (_isLoading) {
      return const Center(
        child: CircularProgressIndicator(color: Colors.blueAccent),
      );
    }
    if (_isAccountNotFound) {
      return _buildCreateAccountButton();
    }
    if (_errorMessage != null) {
      return _buildErrorWidget(_errorMessage!);
    }
    if (_account != null) {
      return _buildAccountInfo(_account!);
    }
    // Trường hợp dự phòng
    return const Center(child: Text("Có lỗi xảy ra."));
  }

  Widget _buildAccountInfo(Account account) {
    return Padding(
      padding: const EdgeInsets.all(8),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  const Text(
                    "Số dư hiện tại",
                    style: TextStyle(fontSize: 13, color: Colors.grey),
                  ),
                  const SizedBox(height: 4),
                ],
              ),
              const Text(
                "💳 Ví của bạn",
                style: TextStyle(
                  fontSize: 17,
                  fontWeight: FontWeight.bold,
                  color: Colors.black87,
                ),
              ),
            ],
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              Row(
                children: [
                  Text(
                    _isObscured ? "••••••" : _formatCurrency(account.balance),
                    style: const TextStyle(
                      fontSize: 17,
                      fontWeight: FontWeight.bold,
                      color: Colors.orangeAccent,
                    ),
                  ),
                  IconButton(
                    icon: Icon(
                      _isObscured
                          ? Icons.visibility_off_outlined
                          : Icons.visibility_outlined,
                      size: 17,
                      color: Colors.grey,
                    ),
                    onPressed: () => setState(() => _isObscured = !_isObscured),
                  ),
                ],
              ),
            ],
          ),

          const Divider(height: 8),
          SizedBox(height: 8),

          // Các hành động
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              GestureDetector(
                onTap: () async {
                  final bool? success = await Navigator.push<bool>(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const DepositPage(),
                    ),
                  );
                  if (success == true) _loadAccountData();
                },
                child: _buildActionItem(Icons.add_circle_outline, "Nạp tiền"),
              ),
              GestureDetector(
                onTap: () {
                  ToastUtils.show("Tính năng sẽ sớm phát triển");
                },
                child: _buildActionItem(Icons.refresh_outlined, "Rút tiền"),
              ),
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => const TransactionHistoryPage(),
                    ),
                  );
                },
                child: _buildActionItem(
                  Icons.account_balance_wallet_outlined,
                  "Chi tiết",
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  /// Mỗi hành động bên dưới ví
  Widget _buildActionItem(IconData icon, String label) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: Colors.orange.shade50,
            shape: BoxShape.circle,
          ),
          child: Icon(icon, color: Colors.orange, size: 17),
        ),
        const SizedBox(height: 6),
        Text(
          label,
          style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w500),
        ),
      ],
    );
  }

  /// Widget khi không tìm thấy tài khoản (lỗi 404)
  Widget _buildCreateAccountButton() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text("Bạn có muốn mở ví ?", style: TextStyle(fontSize: 16)),
          const SizedBox(height: 16),
          ElevatedButton(
            style: ElevatedButton.styleFrom(backgroundColor: Colors.orange),
            // *** SỬA ĐỔI Ở ĐÂY ***
            onPressed: () async {
              // Điều hướng đến trang tạo tài khoản
              final bool? success = await Navigator.push<bool>(
                context,
                MaterialPageRoute(builder: (context) => CreateAccountPage()),
              );

              // Nếu tạo tài khoản thành công (trang mới trả về true),
              // thì tải lại dữ liệu tài khoản
              if (success == true) {
                _loadAccountData();
              }
            },
            child: const Text(
              "Tạo tài khoản ngay",
              style: TextStyle(color: Colors.white),
            ),
          ),
        ],
      ),
    );
  }

  /// Widget khi có lỗi khác
  Widget _buildErrorWidget(String message) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              message,
              style: const TextStyle(color: Colors.red),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _loadAccountData, // Cho phép người dùng thử lại
              child: const Text("Thử lại"),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHorizontalList() {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 10),
      child: SizedBox(
        height: 200,
        child: GridView.builder(
          scrollDirection: Axis.horizontal,
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 2,
            crossAxisSpacing: 10,
            mainAxisSpacing: 10,
            childAspectRatio: 0.8,
          ),
          itemCount: categories.length,
          itemBuilder: (context, index) {
            return GestureDetector(
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder:
                        (context) => TicketBookingPage(
                          // selectedCategory: categories[index].name,
                        ),
                  ),
                );
              },
              child: Column(
                children: [
                  Container(
                    height: 45,
                    width: 45,
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(8),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black26,
                          blurRadius: 6,
                          offset: Offset(3, 3),
                        ),
                      ],
                    ),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(5),
                      child: Padding(
                        padding: const EdgeInsets.all(5.0),
                        child: Image.network(
                          categories[index].image ?? 'error.jpg',
                          fit: BoxFit.cover,
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return Center(
                              child: CircularProgressIndicator(
                                color: Colors.blueAccent,
                              ),
                            );
                          },
                          errorBuilder: (context, error, stackTrace) {
                            return Icon(Icons.directions_bus);
                          },
                        ),
                      ),
                    ),
                  ),
                  SizedBox(height: 8),
                  Text(
                    categories[index].name,
                    style: TextStyle(fontSize: 14, color: Colors.black54),
                    textAlign: TextAlign.center,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildBanner(List<String> imgBanner) {
    double screenWidth = MediaQuery.of(context).size.width;
    double bannerHeight =
        screenWidth > 1024
            ? 350
            : screenWidth > 600
            ? 250
            : 180;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Stack(
        children: [
          CarouselSlider.builder(
            itemCount: imgBanner.length,
            itemBuilder: (context, index, realIndex) {
              return CachedNetworkImage(
                imageUrl: imgBanner[index],
                imageBuilder:
                    (context, imageProvider) => Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(10),
                        image: DecorationImage(
                          image: imageProvider,
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                placeholder:
                    (context, url) => Center(
                      child: CircularProgressIndicator(
                        color: Colors.blueAccent,
                      ),
                    ),
                errorWidget: (context, url, error) => Icon(Icons.error),
              );
            },
            options: CarouselOptions(
              height: bannerHeight,
              enlargeCenterPage: true,
              autoPlay: true,
              autoPlayInterval: Duration(seconds: 3),
              autoPlayAnimationDuration: Duration(milliseconds: 800),
              autoPlayCurve: Curves.easeInOut,
              viewportFraction: 0.8,
              aspectRatio: 16 / 9,
              initialPage: 0,
              onPageChanged: (index, reason) {
                setState(() {
                  _currentIndex = index;
                });
              },
            ),
          ),
          Positioned(
            bottom: 10,
            left: 0,
            right: 0,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(
                imgBanner.length,
                (index) => Container(
                  width: _currentIndex == index ? 17 : 7,
                  height: 8,
                  margin: EdgeInsets.symmetric(horizontal: 4),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color:
                        _currentIndex == index
                            ? Colors.blueAccent
                            : const Color.fromARGB(255, 174, 185, 184),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDynamicAnnouncementSection(
    String title,
    Future<List<Announcement>> future,
  ) {
    return FutureBuilder<List<Announcement>>(
      future: future,
      builder: (context, snapshot) {
        // Trường hợp 1: Đang tải dữ liệu
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Padding(
            padding: EdgeInsets.all(40.0),
            child: Center(
              child: CircularProgressIndicator(color: Colors.blueAccent),
            ),
          );
        }
        // Trường hợp 2: Gặp lỗi
        if (snapshot.hasError) {
          return Padding(
            padding: const EdgeInsets.all(16.0),
            child: Center(child: Text('Lỗi tải tin tức: ${snapshot.error}')),
          );
        }
        // Trường hợp 3: Không có dữ liệu
        if (!snapshot.hasData || snapshot.data!.isEmpty) {
          return Padding(
            padding: const EdgeInsets.all(16.0),
            child: Center(child: Text('Chưa có tin tức nào.')),
          );
        }

        // Trường hợp 4: Tải dữ liệu thành công
        final announcements = snapshot.data!;
        return _buildAnnouncementList(title, announcements);
      },
    );
  }

  // SỬA ĐỔI: `_buildAnnouncementList` giờ nhận vào List<Announcement>
  Widget _buildAnnouncementList(
    String title,
    List<Announcement> announcements,
  ) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 10.0),
      child: Column(
        children: [
          Row(
            // ... phần tiêu đề không đổi
          ),
          ...announcements.map((announcement) {
            // Dùng `announcements` được truyền vào
            return Padding(
              padding: const EdgeInsets.symmetric(vertical: 8.0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  ClipRRect(
                    borderRadius: BorderRadius.circular(10),
                    child: CachedNetworkImage(
                      // Sử dụng CachedNetworkImage cho hiệu năng tốt hơn
                      imageUrl: announcement.imageUrl,
                      width: 80,
                      height: 80,
                      fit: BoxFit.cover,
                      placeholder:
                          (context, url) => Container(
                            width: 80,
                            height: 80,
                            color: Colors.grey[200],
                          ),
                      errorWidget:
                          (context, url, error) =>
                              const Icon(Icons.error, size: 80),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          announcement.title, // Dùng dữ liệu từ API
                          style: const TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                            color: Colors.black,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                        const SizedBox(height: 4),
                        Text(
                          announcement.description, // Dùng dữ liệu từ API
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey[600],
                          ),
                          maxLines: 3,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            );
          }).toList(),
        ],
      ),
    );
  }

  Widget _buildPopularRoutesList() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 10.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                "Tuyến xe khách phổ biến",
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
              ),
              TextButton(
                onPressed: () {},
                child: Text(
                  "Xem tất cả",
                  style: TextStyle(fontSize: 14, color: Colors.blueAccent),
                ),
              ),
            ],
          ),
          SizedBox(
            height: 210,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              itemCount: popularRoutes.length,
              itemBuilder: (context, index) {
                final route = popularRoutes[index];
                return GestureDetector(
                  onTap:
                      () => {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => TicketBookingPage(),
                          ),
                        ),
                      },
                  child: Padding(
                    padding: const EdgeInsets.only(right: 10.0),
                    child: Container(
                      width: 200,
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(10),
                        color: Colors.white,
                        border: Border.all(
                          color: Colors.grey[300]!,
                          width: 1.0,
                        ),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black12,
                            blurRadius: 6,
                            offset: Offset(0, 2),
                          ),
                        ],
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          ClipRRect(
                            borderRadius: BorderRadius.only(
                              topLeft: Radius.circular(10),
                              topRight: Radius.circular(10),
                            ),
                            child: Image.network(
                              route['image']!,
                              height: 120,
                              width: double.infinity,
                              fit: BoxFit.cover,
                              errorBuilder: (context, error, stackTrace) {
                                return Container(
                                  height: 90,
                                  color: Colors.grey[300],
                                  child: Icon(Icons.error),
                                );
                              },
                            ),
                          ),
                          Expanded(
                            child: Padding(
                              padding: const EdgeInsets.all(10.0),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    route['title']!,
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.black,
                                    ),
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                  SizedBox(height: 5),
                                  Text(
                                    route['description']!,
                                    style: TextStyle(
                                      fontSize: 12,
                                      color: Colors.grey[600],
                                    ),
                                    maxLines: 2,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class Category {
  final String name;
  final String? image;

  Category({required this.name, this.image});
}
