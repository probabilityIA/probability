import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/modules/invoicing/app/use_cases.dart';
import 'package:mobile_central/services/modules/invoicing/domain/entities.dart';
import 'package:mobile_central/services/modules/invoicing/domain/ports.dart';
import 'package:mobile_central/services/modules/invoicing/ui/providers/invoicing_provider.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IInvoicingRepository
// ---------------------------------------------------------------------------
class MockInvoicingRepository implements IInvoicingRepository {
  final List<String> calls = [];

  PaginatedResponse<Invoice>? getInvoicesResult;
  Invoice? getInvoiceByIdResult;
  Invoice? createInvoiceResult;
  Invoice? cancelInvoiceResult;
  Invoice? retryInvoiceResult;
  CreditNote? createCreditNoteResult;
  PaginatedResponse<InvoicingConfig>? getConfigsResult;
  InvoicingConfig? getConfigByIdResult;
  InvoicingConfig? createConfigResult;
  InvoicingConfig? updateConfigResult;
  List<SyncLog>? getSyncLogsResult;

  Exception? errorToThrow;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<Invoice>> getInvoices(InvoiceFilters filters) async {
    _trackCall('getInvoices');
    return getInvoicesResult!;
  }

  @override
  Future<Invoice> getInvoiceById(int id) async {
    _trackCall('getInvoiceById');
    return getInvoiceByIdResult!;
  }

  @override
  Future<Invoice> createInvoice(CreateInvoiceDTO data) async {
    _trackCall('createInvoice');
    return createInvoiceResult!;
  }

  @override
  Future<Invoice> cancelInvoice(int id) async {
    _trackCall('cancelInvoice');
    return cancelInvoiceResult!;
  }

  @override
  Future<Invoice> retryInvoice(int id) async {
    _trackCall('retryInvoice');
    return retryInvoiceResult!;
  }

  @override
  Future<CreditNote> createCreditNote(CreateCreditNoteDTO data) async {
    _trackCall('createCreditNote');
    return createCreditNoteResult!;
  }

  @override
  Future<PaginatedResponse<InvoicingConfig>> getConfigs(
      ConfigFilters filters) async {
    _trackCall('getConfigs');
    return getConfigsResult!;
  }

  @override
  Future<InvoicingConfig> getConfigById(int id) async {
    _trackCall('getConfigById');
    return getConfigByIdResult!;
  }

  @override
  Future<InvoicingConfig> createConfig(CreateConfigDTO data) async {
    _trackCall('createConfig');
    return createConfigResult!;
  }

  @override
  Future<InvoicingConfig> updateConfig(int id, UpdateConfigDTO data) async {
    _trackCall('updateConfig');
    return updateConfigResult!;
  }

  @override
  Future<void> deleteConfig(int id) async {
    _trackCall('deleteConfig');
  }

  @override
  Future<void> bulkCreateInvoices(BulkCreateInvoicesDTO data) async {
    _trackCall('bulkCreateInvoices');
  }

