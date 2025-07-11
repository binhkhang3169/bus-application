import 'dart:convert';
import 'package:caoky/ui/login/keys/shipping.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class ManagerAddressPage extends StatefulWidget {
  @override
  _ManagerAddressPageState createState() => _ManagerAddressPageState();
}

class _ManagerAddressPageState extends State<ManagerAddressPage> {
  List<dynamic> provinces = [];
  List<dynamic> districts = [];
  List<dynamic> wards = [];

  Map<String, dynamic>? selectedProvinceData;
  Map<String, dynamic>? selectedDistrictData;
  Map<String, dynamic>? selectedWardData;

  String searchText = '';
  final TextEditingController specificAddressController =
      TextEditingController();

  @override
  void initState() {
    super.initState();
    fetchProvinces();
  }

  Future<void> fetchProvinces() async {
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/province",
    );
    try {
      final response = await http.get(
        url,
        headers: {"Token": Shipping().apiKey},
      );
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          provinces = data['data'] ?? [];
        });
      } else {
        showError("Lỗi khi tải danh sách tỉnh/thành");
      }
    } catch (e) {
      showError("Không thể kết nối đến server");
    }
  }

  Future<void> fetchDistricts(int provinceId) async {
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/district",
    );
    try {
      final response = await http.post(
        url,
        headers: {
          "Content-Type": "application/json",
          "Token": Shipping().apiKey,
        },
        body: jsonEncode({"province_id": provinceId}),
      );
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          districts = data['data'] ?? [];
        });
      } else {
        showError("Lỗi khi tải danh sách quận/huyện");
      }
    } catch (e) {
      showError("Không thể kết nối đến server");
    }
  }

  Future<void> fetchWards(int districtId) async {
    final url = Uri.parse(
      "https://dev-online-gateway.ghn.vn/shiip/public-api/master-data/ward",
    );
    try {
      final response = await http.post(
        url,
        headers: {
          "Content-Type": "application/json",
          "Token": Shipping().apiKey,
        },
        body: jsonEncode({"district_id": districtId}),
      );
      if (response.statusCode == 200) {
        final data = jsonDecode(utf8.decode(response.bodyBytes));
        setState(() {
          wards = data['data'] ?? [];
        });
      } else {
        showError("Lỗi khi tải danh sách phường/xã");
      }
    } catch (e) {
      showError("Không thể kết nối đến server");
    }
  }

  void showError(String message) {
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }

  void completeAddressSelection() {
    final address =
        "${specificAddressController.text}, ${selectedWardData?['WardName']}, ${selectedDistrictData?['DistrictName']}, ${selectedProvinceData?['ProvinceName']}";
    Navigator.pop(context, address);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        centerTitle: true,

        title: Text(
          selectedProvinceData == null
              ? "Chọn tỉnh"
              : selectedDistrictData == null
              ? "Chọn huyện"
              : selectedWardData == null
              ? "Chọn xã"
              : "Hoàn tất",
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Colors.blueAccent,
        iconTheme: IconThemeData(color: Colors.white),
        leading: IconButton(
          icon: Icon(Icons.arrow_back),
          onPressed: () {
            setState(() {
              if (selectedWardData != null) {
                selectedWardData = null;
              } else if (selectedDistrictData != null) {
                selectedDistrictData = null;
                wards = [];
              } else if (selectedProvinceData != null) {
                selectedProvinceData = null;
                districts = [];
              } else {
                Navigator.pop(context);
              }
            });
          },
        ),
      ),
      body: buildAddressSelector(),
    );
  }

  Widget buildAddressSelector() {
    if (selectedProvinceData == null) {
      final filtered =
          provinces.where((province) {
            final name = province['ProvinceName']?.toLowerCase() ?? '';
            return name.contains(searchText.toLowerCase());
          }).toList();

      return buildSelectionList(
        hint: 'Tìm kiếm tỉnh / thành phố ...',
        items: filtered,
        titleKey: 'ProvinceName',
        onTap: (province) async {
          setState(() {
            selectedProvinceData = province;
            searchText = '';
          });
          await fetchDistricts(province['ProvinceID']);
        },
      );
    } else if (selectedDistrictData == null) {
      final filtered =
          districts.where((district) {
            final name = district['DistrictName']?.toLowerCase() ?? '';
            return name.contains(searchText.toLowerCase());
          }).toList();

      return buildSelectionList(
        hint: 'Tìm kiếm quận / huyện ...',
        items: filtered,
        titleKey: 'DistrictName',
        onTap: (district) async {
          setState(() {
            selectedDistrictData = district;
            searchText = '';
          });
          await fetchWards(district['DistrictID']);
        },
      );
    } else {
      final filtered =
          wards.where((ward) {
            final name = ward['WardName']?.toLowerCase() ?? '';
            return name.contains(searchText.toLowerCase());
          }).toList();

      return Column(
        children: [
          buildSearchBar('Tìm kiếm phường / xã ...'),
          Expanded(
            child: ListView.builder(
              itemCount: filtered.length,
              itemBuilder: (context, index) {
                final ward = filtered[index];
                return ListTile(
                  title: Text(ward['WardName'] ?? ''),
                  trailing: Icon(Icons.arrow_forward_ios, size: 16),
                  onTap: () {
                    setState(() {
                      selectedWardData = ward;
                    });
                  },
                );
              },
            ),
          ),
          if (selectedWardData != null) ...[
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 0),
              child: TextField(
                controller: specificAddressController,
                decoration: InputDecoration(
                  labelText: "Địa chỉ",
                  hintText: 'Nhập địa chỉ cụ thể (số nhà, tên đường...)',
                  // Viền các trạng thái
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(color: Colors.grey),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(color: Colors.blue, width: 2),
                  ),
                  errorBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(color: Colors.red),
                  ),
                  focusedErrorBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(color: Colors.red, width: 2),
                  ),

                  // Màu chữ và label khi focus
                  labelStyle: TextStyle(color: Colors.grey), // mặc định
                  floatingLabelStyle: TextStyle(
                    color: Colors.blue,
                  ), // khi focus
                  hintStyle: TextStyle(
                    color: Colors.grey.shade400,
                  ), // màu gợi ý
                  suffixStyle: TextStyle(color: Colors.blue), // màu "VND"
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: ElevatedButton(
                onPressed: () {
                  if (specificAddressController.text.trim().isEmpty) {
                    showError('Vui lòng nhập địa chỉ cụ thể');
                    return;
                  }
                  completeAddressSelection();
                },
                child: Text('Hoàn tất'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.blueAccent,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  minimumSize: Size(double.infinity, 48),
                ),
              ),
            ),
          ],
        ],
      );
    }
  }

  Widget buildSelectionList({
    required String hint,
    required List<dynamic> items,
    required String titleKey,
    required Function(Map<String, dynamic>) onTap,
  }) {
    return Column(
      children: [
        buildSearchBar(hint),
        Expanded(
          child: ListView.builder(
            itemCount: items.length,
            itemBuilder: (context, index) {
              final item = items[index];
              return ListTile(
                title: Text(item[titleKey] ?? ''),
                trailing: Icon(Icons.arrow_forward_ios, size: 16),
                onTap: () => onTap(item),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget buildSearchBar(String hint) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(30),
          boxShadow: [
            BoxShadow(
              color: Colors.grey.withOpacity(0.15),
              blurRadius: 8,
              offset: Offset(0, 4),
            ),
          ],
        ),
        child: TextField(
          decoration: InputDecoration(
            hintText: hint,
            prefixIcon: Icon(Icons.search, color: Colors.blueAccent),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide.none,
            ),
            filled: true,
            fillColor: Colors.white,
            contentPadding: EdgeInsets.symmetric(vertical: 0, horizontal: 16),
          ),
          onChanged: (value) {
            setState(() => searchText = value);
          },
        ),
      ),
    );
  }
}
