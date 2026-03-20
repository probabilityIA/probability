import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/pay/app/use_cases.dart';
import 'package:mobile_central/services/modules/pay/domain/entities.dart';
import 'package:mobile_central/services/modules/pay/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mock for IPayGatewayRepository
// ---------------------------------------------------------------------------
class MockPayGatewayRepository implements IPayGatewayRepository {
  // Configurable return values / errors
  PaymentGatewayTypesResponse? listResult;
  Exception? errorToThrow;
  int callCount = 0;

  @override
  Future<PaymentGatewayTypesResponse> listPaymentGatewayTypes() async {
    callCount++;
    if (errorToThrow != null) throw errorToThrow!;
    return listResult ??
        PaymentGatewayTypesResponse(
          success: true,
          data: [],
        );
  }
}

void main() {
  late MockPayGatewayRepository mockRepo;
  late PayUseCases useCases;

  setUp(() {
    mockRepo = MockPayGatewayRepository();
    useCases = PayUseCases(mockRepo);
  });

  group('listPaymentGatewayTypes', () {
    test('delegates to repository', () async {
      await useCases.listPaymentGatewayTypes();

      expect(mockRepo.callCount, 1);
    });

    test('returns the response from repository', () async {
      final expectedGateways = [
        PaymentGatewayType(
          id: 1,
          name: 'Stripe',
          code: 'stripe',
          isActive: true,
          inDevelopment: false,
        ),
        PaymentGatewayType(
          id: 2,
          name: 'PayU',
          code: 'payu',
          isActive: true,
          inDevelopment: true,
        ),
      ];
      mockRepo.listResult = PaymentGatewayTypesResponse(
        success: true,
        data: expectedGateways,
        message: 'OK',
      );

      final result = await useCases.listPaymentGatewayTypes();

      expect(result.success, true);
      expect(result.data, hasLength(2));
      expect(result.data[0].name, 'Stripe');
      expect(result.data[1].name, 'PayU');
      expect(result.message, 'OK');
    });

    test('returns empty list when repository returns empty', () async {
      mockRepo.listResult = PaymentGatewayTypesResponse(
        success: true,
        data: [],
      );

      final result = await useCases.listPaymentGatewayTypes();

      expect(result.success, true);
      expect(result.data, isEmpty);
    });
  });

  group('error propagation', () {
    test('listPaymentGatewayTypes propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.listPaymentGatewayTypes(),
        throwsA(isA<Exception>()),
      );
    });

    test('listPaymentGatewayTypes propagates specific error message', () async {
      mockRepo.errorToThrow = Exception('unauthorized');

      try {
        await useCases.listPaymentGatewayTypes();
        fail('Should have thrown');
      } catch (e) {
        expect(e.toString(), contains('unauthorized'));
      }
    });
  });
}
