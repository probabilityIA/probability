import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/publicsite_repository.dart';

class PublicSiteProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  PublicBusiness? _business;
  List<PublicProduct> _catalogProducts = [];
  PublicProduct? _selectedProduct;
  Pagination? _catalogPagination;
  bool _isLoading = false;
  String? _error;
  int _catalogPage = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  String _categoryFilter = '';

  PublicSiteProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  PublicBusiness? get business => _business;
  List<PublicProduct> get catalogProducts => _catalogProducts;
  PublicProduct? get selectedProduct => _selectedProduct;
  Pagination? get catalogPagination => _catalogPagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get catalogPage => _catalogPage;

  PublicSiteUseCases get _useCases =>
      PublicSiteUseCases(PublicSiteApiRepository(_apiClient));

  Future<void> fetchBusinessPage(String slug) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _business = await _useCases.getBusinessPage(slug);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchCatalog(String slug) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetPublicCatalogParams(
        page: _catalogPage,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        category: _categoryFilter.isNotEmpty ? _categoryFilter : null,
      );
      final response = await _useCases.getCatalog(slug, params);
      _catalogProducts = response.data;
      _catalogPagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchProduct(String slug, String productId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _selectedProduct = await _useCases.getProduct(slug, productId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Map<String, dynamic>?> submitContact(String slug, ContactFormDTO data) async {
    try {
      return await _useCases.submitContact(slug, data);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  void setCatalogPage(int page) {
    _catalogPage = page;
  }

  void setFilters({String? search, String? category}) {
    _searchFilter = search ?? _searchFilter;
    _categoryFilter = category ?? _categoryFilter;
    _catalogPage = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _categoryFilter = '';
    _catalogPage = 1;
  }
}
