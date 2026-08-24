import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/vehicles_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class VehicleProvider extends ChangeNotifier {
  VehicleProvider({required ApiClient apiClient}) : _apiClient = apiClient {
    list = PagedListController<VehicleInfo>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  late final PagedListController<VehicleInfo> list;

  int? _businessId;
  String? _typeFilter;
  String? _statusFilter;
  String? _error;

  List<VehicleInfo> get vehicles => list.loadedItems;
  bool get isLoading => list.isLoading;
  String? get error => _error ?? list.error;

  VehicleUseCases get _useCases => VehicleUseCases(VehicleApiRepository(_apiClient));

  Future<PaginatedResponse<VehicleInfo>> _fetchPage(int page, int pageSize) {
    return _useCases.getVehicles(GetVehiclesParams(
      page: page,
      pageSize: pageSize,
      businessId: _businessId,
      type: _typeFilter,
      status: _statusFilter,
    ));
  }

  Future<void> fetchVehicles({int? businessId, String? type, String? status}) {
    _businessId = businessId;
    _typeFilter = type;
    _statusFilter = status;
    _error = null;
    return list.refresh();
  }

  Future<VehicleInfo?> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    try { return await _useCases.createVehicle(data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<VehicleInfo?> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    try { return await _useCases.updateVehicle(id, data, businessId: businessId); } catch (e) { _error = parseError(e); notifyListeners(); return null; }
  }

  Future<bool> deleteVehicle(int id, {int? businessId}) async {
    try { await _useCases.deleteVehicle(id, businessId: businessId); return true; } catch (e) { _error = parseError(e); notifyListeners(); return false; }
  }

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
