import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/wallet/domain/entities.dart';

void main() {
  group('Wallet', () {
    test('fromJson parses all fields correctly (snake_case)', () {
      final json = {
        'id': 42,
        'business_id': 5,
        'balance': 150000.50,
      };

      final wallet = Wallet.fromJson(json);

      expect(wallet.id, '42');
      expect(wallet.businessId, 5);
      expect(wallet.balance, 150000.50);
    });

    test('fromJson parses all fields correctly (PascalCase)', () {
      final json = {
        'ID': 42,
        'BusinessID': 5,
        'Balance': 150000.50,
      };

      final wallet = Wallet.fromJson(json);

      expect(wallet.id, '42');
      expect(wallet.businessId, 5);
      expect(wallet.balance, 150000.50);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final wallet = Wallet.fromJson(json);

      expect(wallet.id, '');
      expect(wallet.businessId, 0);
      expect(wallet.balance, 0.0);
    });

    test('fromJson prefers PascalCase over snake_case', () {
      final json = {
        'ID': 10,
        'id': 20,
        'BusinessID': 3,
        'business_id': 5,
        'Balance': 100.0,
        'balance': 200.0,
      };

      final wallet = Wallet.fromJson(json);

      // PascalCase keys are checked first
      expect(wallet.id, '10');
      expect(wallet.businessId, 3);
      expect(wallet.balance, 100.0);
    });

    test('fromJson converts id to string', () {
      final json = {'id': 123};
      final wallet = Wallet.fromJson(json);
      expect(wallet.id, '123');
    });

    test('fromJson handles integer balance', () {
      final json = {'balance': 5000};
      final wallet = Wallet.fromJson(json);
      expect(wallet.balance, 5000.0);
    });
  });

  group('WalletTransaction', () {
    test('fromJson parses all fields correctly (snake_case)', () {
      final json = {
        'id': 1,
        'wallet_id': 5,
        'amount': 25000.0,
        'created_at': '2026-03-01T10:00:00Z',
      };

      final tx = WalletTransaction.fromJson(json);

      expect(tx.id, '1');
      expect(tx.walletId, '5');
      expect(tx.amount, 25000.0);
      expect(tx.createdAt, '2026-03-01T10:00:00Z');
    });

    test('fromJson parses all fields correctly (PascalCase)', () {
      final json = {
        'ID': 1,
        'WalletID': 5,
        'Amount': 25000.0,
        'CreatedAt': '2026-03-01T10:00:00Z',
      };

      final tx = WalletTransaction.fromJson(json);

      expect(tx.id, '1');
      expect(tx.walletId, '5');
      expect(tx.amount, 25000.0);
      expect(tx.createdAt, '2026-03-01T10:00:00Z');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final tx = WalletTransaction.fromJson(json);

      expect(tx.id, '');
      expect(tx.walletId, '');
      expect(tx.amount, 0.0);
      expect(tx.createdAt, '');
    });

    test('fromJson converts ids to string', () {
      final json = {'id': 99, 'wallet_id': 42};
      final tx = WalletTransaction.fromJson(json);
      expect(tx.id, '99');
      expect(tx.walletId, '42');
    });
  });

  group('BusinessSubscription', () {
    test('fromJson parses all fields correctly (snake_case)', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'amount': 99000.0,
        'start_date': '2026-01-01',
        'end_date': '2026-12-31',
        'status': 'active',
        'payment_reference': 'REF-001',
        'notes': 'Annual plan',
        'created_at': '2026-01-01T00:00:00Z',
      };

      final sub = BusinessSubscription.fromJson(json);

      expect(sub.id, 1);
      expect(sub.businessId, 5);
      expect(sub.amount, 99000.0);
      expect(sub.startDate, '2026-01-01');
      expect(sub.endDate, '2026-12-31');
      expect(sub.status, 'active');
      expect(sub.paymentReference, 'REF-001');
      expect(sub.notes, 'Annual plan');
      expect(sub.createdAt, '2026-01-01T00:00:00Z');
    });

    test('fromJson parses all fields correctly (camelCase)', () {
      final json = {
        'id': 1,
        'businessId': 5,
        'amount': 50000.0,
        'startDate': '2026-01-01',
        'endDate': '2026-06-30',
        'status': 'active',
        'paymentReference': 'REF-002',
        'notes': 'Monthly',
        'createdAt': '2026-01-01',
      };

      final sub = BusinessSubscription.fromJson(json);

      expect(sub.businessId, 5);
      expect(sub.startDate, '2026-01-01');
      expect(sub.endDate, '2026-06-30');
      expect(sub.paymentReference, 'REF-002');
      expect(sub.createdAt, '2026-01-01');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final sub = BusinessSubscription.fromJson(json);

      expect(sub.id, isNull);
      expect(sub.businessId, 0);
      expect(sub.amount, 0.0);
      expect(sub.startDate, isNull);
      expect(sub.endDate, isNull);
      expect(sub.status, '');
      expect(sub.paymentReference, isNull);
      expect(sub.notes, isNull);
      expect(sub.createdAt, isNull);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': null,
        'business_id': 1,
        'amount': 100,
        'status': 'pending',
        'start_date': null,
        'end_date': null,
        'payment_reference': null,
        'notes': null,
        'created_at': null,
      };

      final sub = BusinessSubscription.fromJson(json);

      expect(sub.id, isNull);
      expect(sub.startDate, isNull);
      expect(sub.endDate, isNull);
      expect(sub.paymentReference, isNull);
      expect(sub.notes, isNull);
      expect(sub.createdAt, isNull);
    });
  });

  group('BusinessSubscriptionStatus', () {
    test('fromJson parses all fields correctly (snake_case)', () {
      final json = {
        'subscription_status': 'active',
        'subscription_end_date': '2026-12-31',
        'business_name': 'Test Business',
      };

      final status = BusinessSubscriptionStatus.fromJson(json);

      expect(status.subscriptionStatus, 'active');
      expect(status.subscriptionEndDate, '2026-12-31');
      expect(status.businessName, 'Test Business');
    });

    test('fromJson parses all fields correctly (camelCase)', () {
      final json = {
        'subscriptionStatus': 'expired',
        'subscriptionEndDate': '2025-12-31',
        'businessName': 'My Business',
      };

      final status = BusinessSubscriptionStatus.fromJson(json);

      expect(status.subscriptionStatus, 'expired');
      expect(status.subscriptionEndDate, '2025-12-31');
      expect(status.businessName, 'My Business');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final status = BusinessSubscriptionStatus.fromJson(json);

      expect(status.subscriptionStatus, '');
      expect(status.subscriptionEndDate, isNull);
      expect(status.businessName, isNull);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'subscription_status': 'active',
        'subscription_end_date': null,
        'business_name': null,
      };

      final status = BusinessSubscriptionStatus.fromJson(json);

      expect(status.subscriptionStatus, 'active');
      expect(status.subscriptionEndDate, isNull);
      expect(status.businessName, isNull);
    });
  });
}
