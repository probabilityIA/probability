import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/paymentstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/paymentstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/paymentstatus/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mock for IPaymentStatusRepository
// ---------------------------------------------------------------------------
class MockPaymentStatusRepository implements IPaymentStatusRepository {
  List<PaymentStatusInfo>? getPaymentStatusesResult;
  Exception? errorToThrow;
  bool? lastIsActive;

  @override
  Future<List<PaymentStatusInfo>> getPaymentStatuses({bool? isActive}) async {
    lastIsActive = isActive;
    if (errorToThrow != null) throw errorToThrow!;
    return getPaymentStatusesResult ?? [];
  }
}

// ---------------------------------------------------------------------------
// Testable provider that mirrors PaymentStatusProvider logic
// ---------------------------------------------------------------------------
class TestablePaymentStatusProvider {
  final IPaymentStatusRepository _repository;

  List<PaymentStatusInfo> _statuses = [];
  bool _isLoading = false;
  String? _error;

  int _notifyCount = 0;

  TestablePaymentStatusProvider(this._repository);

  List<PaymentStatusInfo> get statuses => _statuses;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get notifyCount => _notifyCount;

  PaymentStatusUseCases get _useCases =>
      PaymentStatusUseCases(_repository);

  void _notifyListeners() {
    _notifyCount++;
  }

  Future<void> fetchStatuses({bool? isActive}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();
    try {
      _statuses =
          await _useCases.getPaymentStatuses(isActive: isActive);
    } catch (e) {
      _error = e.toString();
    }
    _isLoading = false;
    _notifyListeners();
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockPaymentStatusRepository mockRepo;
  late TestablePaymentStatusProvider provider;

  setUp(() {
    mockRepo = MockPaymentStatusRepository();
    provider = TestablePaymentStatusProvider(mockRepo);
  });

  group('initial state', () {
    test('has empty statuses list', () {
      expect(provider.statuses, isEmpty);
    });

    test('isLoading is false', () {
      expect(provider.isLoading, false);
    });

    test('error is null', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchStatuses', () {
    test('updates statuses list on success', () async {
      mockRepo.getPaymentStatusesResult = [
        PaymentStatusInfo(
          id: 1,
          code: 'paid',
          name: 'Paid',
          isActive: true,
        ),
        PaymentStatusInfo(
          id: 2,
          code: 'pending',
          name: 'Pending',
          isActive: true,
        ),
      ];

      await provider.fetchStatuses();

      expect(provider.statuses, hasLength(2));
      expect(provider.statuses[0].code, 'paid');
      expect(provider.statuses[1].code, 'pending');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('server down');

      await provider.fetchStatuses();

      expect(provider.error, contains('server down'));
      expect(provider.statuses, isEmpty);
      expect(provider.isLoading, false);
    });

    test('notifies listeners twice (loading start and end)', () async {
      await provider.fetchStatuses();

      expect(provider.notifyCount, 2);
    });

    test('passes isActive when provided', () async {
      await provider.fetchStatuses(isActive: true);

      expect(mockRepo.lastIsActive, true);
    });

    test('does not pass isActive when not provided', () async {
      await provider.fetchStatuses();

      expect(mockRepo.lastIsActive, isNull);
    });

    test('passes false isActive', () async {
      await provider.fetchStatuses(isActive: false);

      expect(mockRepo.lastIsActive, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('fail');
      await provider.fetchStatuses();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      await provider.fetchStatuses();

      expect(provider.error, isNull);
    });

    test('handles empty response', () async {
      mockRepo.getPaymentStatusesResult = [];

      await provider.fetchStatuses();

      expect(provider.statuses, isEmpty);
      expect(provider.error, isNull);
    });
  });

  group('loading states', () {
    test('isLoading is false before fetch', () {
      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch completes', () async {
      await provider.fetchStatuses();

      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch fails', () async {
      mockRepo.errorToThrow = Exception('fail');

      await provider.fetchStatuses();

      expect(provider.isLoading, false);
    });
  });
}
