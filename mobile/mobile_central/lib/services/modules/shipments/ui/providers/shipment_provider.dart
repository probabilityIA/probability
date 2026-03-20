import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/shipment_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class ShipmentProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<Shipment> _shipments = [];
  List<OriginAddress> _originAddresses = [];
  final List<EnvioClickRate> _quotes = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  ShipmentProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Shipment> get shipments => _shipments;
  List<OriginAddress> get originAddresses => _originAddresses;
  List<EnvioClickRate> get quotes => _quotes;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  ShipmentUseCases get _useCases => ShipmentUseCases(ShipmentApiRepository(_apiClient));

  Future<void> fetchShipments({int? businessId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final params = GetShipmentsParams(page: _page, pageSize: _pageSize, businessId: businessId);
      final response = await _useCases.getShipments(params);
      _shipments = response.data;
      _pagination = response.pagination;
    } catch (e) { _error = parseError(e); }
    _isLoading = false; notifyListeners();
  }

  Future<void> fetchOriginAddresses({int? businessId}) async {
    try {
      _originAddresses = await _useCases.getOriginAddresses(businessId: businessId);
      notifyListeners();
    } catch (e) { _error = parseError(e); notifyListeners(); }
  }

  Future<Map<String, dynamic>?> quoteShipment(EnvioClickQuoteRequest req) async {
    try { return await _useCases.quoteShipment(req); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<Map<String, dynamic>?> generateGuide(EnvioClickQuoteRequest req) async {
    try { return await _useCases.generateGuide(req); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<Map<String, dynamic>?> trackShipment(String trackingNumber) async {
    try { return await _useCases.trackShipment(trackingNumber); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  void setPage(int page) { _page = page; }
}
