import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/business_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class BusinessProvider extends ChangeNotifier {
  final ApiClient _apiClient;
  final BusinessUseCases? _injectedUseCases;

  List<BusinessSimple> _businessesSimple = [];
  List<BusinessType> _businessTypes = [];
  bool _isLoading = false;
  String? _error;
  int? _selectedBusinessId;
  GetBusinessesParams? _params;

  late final PagedListController<Business> list;

  BusinessProvider({required ApiClient apiClient, BusinessUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<Business>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  List<Business> get businesses => list.loadedItems;
  List<BusinessSimple> get businessesSimple => _businessesSimple;
  List<BusinessType> get businessTypes => _businessTypes;
  bool get isLoading => _isLoading || list.isLoading;
  String? get error => _error ?? list.error;
  int? get selectedBusinessId => _selectedBusinessId;

  BusinessUseCases get _useCases =>
      _injectedUseCases ?? BusinessUseCases(BusinessApiRepository(_apiClient));

  void setSelectedBusinessId(int? id) {
    _selectedBusinessId = id;
    notifyListeners();
  }

  Future<PaginatedResponse<Business>> _fetchPage(int page, int pageSize) {
    final base = _params;
    return _useCases.getBusinesses(GetBusinessesParams(
      page: page,
      pageSize: pageSize,
      name: base?.name,
      businessTypeId: base?.businessTypeId,
    ));
  }

  Future<void> fetchBusinesses({GetBusinessesParams? params}) {
    _params = params;
    _error = null;
    return list.refresh();
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
