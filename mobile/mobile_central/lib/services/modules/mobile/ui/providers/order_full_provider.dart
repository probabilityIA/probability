import 'package:flutter/foundation.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/mobile_repository.dart';

class OrderFullProvider extends ChangeNotifier {
  OrderFullProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  final ApiClient _apiClient;

  MobileOrderFull? _orderFull;
  bool _isLoading = false;
  String? _error;

  MobileOrderFull? get orderFull => _orderFull;
  bool get isLoading => _isLoading;
  String? get error => _error;

  MobileUseCases get _useCases => MobileUseCases(MobileApiRepository(_apiClient));

  Future<void> load(String orderId, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _orderFull = await _useCases.getOrderFull(orderId, businessId: businessId);
    } catch (e) {
      _error = parseError(e);
      _orderFull = null;
    }

    _isLoading = false;
    notifyListeners();
  }

  void clear() {
    _orderFull = null;
    _error = null;
  }
}
