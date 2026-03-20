import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/pay/app/use_cases.dart';
import 'package:mobile_central/services/modules/pay/domain/entities.dart';
import 'package:mobile_central/services/modules/pay/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mock for IPayGatewayRepository
// ---------------------------------------------------------------------------
class MockPayGatewayRepository implements IPayGatewayRepository {
  PaymentGatewayTypesResponse? listResult;
  Exception? errorToThrow;

  @override
  Future<PaymentGatewayTypesResponse> listPaymentGatewayTypes() async {
    if (errorToThrow != null) throw errorToThrow!;
    return listResult ??
        PaymentGatewayTypesResponse(
          success: true,
          data: [],
        );
  }
}

// ---------------------------------------------------------------------------
// Testable provider that mirrors PayProvider logic but accepts a repository
// ---------------------------------------------------------------------------
class TestablePayProvider {
  final IPayGatewayRepository _repository;

  List<PaymentGatewayType> _paymentGatewayTypes = [];
  bool _isLoading = false;
  String? _error;

  int _notifyCount = 0;

  TestablePayProvider(this._repository);

  List<PaymentGatewayType> get paymentGatewayTypes => _paymentGatewayTypes;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get notifyCount => _notifyCount;

  PayUseCases get _useCases => PayUseCases(_repository);

  void _notifyListeners() {
    _notifyCount++;
  }

  Future<void> fetchPaymentGatewayTypes() async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.listPaymentGatewayTypes();
      _paymentGatewayTypes = response.data;
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
  late MockPayGatewayRepository mockRepo;
  late TestablePayProvider provider;

  setUp(() {
    mockRepo = MockPayGatewayRepository();
    provider = TestablePayProvider(mockRepo);
  });

  group('initial state', () {
    test('has empty paymentGatewayTypes list', () {
      expect(provider.paymentGatewayTypes, isEmpty);
    });

    test('isLoading is false', () {
      expect(provider.isLoading, false);
    });

    test('error is null', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchPaymentGatewayTypes', () {
    test('updates paymentGatewayTypes list on success', () async {
      mockRepo.listResult = PaymentGatewayTypesResponse(
        success: true,
        data: [
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
        ],
      );

      await provider.fetchPaymentGatewayTypes();

      expect(provider.paymentGatewayTypes, hasLength(2));
      expect(provider.paymentGatewayTypes[0].name, 'Stripe');
      expect(provider.paymentGatewayTypes[1].name, 'PayU');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('server down');

      await provider.fetchPaymentGatewayTypes();

      expect(provider.error, contains('server down'));
      expect(provider.paymentGatewayTypes, isEmpty);
      expect(provider.isLoading, false);
    });

    test('notifies listeners twice (loading start and end)', () async {
      await provider.fetchPaymentGatewayTypes();

      expect(provider.notifyCount, 2);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('fail');
      await provider.fetchPaymentGatewayTypes();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      await provider.fetchPaymentGatewayTypes();

      expect(provider.error, isNull);
    });

    test('handles empty response', () async {
      mockRepo.listResult = PaymentGatewayTypesResponse(
        success: true,
        data: [],
      );

      await provider.fetchPaymentGatewayTypes();

      expect(provider.paymentGatewayTypes, isEmpty);
      expect(provider.error, isNull);
    });
  });

  group('loading states', () {
    test('isLoading is false before fetch', () {
      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch completes', () async {
      await provider.fetchPaymentGatewayTypes();

      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch fails', () async {
      mockRepo.errorToThrow = Exception('fail');

      await provider.fetchPaymentGatewayTypes();

      expect(provider.isLoading, false);
    });
  });
}
