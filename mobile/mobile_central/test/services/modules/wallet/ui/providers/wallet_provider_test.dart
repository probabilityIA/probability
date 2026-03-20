import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/wallet/app/use_cases.dart';
import 'package:mobile_central/services/modules/wallet/domain/entities.dart';
import 'package:mobile_central/services/modules/wallet/domain/ports.dart';

// --- Manual Mock Repository ---

class MockWalletRepository implements IWalletRepository {
  List<Wallet>? getWalletsResult;
  Wallet? getWalletBalanceResult;
  List<dynamic>? getWalletHistoryResult;
  List<dynamic>? getPendingRequestsResult;
  List<dynamic>? getProcessedRequestsResult;
  BusinessSubscription? getMySubscriptionResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedBusinessId;

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
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<dynamic>> getWalletHistory({int? businessId}) async {
    calls.add('getWalletHistory');
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
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> manualDebit({required int businessId, required double amount, required String reference}) async {
    calls.add('manualDebit');
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> adminAdjustBalance({required int businessId, required double amount, required String reference}) async {
    calls.add('adminAdjustBalance');
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> clearRechargeHistory({required int businessId}) async {
    calls.add('clearRechargeHistory');
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> debitForGuide({required double amount, required String trackingNumber, int? businessId}) async {
    calls.add('debitForGuide');
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<BusinessSubscription?> getMySubscription({int? businessId}) async {
    calls.add('getMySubscription');
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
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> disableSubscription({required int businessId}) async {
    calls.add('disableSubscription');
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableWalletProvider {
  final WalletUseCases _useCases;

  Wallet? _wallet;
  List<Wallet> _wallets = [];
  List<dynamic> _history = [];
  BusinessSubscription? _subscription;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableWalletProvider(this._useCases);

  Wallet? get wallet => _wallet;
  List<Wallet> get wallets => _wallets;
  List<dynamic> get history => _history;
  BusinessSubscription? get subscription => _subscription;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchWallets() async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _wallets = await _useCases.getWallets();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchBalance({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _wallet = await _useCases.getWalletBalance(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<bool> rechargeWallet({required double amount, int? businessId}) async {
    try {
      await _useCases.rechargeWallet(amount: amount, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<void> fetchHistory({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _history = await _useCases.getWalletHistory(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchSubscription({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _subscription = await _useCases.getMySubscription(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<bool> processRequest({required String id, required String action}) async {
    try {
      await _useCases.processRequest(id: id, action: action);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
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
  late TestableWalletProvider provider;

  setUp(() {
    mockRepo = MockWalletRepository();
    useCases = WalletUseCases(mockRepo);
    provider = TestableWalletProvider(useCases);
  });

  group('initial state', () {
    test('starts with null wallet', () {
      expect(provider.wallet, isNull);
    });

    test('starts with empty wallets list', () {
      expect(provider.wallets, isEmpty);
    });

    test('starts with empty history', () {
      expect(provider.history, isEmpty);
    });

    test('starts with null subscription', () {
      expect(provider.subscription, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchWallets', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getWalletsResult = [_makeWallet()];

      await provider.fetchWallets();

      expect(provider.notifications.length, 2);
    });

    test('populates wallets on success', () async {
      mockRepo.getWalletsResult = [
        _makeWallet(id: '1', balance: 10000.0),
        _makeWallet(id: '2', balance: 20000.0),
      ];

      await provider.fetchWallets();

      expect(provider.wallets.length, 2);
      expect(provider.wallets[0].balance, 10000.0);
      expect(provider.wallets[1].balance, 20000.0);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchWallets();

      expect(provider.error, contains('Server error'));
      expect(provider.wallets, isEmpty);
    });
  });

  group('fetchBalance', () {
    test('populates wallet on success', () async {
      mockRepo.getWalletBalanceResult = _makeWallet(balance: 75000.0);

      await provider.fetchBalance(businessId: 5);

      expect(provider.wallet, isNotNull);
      expect(provider.wallet!.balance, 75000.0);
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Balance error');

      await provider.fetchBalance();

      expect(provider.error, contains('Balance error'));
      expect(provider.wallet, isNull);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchBalance();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getWalletBalanceResult = _makeWallet();
      await provider.fetchBalance();

      expect(provider.error, isNull);
    });
  });

  group('rechargeWallet', () {
    test('returns true on success', () async {
      final result = await provider.rechargeWallet(amount: 10000.0);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Recharge failed');

      final result = await provider.rechargeWallet(amount: 100.0);

      expect(result, false);
      expect(provider.error, contains('Recharge failed'));
    });
  });

  group('fetchHistory', () {
    test('populates history on success', () async {
      mockRepo.getWalletHistoryResult = [
        {'id': 1, 'amount': 100},
        {'id': 2, 'amount': -50},
      ];

      await provider.fetchHistory();

      expect(provider.history.length, 2);
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('History error');

      await provider.fetchHistory();

      expect(provider.error, contains('History error'));
    });
  });

  group('fetchSubscription', () {
    test('populates subscription on success', () async {
      mockRepo.getMySubscriptionResult = BusinessSubscription(
        id: 1,
        businessId: 5,
        amount: 99000.0,
        status: 'active',
      );

      await provider.fetchSubscription(businessId: 5);

      expect(provider.subscription, isNotNull);
      expect(provider.subscription!.status, 'active');
    });

    test('handles null subscription', () async {
      mockRepo.getMySubscriptionResult = null;

      await provider.fetchSubscription();

      expect(provider.subscription, isNull);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Subscription error');

      await provider.fetchSubscription();

      expect(provider.error, contains('Subscription error'));
    });
  });

  group('processRequest', () {
    test('returns true on success', () async {
      final result = await provider.processRequest(id: 'req-1', action: 'approve');

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Process failed');

      final result = await provider.processRequest(id: 'req-1', action: 'approve');

      expect(result, false);
      expect(provider.error, contains('Process failed'));
    });
  });
}
