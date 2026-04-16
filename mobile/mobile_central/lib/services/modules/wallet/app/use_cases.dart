import '../domain/entities.dart';
import '../domain/ports.dart';

class WalletUseCases {
  final IWalletRepository _repository;

  WalletUseCases(this._repository);

  Future<List<Wallet>> getWallets() {
    return _repository.getWallets();
  }

  Future<Wallet> getWalletBalance({int? businessId}) {
    return _repository.getWalletBalance(businessId: businessId);
  }

  Future<void> rechargeWallet({required double amount, int? businessId}) {
    return _repository.rechargeWallet(amount: amount, businessId: businessId);
  }

  Future<List<dynamic>> getWalletHistory({int? businessId}) {
    return _repository.getWalletHistory(businessId: businessId);
  }

  Future<List<dynamic>> getPendingRequests() {
    return _repository.getPendingRequests();
  }

  Future<List<dynamic>> getProcessedRequests() {
    return _repository.getProcessedRequests();
  }

  Future<void> processRequest({required String id, required String action}) {
    return _repository.processRequest(id: id, action: action);
  }

  Future<void> manualDebit({required int businessId, required double amount, required String reference}) {
    return _repository.manualDebit(businessId: businessId, amount: amount, reference: reference);
  }

  Future<void> adminAdjustBalance({required int businessId, required double amount, required String reference}) {
    return _repository.adminAdjustBalance(businessId: businessId, amount: amount, reference: reference);
  }

  Future<void> clearRechargeHistory({required int businessId}) {
    return _repository.clearRechargeHistory(businessId: businessId);
  }

  Future<void> debitForGuide({required double amount, required String trackingNumber, int? businessId}) {
    return _repository.debitForGuide(amount: amount, trackingNumber: trackingNumber, businessId: businessId);
  }

  Future<BusinessSubscription?> getMySubscription({int? businessId}) {
    return _repository.getMySubscription(businessId: businessId);
  }

  Future<void> registerSubscriptionPayment({
    required int businessId,
    required double amount,
    required int monthsToAdd,
    String? paymentReference,
    String? notes,
  }) {
    return _repository.registerSubscriptionPayment(
      businessId: businessId,
      amount: amount,
      monthsToAdd: monthsToAdd,
      paymentReference: paymentReference,
      notes: notes,
    );
  }

  Future<void> disableSubscription({required int businessId}) {
    return _repository.disableSubscription(businessId: businessId);
  }
}
