// lib/ui/account/transaction_history_page.dart
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:caoky/models/transaction_history_model.dart';
import 'package:caoky/services/account_service.dart';

class TransactionHistoryPage extends StatefulWidget {
  const TransactionHistoryPage({super.key});

  @override
  _TransactionHistoryPageState createState() => _TransactionHistoryPageState();
}

class _TransactionHistoryPageState extends State<TransactionHistoryPage> {
  final AccountService _accountService = AccountService();
  final ScrollController _scrollController = ScrollController();

  final List<TransactionHistory> _transactions = [];
  int _currentPage = 1;
  final int _pageSize = 20; // Số lượng item mỗi lần tải
  bool _isLoading = false;
  bool _isInitialLoad = true; // Cờ cho lần tải đầu tiên
  bool _hasMore = true; // Còn dữ liệu để tải thêm không
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _fetchTransactions(); // Tải dữ liệu lần đầu
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _fetchTransactions({bool isRefresh = false}) async {
    if (_isLoading) return;

    setState(() {
      _isLoading = true;
      if (isRefresh) {
        _isInitialLoad = true;
        _transactions.clear();
        _currentPage = 1;
        _hasMore = true;
        _errorMessage = null;
      }
    });

    try {
      final newTransactions = await _accountService.getTransactionHistory(
        pageId: _currentPage,
        pageSize: _pageSize,
      );

      setState(() {
        _transactions.addAll(newTransactions);
        _currentPage++;
        // Nếu kết quả trả về ít hơn page size, nghĩa là đã hết dữ liệu
        if (newTransactions.length < _pageSize) {
          _hasMore = false;
        }
      });
    } catch (e) {
      setState(() {
        _errorMessage = e.toString();
      });
    } finally {
      setState(() {
        _isLoading = false;
        _isInitialLoad = false;
      });
    }
  }

  void _onScroll() {
    // Tải thêm khi người dùng cuộn gần đến cuối danh sách
    if (_scrollController.position.pixels >=
            _scrollController.position.maxScrollExtent * 0.9 &&
        _hasMore &&
        !_isLoading) {
      _fetchTransactions();
    }
  }

  /// *** NÂNG CẤP TẠI ĐÂY ***
  /// Helper để lấy thông tin hiển thị cho từng loại giao dịch dựa trên các hằng số từ backend.
  Map<String, dynamic> _getTransactionUIData(String type) {
    switch (type) {
      case 'CREATE_ACCOUNT':
        return {
          'icon': Icons.person_add_alt_1,
          'color': Colors.blue,
          'label': 'Tạo tài khoản thành công',
          'isAmountBased': false, // Không hiển thị số tiền
        };
      case 'DEPOSIT':
        return {
          'icon': Icons.arrow_upward,
          'color': Colors.green,
          'label': 'Nạp tiền',
          'prefix': '+ ',
          'isAmountBased': true,
        };
      case 'PAYMENT':
        return {
          'icon': Icons.arrow_downward,
          'color': Colors.red,
          'label': 'Thanh toán',
          'prefix': '- ',
          'isAmountBased': true,
        };
      case 'CLOSE_ACCOUNT':
        return {
          'icon': Icons.person_remove,
          'color': Colors.grey[700],
          'label': 'Đóng tài khoản',
          'isAmountBased': false,
        };
      default: // Trường hợp dự phòng cho các loại giao dịch không xác định
        return {
          'icon': Icons.history,
          'color': Colors.grey,
          'label': 'Giao dịch khác',
          'prefix': '',
          'isAmountBased': true,
        };
    }
  }

  String _formatCurrency(int amount) {
    final format = NumberFormat.currency(locale: 'vi_VN', symbol: 'đ');
    return format.format(amount);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        title: const Text("Lịch sử giao dịch"),
        backgroundColor: Colors.blueAccent,
        foregroundColor: Colors.white,
        centerTitle: true,
      ),
      body: RefreshIndicator(
        onRefresh: () => _fetchTransactions(isRefresh: true),
        child: _buildBody(),
      ),
    );
  }

  Widget _buildBody() {
    if (_isInitialLoad && _isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_errorMessage != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text("Lỗi: $_errorMessage", textAlign: TextAlign.center),
              const SizedBox(height: 8),
              ElevatedButton(
                onPressed: () => _fetchTransactions(isRefresh: true),
                child: const Text("Thử lại"),
              ),
            ],
          ),
        ),
      );
    }

    if (_transactions.isEmpty && !_isLoading) {
      return const Center(
        child: Text(
          "Bạn chưa có giao dịch nào.",
          style: TextStyle(fontSize: 16, color: Colors.grey),
        ),
      );
    }

    return ListView.builder(
      controller: _scrollController,
      itemCount: _transactions.length + (_hasMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == _transactions.length) {
          // Widget loading ở cuối danh sách
          return _isLoading
              ? const Padding(
                padding: EdgeInsets.symmetric(vertical: 16.0),
                child: Center(child: CircularProgressIndicator()),
              )
              : const SizedBox.shrink();
        }

        final tx = _transactions[index];
        final uiData = _getTransactionUIData(tx.transactionType);

        // Xác định widget hiển thị ở cuối (trailing)
        Widget trailingWidget;
        if (uiData['isAmountBased'] as bool) {
          // Hiển thị số tiền cho các giao dịch nạp, rút, thanh toán...
          trailingWidget = Text(
            tx.amount != null
                ? '${uiData['prefix']}${_formatCurrency(tx.amount!)}'
                : 'N/A',
            style: TextStyle(
              color: uiData['color'],
              fontWeight: FontWeight.bold,
              fontSize: 15,
            ),
          );
        } else {
          // Hiển thị icon cho các giao dịch không có số tiền (tạo, đóng tài khoản)
          trailingWidget = Icon(uiData['icon'], color: uiData['color']);
        }

        return Card(
          color: Colors.white,
          margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
          elevation: 2,
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: (uiData['color'] as Color).withOpacity(0.1),
              child: Icon(uiData['icon'], color: uiData['color']),
            ),
            title: Text(
              uiData['label'],
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            subtitle: Text(
              DateFormat(
                'dd/MM/yyyy HH:mm',
              ).format(tx.transactionTimestamp.toLocal()),
              style: const TextStyle(color: Colors.grey),
            ),
            trailing: trailingWidget,
          ),
        );
      },
    );
  }
}
