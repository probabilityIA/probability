import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class ShipmentUseCases {
  final IShipmentRepository _repository;
  ShipmentUseCases(this._repository);

  Future<PaginatedResponse<Shipment>> getShipments(GetShipmentsParams? params) => _repository.getShipments(params);
  Future<Map<String, dynamic>> quoteShipment(EnvioClickQuoteRequest req) => _repository.quoteShipment(req);
  Future<Map<String, dynamic>> generateGuide(EnvioClickQuoteRequest req) => _repository.generateGuide(req);
  Future<Map<String, dynamic>> trackShipment(String trackingNumber) => _repository.trackShipment(trackingNumber);
  Future<Map<String, dynamic>> cancelShipment(String id) => _repository.cancelShipment(id);
  Future<Shipment> createShipment(Map<String, dynamic> data) => _repository.createShipment(data);
  Future<List<OriginAddress>> getOriginAddresses({int? businessId}) => _repository.getOriginAddresses(businessId: businessId);
  Future<OriginAddress> createOriginAddress(CreateOriginAddressDTO data, {int? businessId}) => _repository.createOriginAddress(data, businessId: businessId);
  Future<OriginAddress> updateOriginAddress(int id, Map<String, dynamic> data, {int? businessId}) => _repository.updateOriginAddress(id, data, businessId: businessId);
  Future<void> deleteOriginAddress(int id, {int? businessId}) => _repository.deleteOriginAddress(id, businessId: businessId);
}
