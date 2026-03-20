import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/wallet/app/use_cases.dart';
import 'package:mobile_central/services/modules/wallet/domain/entities.dart';
import 'package:mobile_central/services/modules/wallet/domain/ports.dart';

// --- Manual Mock ---

class MockWalletRepository implements IWalletRepository {
  final List<String> calls = [];

  List<Wallet>? getWalletsResult;
  Wallet? getWalletBalanceResult;
  List<dynamic>? getWalletHistoryResult;
  List<dynamic>? getPendingRequestsResult;
  List<dynamic>? getProcessedRequestsResult;
  BusinessSubscription? getMySubscriptionResult;

  Exception? errorToThrow;

  int? capturedBusinessId;
  double? capturedAmount;
  String? capturedTrackingNumber;
  String? capturedId;
  String? capturedAction;
  String? capturedReference;
  int? capturedMonthsToAdd;
  String? capturedPaymentReference;
  String? capturedNotes;

  @override
  Future<List<Wallet>> getWallets() async {
    calls.add('getWallets');
    if (errorToThrow != null) throw errorToThrow!;
    return getWalletsResult!;
  }

  @override
  Future<Wallet> getWalletBalance({int? businessId}) async {
    calls.add('getWalletBalance');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getWalletBalanceResult!;
  }

