import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/publicsite/app/use_cases.dart';
import 'package:mobile_central/services/modules/publicsite/domain/entities.dart';
import 'package:mobile_central/services/modules/publicsite/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockPublicSiteRepository implements IPublicSiteRepository {
  final List<String> calls = [];

  PublicBusiness? getBusinessPageResult;
  PaginatedResponse<PublicProduct>? getCatalogResult;
  PublicProduct? getProductResult;
  Map<String, dynamic>? submitContactResult;

  Exception? errorToThrow;

  String? capturedSlug;
  GetPublicCatalogParams? capturedCatalogParams;
  String? capturedProductId;
  ContactFormDTO? capturedContactData;

  @override
  Future<PublicBusiness> getBusinessPage(String slug) async {
    calls.add('getBusinessPage');
    capturedSlug = slug;
    if (errorToThrow != null) throw errorToThrow!;
    return getBusinessPageResult!;
  }

  @override
  Future<PaginatedResponse<PublicProduct>> getCatalog(String slug, GetPublicCatalogParams? params) async {
    calls.add('getCatalog');
    capturedSlug = slug;
    capturedCatalogParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getCatalogResult!;
  }

  @override
  Future<PublicProduct> getProduct(String slug, String productId) async {
    calls.add('getProduct');
    capturedSlug = slug;
    capturedProductId = productId;
    if (errorToThrow != null) throw errorToThrow!;
    return getProductResult!;
  }

  @override
  Future<Map<String, dynamic>> submitContact(String slug, ContactFormDTO data) async {
    calls.add('submitContact');
    capturedSlug = slug;
    capturedContactData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return submitContactResult!;
  }
}

// --- Helpers ---

PublicBusiness _makeBusiness({int id = 1, String name = 'TestStore'}) {
  return PublicBusiness(
    id: id, name: name, code: 'teststore', description: 'A store',
    logoUrl: '', primaryColor: '', secondaryColor: '', tertiaryColor: '',
    quaternaryColor: '', navbarImageUrl: '', featuredProducts: [],
  );
}

PublicProduct _makeProduct({String id = '1', String name = 'TestProduct'}) {
  return PublicProduct(
    id: id, name: name, description: '', shortDescription: '', price: 10.0,
    currency: 'COP', imageUrl: '', sku: 'SKU', stockQuantity: 5,
    category: '', brand: '', isFeatured: false, createdAt: '',
  );
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1, perPage: 10, total: 1, lastPage: 1,
    hasNext: false, hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockPublicSiteRepository mockRepo;
  late PublicSiteUseCases useCases;

  setUp(() {
    mockRepo = MockPublicSiteRepository();
    useCases = PublicSiteUseCases(mockRepo);
  });

  group('getBusinessPage', () {
    test('delegates to repository with correct slug', () async {
      mockRepo.getBusinessPageResult = _makeBusiness(id: 1, name: 'My Store');

      final result = await useCases.getBusinessPage('mystore');

      expect(result.id, 1);
      expect(result.name, 'My Store');
      expect(mockRepo.capturedSlug, 'mystore');
      expect(mockRepo.calls, ['getBusinessPage']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.getBusinessPage('invalid'), throwsException);
    });
  });

  group('getCatalog', () {
    test('delegates to repository with correct slug and params', () async {
      final expected = PaginatedResponse<PublicProduct>(
        data: [_makeProduct()],
        pagination: _makePagination(),
      );
      mockRepo.getCatalogResult = expected;
      final params = GetPublicCatalogParams(page: 1, pageSize: 20);

      final result = await useCases.getCatalog('mystore', params);

      expect(result.data.length, 1);
      expect(mockRepo.capturedSlug, 'mystore');
      expect(mockRepo.capturedCatalogParams, params);
      expect(mockRepo.calls, ['getCatalog']);
    });

    test('passes null params to repository', () async {
      mockRepo.getCatalogResult = PaginatedResponse<PublicProduct>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getCatalog('mystore', null);

      expect(mockRepo.capturedCatalogParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getCatalog('mystore', null), throwsException);
    });
  });

  group('getProduct', () {
    test('delegates to repository with correct slug and productId', () async {
      mockRepo.getProductResult = _makeProduct(id: '42', name: 'Widget');

      final result = await useCases.getProduct('mystore', '42');

      expect(result.id, '42');
      expect(result.name, 'Widget');
      expect(mockRepo.capturedSlug, 'mystore');
      expect(mockRepo.capturedProductId, '42');
      expect(mockRepo.calls, ['getProduct']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.getProduct('mystore', '999'), throwsException);
    });
  });

  group('submitContact', () {
    test('delegates to repository with correct slug and data', () async {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');
      mockRepo.submitContactResult = {'success': true};

      final result = await useCases.submitContact('mystore', dto);

      expect(result['success'], true);
      expect(mockRepo.capturedSlug, 'mystore');
      expect(mockRepo.capturedContactData, dto);
      expect(mockRepo.calls, ['submitContact']);
    });

    test('propagates repository errors', () async {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');
      mockRepo.errorToThrow = Exception('Submit failed');

      expect(() => useCases.submitContact('mystore', dto), throwsException);
    });
  });
}
