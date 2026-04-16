import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/domain/ports.dart';

// --- Manual Mock Repository ---

class MockFulfillmentStatusRepository implements IFulfillmentStatusRepository {
  List<FulfillmentStatusInfo>? getFulfillmentStatusesResult;
  Exception? errorToThrow;

  final List<String> calls = [];

  @override
  Future<List<FulfillmentStatusInfo>> getFulfillmentStatuses() async {
    calls.add('getFulfillmentStatuses');
    if (errorToThrow != null) throw errorToThrow!;
    return getFulfillmentStatusesResult!;
  }
}

// --- Testable Provider ---

class TestableFulfillmentStatusProvider {
  final FulfillmentStatusUseCases _useCases;

  List<FulfillmentStatusInfo> _statuses = [];
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableFulfillmentStatusProvider(this._useCases);

  List<FulfillmentStatusInfo> get statuses => _statuses;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchStatuses() async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _statuses = await _useCases.getFulfillmentStatuses();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }
}

// --- Helpers ---

FulfillmentStatusInfo _makeStatus({int id = 1, String code = 'pending', String name = 'Pending'}) {
  return FulfillmentStatusInfo(id: id, code: code, name: name);
}

// --- Tests ---

void main() {
  late MockFulfillmentStatusRepository mockRepo;
  late FulfillmentStatusUseCases useCases;
  late TestableFulfillmentStatusProvider provider;

  setUp(() {
    mockRepo = MockFulfillmentStatusRepository();
    useCases = FulfillmentStatusUseCases(mockRepo);
    provider = TestableFulfillmentStatusProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty statuses list', () {
      expect(provider.statuses, isEmpty);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchStatuses', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getFulfillmentStatusesResult = [_makeStatus()];

      await provider.fetchStatuses();

      expect(provider.notifications.length, 2);
    });

    test('populates statuses on success', () async {
      mockRepo.getFulfillmentStatusesResult = [
        _makeStatus(id: 1, code: 'pending', name: 'Pending'),
        _makeStatus(id: 2, code: 'shipped', name: 'Shipped'),
        _makeStatus(id: 3, code: 'delivered', name: 'Delivered'),
      ];

      await provider.fetchStatuses();

      expect(provider.statuses.length, 3);
      expect(provider.statuses[0].code, 'pending');
      expect(provider.statuses[1].code, 'shipped');
      expect(provider.statuses[2].code, 'delivered');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchStatuses();

      expect(provider.error, contains('Server error'));
      expect(provider.statuses, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchStatuses();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getFulfillmentStatusesResult = [];
      await provider.fetchStatuses();

      expect(provider.error, isNull);
    });

    test('handles empty list from repository', () async {
      mockRepo.getFulfillmentStatusesResult = [];

      await provider.fetchStatuses();

      expect(provider.statuses, isEmpty);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('replaces previous statuses on subsequent fetch', () async {
      mockRepo.getFulfillmentStatusesResult = [
        _makeStatus(id: 1, code: 'pending', name: 'Pending'),
      ];
      await provider.fetchStatuses();
      expect(provider.statuses.length, 1);

      mockRepo.getFulfillmentStatusesResult = [
        _makeStatus(id: 1, code: 'pending', name: 'Pending'),
        _makeStatus(id: 2, code: 'shipped', name: 'Shipped'),
      ];
      await provider.fetchStatuses();

      expect(provider.statuses.length, 2);
    });
  });
}
