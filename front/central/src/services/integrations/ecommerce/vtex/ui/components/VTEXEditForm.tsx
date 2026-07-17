'use client';

import { VTEXCredentials } from '../../domain/types';
import { VTEXConfigForm } from './VTEXConfigForm';

interface VTEXEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: any;
        credentials?: VTEXCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function VTEXEditForm({ integrationId, initialData, onSuccess, onCancel }: VTEXEditFormProps) {
    return (
        <VTEXConfigForm
            isEdit
            integrationId={integrationId}
            initialData={initialData}
            onSuccess={onSuccess}
            onCancel={onCancel}
        />
    );
}
