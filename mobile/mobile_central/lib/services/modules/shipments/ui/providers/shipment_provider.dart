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

  bool _isLoadingMore = false;
  String _statusFilter = '';
  String _searchFilter = '';

  bool get isLoadingMore => _isLoadingMore;
  bool get hasMore => _pagination?.hasNext ?? false;
  String get statusFilter => _statusFilter;

  void setShipmentFilters({String? status, String? search}) {
    _statusFilter = status ?? _statusFilter;
    _searchFilter = search ?? _searchFilter;
    _page = 1;
  }

  GetShipmentsParams _params(int? businessId, int page) => GetShipmentsParams(
        page: page,
        pageSize: _pageSize,
        businessId: businessId,
        status: _statusFilter.isNotEmpty ? _statusFilter : null,
        trackingNumber: _searchFilter.isNotEmpty ? _searchFilter : null,
      );

  Future<void> fetchShipments({int? businessId}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final response = await _useCases.getShipments(_params(businessId, _page));
      _shipments = response.data;
      _pagination = response.pagination;
    } catch (e) { _error = parseError(e); }
    _isLoading = false; notifyListeners();
  }

  Future<void> loadMoreShipments({int? businessId}) async {
    if (_isLoading || _isLoadingMore || !hasMore) return;
    _isLoadingMore = true; notifyListeners();
    try {
      final response = await _useCases.getShipments(_params(businessId, _page + 1));
      _shipments = [..._shipments, ...response.data];
      _pagination = response.pagination;
      _page += 1;
    } catch (e) { _error = parseError(e); }
    _isLoadingMore = false; notifyListeners();
  }

  Future<bool> cancelShipment(int id, {int? businessId}) async {
    _error = null;
    try {
      await _useCases.cancelShipment(id.toString());
      await fetchShipments(businessId: businessId);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Shipment? shipmentById(int id) {
    for (final shipment in _shipments) {
      if (shipment.id == id) return shipment;
    }
    return null;
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
