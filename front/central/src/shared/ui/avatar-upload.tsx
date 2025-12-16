'use client';

import { useRef, useState, useEffect, useCallback } from 'react';
import { CameraIcon } from '@heroicons/react/24/outline';

interface AvatarUploadProps {
  currentAvatarUrl?: string | null;
  onFileSelect: (file: File | null) => void;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export function AvatarUpload({
  currentAvatarUrl,
  onFileSelect,
  size = 'md',
  className = '',
}: AvatarUploadProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [localPreview, setLocalPreview] = useState<string | null>(null);
  const [hasError, setHasError] = useState(false);

  // Tamaños del avatar
  const sizeClasses = {
    sm: 'w-16 h-16',
    md: 'w-24 h-24',
    lg: 'w-32 h-32',
  };

  const iconSizes = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
  };

  // Limpiar blob URL al desmontar
  useEffect(() => {
    return () => {
      if (localPreview) {
        URL.revokeObjectURL(localPreview);
      }
    };
  }, [localPreview]);

  const handleClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    
    if (file) {
      // Limpiar preview anterior
      if (localPreview) {
        URL.revokeObjectURL(localPreview);
      }
      const url = URL.createObjectURL(file);
      setLocalPreview(url);
      setHasError(false);
      onFileSelect(file);
    }
  }, [localPreview, onFileSelect]);

  const handleRemove = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    if (localPreview) {
      URL.revokeObjectURL(localPreview);
    }
    setLocalPreview(null);
    setHasError(false);
    onFileSelect(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  }, [localPreview, onFileSelect]);

  const handleImageError = useCallback(() => {
    console.error('AvatarUpload - Failed to load image');
    setHasError(true);
  }, []);

  const handleImageLoad = useCallback(() => {
    setHasError(false);
  }, []);

  // Determinar qué imagen mostrar:
  // Si hay un preview local (archivo seleccionado), usarlo
  // Si no, usar la URL del avatar actual
  const displayUrl = localPreview || currentAvatarUrl;
  const showImage = displayUrl && !hasError;

  return (
    <div className={`flex flex-col items-center gap-2 ${className}`}>
      <div className="relative group">
        <button
          type="button"
          onClick={handleClick}
          className={`
            ${sizeClasses[size]}
            rounded-full overflow-hidden
            border-2 border-gray-300
            bg-gray-100
            flex items-center justify-center
            cursor-pointer
            transition-all duration-200
            hover:border-blue-500 hover:shadow-md
            focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
            relative
          `}
          aria-label="Cambiar foto de perfil"
        >
          {showImage ? (
            <img
              src={displayUrl}
              alt="Avatar"
              className="w-full h-full object-cover"
              onError={handleImageError}
              onLoad={handleImageLoad}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-200 to-gray-300">
              <CameraIcon className={`${iconSizes[size]} text-gray-500`} />
            </div>
          )}
        </button>

        {/* Botón para eliminar si hay imagen */}
        {showImage && (
          <button
            type="button"
            onClick={handleRemove}
            className="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center hover:bg-red-600 transition-colors shadow-md z-10"
            aria-label="Eliminar foto"
            title="Eliminar foto"
          >
            <span className="text-xs font-bold">×</span>
          </button>
        )}
      </div>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={handleFileChange}
      />

      <p className="text-xs text-gray-500 text-center">
        Haz clic para cambiar
      </p>
    </div>
  );
}
