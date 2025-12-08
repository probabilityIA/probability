import { useState, useEffect } from 'react';
import { createUserAction, updateUserAction } from '../../infra/actions';
import { User, CreateUserDTO, UpdateUserDTO } from '../../domain/types';

export const useUserForm = (initialData?: User, onSuccess?: () => void) => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);

    const [formData, setFormData] = useState<Partial<CreateUserDTO>>({
        name: '',
        email: '',
        phone: '',
        is_active: true,
        business_ids: [],
    });

    // Separate state for file input as it's not part of initialData usually
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [removeAvatar, setRemoveAvatar] = useState(false);

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                email: initialData.email,
                phone: initialData.phone || '',
                is_active: initialData.is_active,
                business_ids: initialData.business_role_assignments?.map(a => a.business_id) || [],
            });
        }
    }, [initialData]);

    const handleChange = (field: keyof CreateUserDTO, value: string | number | boolean | number[] | null) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const handleFileChange = (file: File | null) => {
        setAvatarFile(file);
        // Si se elimina el archivo y hay datos iniciales con avatar, marcar para eliminar avatar
        if (!file && initialData?.avatar_url) {
            setRemoveAvatar(true);
        } else {
            // Si se selecciona un nuevo archivo, no eliminar el avatar existente
            setRemoveAvatar(false);
        }
    };

    const submit = async () => {
        setLoading(true);
        setError(null);
        setSuccessMessage(null);

        try {
            let dataToSubmit: CreateUserDTO | UpdateUserDTO = { ...formData };
            
            if (avatarFile) {
                dataToSubmit.avatarFile = avatarFile;
            }
            
            // Si estamos actualizando y se eliminó el avatar, agregar flag
            if (initialData) {
                const updateData = dataToSubmit as UpdateUserDTO;
                if (removeAvatar) {
                    updateData.remove_avatar = true;
                }
                dataToSubmit = updateData;
            }

            let response;
            if (initialData) {
                response = await updateUserAction(initialData.id, dataToSubmit as UpdateUserDTO);
            } else {
                response = await createUserAction(dataToSubmit as CreateUserDTO);
            }

            if (response.success) {
                if (!initialData && response.password) {
                    setSuccessMessage(`User created successfully. Password: ${response.password}`);
                    // Don't close modal immediately if we need to show password
                } else {
                    if (onSuccess) onSuccess();
                }
                return true;
            }
            return false;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error saving user';
            setError(errorMessage);
            return false;
        } finally {
            setLoading(false);
        }
    };

    return {
        formData,
        avatarFile,
        loading,
        error,
        successMessage,
        handleChange,
        handleFileChange,
        submit,
        setError,
        setSuccessMessage
    };
};
