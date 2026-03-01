"use client";

import { useState, useEffect } from "react";
import { NotificationType } from "../../domain/types";
import { useToast } from "@/shared/providers/toast-provider";
import {
  getNotificationTypesAction,
  deleteNotificationTypeAction,
} from "../../infra/actions";
import { ConfirmModal } from "@/shared/ui/confirm-modal";

interface NotificationTypeListProps {
  onEdit: (type: NotificationType) => void;
  onCreate: () => void;
  refreshKey: number;
}

export function NotificationTypeList({
  onEdit,
  onCreate,
  refreshKey,
}: NotificationTypeListProps) {
  const [types, setTypes] = useState<NotificationType[]>([]);
  const [loading, setLoading] = useState(false);
  const [deleteModal, setDeleteModal] = useState<{
    isOpen: boolean;
    type?: NotificationType;
  }>({ isOpen: false });
  const { showToast } = useToast();

  const fetchTypes = async () => {
    setLoading(true);
    try {
      const result = await getNotificationTypesAction();
      if (result.success) {
        setTypes(result.data);
      } else {
        showToast("Error al cargar tipos", "error");
      }
    } catch (error) {
      showToast("Error al cargar tipos", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTypes();
  }, [refreshKey]);

  const handleDelete = async () => {
    if (!deleteModal.type) return;

    try {
      const result = await deleteNotificationTypeAction(deleteModal.type.id);
      if (result.success) {
        showToast("Tipo eliminado exitosamente", "success");
        setDeleteModal({ isOpen: false });
        fetchTypes();
      } else {
        showToast(result.error || "Error al eliminar", "error");
      }
    } catch (error: any) {
      showToast(error.message || "Error al eliminar", "error");
    }
  };

  if (loading) {
    return <div className="text-center py-8">Cargando...</div>;
  }

  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden">
      <div className="p-4 border-b flex justify-between items-center">
        <h2 className="text-lg font-semibold">Tipos de Notificación</h2>
        <button
          type="button"
          onClick={onCreate}
          className="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
          title="Nuevo Tipo"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Nombre
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Código
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Descripción
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {types.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                  No hay tipos de notificación
                </td>
              </tr>
            ) : (
              types.map((type) => (
                <tr key={type.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      {type.icon && (
                        <span className="mr-2 text-gray-400">{type.icon}</span>
                      )}
                      <div className="text-sm font-medium text-gray-900">
                        {type.name}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <code className="text-sm text-gray-600 bg-gray-100 px-2 py-1 rounded">
                      {type.code}
                    </code>
                  </td>
                  <td className="px-6 py-4">
                    <div className="text-sm text-gray-500">
                      {type.description || "-"}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        type.is_active
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {type.is_active ? "Activo" : "Inactivo"}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => onEdit(type)}
                        className="p-1.5 rounded-md bg-amber-50 text-amber-600 hover:bg-amber-100 transition-colors"
                        title="Editar"
                      >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                      </button>
                      <button
                        type="button"
                        onClick={() => setDeleteModal({ isOpen: true, type })}
                        className="p-1.5 rounded-md bg-red-50 text-red-500 hover:bg-red-100 transition-colors"
                        title="Eliminar"
                      >
                        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <ConfirmModal
        isOpen={deleteModal.isOpen}
        onClose={() => setDeleteModal({ isOpen: false })}
        onConfirm={handleDelete}
        title="Eliminar Tipo de Notificación"
        message={`¿Estás seguro de eliminar el tipo "${deleteModal.type?.name}"? Esta acción no se puede deshacer.`}
        confirmText="Eliminar"
        cancelText="Cancelar"
      />
    </div>
  );
}
