import { useState, useEffect } from 'react';
import { createBusinessAction, updateBusinessAction } from '../../infra/actions';
import { Business, CreateBusinessDTO } from '../../domain/types';

export const useBusinessForm = (initialData?: Business, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [formData, setFormData] = useState<Partial<CreateBusinessDTO>>({
        name: '',
        code: '',
        business_type_id: undefined,
        timezone: '',
        address: '',
        description: '',
        primary_color: '#000000',
        secondary_color: '#ffffff',
        tertiary_color: '#cccccc',
        quaternary_color: '#eeeeee',
        is_active: true,
        enable_delivery: false,
        enable_pickup: false,
        enable_reservations: false,
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                code: initialData.code,
                business_type_id: initialData.business_type_id,
                timezone: initialData.timezone,
                address: initialData.address,
                description: initialData.description,
                primary_color: initialData.primary_color,
                secondary_color: initialData.secondary_color,
                tertiary_color: initialData.tertiary_color,
                quaternary_color: initialData.quaternary_color,
                is_active: initialData.is_active,
                enable_delivery: initialData.enable_delivery,
                enable_pickup: initialData.enable_pickup,
                enable_reservations: initialData.enable_reservations,
                // No incluir logo_file ni navbar_image_file aquí, solo las URLs
                // Los archivos se manejan cuando el usuario selecciona nuevos archivos
            });
        }
    }, [initialData]);

    const handleChange = (field: keyof CreateBusinessDTO, value: any) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const submit = async () => {
        setLoading(true);
        setError(null);

        try {
            console.log('Submitting Business Data:', formData);
            if (initialData) {
                await updateBusinessAction(initialData.id, formData);
            } else {
                await createBusinessAction(formData as CreateBusinessDTO);
            }
            if (onSuccess) onSuccess();
            return true;
        } catch (err: any) {
            console.error('Error in useBusinessForm submit:', err);
            setError(err.message || 'Error saving business');
            return false;
        } finally {
            setLoading(false);
        }
    };

    return {
        formData,
        loading,
        error,
        handleChange,
        submit,
        setError
    };
};
