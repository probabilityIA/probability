import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/pay_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class PayProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<PaymentGatewayType> _paymentGatewayTypes = [];
  bool _isLoading = false;
  String? _error;

  PayProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<PaymentGatewayType> get paymentGatewayTypes => _paymentGatewayTypes;
  bool get isLoading => _isLoading;
  String? get error => _error;

  PayUseCases get _useCases =>
      PayUseCases(PayGatewayApiRepository(_apiClient));

  Future<void> fetchPaymentGatewayTypes() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.listPaymentGatewayTypes();
      _paymentGatewayTypes = response.data;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }
}