  @override
  Future<List<SyncLog>> getSyncLogs(int invoiceId) async {
    _trackCall('getSyncLogs');
    return getSyncLogsResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
Pagination _makePagination({int total = 50}) {
  return Pagination(
    currentPage: 1,
    perPage: 20,
    total: total,
    lastPage: 3,
    hasNext: true,
    hasPrev: false,
  );
}

Invoice _makeInvoice({int id = 1, String status = 'completed'}) {
  return Invoice(
    id: id,
    orderId: 'ord-$id',
    businessId: 1,
    integrationId: 1,
    invoicingProviderId: 1,
    invoiceNumber: 'INV-$id',
    status: status,
    totalAmount: 1000.0,
    subtotal: 840.0,
    tax: 160.0,
    discount: 0,
    currency: 'COP',
    customerName: 'Test',
    createdAt: '',
    updatedAt: '',
  );
}

InvoicingConfig _makeConfig({int id = 1}) {
  return InvoicingConfig(
    id: id,
    businessId: 1,
    integrationIds: [1],
    enabled: true,
    autoInvoice: false,
    createdAt: '',
    updatedAt: '',
  );
}

InvoicingProvider _createProvider(MockInvoicingRepository mockRepo) {
  final apiClient = ApiClient();
  final useCases = InvoicingUseCases(mockRepo);
  return InvoicingProvider(apiClient: apiClient, useCases: useCases);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockInvoicingRepository mockRepo;
  late InvoicingProvider provider;

  setUp(() {
    mockRepo = MockInvoicingRepository();
    provider = _createProvider(mockRepo);
  });

  group('Initial state', () {
    test('has empty invoices list', () {
      expect(provider.invoices, isEmpty);
    });

    test('has empty configs list', () {
      expect(provider.configs, isEmpty);
    });

    test('has null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('is not loading', () {
      expect(provider.isLoading, false);
    });

    test('has no error', () {
      expect(provider.error, isNull);
    });

    test('has default page 1', () {
      expect(provider.page, 1);
    });
  });

  group('fetchInvoices', () {
    test('updates invoices and pagination on success', () async {
      mockRepo.getInvoicesResult = PaginatedResponse<Invoice>(
        data: [_makeInvoice(id: 1), _makeInvoice(id: 2)],
        pagination: _makePagination(total: 2),
      );

      await provider.fetchInvoices();

      expect(provider.invoices.length, 2);
      expect(provider.invoices[0].id, 1);
      expect(provider.invoices[1].id, 2);
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets isLoading during fetch and clears after', () async {
      final loadingStates = <bool>[];
      provider.addListener(() {
        loadingStates.add(provider.isLoading);
      });

      mockRepo.getInvoicesResult = PaginatedResponse<Invoice>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchInvoices();

      expect(loadingStates, [true, false]);
    });

    test('clears previous error before fetching', () async {
      mockRepo.errorToThrow = Exception('first error');
      await provider.fetchInvoices();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getInvoicesResult = PaginatedResponse<Invoice>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchInvoices();

      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch failed');

      await provider.fetchInvoices();

      expect(provider.error, contains('fetch failed'));
      expect(provider.isLoading, false);
    });
  });

  group('fetchConfigs', () {
    test('updates configs on success', () async {
      mockRepo.getConfigsResult = PaginatedResponse<InvoicingConfig>(
        data: [_makeConfig(id: 1), _makeConfig(id: 2)],
        pagination: _makePagination(),
      );

      await provider.fetchConfigs();

      expect(provider.configs.length, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('configs failed');

      await provider.fetchConfigs();

      expect(provider.error, contains('configs failed'));
      expect(provider.isLoading, false);
    });
  });

  group('createInvoice', () {
    test('returns invoice on success', () async {
      mockRepo.createInvoiceResult = _makeInvoice(id: 99);

      final dto = CreateInvoiceDTO(
        orderId: 'ord-1',
        businessId: 1,
        integrationId: 2,
      );
      final result = await provider.createInvoice(dto);

      expect(result, isNotNull);
      expect(result!.id, 99);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('create failed');

      final dto = CreateInvoiceDTO(
        orderId: 'ord-1',
        businessId: 1,
        integrationId: 2,
      );
      final result = await provider.createInvoice(dto);

      expect(result, isNull);
      expect(provider.error, contains('create failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('create error');

      await provider.createInvoice(
        CreateInvoiceDTO(orderId: 'o', businessId: 1, integrationId: 1),
      );

      expect(notified, true);
    });
  });

  group('cancelInvoice', () {
    test('returns true on success', () async {
      mockRepo.cancelInvoiceResult =
          _makeInvoice(id: 5, status: 'cancelled');

      final result = await provider.cancelInvoice(5);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('cancel failed');

      final result = await provider.cancelInvoice(5);

      expect(result, false);
      expect(provider.error, contains('cancel failed'));
    });
  });

  group('retryInvoice', () {
    test('returns true on success', () async {
      mockRepo.retryInvoiceResult = _makeInvoice(id: 5, status: 'pending');

      final result = await provider.retryInvoice(5);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('retry failed');

      final result = await provider.retryInvoice(5);

      expect(result, false);
      expect(provider.error, contains('retry failed'));
    });
  });

  group('bulkCreateInvoices', () {
    test('returns true on success', () async {
      final dto = BulkCreateInvoicesDTO(orderIds: ['o1', 'o2']);

      final result = await provider.bulkCreateInvoices(dto);

      expect(result, true);
      expect(mockRepo.calls, ['bulkCreateInvoices']);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('bulk failed');

      final result = await provider.bulkCreateInvoices(
        BulkCreateInvoicesDTO(orderIds: ['o1']),
      );

      expect(result, false);
      expect(provider.error, contains('bulk failed'));
    });
  });

  group('setPage', () {
    test('updates page', () {
      provider.setPage(3);
      expect(provider.page, 3);
    });
  });

  group('setFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(status: 'pending');
      expect(provider.page, 1);
    });
  });

  group('resetFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.resetFilters();
      expect(provider.page, 1);
    });
  });
}
