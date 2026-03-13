import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/vehicles_repository.dart';

class VehicleProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  List<VehicleInfo> _vehicles = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  VehicleProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<VehicleInfo> get vehicles => _vehicles;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  VehicleUseCases get _useCases => VehicleUseCases(VehicleApiRepository(_apiClient));

  Future<void> fetchVehicles({int? businessId, String? type, String? status}) async {
    _isLoading = true; _error = null; notifyListeners();
    try {
      final params = GetVehiclesParams(page: _page, pageSize: _pageSize, businessId: businessId, type: type, status: status);
      final response = await _useCases.getVehicles(params);
      _vehicles = response.data; _pagination = response.pagination;
    } catch (e) { _error = e.toString(); }
    _isLoading = false; notifyListeners();
  }

  Future<VehicleInfo?> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    try { return await _useCases.createVehicle(data, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
  }

  Future<VehicleInfo?> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    try { return await _useCases.updateVehicle(id, data, businessId: businessId); } catch (e) { _error = e.toString(); notifyListeners(); return null; }
  }

  Future<bool> deleteVehicle(int id, {int? businessId}) async {
    try { await _useCases.deleteVehicle(id, businessId: businessId); return true; } catch (e) { _error = e.toString(); notifyListeners(); return false; }
  }

  void setPage(int page) { _page = page; }
}
