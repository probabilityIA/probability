export interface ShippingBox {
    name: string;
    weight: number | null;
    length: number | null;
    width: number | null;
    height: number | null;
    max_items: number;
}

export interface DirectIntegration {
    enabled: boolean;
    integration_id: number | null;
    status: 'unavailable' | 'pending' | 'active';
}

export interface CarrierSetting {
    code: string;
    enabled: boolean;
    allow_cod: boolean;
    allow_prepaid: boolean;
    direct: DirectIntegration;
}

export interface ShippingConfig {
    id: number;
    business_id: number;
    warehouse_id: number | null;
    package_strategy: 'product_dimensions' | 'standard_box';
    boxes: ShippingBox[];
    carriers: CarrierSetting[];
    always_insure: boolean;
}

export interface WarehouseOrigin {
    id: number;
    name: string;
    address: string;
    city: string;
    state: string;
    phone: string;
    is_default: boolean;
    is_active: boolean;
    has_config: boolean;
}

export interface CarrierCatalogItem {
    code: string;
    name: string;
    direct_available: boolean;
}

export interface ShippingConfigOverview {
    business: ShippingConfig;
    warehouses: WarehouseOrigin[];
    overrides: ShippingConfig[];
    carriers: CarrierCatalogItem[];
}

export interface SaveShippingConfigRequest {
    package_strategy: 'product_dimensions' | 'standard_box';
    boxes: ShippingBox[];
    carriers: CarrierSetting[];
    always_insure: boolean;
}
