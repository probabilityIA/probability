import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/paymentstatus/domain/entities.dart';

void main() {
  group('PaymentStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'paid',
        'name': 'Paid',
        'description': 'Payment completed',
        'category': 'completed',
        'color': '#00FF00',
        'icon': 'check_circle',
        'is_active': true,
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.id, 1);
      expect(status.code, 'paid');
      expect(status.name, 'Paid');
      expect(status.description, 'Payment completed');
      expect(status.category, 'completed');
      expect(status.color, '#00FF00');
      expect(status.icon, 'check_circle');
      expect(status.isActive, true);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 2,
        'code': 'pending',
        'name': 'Pending',
        'is_active': true,
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.id, 2);
      expect(status.code, 'pending');
      expect(status.name, 'Pending');
      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
      expect(status.icon, isNull);
      expect(status.isActive, true);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
      expect(status.isActive, true);
    });

    test('fromJson handles false is_active', () {
      final json = {
        'id': 3,
        'code': 'cancelled',
        'name': 'Cancelled',
        'is_active': false,
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.isActive, false);
    });

    test('fromJson handles all null optional fields explicitly', () {
      final json = {
        'id': 4,
        'code': 'refunded',
        'name': 'Refunded',
        'description': null,
        'category': null,
        'color': null,
        'icon': null,
        'is_active': true,
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
      expect(status.icon, isNull);
    });
  });

  group('GetPaymentStatusesParams', () {
    test('toQueryParams includes isActive when provided', () {
      final params = GetPaymentStatusesParams(isActive: true);

      final query = params.toQueryParams();

      expect(query['is_active'], true);
    });

    test('toQueryParams returns empty map when isActive is null', () {
      final params = GetPaymentStatusesParams();

      final query = params.toQueryParams();

      expect(query, isEmpty);
    });

    test('toQueryParams includes false isActive', () {
      final params = GetPaymentStatusesParams(isActive: false);

      final query = params.toQueryParams();

      expect(query['is_active'], false);
    });
  });
}
