class Wallet {
  final String id;
  final int businessId;
  final double balance;

  Wallet({
    required this.id,
    required this.businessId,
    required this.balance,
  });

  factory Wallet.fromJson(Map<String, dynamic> json) {
    return Wallet(
      id: json['ID']?.toString() ?? json['id']?.toString() ?? '',
      businessId: json['BusinessID'] ?? json['business_id'] ?? 0,
      balance: (json['Balance'] ?? json['balance'] ?? 0).toDouble(),
    );
  }
}

class WalletTransaction {
  final String id;
  final String walletId;
  final double amount;
  final String createdAt;

  WalletTransaction({
    required this.id,
    required this.walletId,
    required this.amount,
    required this.createdAt,
  });

  factory WalletTransaction.fromJson(Map<String, dynamic> json) {
    return WalletTransaction(
      id: json['ID']?.toString() ?? json['id']?.toString() ?? '',
      walletId: json['WalletID']?.toString() ?? json['wallet_id']?.toString() ?? '',
      amount: (json['Amount'] ?? json['amount'] ?? 0).toDouble(),
      createdAt: json['CreatedAt'] ?? json['created_at'] ?? '',
    );
  }
}

class BusinessSubscription {
  final int? id;
  final int businessId;
  final double amount;
  final String? startDate;
  final String? endDate;
  final String status;
  final String? paymentReference;
  final String? notes;
  final String? createdAt;

  BusinessSubscription({
    this.id,
    required this.businessId,
    required this.amount,
    this.startDate,
    this.endDate,
    required this.status,
    this.paymentReference,
    this.notes,
    this.createdAt,
  });

  factory BusinessSubscription.fromJson(Map<String, dynamic> json) {
    return BusinessSubscription(
      id: json['id'],
      businessId: json['businessId'] ?? json['business_id'] ?? 0,
      amount: (json['amount'] ?? 0).toDouble(),
      startDate: json['startDate'] ?? json['start_date'],
      endDate: json['endDate'] ?? json['end_date'],
      status: json['status'] ?? '',
      paymentReference: json['paymentReference'] ?? json['payment_reference'],
      notes: json['notes'],
      createdAt: json['createdAt'] ?? json['created_at'],
    );
  }
}

class BusinessSubscriptionStatus {
  final String subscriptionStatus;
  final String? subscriptionEndDate;
  final String? businessName;

  BusinessSubscriptionStatus({
    required this.subscriptionStatus,
    this.subscriptionEndDate,
    this.businessName,
  });

  factory BusinessSubscriptionStatus.fromJson(Map<String, dynamic> json) {
    return BusinessSubscriptionStatus(
      subscriptionStatus: json['subscriptionStatus'] ?? json['subscription_status'] ?? '',
      subscriptionEndDate: json['subscriptionEndDate'] ?? json['subscription_end_date'],
      businessName: json['businessName'] ?? json['business_name'],
    );
  }
}
