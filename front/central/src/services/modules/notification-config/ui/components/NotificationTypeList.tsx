"use client";

import { useState, useEffect } from "react";
import { NotificationType } from "../../domain/types";
import { Button } from "@/shared/ui/button";
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
        <Button onClick={onCreate}>+ Nuevo Tipo</Button>
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
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => onEdit(type)}
                      className="mr-2"
                    >
                      Editar
                    </Button>
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => setDeleteModal({ isOpen: true, type })}
                    >
                      Eliminar
                    </Button>
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
