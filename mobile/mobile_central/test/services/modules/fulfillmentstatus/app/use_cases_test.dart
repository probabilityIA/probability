import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/fulfillmentstatus/domain/ports.dart';

// --- Manual Mock ---

class MockFulfillmentStatusRepository implements IFulfillmentStatusRepository {
  final List<String> calls = [];

  List<FulfillmentStatusInfo>? getFulfillmentStatusesResult;
  Exception? errorToThrow;

  @override
  Future<List<FulfillmentStatusInfo>> getFulfillmentStatuses() async {
    calls.add('getFulfillmentStatuses');
    if (errorToThrow != null) throw errorToThrow!;
    return getFulfillmentStatusesResult!;
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

  setUp(() {
    mockRepo = MockFulfillmentStatusRepository();
    useCases = FulfillmentStatusUseCases(mockRepo);
  });

  group('getFulfillmentStatuses', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getFulfillmentStatusesResult = [
        _makeStatus(id: 1, code: 'pending', name: 'Pending'),
        _makeStatus(id: 2, code: 'shipped', name: 'Shipped'),
        _makeStatus(id: 3, code: 'delivered', name: 'Delivered'),
      ];

      final result = await useCases.getFulfillmentStatuses();

      expect(result.length, 3);
      expect(result[0].code, 'pending');
      expect(result[1].code, 'shipped');
      expect(result[2].code, 'delivered');
      expect(mockRepo.calls, ['getFulfillmentStatuses']);
    });

    test('returns empty list when repository returns empty', () async {
      mockRepo.getFulfillmentStatusesResult = [];

      final result = await useCases.getFulfillmentStatuses();

      expect(result, isEmpty);
      expect(mockRepo.calls, ['getFulfillmentStatuses']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getFulfillmentStatuses(), throwsException);
    });

    test('calls repository exactly once', () async {
      mockRepo.getFulfillmentStatusesResult = [_makeStatus()];

      await useCases.getFulfillmentStatuses();

      expect(mockRepo.calls.length, 1);
    });
  });
}
