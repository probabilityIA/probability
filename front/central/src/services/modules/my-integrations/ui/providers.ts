import type { ComponentType } from 'react';
import { syncShopifyInventoryAction, reconcileShopifyProductsAction, associateShopifyProductsAction } from '@/services/integrations/ecommerce/shopify/infra/actions';
import { syncMeliInventoryAction, reconcileMeliProductsAction, associateMeliProductsAction } from '@/services/integrations/ecommerce/mercadolibre/infra/actions';
import { syncWooInventoryAction, reconcileWooProductsAction, associateWooProductsAction } from '@/services/integrations/ecommerce/woocommerce/infra/actions';
import { syncJumpsellerInventoryAction, reconcileJumpsellerProductsAction, associateJumpsellerProductsAction } from '@/services/integrations/ecommerce/jumpseller/infra/actions';
import { ShopifyProductSyncModal } from '@/services/integrations/ecommerce/shopify/ui/components/ShopifyProductSyncModal';
import { MercadoLibreProductSyncModal } from '@/services/integrations/ecommerce/mercadolibre/ui/components/MercadoLibreProductSyncModal';
import { WooProductSyncModal } from '@/services/integrations/ecommerce/woocommerce/ui/components/WooProductSyncModal';
import { JumpsellerProductSyncModal } from '@/services/integrations/ecommerce/jumpseller/ui/components/JumpsellerProductSyncModal';

export interface ProductSyncModalProps {
    isOpen: boolean;
    onClose: () => void;
    integrationId: number;
    businessId: number | null;
    onCompleted?: () => void;
}

export interface SyncProvider {
    typeId: number;
    key: string;
    label: string;
    inventoryEventPrefix: string;
    syncInventory: (integrationId: number, businessId?: number) => Promise<unknown>;
    reconcileProducts: (integrationId: number, businessId?: number) => Promise<unknown>;
    associateProducts: (integrationId: number, businessId?: number, skus?: string[]) => Promise<unknown>;
    onlyInChannelField: string;
    channelNoSkuField: string;
    ProductSyncModal: ComponentType<ProductSyncModalProps>;
}

export const SYNC_PROVIDERS: Record<number, SyncProvider> = {
    1: {
        typeId: 1,
        key: 'shopify',
        label: 'Shopify',
        inventoryEventPrefix: 'shopify',
        syncInventory: syncShopifyInventoryAction,
        reconcileProducts: reconcileShopifyProductsAction,
        associateProducts: associateShopifyProductsAction,
        onlyInChannelField: 'only_in_shopify',
        channelNoSkuField: 'shopify_no_sku',
        ProductSyncModal: ShopifyProductSyncModal,
    },
    3: {
        typeId: 3,
        key: 'meli',
        label: 'Mercado Libre',
        inventoryEventPrefix: 'meli',
        syncInventory: syncMeliInventoryAction,
        reconcileProducts: reconcileMeliProductsAction,
        associateProducts: associateMeliProductsAction,
        onlyInChannelField: 'only_in_meli',
        channelNoSkuField: 'meli_no_sku',
        ProductSyncModal: MercadoLibreProductSyncModal,
    },
    4: {
        typeId: 4,
        key: 'woocommerce',
        label: 'WooCommerce',
        inventoryEventPrefix: 'woo',
        syncInventory: syncWooInventoryAction,
        reconcileProducts: reconcileWooProductsAction,
        associateProducts: associateWooProductsAction,
        onlyInChannelField: 'only_in_woocommerce',
        channelNoSkuField: 'woocommerce_no_sku',
        ProductSyncModal: WooProductSyncModal,
    },
    33: {
        typeId: 33,
        key: 'jumpseller',
        label: 'Jumpseller',
        inventoryEventPrefix: 'jumpseller',
        syncInventory: syncJumpsellerInventoryAction,
        reconcileProducts: reconcileJumpsellerProductsAction,
        associateProducts: associateJumpsellerProductsAction,
        onlyInChannelField: 'only_in_jumpseller',
        channelNoSkuField: 'jumpseller_no_sku',
        ProductSyncModal: JumpsellerProductSyncModal,
    },
};

const INVENTORY_EVENT_SUFFIXES = ['started', 'item', 'progress', 'completed'];

export const GLOBAL_INVENTORY_EVENT_TYPES = Object.values(SYNC_PROVIDERS).flatMap(p =>
    INVENTORY_EVENT_SUFFIXES.map(s => `${p.inventoryEventPrefix}.inventory.sync.${s}`)
);

export function getSyncProvider(integrationTypeId: number | string | undefined): SyncProvider | null {
    if (integrationTypeId === undefined || integrationTypeId === null) return null;
    return SYNC_PROVIDERS[Number(integrationTypeId)] || null;
}
