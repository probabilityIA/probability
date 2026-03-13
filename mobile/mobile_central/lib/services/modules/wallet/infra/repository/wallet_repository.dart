import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class WalletApiRepository implements IWalletRepository {
  final ApiClient _client;

  WalletApiRepository(this._client);

  @override
  Future<List<Wallet>> getWallets() async {
    final response = await _client.get('/pay/wallet/all');
    final data = response.data;
    if (data is List) {
      return data.map((e) => Wallet.fromJson(e)).toList();
    }
    return [];
  }

  @override
  Future<Wallet> getWalletBalance({int? businessId}) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['business_id'] = businessId;

    final response = await _client.get(
      '/pay/wallet/balance',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    return Wallet.fromJson(response.data);
  }

  @override
  Future<void> rechargeWallet({required double amount, int? businessId}) async {
    final body = <String, dynamic>{'amount': amount};
    if (businessId != null) body['business_id'] = businessId;

    await _client.post('/pay/wallet/recharge', data: body);
  }

  @override
  Future<List<dynamic>> getWalletHistory({int? businessId}) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['business_id'] = businessId;

    final response = await _client.get(
      '/pay/wallet/history',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    final data = response.data;
    if (data is List) return data;
    return [];
  }

  @override
  Future<List<dynamic>> getPendingRequests() async {
    final response = await _client.get('/pay/wallet/admin/pending-requests');
    final data = response.data;
    if (data is List) return data;
    return [];
  }

  @override
  Future<List<dynamic>> getProcessedRequests() async {
    final response = await _client.get('/pay/wallet/admin/processed-requests');
    final data = response.data;
    if (data is List) return data;
    return [];
  }

  @override
  Future<void> processRequest({required String id, required String action}) async {
    await _client.post('/pay/wallet/admin/requests/$id/$action');
  }

  @override
  Future<void> manualDebit({required int businessId, required double amount, required String reference}) async {
    await _client.post('/pay/wallet/admin/manual-debit', data: {
      'business_id': businessId,
      'amount': amount,
      'reference': reference,
    });
  }

  @override
  Future<void> adminAdjustBalance({required int businessId, required double amount, required String reference}) async {
    await _client.post('/pay/wallet/admin/adjust-balance', data: {
      'business_id': businessId,
      'amount': amount,
      'reference': reference,
    });
  }

  @override
  Future<void> clearRechargeHistory({required int businessId}) async {
    await _client.delete('/pay/wallet/admin/history/$businessId');
  }

  @override
  Future<void> debitForGuide({required double amount, required String trackingNumber, int? businessId}) async {
    final body = <String, dynamic>{
      'amount': amount,
      'tracking_number': trackingNumber,
    };
    if (businessId != null) body['business_id'] = businessId;

    await _client.post('/pay/wallet/debit-guide', data: body);
  }

  @override
  Future<BusinessSubscription?> getMySubscription({int? businessId}) async {
    final queryParams = <String, dynamic>{};
    if (businessId != null) queryParams['businessId'] = businessId;

    final response = await _client.get(
      '/api/v1/subscriptions/me',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );
    final data = response.data;
    if (data['data'] != null) {
      return BusinessSubscription.fromJson(data['data']);
    }
    return null;
  }

  @override
  Future<void> registerSubscriptionPayment({
    required int businessId,
    required double amount,
    required int monthsToAdd,
    String? paymentReference,
    String? notes,
  }) async {
    final body = <String, dynamic>{
      'businessId': businessId,
      'amount': amount,
      'monthsToAdd': monthsToAdd,
    };
    if (paymentReference != null) body['paymentReference'] = paymentReference;
    if (notes != null) body['notes'] = notes;

    await _client.post('/api/v1/subscriptions/register-payment', data: body);
  }

  @override
  Future<void> disableSubscription({required int businessId}) async {
    await _client.post(
      '/api/v1/subscriptions/disable',
      queryParameters: {'businessId': businessId},
    );
  }
}
