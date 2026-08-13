import type { ComponentType } from 'react';
import { syncShopifyInventoryAction, reconcileShopifyProductsAction, associateShopifyProductsAction, applyShopifyProductsAction } from '@/services/integrations/ecommerce/shopify/infra/actions';
import { syncMeliInventoryAction, compareMeliInventoryAction, reconcileMeliProductsAction, associateMeliProductsAction, applyMeliProductsAction } from '@/services/integrations/ecommerce/mercadolibre/infra/actions';
import { syncWooInventoryAction, compareWooInventoryAction, reconcileWooProductsAction, associateWooProductsAction, applyWooProductsAction } from '@/services/integrations/ecommerce/woocommerce/infra/actions';
import { syncJumpsellerInventoryAction, reconcileJumpsellerProductsAction, associateJumpsellerProductsAction, applyJumpsellerProductsAction } from '@/services/integrations/ecommerce/jumpseller/infra/actions';
import { syncSiigoInventoryAction, compareSiigoInventoryAction, startSiigoProductsReconcileAction, applySiigoProductsAction } from '@/services/integrations/invoicing/siigo/infra/actions';
import { syncVTEXInventoryAction, reconcileVTEXProductsAction, associateVTEXProductsAction, applyVTEXProductsAction } from '@/services/integrations/ecommerce/vtex/infra/actions';
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

export type ProductApply = (integrationId: number, businessId?: number, skus?: string[]) => Promise<unknown>;

export interface ProductApplyActions {
    createInChannel?: ProductApply;
    createInProbability?: ProductApply;
    updateInProbability?: ProductApply;
}

export interface SyncProvider {
    typeId: number;
    key: string;
    label: string;
    inventoryEventPrefix: string;
    productEventPrefix?: string;
    syncInventory: (integrationId: number, businessId?: number, skus?: string[]) => Promise<unknown>;
    compareInventory?: (integrationId: number, businessId?: number, page?: number, pageSize?: number, skus?: string[]) => Promise<unknown>;
    reconcileProducts: (integrationId: number, businessId?: number) => Promise<unknown>;
    associateProducts?: (integrationId: number, businessId?: number, skus?: string[]) => Promise<unknown>;
    apply: ProductApplyActions;
    onlyInChannelField: string;
    channelNoSkuField: string;
    ProductSyncModal?: ComponentType<ProductSyncModalProps>;
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
        apply: {
            createInChannel: (id, bid, skus) => applyShopifyProductsAction(id, 'to_shopify', bid, skus),
            createInProbability: (id, bid, skus) => applyShopifyProductsAction(id, 'to_probability', bid, skus),
        },
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
        compareInventory: compareMeliInventoryAction,
        reconcileProducts: reconcileMeliProductsAction,
        associateProducts: associateMeliProductsAction,
        apply: {
            createInChannel: (id, bid, skus) => applyMeliProductsAction(id, 'to_meli', bid, skus),
            createInProbability: (id, bid, skus) => applyMeliProductsAction(id, 'to_probability', bid, skus),
        },
        onlyInChannelField: 'only_in_meli',
        channelNoSkuField: 'meli_no_sku',
        ProductSyncModal: MercadoLibreProductSyncModal,
    },
    4: {
        typeId: 4,
        key: 'woocommerce',
        label: 'WooCommerce',
        inventoryEventPrefix: 'woo',
        productEventPrefix: 'woocommerce',
        syncInventory: syncWooInventoryAction,
        compareInventory: compareWooInventoryAction,
        reconcileProducts: reconcileWooProductsAction,
        associateProducts: associateWooProductsAction,
        apply: {
            createInChannel: (id, bid, skus) => applyWooProductsAction(id, 'to_woo', bid, skus),
            createInProbability: (id, bid, skus) => applyWooProductsAction(id, 'to_probability', bid, skus),
        },
        onlyInChannelField: 'only_in_woocommerce',
        channelNoSkuField: 'woocommerce_no_sku',
        ProductSyncModal: WooProductSyncModal,
    },
    16: {
        typeId: 16,
        key: 'vtex',
        label: 'VTEX',
        inventoryEventPrefix: 'vtex',
        syncInventory: syncVTEXInventoryAction,
        reconcileProducts: reconcileVTEXProductsAction,
        associateProducts: associateVTEXProductsAction,
        apply: {
            createInProbability: (id, bid, skus) => applyVTEXProductsAction(id, bid, 'create', skus),
            updateInProbability: (id, bid, skus) => applyVTEXProductsAction(id, bid, 'update', skus),
        },
        onlyInChannelField: 'only_in_vtex',
        channelNoSkuField: 'vtex_no_sku',
    },
    8: {
        typeId: 8,
        key: 'siigo',
        label: 'Siigo',
        inventoryEventPrefix: 'siigo',
        syncInventory: syncSiigoInventoryAction,
        compareInventory: compareSiigoInventoryAction,
        reconcileProducts: startSiigoProductsReconcileAction,
        apply: {
            createInProbability: (id, bid, skus) => applySiigoProductsAction(id, bid, skus),
        },
        onlyInChannelField: 'only_in_siigo',
        channelNoSkuField: 'siigo_no_sku',
    },
    33: {
        typeId: 33,
        key: 'jumpseller',
        label: 'Jumpseller',
        inventoryEventPrefix: 'jumpseller',
        syncInventory: syncJumpsellerInventoryAction,
        reconcileProducts: reconcileJumpsellerProductsAction,
        associateProducts: associateJumpsellerProductsAction,
        apply: {
            createInChannel: (id, bid, skus) => applyJumpsellerProductsAction(id, 'to_jumpseller', bid, 'create', skus),
            createInProbability: (id, bid, skus) => applyJumpsellerProductsAction(id, 'to_probability', bid, 'create', skus),
            updateInProbability: (id, bid, skus) => applyJumpsellerProductsAction(id, 'to_probability', bid, 'update', skus),
        },
        onlyInChannelField: 'only_in_jumpseller',
        channelNoSkuField: 'jumpseller_no_sku',
        ProductSyncModal: JumpsellerProductSyncModal,
    },
};

const INVENTORY_EVENT_SUFFIXES = ['started', 'item', 'progress', 'completed'];

const PRODUCT_SYNC_SUFFIXES = ['started', 'progress', 'completed'];

export const GLOBAL_INVENTORY_EVENT_TYPES = Object.values(SYNC_PROVIDERS).flatMap(p => [
    ...INVENTORY_EVENT_SUFFIXES.map(s => `${p.inventoryEventPrefix}.inventory.sync.${s}`),
    `${p.inventoryEventPrefix}.product.reconcile.started`,
    `${p.inventoryEventPrefix}.product.reconcile.completed`,
    ...PRODUCT_SYNC_SUFFIXES.map(s => `${p.productEventPrefix ?? p.inventoryEventPrefix}.product.sync.${s}`),
]);

export const SIIGO_TYPE_ID = 8;

export function getSyncProvider(integrationTypeId: number | string | undefined): SyncProvider | null {
    if (integrationTypeId === undefined || integrationTypeId === null) return null;
    return SYNC_PROVIDERS[Number(integrationTypeId)] || null;
}
