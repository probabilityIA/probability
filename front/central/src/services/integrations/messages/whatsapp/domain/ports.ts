import {
    WhatsAppAddNumberValues,
    WhatsAppConnectionResponse,
    WhatsAppConnectionValues,
    WhatsAppEmbeddedSignupConfigResponse,
    WhatsAppEmbeddedSignupPayload,
    WhatsAppEmbeddedSignupResponse,
    WhatsAppNumberResponse,
    WhatsAppBusinessProfileResponse,
    WhatsAppBusinessProfileValues,
    WhatsAppProvisionResponse,
    WhatsAppTemplatesResponse,
} from './types';

export interface IWhatsAppRepository {
    getTemplatesStatus(businessId?: number, refresh?: boolean): Promise<WhatsAppTemplatesResponse>;
    provisionTemplates(businessId?: number): Promise<WhatsAppProvisionResponse>;
    saveConnection(values: WhatsAppConnectionValues, businessId?: number): Promise<WhatsAppConnectionResponse>;
    getNumberState(businessId?: number): Promise<WhatsAppNumberResponse>;
    getBusinessProfile(businessId?: number): Promise<WhatsAppBusinessProfileResponse>;
    updateBusinessProfile(values: WhatsAppBusinessProfileValues, businessId?: number): Promise<WhatsAppBusinessProfileResponse>;
    updateBusinessProfilePhoto(file: File, businessId?: number): Promise<WhatsAppBusinessProfileResponse>;
    addNumber(values: WhatsAppAddNumberValues, businessId?: number): Promise<WhatsAppNumberResponse>;
    requestNumberCode(method: string, businessId?: number): Promise<WhatsAppNumberResponse>;
    verifyNumberCode(code: string, businessId?: number): Promise<WhatsAppNumberResponse>;
    registerNumber(businessId?: number): Promise<WhatsAppNumberResponse>;
    getEmbeddedSignupConfig(): Promise<WhatsAppEmbeddedSignupConfigResponse>;
    completeEmbeddedSignup(
        payload: WhatsAppEmbeddedSignupPayload,
        businessId?: number
    ): Promise<WhatsAppEmbeddedSignupResponse>;
}
