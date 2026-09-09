class DeviceRegistration {
  final String token;
  final String platform;
  final String? appVersion;
  final String? deviceName;
  final int? businessId;

  const DeviceRegistration({
    required this.token,
    required this.platform,
    this.appVersion,
    this.deviceName,
    this.businessId,
  });

  Map<String, dynamic> toJson() => {
        'token': token,
        'platform': platform,
        if (appVersion != null) 'app_version': appVersion,
        if (deviceName != null) 'device_name': deviceName,
        if (businessId != null) 'business_id': businessId,
      };
}

class PushPayload {
  final String? title;
  final String? body;
  final Map<String, String> data;

  const PushPayload({this.title, this.body, this.data = const {}});

  String? get route => data['route'];
  String? get orderId => data['order_id'];
  String? get orderNumber => data['order_number'];
}
