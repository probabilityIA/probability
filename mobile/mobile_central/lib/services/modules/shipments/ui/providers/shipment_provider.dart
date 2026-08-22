import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/shipment_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class ShipmentProvider extends ChangeNotifier {
  ShipmentProvider({required ApiClient apiClient, ShipmentUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<Shipment>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final ShipmentUseCases? _injectedUseCases;
  late final PagedListController<Shipment> list;

  List<OriginAddress> _originAddresses = [];
  List<EnvioClickRate> _quotes = [];
  final bool _isBusy = false;
  String? _error;
  int? _businessId;
  String _statusFilter = '';
  String _carrierFilter = '';
  String _searchField = 'tracking_number';
  String _searchTerm = '';
  int _unfilteredTotal = 0;

  List<Shipment> get shipments => list.loadedItems;
  List<OriginAddress> get originAddresses => _originAddresses;
  List<EnvioClickRate> get quotes => _quotes;
  bool get isLoading => list.isLoading || _isBusy;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String get statusFilter => _statusFilter;
  String? get error => _error ?? list.error;
  int get unfilteredTotal => _unfilteredTotal;

  bool get hasFilters =>
      _searchTerm.isNotEmpty ||
      _statusFilter.isNotEmpty ||
      _carrierFilter.isNotEmpty;

  ShipmentUseCases get _useCases =>
      _injectedUseCases ?? ShipmentUseCases(ShipmentApiRepository(_apiClient));

  void setSearch({String? field, String? term}) {
    if (field != null) _searchField = field;
    if (term != null) _searchTerm = term.trim();
  }

  void applyFilters({String? status, String? carrier}) {
    _statusFilter = status ?? '';
    _carrierFilter = carrier ?? '';
  }

  String? _termFor(String field) =>
      _searchField == field && _searchTerm.isNotEmpty ? _searchTerm : null;

  Future<PaginatedResponse<Shipment>> _fetchPage(int page, int pageSize) async {
    final response = await _useCases.getShipments(GetShipmentsParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      status: _statusFilter.isNotEmpty ? _statusFilter : null,
      carrier: _carrierFilter.isNotEmpty ? _carrierFilter : null,
      trackingNumber: _termFor('tracking_number'),
      customerName: _termFor('customer_name'),
      orderId: _termFor('order_id'),
    ));

    if (!hasFilters) _unfilteredTotal = response.pagination.total;
    return response;
  }

  Future<void> fetchShipments({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  Future<void> loadMoreShipments({int? businessId}) {
    if (businessId != null) _businessId = businessId;
    return list.loadMore();
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
    for (final shipment in list.loadedItems) {
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
    _error = null;
    try {
      final result = await _useCases.quoteShipment(req);
      final payload = result['data'] is Map ? result['data'] : result;
      final rawRates = payload is Map ? payload['rates'] : null;
      _quotes = rawRates is List
          ? rawRates
              .whereType<Map>()
              .map((e) => EnvioClickRate.fromJson(Map<String, dynamic>.from(e)))
              .toList()
          : <EnvioClickRate>[];
      notifyListeners();
      return result;
    } catch (e) {
      _error = parseError(e);
      _quotes = <EnvioClickRate>[];
      notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> generateGuide(EnvioClickQuoteRequest req) async {
    try { return await _useCases.generateGuide(req); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<Map<String, dynamic>?> trackShipment(String trackingNumber) async {
    try { return await _useCases.trackShipment(trackingNumber); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
