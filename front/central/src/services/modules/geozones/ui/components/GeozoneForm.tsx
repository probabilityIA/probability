'use client';

import { useState } from 'react';
import dynamic from 'next/dynamic';
import { Geozone, GeozoneType, CreateGeozoneDTO } from '../../domain/types';
import { createGeozoneAction } from '../../infra/actions';
import { Button, Alert } from '@/shared/ui';

const DrawMap = dynamic(() => import('./GeozoneDrawMap'), { ssr: false });
import { pointsToPolygon } from './GeozoneDrawMap';

interface GeozoneFormProps {
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
    contextLayers?: Geozone[];
}

type Mode = 'draw' | 'paste';

export default function GeozoneForm({ onSuccess, onCancel, businessId, contextLayers }: GeozoneFormProps) {
    const [mode, setMode] = useState<Mode>('draw');
    const [name, setName] = useState('');
    const [type, setType] = useState<GeozoneType>('custom');
    const [code, setCode] = useState('');
    const [parentIdRaw, setParentIdRaw] = useState('');
    const [points, setPoints] = useState<Array<[number, number]>>([]);
    const [pasted, setPasted] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);

        if (!name.trim()) { setError('El nombre es obligatorio'); return; }

        let geometry;
        if (mode === 'draw') {
            geometry = pointsToPolygon(points);
            if (!geometry) { setError('Dibuja al menos 3 vertices en el mapa'); return; }
        } else {
            try {
                const parsed = JSON.parse(pasted);
                geometry = parsed.geometry || parsed;
                if (!geometry?.type || !geometry?.coordinates) {
                    throw new Error('GeoJSON invalido (falta type o coordinates)');
                }
            } catch (err: any) {
                setError('GeoJSON invalido: ' + err.message);
                return;
            }
        }

        const parent_id = parentIdRaw.trim() ? Number(parentIdRaw.trim()) : null;
        if (parentIdRaw.trim() && Number.isNaN(parent_id as number)) {
            setError('parent_id debe ser un numero');
            return;
        }

        const dto: CreateGeozoneDTO = {
            type,
            name: name.trim(),
            code: code.trim() || null,
            parent_id,
            geometry,
            properties: { source: 'manual' },
        };

        setSubmitting(true);
        try {
            await createGeozoneAction(dto, businessId);
            onSuccess();
        } catch (err: any) {
            setError(err.message || 'Error al crear la geozona');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Nombre *</label>
                    <input
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="Ej: Zona Norte Bogota"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    />
                </div>
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Tipo</label>
                    <select
                        value={type}
                        onChange={(e) => setType(e.target.value as GeozoneType)}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500"
                    >
                        <option value="custom">Personalizada</option>
                        <option value="neighborhood">Barrio</option>
                        <option value="locality">Localidad</option>
                        <option value="city">Municipio</option>
                        <option value="state">Departamento</option>
                        <option value="country">Pais</option>
                    </select>
                </div>
                <div>
                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Codigo (opcional)</label>
                    <input
                        value={code}
                        onChange={(e) => setCode(e.target.value)}
                        placeholder="Ej: ZN-001"
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500"
                    />
                </div>
            </div>

            <div>
                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">
                    Geozona padre (id, opcional)
                </label>
                <input
                    value={parentIdRaw}
                    onChange={(e) => setParentIdRaw(e.target.value)}
                    placeholder="Ej: 183 (Bogota DC ciudad)"
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500"
                />
            </div>

            <div className="flex gap-1 p-1 bg-gray-100 dark:bg-gray-700 rounded-lg">
                <button
                    type="button"
                    onClick={() => setMode('draw')}
                    className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-all ${
                        mode === 'draw'
                            ? 'bg-white dark:bg-gray-800 text-purple-700 dark:text-purple-300 shadow'
                            : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
                    }`}
                >
                    🎨 Dibujar en el mapa
                </button>
                <button
                    type="button"
                    onClick={() => setMode('paste')}
                    className={`flex-1 px-4 py-2 text-sm font-medium rounded-md transition-all ${
                        mode === 'paste'
                            ? 'bg-white dark:bg-gray-800 text-purple-700 dark:text-purple-300 shadow'
                            : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white'
                    }`}
                >
                    📋 Pegar GeoJSON
                </button>
            </div>

            {mode === 'draw' ? (
                <DrawMap points={points} onChange={setPoints} contextLayers={contextLayers} height="380px" />
            ) : (
                <textarea
                    value={pasted}
                    onChange={(e) => setPasted(e.target.value)}
                    placeholder='{"type":"Polygon","coordinates":[[[-74.10,4.70],[-74.00,4.70],[-74.00,4.80],[-74.10,4.80],[-74.10,4.70]]]}'
                    rows={10}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-xs font-mono bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500"
                />
            )}

            <div className="flex justify-end gap-2 pt-2 border-t border-gray-200 dark:border-gray-700">
                <Button variant="secondary" onClick={onCancel} type="button">Cancelar</Button>
                <Button variant="purple" type="submit" disabled={submitting}>
                    {submitting ? 'Creando...' : 'Crear geozona'}
                </Button>
            </div>
        </form>
    );
}
