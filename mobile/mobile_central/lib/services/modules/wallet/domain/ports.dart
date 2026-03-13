import 'entities.dart';

abstract class IWalletRepository {
  Future<List<Wallet>> getWallets();
  Future<Wallet> getWalletBalance({int? businessId});
  Future<void> rechargeWallet({required double amount, int? businessId});
  Future<List<dynamic>> getWalletHistory({int? businessId});
  Future<List<dynamic>> getPendingRequests();
  Future<List<dynamic>> getProcessedRequests();
  Future<void> processRequest({required String id, required String action});
  Future<void> manualDebit({required int businessId, required double amount, required String reference});
  Future<void> adminAdjustBalance({required int businessId, required double amount, required String reference});
  Future<void> clearRechargeHistory({required int businessId});
  Future<void> debitForGuide({required double amount, required String trackingNumber, int? businessId});
  Future<BusinessSubscription?> getMySubscription({int? businessId});
  Future<void> registerSubscriptionPayment({
    required int businessId,
    required double amount,
    required int monthsToAdd,
    String? paymentReference,
    String? notes,
  });
  Future<void> disableSubscription({required int businessId});
}
