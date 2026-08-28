import type { TourDefinition } from '../domain/types';
import { homeTour } from './home';
import { ordersTour } from './orders';
import { productsTour } from './products';
import {
    shipmentsTour,
    shipmentsCodTour,
    shipmentsQuotesTour,
    shippingMarginsTour,
} from './shipments';
import { inventoryTour, warehousesTour } from './inventory';
import { walletTour, accountingTour, invoicingTour, subscriptionTour } from './finanzas';
import {
    customersTour,
    integrationsTour,
    deliveryTour,
    notificationsTour,
    storefrontTour,
} from './operacion';
import { iamTour, ticketsTour, announcementsTour, commercialTour } from './administracion';

export const TOUR_LIST: TourDefinition[] = [
    homeTour,
    ordersTour,
    productsTour,
    shipmentsTour,
    shipmentsCodTour,
    shipmentsQuotesTour,
    shippingMarginsTour,
    inventoryTour,
    warehousesTour,
    walletTour,
    accountingTour,
    invoicingTour,
    subscriptionTour,
    customersTour,
    integrationsTour,
    deliveryTour,
    notificationsTour,
    storefrontTour,
    iamTour,
    ticketsTour,
    announcementsTour,
    commercialTour,
];

export const TOUR_REGISTRY: Record<string, TourDefinition> = TOUR_LIST.reduce(
    (acc, tour) => ({ ...acc, [tour.key]: tour }),
    {} as Record<string, TourDefinition>,
);
