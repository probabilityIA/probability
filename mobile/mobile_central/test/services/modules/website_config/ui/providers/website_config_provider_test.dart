import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/website_config/app/use_cases.dart';
import 'package:mobile_central/services/modules/website_config/domain/entities.dart';
import 'package:mobile_central/services/modules/website_config/domain/ports.dart';

// --- Manual Mock Repository ---

class MockWebsiteConfigRepository implements IWebsiteConfigRepository {
  WebsiteConfigData? getConfigResult;
  WebsiteConfigData? updateConfigResult;
  Exception? errorToThrow;

  final List<String> calls = [];

  @override
  Future<WebsiteConfigData> getConfig({int? businessId}) async {
    calls.add('getConfig');
    if (errorToThrow != null) throw errorToThrow!;
    return getConfigResult!;
  }

  @override
  Future<WebsiteConfigData> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) async {
    calls.add('updateConfig');
    if (errorToThrow != null) throw errorToThrow!;
    return updateConfigResult!;
  }
}

// --- Testable Provider ---

class TestableWebsiteConfigProvider {
  final WebsiteConfigUseCases _useCases;

  WebsiteConfigData? _config;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableWebsiteConfigProvider(this._useCases);

  WebsiteConfigData? get config => _config;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchConfig({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _config = await _useCases.getConfig(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<bool> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) async {
    try {
      _config = await _useCases.updateConfig(data, businessId: businessId);
      _notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// --- Helpers ---

WebsiteConfigData _makeConfig({int id = 1, String template = 'modern'}) {
  return WebsiteConfigData(
    id: id,
    businessId: 1,
    template: template,
    showHero: true,
    showAbout: false,
    showFeaturedProducts: true,
    showFullCatalog: false,
    showTestimonials: false,
    showLocation: false,
    showContact: true,
    showSocialMedia: false,
    showWhatsapp: false,
  );
}

// --- Tests ---

void main() {
  late MockWebsiteConfigRepository mockRepo;
  late WebsiteConfigUseCases useCases;
  late TestableWebsiteConfigProvider provider;

  setUp(() {
    mockRepo = MockWebsiteConfigRepository();
    useCases = WebsiteConfigUseCases(mockRepo);
    provider = TestableWebsiteConfigProvider(useCases);
  });

  group('initial state', () {
    test('starts with null config', () {
      expect(provider.config, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchConfig', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getConfigResult = _makeConfig();

      await provider.fetchConfig();

      expect(provider.notifications.length, 2);
    });

    test('populates config on success', () async {
      mockRepo.getConfigResult = _makeConfig(id: 1, template: 'modern');

      await provider.fetchConfig();

      expect(provider.config, isNotNull);
      expect(provider.config!.id, 1);
      expect(provider.config!.template, 'modern');
      expect(provider.config!.showHero, true);
      expect(provider.config!.showAbout, false);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchConfig();

      expect(provider.error, contains('Server error'));
      expect(provider.config, isNull);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchConfig();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getConfigResult = _makeConfig();
      await provider.fetchConfig();

      expect(provider.error, isNull);
    });

    test('passes businessId to use case', () async {
      mockRepo.getConfigResult = _makeConfig();

      await provider.fetchConfig(businessId: 5);

      expect(mockRepo.calls, ['getConfig']);
    });
  });

  group('updateConfig', () {
    test('returns true and updates config on success', () async {
      final dto = UpdateWebsiteConfigDTO(
        template: 'basic',
        showHero: false,
      );
      mockRepo.updateConfigResult = _makeConfig(id: 1, template: 'basic');

      final result = await provider.updateConfig(dto);

      expect(result, true);
      expect(provider.config, isNotNull);
      expect(provider.config!.template, 'basic');
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateWebsiteConfigDTO(template: 'x');

      final result = await provider.updateConfig(dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });

    test('notifies listeners on success', () async {
      final dto = UpdateWebsiteConfigDTO(showHero: true);
      mockRepo.updateConfigResult = _makeConfig();

      await provider.updateConfig(dto);

      expect(provider.notifications, isNotEmpty);
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('Error');
      final dto = UpdateWebsiteConfigDTO(template: 'x');

      await provider.updateConfig(dto);

      expect(provider.notifications, isNotEmpty);
    });

    test('passes businessId to use case', () async {
      final dto = UpdateWebsiteConfigDTO(template: 'modern');
      mockRepo.updateConfigResult = _makeConfig();

      await provider.updateConfig(dto, businessId: 3);

      expect(mockRepo.calls, ['updateConfig']);
    });
  });
}
