import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/resource_repository.dart';

class ResourceProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<Resource> _resources = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  ResourceProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Resource> get resources => _resources;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  ResourceUseCases get _useCases =>
      ResourceUseCases(ResourceApiRepository(_apiClient));

  Future<void> fetchResources({GetResourcesParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.getResources(params);
      _resources = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Resource?> createResource(CreateResourceDTO data) async {
    try {
      return await _useCases.createResource(data);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateResource(int id, UpdateResourceDTO data) async {
    try {
      await _useCases.updateResource(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteResource(int id) async {
    try {
      await _useCases.deleteResource(id);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}
