import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/products/app/use_cases.dart';
import 'package:mobile_central/services/modules/products/domain/entities.dart';
import 'package:mobile_central/services/modules/products/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockProductRepository implements IProductRepository {
  final List<String> calls = [];

  PaginatedResponse<Product>? getProductsResult;
  Product? getProductByIdResult;
  Product? createProductResult;
  Product? updateProductResult;
  List<ProductIntegration>? getProductIntegrationsResult;

  Exception? errorToThrow;

  GetProductsParams? capturedGetProductsParams;
  String? capturedId;
  int? capturedBusinessId;
  CreateProductDTO? capturedCreateData;
  UpdateProductDTO? capturedUpdateData;
  String? capturedDeleteId;
  String? capturedIntegrationProductId;
  AddProductIntegrationDTO? capturedAddIntegrationData;
  int? capturedRemoveIntegrationId;

  @override
  Future<PaginatedResponse<Product>> getProducts(GetProductsParams? params) async {
    calls.add('getProducts');
    capturedGetProductsParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getProductsResult!;
  }

  @override
  Future<Product> getProductById(String id, {int? businessId}) async {
    calls.add('getProductById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getProductByIdResult!;
  }

  @override
  Future<Product> createProduct(CreateProductDTO data, {int? businessId}) async {
    calls.add('createProduct');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createProductResult!;
  }

  @override
  Future<Product> updateProduct(String id, UpdateProductDTO data, {int? businessId}) async {
    calls.add('updateProduct');
    capturedId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateProductResult!;
  }

  @override
  Future<void> deleteProduct(String id, {int? businessId}) async {
    calls.add('deleteProduct');
    capturedDeleteId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<ProductIntegration>> getProductIntegrations(String productId, {int? businessId}) async {
    calls.add('getProductIntegrations');
    capturedIntegrationProductId = productId;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getProductIntegrationsResult!;
  }

  @override
  Future<void> addProductIntegration(String productId, AddProductIntegrationDTO data, {int? businessId}) async {
    calls.add('addProductIntegration');
    capturedIntegrationProductId = productId;
    capturedAddIntegrationData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> removeProductIntegration(String productId, int integrationId, {int? businessId}) async {
    calls.add('removeProductIntegration');
    capturedIntegrationProductId = productId;
    capturedRemoveIntegrationId = integrationId;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Product _makeProduct({String id = '1', String name = 'TestProduct'}) {
  return Product(
    id: id,
    createdAt: '2026-01-01',
    updatedAt: '2026-01-01',
    businessId: 1,
    sku: 'SKU-$id',
    name: name,
    price: 10.0,
    currency: 'COP',
    stock: 5,
    manageStock: true,
    status: 'active',
    isActive: true,
  );
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 10,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockProductRepository mockRepo;
  late ProductUseCases useCases;

  setUp(() {
    mockRepo = MockProductRepository();
    useCases = ProductUseCases(mockRepo);
  });

  group('getProducts', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Product>(
        data: [_makeProduct()],
        pagination: _makePagination(),
      );
      mockRepo.getProductsResult = expected;
      final params = GetProductsParams(page: 1, pageSize: 10);

      final result = await useCases.getProducts(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'TestProduct');
      expect(mockRepo.calls, ['getProducts']);
      expect(mockRepo.capturedGetProductsParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getProductsResult = PaginatedResponse<Product>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getProducts(null);

      expect(mockRepo.capturedGetProductsParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getProducts(null), throwsException);
    });
  });

  group('getProductById', () {
    test('delegates to repository with correct id and businessId', () async {
      mockRepo.getProductByIdResult = _makeProduct(id: '42', name: 'Found');

      final result = await useCases.getProductById('42', businessId: 5);

      expect(result.id, '42');
      expect(result.name, 'Found');
      expect(mockRepo.capturedId, '42');
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getProductById']);
    });
  });

  group('createProduct', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateProductDTO(businessId: 1, sku: 'S', name: 'N', price: 10, stock: 5);
      mockRepo.createProductResult = _makeProduct(id: '99', name: 'N');

      final result = await useCases.createProduct(dto, businessId: 1);

      expect(result.id, '99');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.capturedBusinessId, 1);
      expect(mockRepo.calls, ['createProduct']);
    });

    test('propagates repository errors', () async {
      final dto = CreateProductDTO(businessId: 1, sku: 'S', name: 'N', price: 10, stock: 5);
      mockRepo.errorToThrow = Exception('Create failed');

      expect(() => useCases.createProduct(dto), throwsException);
    });
  });

  group('updateProduct', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateProductDTO(name: 'Updated');
      mockRepo.updateProductResult = _makeProduct(id: '5', name: 'Updated');

      final result = await useCases.updateProduct('5', dto, businessId: 2);

      expect(result.id, '5');
      expect(result.name, 'Updated');
      expect(mockRepo.capturedId, '5');
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.capturedBusinessId, 2);
      expect(mockRepo.calls, ['updateProduct']);
    });
  });

  group('deleteProduct', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteProduct('7', businessId: 3);

      expect(mockRepo.capturedDeleteId, '7');
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['deleteProduct']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      expect(() => useCases.deleteProduct('7'), throwsException);
    });
  });

  group('getProductIntegrations', () {
    test('delegates to repository with correct productId', () async {
      mockRepo.getProductIntegrationsResult = [
        ProductIntegration(
          id: 1,
          productId: '42',
          integrationId: 10,
          externalProductId: 'ext-1',
          createdAt: '2026-01-01',
          updatedAt: '2026-01-01',
        ),
      ];

      final result = await useCases.getProductIntegrations('42', businessId: 5);

      expect(result.length, 1);
      expect(result[0].productId, '42');
      expect(mockRepo.capturedIntegrationProductId, '42');
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getProductIntegrations']);
    });
  });

  group('addProductIntegration', () {
    test('delegates to repository with correct data', () async {
      final dto = AddProductIntegrationDTO(integrationId: 10, externalProductId: 'ext-1');

      await useCases.addProductIntegration('42', dto, businessId: 5);

      expect(mockRepo.capturedIntegrationProductId, '42');
      expect(mockRepo.capturedAddIntegrationData, dto);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['addProductIntegration']);
    });
  });

  group('removeProductIntegration', () {
    test('delegates to repository with correct ids', () async {
      await useCases.removeProductIntegration('42', 10, businessId: 5);

      expect(mockRepo.capturedIntegrationProductId, '42');
      expect(mockRepo.capturedRemoveIntegrationId, 10);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['removeProductIntegration']);
    });
  });
}