  @override
  Future<void> rechargeWallet({required double amount, int? businessId}) async {
    calls.add('rechargeWallet');
    capturedAmount = amount;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<dynamic>> getWalletHistory({int? businessId}) async {
    calls.add('getWalletHistory');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getWalletHistoryResult!;
  }

  @override
  Future<List<dynamic>> getPendingRequests() async {
    calls.add('getPendingRequests');
    if (errorToThrow != null) throw errorToThrow!;
    return getPendingRequestsResult!;
  }

  @override
  Future<List<dynamic>> getProcessedRequests() async {
    calls.add('getProcessedRequests');
    if (errorToThrow != null) throw errorToThrow!;
    return getProcessedRequestsResult!;
  }

  @override
  Future<void> processRequest({required String id, required String action}) async {
    calls.add('processRequest');
    capturedId = id;
    capturedAction = action;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> manualDebit({required int businessId, required double amount, required String reference}) async {
    calls.add('manualDebit');
    capturedBusinessId = businessId;
    capturedAmount = amount;
    capturedReference = reference;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> adminAdjustBalance({required int businessId, required double amount, required String reference}) async {
    calls.add('adminAdjustBalance');
    capturedBusinessId = businessId;
    capturedAmount = amount;
    capturedReference = reference;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> clearRechargeHistory({required int businessId}) async {
    calls.add('clearRechargeHistory');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> debitForGuide({required double amount, required String trackingNumber, int? businessId}) async {
    calls.add('debitForGuide');
    capturedAmount = amount;
    capturedTrackingNumber = trackingNumber;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<BusinessSubscription?> getMySubscription({int? businessId}) async {
    calls.add('getMySubscription');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getMySubscriptionResult;
  }

  @override
  Future<void> registerSubscriptionPayment({
    required int businessId,
    required double amount,
    required int monthsToAdd,
    String? paymentReference,
    String? notes,
  }) async {
    calls.add('registerSubscriptionPayment');
    capturedBusinessId = businessId;
    capturedAmount = amount;
    capturedMonthsToAdd = monthsToAdd;
    capturedPaymentReference = paymentReference;
    capturedNotes = notes;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> disableSubscription({required int businessId}) async {
    calls.add('disableSubscription');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Wallet _makeWallet({String id = '1', double balance = 50000.0}) {
  return Wallet(id: id, businessId: 1, balance: balance);
}

// --- Tests ---

void main() {
  late MockWalletRepository mockRepo;
  late WalletUseCases useCases;

  setUp(() {
    mockRepo = MockWalletRepository();
    useCases = WalletUseCases(mockRepo);
  });

  group('getWallets', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getWalletsResult = [_makeWallet(), _makeWallet(id: '2')];

      final result = await useCases.getWallets();

      expect(result.length, 2);
      expect(mockRepo.calls, ['getWallets']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getWallets(), throwsException);
    });
  });

  group('getWalletBalance', () {
    test('delegates to repository with businessId', () async {
      mockRepo.getWalletBalanceResult = _makeWallet(balance: 75000.0);

      final result = await useCases.getWalletBalance(businessId: 5);

      expect(result.balance, 75000.0);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getWalletBalance']);
    });

    test('delegates without businessId', () async {
      mockRepo.getWalletBalanceResult = _makeWallet();

      await useCases.getWalletBalance();

      expect(mockRepo.capturedBusinessId, isNull);
    });
  });

  group('rechargeWallet', () {
    test('delegates to repository with correct amount', () async {
      await useCases.rechargeWallet(amount: 10000.0, businessId: 3);

      expect(mockRepo.capturedAmount, 10000.0);
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['rechargeWallet']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Recharge failed');

      expect(
        () => useCases.rechargeWallet(amount: 100.0),
        throwsException,
      );
    });
  });

  group('getWalletHistory', () {
    test('delegates to repository', () async {
      mockRepo.getWalletHistoryResult = [
        {'id': 1, 'amount': 100},
        {'id': 2, 'amount': 200},
      ];

      final result = await useCases.getWalletHistory(businessId: 2);

      expect(result.length, 2);
      expect(mockRepo.capturedBusinessId, 2);
      expect(mockRepo.calls, ['getWalletHistory']);
    });
  });

  group('getPendingRequests', () {
    test('delegates to repository', () async {
      mockRepo.getPendingRequestsResult = [{'id': 'req-1'}];

      final result = await useCases.getPendingRequests();

      expect(result.length, 1);
      expect(mockRepo.calls, ['getPendingRequests']);
    });
  });

  group('getProcessedRequests', () {
    test('delegates to repository', () async {
      mockRepo.getProcessedRequestsResult = [];

      final result = await useCases.getProcessedRequests();

      expect(result, isEmpty);
      expect(mockRepo.calls, ['getProcessedRequests']);
    });
  });

  group('processRequest', () {
    test('delegates to repository with correct id and action', () async {
      await useCases.processRequest(id: 'req-1', action: 'approve');

      expect(mockRepo.capturedId, 'req-1');
      expect(mockRepo.capturedAction, 'approve');
      expect(mockRepo.calls, ['processRequest']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Process failed');

      expect(
        () => useCases.processRequest(id: 'x', action: 'approve'),
        throwsException,
      );
    });
  });

  group('manualDebit', () {
    test('delegates to repository with correct data', () async {
      await useCases.manualDebit(
        businessId: 1,
        amount: 5000.0,
        reference: 'REF-001',
      );

      expect(mockRepo.capturedBusinessId, 1);
      expect(mockRepo.capturedAmount, 5000.0);
      expect(mockRepo.capturedReference, 'REF-001');
      expect(mockRepo.calls, ['manualDebit']);
    });
  });

  group('adminAdjustBalance', () {
    test('delegates to repository with correct data', () async {
      await useCases.adminAdjustBalance(
        businessId: 2,
        amount: -3000.0,
        reference: 'ADJ-001',
      );

      expect(mockRepo.capturedBusinessId, 2);
      expect(mockRepo.capturedAmount, -3000.0);
      expect(mockRepo.capturedReference, 'ADJ-001');
      expect(mockRepo.calls, ['adminAdjustBalance']);
    });
  });

  group('clearRechargeHistory', () {
    test('delegates to repository with correct businessId', () async {
      await useCases.clearRechargeHistory(businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['clearRechargeHistory']);
    });
  });

  group('debitForGuide', () {
    test('delegates to repository with correct data', () async {
      await useCases.debitForGuide(
        amount: 8000.0,
        trackingNumber: 'TRACK-001',
        businessId: 4,
      );

      expect(mockRepo.capturedAmount, 8000.0);
      expect(mockRepo.capturedTrackingNumber, 'TRACK-001');
      expect(mockRepo.capturedBusinessId, 4);
      expect(mockRepo.calls, ['debitForGuide']);
    });
  });

  group('getMySubscription', () {
    test('delegates to repository and returns subscription', () async {
      mockRepo.getMySubscriptionResult = BusinessSubscription(
        id: 1,
        businessId: 5,
        amount: 99000.0,
        status: 'active',
      );

      final result = await useCases.getMySubscription(businessId: 5);

      expect(result, isNotNull);
      expect(result!.status, 'active');
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getMySubscription']);
    });

    test('returns null when no subscription exists', () async {
      mockRepo.getMySubscriptionResult = null;

      final result = await useCases.getMySubscription();

      expect(result, isNull);
    });
  });

  group('registerSubscriptionPayment', () {
    test('delegates to repository with all parameters', () async {
      await useCases.registerSubscriptionPayment(
        businessId: 1,
        amount: 99000.0,
        monthsToAdd: 12,
        paymentReference: 'PAY-001',
        notes: 'Annual subscription',
      );

      expect(mockRepo.capturedBusinessId, 1);
      expect(mockRepo.capturedAmount, 99000.0);
      expect(mockRepo.capturedMonthsToAdd, 12);
      expect(mockRepo.capturedPaymentReference, 'PAY-001');
      expect(mockRepo.capturedNotes, 'Annual subscription');
      expect(mockRepo.calls, ['registerSubscriptionPayment']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Payment failed');

      expect(
        () => useCases.registerSubscriptionPayment(
          businessId: 1,
          amount: 100.0,
          monthsToAdd: 1,
        ),
        throwsException,
      );
    });
  });

  group('disableSubscription', () {
    test('delegates to repository with correct businessId', () async {
      await useCases.disableSubscription(businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
      expect(mockRepo.calls, ['disableSubscription']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Disable failed');

      expect(
        () => useCases.disableSubscription(businessId: 1),
        throwsException,
      );
    });
  });
}
