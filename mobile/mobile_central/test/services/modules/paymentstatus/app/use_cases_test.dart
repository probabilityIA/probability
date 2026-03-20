import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/paymentstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/paymentstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/paymentstatus/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mock for IPaymentStatusRepository
// ---------------------------------------------------------------------------
class MockPaymentStatusRepository implements IPaymentStatusRepository {
  // Captured arguments
  bool? lastGetPaymentStatusesIsActive;
  int callCount = 0;

  // Configurable return values / errors
  List<PaymentStatusInfo>? getPaymentStatusesResult;
  Exception? errorToThrow;

  @override
  Future<List<PaymentStatusInfo>> getPaymentStatuses({bool? isActive}) async {
    callCount++;
    lastGetPaymentStatusesIsActive = isActive;
    if (errorToThrow != null) throw errorToThrow!;
    return getPaymentStatusesResult ?? [];
  }
}

void main() {
  late MockPaymentStatusRepository mockRepo;
  late PaymentStatusUseCases useCases;

  setUp(() {
    mockRepo = MockPaymentStatusRepository();
    useCases = PaymentStatusUseCases(mockRepo);
  });

  group('getPaymentStatuses', () {
    test('delegates to repository', () async {
      await useCases.getPaymentStatuses();

      expect(mockRepo.callCount, 1);
    });

    test('delegates to repository with isActive param', () async {
      await useCases.getPaymentStatuses(isActive: true);

      expect(mockRepo.lastGetPaymentStatusesIsActive, true);
    });

    test('delegates to repository without isActive param', () async {
      await useCases.getPaymentStatuses();

      expect(mockRepo.lastGetPaymentStatusesIsActive, isNull);
    });

    test('delegates with isActive false', () async {
      await useCases.getPaymentStatuses(isActive: false);

      expect(mockRepo.lastGetPaymentStatusesIsActive, false);
    });

    test('returns the list from repository', () async {
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
        PaymentStatusInfo(
          id: 3,
          code: 'failed',
          name: 'Failed',
          isActive: false,
        ),
      ];

      final result = await useCases.getPaymentStatuses();

      expect(result, hasLength(3));
      expect(result[0].code, 'paid');
      expect(result[1].code, 'pending');
      expect(result[2].code, 'failed');
    });

    test('returns empty list when repository returns empty', () async {
      mockRepo.getPaymentStatusesResult = [];

      final result = await useCases.getPaymentStatuses();

      expect(result, isEmpty);
    });
  });

  group('error propagation', () {
    test('getPaymentStatuses propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getPaymentStatuses(),
        throwsA(isA<Exception>()),
      );
    });

    test('getPaymentStatuses propagates specific error message', () async {
      mockRepo.errorToThrow = Exception('unauthorized');

      try {
        await useCases.getPaymentStatuses();
        fail('Should have thrown');
      } catch (e) {
        expect(e.toString(), contains('unauthorized'));
      }
    });

    test('getPaymentStatuses with isActive propagates exceptions', () {
      mockRepo.errorToThrow = Exception('server error');

      expect(
        () => useCases.getPaymentStatuses(isActive: true),
        throwsA(isA<Exception>()),
      );
    });
  });
}
