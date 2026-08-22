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
  String _searchFilter = '';

  List<Shipment> get shipments => list.loadedItems;
  List<OriginAddress> get originAddresses => _originAddresses;
  List<EnvioClickRate> get quotes => _quotes;
  bool get isLoading => list.isLoading || _isBusy;
  bool get isLoadingMore => list.isLoadingMore;
  bool get hasMore => list.hasMore;
  String get statusFilter => _statusFilter;
  String? get error => _error ?? list.error;

  ShipmentUseCases get _useCases =>
      _injectedUseCases ?? ShipmentUseCases(ShipmentApiRepository(_apiClient));

  void setShipmentFilters({String? status, String? search}) {
    _statusFilter = status ?? _statusFilter;
    _searchFilter = search ?? _searchFilter;
  }

  Future<PaginatedResponse<Shipment>> _fetchPage(int page, int pageSize) {
    return _useCases.getShipments(GetShipmentsParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      status: _statusFilter.isNotEmpty ? _statusFilter : null,
      trackingNumber: _searchFilter.isNotEmpty ? _searchFilter : null,
    ));
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
