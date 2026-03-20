import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/invoicing/app/use_cases.dart';
import 'package:mobile_central/services/modules/invoicing/domain/entities.dart';
import 'package:mobile_central/services/modules/invoicing/domain/ports.dart';
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

  // Captured args
  InvoiceFilters? lastInvoiceFilters;
  int? lastId;
  CreateInvoiceDTO? lastCreateInvoiceDTO;
  CreateCreditNoteDTO? lastCreateCreditNoteDTO;
  ConfigFilters? lastConfigFilters;
  CreateConfigDTO? lastCreateConfigDTO;
  UpdateConfigDTO? lastUpdateConfigDTO;
  BulkCreateInvoicesDTO? lastBulkDTO;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<Invoice>> getInvoices(InvoiceFilters filters) async {
    lastInvoiceFilters = filters;
    _trackCall('getInvoices');
    return getInvoicesResult!;
  }

  @override
  Future<Invoice> getInvoiceById(int id) async {
    lastId = id;
    _trackCall('getInvoiceById');
    return getInvoiceByIdResult!;
  }

  @override
  Future<Invoice> createInvoice(CreateInvoiceDTO data) async {
    lastCreateInvoiceDTO = data;
    _trackCall('createInvoice');
    return createInvoiceResult!;
  }

  @override
  Future<Invoice> cancelInvoice(int id) async {
    lastId = id;
    _trackCall('cancelInvoice');
    return cancelInvoiceResult!;
  }

  @override
  Future<Invoice> retryInvoice(int id) async {
    lastId = id;
    _trackCall('retryInvoice');
    return retryInvoiceResult!;
  }

  @override
  Future<CreditNote> createCreditNote(CreateCreditNoteDTO data) async {
    lastCreateCreditNoteDTO = data;
    _trackCall('createCreditNote');
    return createCreditNoteResult!;
  }

  @override
  Future<PaginatedResponse<InvoicingConfig>> getConfigs(
      ConfigFilters filters) async {
    lastConfigFilters = filters;
    _trackCall('getConfigs');
    return getConfigsResult!;
  }

  @override
  Future<InvoicingConfig> getConfigById(int id) async {
    lastId = id;
    _trackCall('getConfigById');
    return getConfigByIdResult!;
  }

  @override
  Future<InvoicingConfig> createConfig(CreateConfigDTO data) async {
    lastCreateConfigDTO = data;
    _trackCall('createConfig');
    return createConfigResult!;
  }

  @override
  Future<InvoicingConfig> updateConfig(int id, UpdateConfigDTO data) async {
    lastId = id;
    lastUpdateConfigDTO = data;
    _trackCall('updateConfig');
    return updateConfigResult!;
  }

  @override
  Future<void> deleteConfig(int id) async {
    lastId = id;
    _trackCall('deleteConfig');
  }

  @override
  Future<void> bulkCreateInvoices(BulkCreateInvoicesDTO data) async {
    lastBulkDTO = data;
    _trackCall('bulkCreateInvoices');
  }

  @override
  Future<List<SyncLog>> getSyncLogs(int invoiceId) async {
    lastId = invoiceId;
    _trackCall('getSyncLogs');
    return getSyncLogsResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
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
    subtotal: 840.34,
    tax: 159.66,
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

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockInvoicingRepository mockRepo;
  late InvoicingUseCases useCases;

  setUp(() {
    mockRepo = MockInvoicingRepository();
    useCases = InvoicingUseCases(mockRepo);
  });

  group('getInvoices', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Invoice>(
        data: [_makeInvoice()],
        pagination: _makePagination(),
      );
      mockRepo.getInvoicesResult = expected;

      final filters = InvoiceFilters(page: 1, pageSize: 10);
      final result = await useCases.getInvoices(filters);

      expect(mockRepo.calls, ['getInvoices']);
      expect(mockRepo.lastInvoiceFilters, filters);
      expect(result.data.length, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getInvoices(InvoiceFilters()),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getInvoiceById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getInvoiceByIdResult = _makeInvoice(id: 42);

      final result = await useCases.getInvoiceById(42);

      expect(mockRepo.calls, ['getInvoiceById']);
      expect(mockRepo.lastId, 42);
      expect(result.id, 42);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getInvoiceById(999),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('createInvoice', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateInvoiceDTO(
        orderId: 'ord-1',
        businessId: 1,
        integrationId: 2,
      );
      mockRepo.createInvoiceResult = _makeInvoice(id: 10);

      final result = await useCases.createInvoice(dto);

      expect(mockRepo.calls, ['createInvoice']);
      expect(mockRepo.lastCreateInvoiceDTO, dto);
      expect(result.id, 10);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('create error');

      expect(
        () => useCases.createInvoice(
          CreateInvoiceDTO(orderId: 'o', businessId: 1, integrationId: 1),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('cancelInvoice', () {
    test('delegates to repository with correct id', () async {
      mockRepo.cancelInvoiceResult =
          _makeInvoice(id: 5, status: 'cancelled');

      final result = await useCases.cancelInvoice(5);

      expect(mockRepo.calls, ['cancelInvoice']);
      expect(mockRepo.lastId, 5);
      expect(result.status, 'cancelled');
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('cancel error');

      expect(
        () => useCases.cancelInvoice(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('retryInvoice', () {
    test('delegates to repository with correct id', () async {
      mockRepo.retryInvoiceResult = _makeInvoice(id: 5, status: 'pending');

      final result = await useCases.retryInvoice(5);

      expect(mockRepo.calls, ['retryInvoice']);
      expect(mockRepo.lastId, 5);
      expect(result.id, 5);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('retry error');

      expect(
        () => useCases.retryInvoice(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('createCreditNote', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateCreditNoteDTO(
        invoiceId: 10,
        amount: 500.0,
        reason: 'return',
        noteType: 'full',
      );
      mockRepo.createCreditNoteResult = CreditNote(
        id: 1,
        invoiceId: 10,
        creditNoteNumber: 'CN-001',
        amount: 500.0,
        reason: 'return',
        noteType: 'full',
        status: 'completed',
        createdAt: '',
      );

      final result = await useCases.createCreditNote(dto);

      expect(mockRepo.calls, ['createCreditNote']);
      expect(mockRepo.lastCreateCreditNoteDTO, dto);
      expect(result.id, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('credit note error');

      expect(
        () => useCases.createCreditNote(
          CreateCreditNoteDTO(
            invoiceId: 1,
            amount: 100,
            reason: 'r',
            noteType: 'full',
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getConfigs', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getConfigsResult = PaginatedResponse<InvoicingConfig>(
        data: [_makeConfig()],
        pagination: _makePagination(),
      );

      final filters = ConfigFilters(businessId: 1);
      final result = await useCases.getConfigs(filters);

      expect(mockRepo.calls, ['getConfigs']);
      expect(mockRepo.lastConfigFilters, filters);
      expect(result.data.length, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('configs error');

      expect(
        () => useCases.getConfigs(ConfigFilters()),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getConfigById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getConfigByIdResult = _makeConfig(id: 7);

      final result = await useCases.getConfigById(7);

      expect(mockRepo.calls, ['getConfigById']);
      expect(mockRepo.lastId, 7);
      expect(result.id, 7);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('config not found');

      expect(
        () => useCases.getConfigById(999),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('createConfig', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationIds: [1, 2],
        invoicingIntegrationId: 5,
      );
      mockRepo.createConfigResult = _makeConfig(id: 99);

      final result = await useCases.createConfig(dto);

      expect(mockRepo.calls, ['createConfig']);
      expect(mockRepo.lastCreateConfigDTO, dto);
      expect(result.id, 99);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('create config error');

      expect(
        () => useCases.createConfig(
          CreateConfigDTO(
            businessId: 1,
            integrationIds: [1],
            invoicingIntegrationId: 1,
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('updateConfig', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateConfigDTO(enabled: false);
      mockRepo.updateConfigResult = _makeConfig(id: 5);

      final result = await useCases.updateConfig(5, dto);

      expect(mockRepo.calls, ['updateConfig']);
      expect(mockRepo.lastId, 5);
      expect(mockRepo.lastUpdateConfigDTO, dto);
      expect(result.id, 5);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('update config error');

      expect(
        () => useCases.updateConfig(1, UpdateConfigDTO()),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('deleteConfig', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteConfig(8);

      expect(mockRepo.calls, ['deleteConfig']);
      expect(mockRepo.lastId, 8);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('delete error');

      expect(
        () => useCases.deleteConfig(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('bulkCreateInvoices', () {
    test('delegates to repository with correct DTO', () async {
      final dto = BulkCreateInvoicesDTO(orderIds: ['o1', 'o2']);

      await useCases.bulkCreateInvoices(dto);

      expect(mockRepo.calls, ['bulkCreateInvoices']);
      expect(mockRepo.lastBulkDTO, dto);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('bulk error');

      expect(
        () => useCases.bulkCreateInvoices(
          BulkCreateInvoicesDTO(orderIds: ['o1']),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getSyncLogs', () {
    test('delegates to repository with correct invoiceId', () async {
      mockRepo.getSyncLogsResult = [
        SyncLog(
          id: 1,
          invoiceId: 10,
          operationType: 'create',
          status: 'completed',
          retryCount: 0,
          maxRetries: 3,
          triggeredBy: 'auto',
          startedAt: '',
          createdAt: '',
        ),
      ];

      final result = await useCases.getSyncLogs(10);

      expect(mockRepo.calls, ['getSyncLogs']);
      expect(mockRepo.lastId, 10);
      expect(result.length, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('sync logs error');

      expect(
        () => useCases.getSyncLogs(1),
        throwsA(isA<Exception>()),
      );
    });
  });
}
