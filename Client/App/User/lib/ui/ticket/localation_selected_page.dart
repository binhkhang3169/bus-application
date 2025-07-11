import 'dart:convert'; // Keep if AddressInfo or other models might use it indirectly, otherwise can be removed.

import 'package:caoky/global/toast.dart';
// import 'package:caoky/models/LocationGroup.dart'; // Likely not needed anymore
import 'package:caoky/models/trip/address_info.dart';
// import 'package:caoky/models/trip/location.dart'; // Likely not needed anymore
// import 'package:caoky/models/trip/location_group.dart'; // Likely not needed anymore
// import 'package:caoky/models/trip/station.dart'; // Not needed anymore
import 'package:caoky/services/api_trip_service.dart';
// import 'package:caoky/ui/login/keys/shipping.dart'; // Check if still used elsewhere or can be removed
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
// import 'package:http/http.dart' as http show get; // Not used directly here

class LocationSelectPage extends StatefulWidget {
  // The 'province' parameter seems to have been for pre-selecting a province
  // to then show stations. If we are only selecting a province, this might
  // not be needed unless you have a use case for an initial highlighted province.
  // For now, I'll remove it to simplify, assuming the goal is purely to select one province from the list.
  // final int? province;
  // LocationSelectPage({this.province});

  LocationSelectPage(); // Constructor without the province parameter

  @override
  _LocationSelectPageState createState() => _LocationSelectPageState();
}

class _LocationSelectPageState extends State<LocationSelectPage> {
  // late Future<List<Station>> futureLocations; // Not needed for stations anymore
  // List<Station> allData = []; // Not needed for stations anymore
  // int? selectedProvince; // Not needed if we pop immediately
  String searchText = '';
  List<AddressInfo> provinces = [];
  bool _isLoadingProvinces = true; // To manage loading state for provinces
  late ApiTripService apiService;

  @override
  void initState() {
    super.initState();
    apiService = ApiTripService(Dio());
    fetchProvinces();
  }

  Future<void> fetchProvinces() async {
    setState(() {
      _isLoadingProvinces = true;
    });
    try {
      final response = await apiService.getListAddress();
      if (response.code == 200) {
        setState(() {
          provinces = response.data ?? []; // Use ?? [] for safety
        });
      } else {
        ToastUtils.show("Lỗi khi tải danh sách tỉnh/thành: ${response.message}");
      }
    } catch (e) {
      print("Lỗi khi gọi API provinces: $e");
      ToastUtils.show("Lỗi khi tải danh sách tỉnh/thành. Vui lòng thử lại.");
    } finally {
      setState(() {
        _isLoadingProvinces = false;
      });
    }
  }

  // _filterData() and fetchLocations() are no longer needed as we are not dealing with stations.

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        // Title is now static as we are always choosing a province on this screen
        title: Text(
          "Chọn tỉnh/thành phố",
          style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
        backgroundColor: Colors.blueAccent,
        iconTheme: IconThemeData(color: Colors.white),
      ),
      body: _isLoadingProvinces
          ? Center(
              child: CircularProgressIndicator(
                color: Colors.blueAccent,
              ),
            )
          : provinces.isEmpty && !_isLoadingProvinces
              ? Center(child: Text("Không có dữ liệu tỉnh/thành phố."))
              : buildProvinceList(), // Directly build province list
    );
  }

  Widget buildProvinceList() {
    final filteredProvinces = provinces.where((province) {
      final name = province.name.toLowerCase();
      return name.contains(searchText.toLowerCase());
    }).toList();

    if (filteredProvinces.isEmpty && searchText.isNotEmpty) {
       return Column(
         children: [
           buildSearchBar('Tìm kiếm tỉnh / thành phố ...'),
           Expanded(
             child: Center(child: Text("Không tìm thấy tỉnh/thành phố phù hợp.")),
           ),
         ],
       );
    }

    return Column(
      children: [
        buildSearchBar('Tìm kiếm tỉnh / thành phố ...'),
        Expanded(
          child: ListView.builder(
            itemCount: filteredProvinces.length,
            itemBuilder: (context, index) {
              final province = filteredProvinces[index];
              return Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  ListTile(
                    title: Text(province.name),
                    trailing: Icon(Icons.arrow_forward_ios, size: 16),
                    onTap: () {
                      // When a province is tapped, pop the page and return the selected province
                      Navigator.pop(context, province);
                    },
                  ),
                  // Divider for better visual separation, optional
                  if (index < filteredProvinces.length - 1)
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16.0),
                      child: Divider(height: 1, color: Colors.grey[300]),
                    ),
                ],
              );
            },
          ),
        ),
      ],
    );
  }

  // buildLocationList is no longer needed.

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