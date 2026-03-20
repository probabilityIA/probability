import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/website_config_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class WebsiteConfigProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  WebsiteConfigData? _config;
  bool _isLoading = false;
  String? _error;

  WebsiteConfigProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  WebsiteConfigData? get config => _config;
  bool get isLoading => _isLoading;
  String? get error => _error;

  WebsiteConfigUseCases get _useCases =>
      WebsiteConfigUseCases(WebsiteConfigApiRepository(_apiClient));

  Future<void> fetchConfig({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _config = await _useCases.getConfig(businessId: businessId);
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _config = await _useCases.updateConfig(data, businessId: businessId);
      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = parseError(e);
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }
}
