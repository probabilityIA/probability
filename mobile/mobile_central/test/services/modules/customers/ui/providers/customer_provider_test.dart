import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/customers/app/use_cases.dart';
import 'package:mobile_central/services/modules/customers/domain/entities.dart';
import 'package:mobile_central/services/modules/customers/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockCustomerRepository implements ICustomerRepository {
  PaginatedResponse<CustomerInfo>? getCustomersResult;
  CustomerDetail? getCustomerByIdResult;
  CustomerInfo? createCustomerResult;
  CustomerInfo? updateCustomerResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;
  int? capturedDeleteBusinessId;

  @override
  Future<PaginatedResponse<CustomerInfo>> getCustomers(GetCustomersParams? params) async {
    calls.add('getCustomers');
    if (errorToThrow != null) throw errorToThrow!;
    return getCustomersResult!;
  }

  @override
  Future<CustomerDetail> getCustomerById(int id, {int? businessId}) async {
    calls.add('getCustomerById');
    if (errorToThrow != null) throw errorToThrow!;
    return getCustomerByIdResult!;
  }

  @override
  Future<CustomerInfo> createCustomer(CreateCustomerDTO data, {int? businessId}) async {
    calls.add('createCustomer');
    if (errorToThrow != null) throw errorToThrow!;
    return createCustomerResult!;
  }

  @override
  Future<CustomerInfo> updateCustomer(int id, UpdateCustomerDTO data, {int? businessId}) async {
    calls.add('updateCustomer');
    if (errorToThrow != null) throw errorToThrow!;
    return updateCustomerResult!;
  }

  @override
  Future<void> deleteCustomer(int id, {int? businessId}) async {
    calls.add('deleteCustomer');
    capturedDeleteId = id;
    capturedDeleteBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableCustomerProvider {
  final CustomerUseCases _useCases;

  List<CustomerInfo> _customers = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _searchFilter = '';

  final List<String> notifications = [];

  TestableCustomerProvider(this._useCases);

  List<CustomerInfo> get customers => _customers;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchCustomers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final params = GetCustomersParams(
        page: _page,
        pageSize: _pageSize,
        businessId: businessId,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
      );
      final response = await _useCases.getCustomers(params);
      _customers = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<CustomerInfo?> createCustomer(CreateCustomerDTO data) async {
    try {
      return await _useCases.createCustomer(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateCustomer(int id, UpdateCustomerDTO data) async {
    try {
      await _useCases.updateCustomer(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteCustomer(int id) async {
    try {
      await _useCases.deleteCustomer(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  void setPage(int page) {
    _page = page;
  }

  void setSearch(String search) {
    _searchFilter = search;
    _page = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _page = 1;
  }
}

// --- Helpers ---

Pagination _makePagination({
  int currentPage = 1,
  int total = 5,
  int lastPage = 1,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 20,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
  );
}

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

// --- Tests ---

void main() {
  late MockCustomerRepository mockRepo;
  late CustomerUseCases useCases;
  late TestableCustomerProvider provider;

  setUp(() {
    mockRepo = MockCustomerRepository();
    useCases = CustomerUseCases(mockRepo);
    provider = TestableCustomerProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty customers list', () {
      expect(provider.customers, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });

    test('starts at page 1', () {
      expect(provider.page, 1);
    });
  });

  group('fetchCustomers', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getCustomersResult = PaginatedResponse<CustomerInfo>(
        data: [_makeCustomer()],
        pagination: _makePagination(),
      );

      await provider.fetchCustomers();

      expect(provider.notifications.length, 2);
    });

    test('populates customers and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getCustomersResult = PaginatedResponse<CustomerInfo>(
        data: [_makeCustomer(id: 1), _makeCustomer(id: 2, name: 'Jane')],
        pagination: pagination,
      );

      await provider.fetchCustomers();

      expect(provider.customers.length, 2);
      expect(provider.customers[0].id, 1);
      expect(provider.customers[1].name, 'Jane');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchCustomers();

      expect(provider.error, contains('Server error'));
      expect(provider.customers, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchCustomers();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getCustomersResult = PaginatedResponse<CustomerInfo>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchCustomers();

      expect(provider.error, isNull);
    });
  });

  group('createCustomer', () {
    test('returns created customer on success', () async {
      final dto = CreateCustomerDTO(name: 'New');
      mockRepo.createCustomerResult = _makeCustomer(id: 10, name: 'New');

      final result = await provider.createCustomer(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'New');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateCustomerDTO(name: 'Fail');

      final result = await provider.createCustomer(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateCustomer', () {
    test('returns true on success', () async {
      final dto = UpdateCustomerDTO(name: 'Updated');
      mockRepo.updateCustomerResult = _makeCustomer(id: 5, name: 'Updated');

      final result = await provider.updateCustomer(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateCustomerDTO(name: 'Fail');

      final result = await provider.updateCustomer(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteCustomer', () {
    test('returns true on success', () async {
      final result = await provider.deleteCustomer(7);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteCustomer(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('setPage', () {
    test('updates page value', () {
      provider.setPage(3);

      expect(provider.page, 3);
    });
  });

  group('setSearch', () {
    test('updates search and resets page to 1', () {
      provider.setPage(5);
      provider.setSearch('query');

      expect(provider.page, 1);
    });
  });

  group('resetFilters', () {
    test('resets search and page', () {
      provider.setPage(3);
      provider.setSearch('test');
      provider.resetFilters();

      expect(provider.page, 1);
    });
  });
}
