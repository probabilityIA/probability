import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/wallet_repository.dart';

class WalletProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  Wallet? _wallet;
  List<Wallet> _wallets = [];
  List<dynamic> _history = [];
  List<dynamic> _pendingRequests = [];
  List<dynamic> _processedRequests = [];
  BusinessSubscription? _subscription;
  bool _isLoading = false;
  String? _error;

  WalletProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  Wallet? get wallet => _wallet;
  List<Wallet> get wallets => _wallets;
  List<dynamic> get history => _history;
  List<dynamic> get pendingRequests => _pendingRequests;
  List<dynamic> get processedRequests => _processedRequests;
  BusinessSubscription? get subscription => _subscription;
  bool get isLoading => _isLoading;
  String? get error => _error;

  WalletUseCases get _useCases =>
      WalletUseCases(WalletApiRepository(_apiClient));

  Future<void> fetchWallets() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _wallets = await _useCases.getWallets();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchBalance({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _wallet = await _useCases.getWalletBalance(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> rechargeWallet({required double amount, int? businessId}) async {
    try {
      await _useCases.rechargeWallet(amount: amount, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<void> fetchHistory({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _history = await _useCases.getWalletHistory(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchPendingRequests() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _pendingRequests = await _useCases.getPendingRequests();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<void> fetchProcessedRequests() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _processedRequests = await _useCases.getProcessedRequests();
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> processRequest({required String id, required String action}) async {
    try {
      await _useCases.processRequest(id: id, action: action);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> manualDebit({required int businessId, required double amount, required String reference}) async {
    try {
      await _useCases.manualDebit(businessId: businessId, amount: amount, reference: reference);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> adminAdjustBalance({required int businessId, required double amount, required String reference}) async {
    try {
      await _useCases.adminAdjustBalance(businessId: businessId, amount: amount, reference: reference);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> clearRechargeHistory({required int businessId}) async {
    try {
      await _useCases.clearRechargeHistory(businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> debitForGuide({required double amount, required String trackingNumber, int? businessId}) async {
    try {
      await _useCases.debitForGuide(amount: amount, trackingNumber: trackingNumber, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<void> fetchSubscription({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _subscription = await _useCases.getMySubscription(businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<bool> registerSubscriptionPayment({
    required int businessId,
    required double amount,
    required int monthsToAdd,
    String? paymentReference,
    String? notes,
  }) async {
    try {
      await _useCases.registerSubscriptionPayment(
        businessId: businessId,
        amount: amount,
        monthsToAdd: monthsToAdd,
        paymentReference: paymentReference,
        notes: notes,
      );
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> disableSubscription({required int businessId}) async {
    try {
      await _useCases.disableSubscription(businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}
