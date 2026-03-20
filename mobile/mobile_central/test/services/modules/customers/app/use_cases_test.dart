import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/customers/app/use_cases.dart';
import 'package:mobile_central/services/modules/customers/domain/entities.dart';
import 'package:mobile_central/services/modules/customers/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockCustomerRepository implements ICustomerRepository {
  final List<String> calls = [];

  PaginatedResponse<CustomerInfo>? getCustomersResult;
  CustomerDetail? getCustomerByIdResult;
  CustomerInfo? createCustomerResult;
  CustomerInfo? updateCustomerResult;

  Exception? errorToThrow;

  GetCustomersParams? capturedGetCustomersParams;
  int? capturedId;
  int? capturedBusinessId;
  CreateCustomerDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateCustomerDTO? capturedUpdateData;
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<CustomerInfo>> getCustomers(GetCustomersParams? params) async {
    calls.add('getCustomers');
    capturedGetCustomersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getCustomersResult!;
  }

  @override
  Future<CustomerDetail> getCustomerById(int id, {int? businessId}) async {
    calls.add('getCustomerById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getCustomerByIdResult!;
  }

  @override
  Future<CustomerInfo> createCustomer(CreateCustomerDTO data, {int? businessId}) async {
    calls.add('createCustomer');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createCustomerResult!;
  }

  @override
  Future<CustomerInfo> updateCustomer(int id, UpdateCustomerDTO data, {int? businessId}) async {
    calls.add('updateCustomer');
    capturedUpdateId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateCustomerResult!;
  }

  @override
  Future<void> deleteCustomer(int id, {int? businessId}) async {
    calls.add('deleteCustomer');
    capturedDeleteId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

CustomerInfo _makeCustomer({int id = 1, String name = 'TestCustomer'}) {
  return CustomerInfo(
    id: id,
    businessId: 1,
    name: name,
    phone: '555-1234',
    createdAt: '2026-01-01',
    updatedAt: '2026-01-02',
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
  late MockCustomerRepository mockRepo;
  late CustomerUseCases useCases;

  setUp(() {
    mockRepo = MockCustomerRepository();
    useCases = CustomerUseCases(mockRepo);
  });

  group('getCustomers', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<CustomerInfo>(
        data: [_makeCustomer()],
        pagination: _makePagination(),
      );
      mockRepo.getCustomersResult = expected;
      final params = GetCustomersParams(page: 1, pageSize: 10);

      final result = await useCases.getCustomers(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'TestCustomer');
      expect(mockRepo.calls, ['getCustomers']);
      expect(mockRepo.capturedGetCustomersParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getCustomersResult = PaginatedResponse<CustomerInfo>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getCustomers(null);

      expect(mockRepo.capturedGetCustomersParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getCustomers(null), throwsException);
    });
  });

  group('getCustomerById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getCustomerByIdResult = CustomerDetail(
        id: 42,
        businessId: 1,
        name: 'Found',
        phone: '555',
        createdAt: '',
        updatedAt: '',
        orderCount: 5,
        totalSpent: 1000.0,
      );

      final result = await useCases.getCustomerById(42);

      expect(result.id, 42);
      expect(result.name, 'Found');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getCustomerById']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getCustomerByIdResult = CustomerDetail(
        id: 1,
        businessId: 5,
        name: 'Test',
        phone: '555',
        createdAt: '',
        updatedAt: '',
        orderCount: 0,
        totalSpent: 0,
      );

      await useCases.getCustomerById(1, businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });
  });

  group('createCustomer', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateCustomerDTO(name: 'NewCustomer', email: 'new@test.com');
      mockRepo.createCustomerResult = _makeCustomer(id: 99, name: 'NewCustomer');

      final result = await useCases.createCustomer(dto);

      expect(result.id, 99);
      expect(result.name, 'NewCustomer');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createCustomer']);
    });

    test('passes businessId to repository', () async {
      final dto = CreateCustomerDTO(name: 'Test');
      mockRepo.createCustomerResult = _makeCustomer();

      await useCases.createCustomer(dto, businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
    });
  });

  group('updateCustomer', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateCustomerDTO(name: 'Updated');
      mockRepo.updateCustomerResult = _makeCustomer(id: 5, name: 'Updated');

      final result = await useCases.updateCustomer(5, dto);

      expect(result.id, 5);
      expect(result.name, 'Updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateCustomer']);
    });

    test('passes businessId to repository', () async {
      final dto = UpdateCustomerDTO(name: 'Test');
      mockRepo.updateCustomerResult = _makeCustomer();

      await useCases.updateCustomer(1, dto, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });
  });

  group('deleteCustomer', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteCustomer(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteCustomer']);
    });

    test('passes businessId to repository', () async {
      await useCases.deleteCustomer(7, businessId: 2);

      expect(mockRepo.capturedBusinessId, 2);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteCustomer(7), throwsException);
    });
  });
}
