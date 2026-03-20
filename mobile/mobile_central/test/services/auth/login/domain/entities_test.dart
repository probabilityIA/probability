import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/login/domain/entities.dart';

void main() {
  group('UserInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'name': 'John Doe',
        'email': 'john@example.com',
        'phone': '+573001234567',
        'avatar_url': 'https://example.com/avatar.png',
        'is_active': true,
        'last_login_at': '2026-03-19T10:00:00Z',
      };

      final user = UserInfo.fromJson(json);

      expect(user.id, 42);
      expect(user.name, 'John Doe');
      expect(user.email, 'john@example.com');
      expect(user.phone, '+573001234567');
      expect(user.avatarUrl, 'https://example.com/avatar.png');
      expect(user.isActive, true);
      expect(user.lastLoginAt, '2026-03-19T10:00:00Z');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'name': 'Jane',
        'email': 'jane@example.com',
        'is_active': true,
      };

      final user = UserInfo.fromJson(json);

      expect(user.id, 1);
      expect(user.name, 'Jane');
      expect(user.email, 'jane@example.com');
      expect(user.phone, isNull);
      expect(user.avatarUrl, isNull);
      expect(user.isActive, true);
      expect(user.lastLoginAt, isNull);
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final user = UserInfo.fromJson(json);

      expect(user.id, 0);
      expect(user.name, '');
      expect(user.email, '');
      expect(user.isActive, false);
    });

    test('fromJson handles is_active false', () {
      final json = {
        'id': 5,
        'name': 'Inactive User',
        'email': 'inactive@example.com',
        'is_active': false,
      };

      final user = UserInfo.fromJson(json);

      expect(user.isActive, false);
    });
  });

  group('BusinessInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 10,
        'name': 'Test Business',
        'logo_url': 'https://example.com/logo.png',
        'primary_color': '#FF0000',
        'secondary_color': '#00FF00',
        'accent_color': '#0000FF',
      };

      final business = BusinessInfo.fromJson(json);

      expect(business.id, 10);
      expect(business.name, 'Test Business');
      expect(business.logoUrl, 'https://example.com/logo.png');
      expect(business.primaryColor, '#FF0000');
      expect(business.secondaryColor, '#00FF00');
      expect(business.accentColor, '#0000FF');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 3,
        'name': 'Minimal Business',
      };

      final business = BusinessInfo.fromJson(json);

      expect(business.id, 3);
      expect(business.name, 'Minimal Business');
      expect(business.logoUrl, isNull);
      expect(business.primaryColor, isNull);
      expect(business.secondaryColor, isNull);
      expect(business.accentColor, isNull);
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final business = BusinessInfo.fromJson(json);

      expect(business.id, 0);
      expect(business.name, '');
    });
  });

  group('LoginResponse', () {
    test('fromJson parses complete response', () {
      final json = {
        'user': {
          'id': 1,
          'name': 'Admin',
          'email': 'admin@example.com',
          'is_active': true,
        },
        'token': 'jwt-token-abc123',
        'require_password_change': false,
        'businesses': [
          {'id': 1, 'name': 'Business One'},
          {'id': 2, 'name': 'Business Two'},
        ],
        'scope': 'full',
        'is_super_admin': true,
      };

      final response = LoginResponse.fromJson(json);

      expect(response.user.id, 1);
      expect(response.user.name, 'Admin');
      expect(response.user.email, 'admin@example.com');
      expect(response.token, 'jwt-token-abc123');
      expect(response.requirePasswordChange, false);
      expect(response.businesses.length, 2);
      expect(response.businesses[0].name, 'Business One');
      expect(response.businesses[1].name, 'Business Two');
      expect(response.scope, 'full');
      expect(response.isSuperAdmin, true);
    });

    test('fromJson handles empty businesses list', () {
      final json = {
        'user': {'id': 1, 'name': 'User', 'email': 'u@e.com', 'is_active': true},
        'token': 'tok',
        'require_password_change': false,
        'businesses': [],
        'is_super_admin': false,
      };

      final response = LoginResponse.fromJson(json);

      expect(response.businesses, isEmpty);
    });

    test('fromJson handles null businesses', () {
      final json = {
        'user': {'id': 1, 'name': 'User', 'email': 'u@e.com', 'is_active': true},
        'token': 'tok',
        'require_password_change': false,
        'is_super_admin': false,
      };

      final response = LoginResponse.fromJson(json);

      expect(response.businesses, isEmpty);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = LoginResponse.fromJson(json);

      expect(response.user.id, 0);
      expect(response.token, '');
      expect(response.requirePasswordChange, false);
      expect(response.businesses, isEmpty);
      expect(response.scope, isNull);
      expect(response.isSuperAdmin, false);
    });

    test('fromJson handles require_password_change true', () {
      final json = {
        'user': {'id': 1, 'name': 'New', 'email': 'new@e.com', 'is_active': true},
        'token': 'temp-token',
        'require_password_change': true,
        'is_super_admin': false,
      };

      final response = LoginResponse.fromJson(json);

      expect(response.requirePasswordChange, true);
    });
  });

  group('LoginSuccessResponse', () {
    test('fromJson parses success response', () {
      final json = {
        'success': true,
        'data': {
          'user': {
            'id': 5,
            'name': 'Test User',
            'email': 'test@example.com',
            'is_active': true,
          },
          'token': 'my-token',
          'require_password_change': false,
          'businesses': [
            {'id': 1, 'name': 'Biz'},
          ],
          'is_super_admin': false,
        },
      };

      final response = LoginSuccessResponse.fromJson(json);

      expect(response.success, true);
      expect(response.data.user.id, 5);
      expect(response.data.user.name, 'Test User');
      expect(response.data.token, 'my-token');
      expect(response.data.businesses.length, 1);
    });

    test('fromJson handles success false', () {
      final json = {
        'success': false,
        'data': {
          'user': {'id': 0, 'name': '', 'email': '', 'is_active': false},
          'token': '',
          'require_password_change': false,
          'is_super_admin': false,
        },
      };

      final response = LoginSuccessResponse.fromJson(json);

      expect(response.success, false);
    });

    test('fromJson falls back to root json when data is missing', () {
      // When 'data' key is missing, fromJson uses json itself as the data
      final json = {
        'success': true,
        'user': {
          'id': 7,
          'name': 'Fallback',
          'email': 'fb@e.com',
          'is_active': true,
        },
        'token': 'fallback-token',
        'require_password_change': false,
        'is_super_admin': false,
      };

      final response = LoginSuccessResponse.fromJson(json);

      expect(response.success, true);
      expect(response.data.user.id, 7);
      expect(response.data.token, 'fallback-token');
    });

    test('fromJson defaults success to false when missing', () {
      final json = {
        'data': {
          'user': {'id': 1, 'name': 'X', 'email': 'x@e.com', 'is_active': true},
          'token': 't',
          'require_password_change': false,
          'is_super_admin': false,
        },
      };

      final response = LoginSuccessResponse.fromJson(json);

      expect(response.success, false);
    });
  });

  group('UserRolesPermissionsResponse', () {
    test('fromJson parses complete response', () {
      final json = {
        'is_super': true,
        'business_id': 5,
        'business_name': 'My Business',
        'business_type_id': 2,
        'business_type_name': 'E-commerce',
        'role': 'admin',
        'resources': [
          {
            'resource': 'orders',
            'actions': ['read', 'write', 'delete'],
          },
          {
            'resource': 'products',
            'actions': ['read'],
          },
        ],
        'subscription_status': 'active',
      };

      final response = UserRolesPermissionsResponse.fromJson(json);

      expect(response.isSuper, true);
      expect(response.businessId, 5);
      expect(response.businessName, 'My Business');
      expect(response.businessTypeId, 2);
      expect(response.businessTypeName, 'E-commerce');
      expect(response.role, 'admin');
      expect(response.resources.length, 2);
      expect(response.resources[0].resource, 'orders');
      expect(response.resources[0].actions, ['read', 'write', 'delete']);
      expect(response.resources[1].resource, 'products');
      expect(response.resources[1].actions, ['read']);
      expect(response.subscriptionStatus, 'active');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'is_super': false,
        'business_id': 1,
        'resources': [],
      };

      final response = UserRolesPermissionsResponse.fromJson(json);

      expect(response.isSuper, false);
      expect(response.businessId, 1);
      expect(response.businessName, isNull);
      expect(response.businessTypeId, isNull);
      expect(response.businessTypeName, isNull);
      expect(response.role, isNull);
      expect(response.resources, isEmpty);
      expect(response.subscriptionStatus, isNull);
    });

    test('fromJson handles null resources', () {
      final json = {
        'is_super': false,
        'business_id': 1,
      };

      final response = UserRolesPermissionsResponse.fromJson(json);

      expect(response.resources, isEmpty);
    });

    test('fromJson uses defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final response = UserRolesPermissionsResponse.fromJson(json);

      expect(response.isSuper, false);
      expect(response.businessId, 0);
      expect(response.resources, isEmpty);
    });
  });

  group('ResourcePermission', () {
    test('fromJson parses correctly', () {
      final json = {
        'resource': 'invoices',
        'actions': ['read', 'write', 'export'],
      };

      final permission = ResourcePermission.fromJson(json);

      expect(permission.resource, 'invoices');
      expect(permission.actions, ['read', 'write', 'export']);
    });

    test('fromJson handles empty actions', () {
      final json = {
        'resource': 'dashboard',
        'actions': [],
      };

      final permission = ResourcePermission.fromJson(json);

      expect(permission.resource, 'dashboard');
      expect(permission.actions, isEmpty);
    });

    test('fromJson handles null actions', () {
      final json = {
        'resource': 'settings',
      };

      final permission = ResourcePermission.fromJson(json);

      expect(permission.resource, 'settings');
      expect(permission.actions, isEmpty);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final permission = ResourcePermission.fromJson(json);

      expect(permission.resource, '');
      expect(permission.actions, isEmpty);
    });

    test('fromJson handles single action', () {
      final json = {
        'resource': 'reports',
        'actions': ['read'],
      };

      final permission = ResourcePermission.fromJson(json);

      expect(permission.actions.length, 1);
      expect(permission.actions.first, 'read');
    });
  });

  group('ChangePasswordResponse', () {
    test('fromJson parses success response', () {
      final json = {
        'success': true,
        'message': 'Password changed successfully',
      };

      final response = ChangePasswordResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, 'Password changed successfully');
    });

    test('fromJson parses failure response', () {
      final json = {
        'success': false,
        'message': 'Current password is incorrect',
      };

      final response = ChangePasswordResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, 'Current password is incorrect');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = ChangePasswordResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, '');
    });
  });

  group('GeneratePasswordResponse', () {
    test('fromJson parses success response with password', () {
      final json = {
        'success': true,
        'message': 'Password generated',
        'password': 'xK9#mP2!nQ4',
      };

      final response = GeneratePasswordResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, 'Password generated');
      expect(response.password, 'xK9#mP2!nQ4');
    });

    test('fromJson handles null password', () {
      final json = {
        'success': true,
        'message': 'Password sent by email',
      };

      final response = GeneratePasswordResponse.fromJson(json);

      expect(response.success, true);
      expect(response.message, 'Password sent by email');
      expect(response.password, isNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final response = GeneratePasswordResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, '');
      expect(response.password, isNull);
    });

    test('fromJson parses failure response', () {
      final json = {
        'success': false,
        'message': 'User not found',
      };

      final response = GeneratePasswordResponse.fromJson(json);

      expect(response.success, false);
      expect(response.message, 'User not found');
      expect(response.password, isNull);
    });
  });
}
