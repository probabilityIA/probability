import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/website_config/app/use_cases.dart';
import 'package:mobile_central/services/modules/website_config/domain/entities.dart';
import 'package:mobile_central/services/modules/website_config/domain/ports.dart';

// --- Manual Mock ---

class MockWebsiteConfigRepository implements IWebsiteConfigRepository {
  final List<String> calls = [];

  WebsiteConfigData? getConfigResult;
  WebsiteConfigData? updateConfigResult;

  Exception? errorToThrow;

  int? capturedBusinessId;
  UpdateWebsiteConfigDTO? capturedUpdateData;

  @override
  Future<WebsiteConfigData> getConfig({int? businessId}) async {
    calls.add('getConfig');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getConfigResult!;
  }

  @override
  Future<WebsiteConfigData> updateConfig(UpdateWebsiteConfigDTO data, {int? businessId}) async {
    calls.add('updateConfig');
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateConfigResult!;
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

  setUp(() {
    mockRepo = MockWebsiteConfigRepository();
    useCases = WebsiteConfigUseCases(mockRepo);
  });

  group('getConfig', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getConfigResult = _makeConfig(id: 1, template: 'modern');

      final result = await useCases.getConfig();

      expect(result.id, 1);
      expect(result.template, 'modern');
      expect(mockRepo.calls, ['getConfig']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getConfigResult = _makeConfig();

      await useCases.getConfig(businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });

    test('passes null businessId when not provided', () async {
      mockRepo.getConfigResult = _makeConfig();

      await useCases.getConfig();

      expect(mockRepo.capturedBusinessId, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getConfig(), throwsException);
    });
  });

  group('updateConfig', () {
    test('delegates to repository with correct data', () async {
      final dto = UpdateWebsiteConfigDTO(
        template: 'basic',
        showHero: false,
      );
      mockRepo.updateConfigResult = _makeConfig(id: 1, template: 'basic');

      final result = await useCases.updateConfig(dto);

      expect(result.template, 'basic');
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateConfig']);
    });

    test('passes businessId to repository', () async {
      final dto = UpdateWebsiteConfigDTO(template: 'modern');
      mockRepo.updateConfigResult = _makeConfig();

      await useCases.updateConfig(dto, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateWebsiteConfigDTO(template: 'x');

      expect(() => useCases.updateConfig(dto), throwsException);
    });
  });
}
