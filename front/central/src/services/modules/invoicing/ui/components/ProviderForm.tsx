/**
 * Formulario para crear/editar proveedores de facturación
 */

'use client';

import { useState, useEffect } from 'react';
import { FormModal } from '@/shared/ui/form-modal';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Label } from '@/shared/ui/label';
import { Checkbox } from '@/shared/ui/checkbox';
import { useToast } from '@/shared/providers/toast-provider';
import {
  createProviderAction,
  updateProviderAction,
  getProviderTypesAction,
  testProviderConnectionAction,
} from '../../infra/actions';
import type { InvoicingProvider, InvoicingProviderType, CreateProviderDTO, UpdateProviderDTO } from '../../domain/types';

interface ProviderFormProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  provider?: InvoicingProvider;
  businessId: number;
}

export function ProviderForm({ isOpen, onClose, onSuccess, provider, businessId }: ProviderFormProps) {
  const { showToast } = useToast();
  const [loading, setLoading] = useState(false);
  const [testLoading, setTestLoading] = useState(false);
  const [providerTypes, setProviderTypes] = useState<InvoicingProviderType[]>([]);
  const [formData, setFormData] = useState({
    name: '',
    provider_type_code: '',
    description: '',
    is_active: true,
    is_default: false,
    referer: '',
    branch_code: '',
    api_key: '',
    api_secret: '',
  });

  useEffect(() => {
    if (isOpen) {
      loadProviderTypes();
      if (provider) {
        setFormData({
          name: provider.name,
          provider_type_code: provider.provider_type_code,
          description: provider.description || '',
          is_active: provider.is_active,
          is_default: provider.is_default,
          referer: provider.config?.referer || '',
          branch_code: provider.config?.branch_code || '',
          api_key: provider.credentials?.api_key || '',
          api_secret: provider.credentials?.api_secret || '',
        });
      } else {
        resetForm();
      }
    }
  }, [isOpen, provider]);

  const loadProviderTypes = async () => {
    try {
      const types = await getProviderTypesAction();
      setProviderTypes(types);
      if (!provider && types.length > 0) {
        setFormData(prev => ({ ...prev, provider_type_code: types[0].code }));
      }
    } catch (error: any) {
      showToast('Error al cargar tipos de proveedores', 'error');
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      provider_type_code: providerTypes[0]?.code || '',
      description: '',
      is_active: true,
      is_default: false,
      referer: '',
      branch_code: '',
      api_key: '',
      api_secret: '',
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const data: CreateProviderDTO | UpdateProviderDTO = {
        name: formData.name,
        description: formData.description || undefined,
        config: {
          referer: formData.referer,
          branch_code: formData.branch_code,
        },
        credentials: {
          api_key: formData.api_key,
          api_secret: formData.api_secret,
        },
        is_active: formData.is_active,
        is_default: formData.is_default,
      };

      if (provider) {
        await updateProviderAction(provider.id, data);
        showToast('Proveedor actualizado exitosamente', 'success');
      } else {
        await createProviderAction({
          ...data,
          business_id: businessId,
          provider_type_code: formData.provider_type_code,
        } as CreateProviderDTO);
        showToast('Proveedor creado exitosamente', 'success');
      }

      onSuccess();
      onClose();
    } catch (error: any) {
      showToast('Error: ' + error.message, 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = async () => {
    if (!provider) {
      showToast('Primero debes guardar el proveedor', 'warning');
      return;
    }

    setTestLoading(true);
    try {
      const result = await testProviderConnectionAction(provider.id);
      if (result.success) {
        showToast('✅ Conexión exitosa', 'success');
      } else {
        showToast('❌ ' + result.message, 'error');
      }
    } catch (error: any) {
      showToast('Error al probar conexión: ' + error.message, 'error');
    } finally {
      setTestLoading(false);
    }
  };

  return (
    <FormModal
      isOpen={isOpen}
      onClose={onClose}
      title={provider ? 'Editar Proveedor' : 'Nuevo Proveedor'}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <Label htmlFor="name">Nombre *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Ej: Softpymes - Mi Negocio"
            required
          />
        </div>

        {!provider && (
          <div>
            <Label htmlFor="provider_type">Tipo de Proveedor *</Label>
            <Select
              id="provider_type"
              value={formData.provider_type_code}
              onChange={(e) => setFormData({ ...formData, provider_type_code: e.target.value })}
              options={providerTypes.map((type) => ({ value: type.code, label: type.name }))}
              required
            />
          </div>
        )}

        <div>
          <Label htmlFor="description">Descripción</Label>
          <Input
            id="description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Descripción opcional"
          />
        </div>

        <div className="border-t pt-4">
          <h4 className="font-semibold mb-3">Configuración de Softpymes</h4>

          <div className="space-y-3">
            <div>
              <Label htmlFor="referer">NIT / Referer *</Label>
              <Input
                id="referer"
                value={formData.referer}
                onChange={(e) => setFormData({ ...formData, referer: e.target.value })}
                placeholder="900123456"
                required
              />
            </div>

            <div>
              <Label htmlFor="branch_code">Código de Sucursal *</Label>
              <Input
                id="branch_code"
                value={formData.branch_code}
                onChange={(e) => setFormData({ ...formData, branch_code: e.target.value })}
                placeholder="001"
                required
              />
            </div>
          </div>
        </div>

        <div className="border-t pt-4">
          <h4 className="font-semibold mb-3">Credenciales de API</h4>

          <div className="space-y-3">
            <div>
              <Label htmlFor="api_key">API Key *</Label>
              <Input
                id="api_key"
                type="password"
                value={formData.api_key}
                onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                placeholder="Tu API Key de Softpymes"
                required
              />
            </div>

            <div>
              <Label htmlFor="api_secret">API Secret *</Label>
              <Input
                id="api_secret"
                type="password"
                value={formData.api_secret}
                onChange={(e) => setFormData({ ...formData, api_secret: e.target.value })}
                placeholder="Tu API Secret de Softpymes"
                required
              />
            </div>
          </div>
        </div>

        <div className="border-t pt-4 space-y-2">
          <label className="flex items-center gap-2">
            <Checkbox
              id="is_active"
              checked={formData.is_active}
              onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
            />
            <span className="text-sm text-gray-700">Activo</span>
          </label>

          <label className="flex items-center gap-2">
            <Checkbox
              id="is_default"
              checked={formData.is_default}
              onChange={(e) => setFormData({ ...formData, is_default: e.target.checked })}
            />
            <span className="text-sm text-gray-700">Proveedor por defecto</span>
          </label>
        </div>

        <div className="flex gap-2 pt-4">
          {provider && (
            <button
              type="button"
              onClick={handleTestConnection}
              disabled={testLoading}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {testLoading ? 'Probando...' : 'Probar Conexión'}
            </button>
          )}

          <button
            type="submit"
            disabled={loading}
            className="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
          >
            {loading ? 'Guardando...' : provider ? 'Actualizar' : 'Crear Proveedor'}
          </button>

          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
          >
            Cancelar
          </button>
        </div>
      </form>
    </FormModal>
  );
}
