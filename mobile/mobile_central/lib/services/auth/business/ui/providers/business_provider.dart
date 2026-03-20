import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/business_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class BusinessProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final BusinessUseCases? _injectedUseCases;

  List<Business> _businesses = [];
  List<BusinessSimple> _businessesSimple = [];
  List<BusinessType> _businessTypes = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int? _selectedBusinessId;

  BusinessProvider({required ApiClient apiClient, BusinessUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases;

  List<Business> get businesses => _businesses;
  List<BusinessSimple> get businessesSimple => _businessesSimple;
  List<BusinessType> get businessTypes => _businessTypes;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int? get selectedBusinessId => _selectedBusinessId;

  BusinessUseCases get _useCases =>
      _injectedUseCases ?? BusinessUseCases(BusinessApiRepository(_apiClient));

  void setSelectedBusinessId(int? id) {
    _selectedBusinessId = id;
    notifyListeners();
  }

  Future<void> fetchBusinesses({GetBusinessesParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.getBusinesses(params);
      _businesses = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchBusinessesSimple() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _businessesSimple = await _useCases.getBusinessesSimple();
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchBusinessTypes() async {
    try {
      _businessTypes = await _useCases.getBusinessTypes();
      notifyListeners();
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
    }
  }

  Future<Business?> createBusiness(CreateBusinessDTO data) async {
    try {
      return await _useCases.createBusiness(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateBusiness(int id, UpdateBusinessDTO data) async {
    try {
      await _useCases.updateBusiness(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteBusiness(int id) async {
    try {
      await _useCases.deleteBusiness(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> activateBusiness(int id) async {
    try {
      await _useCases.activateBusiness(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deactivateBusiness(int id) async {
    try {
      await _useCases.deactivateBusiness(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }
}
