'use client';

import { useRef, useState, useEffect } from 'react';
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
  const [preview, setPreview] = useState<string | null>(null);

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

  // Limpiar preview cuando se desmonta
  useEffect(() => {
    return () => {
      if (preview && preview.startsWith('blob:')) {
        URL.revokeObjectURL(preview);
      }
    };
  }, [preview]);

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    
    if (file) {
      // Crear preview de la nueva imagen
      const url = URL.createObjectURL(file);
      setPreview(url);
      onFileSelect(file);
    } else {
      setPreview(null);
      onFileSelect(null);
    }
  };

  const handleRemove = (e: React.MouseEvent) => {
    e.stopPropagation();
    // Limpiar preview si es un blob URL
    if (preview && preview.startsWith('blob:')) {
      URL.revokeObjectURL(preview);
    }
    setPreview(null);
    onFileSelect(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // Determinar qué imagen mostrar
  const imageUrl = preview || currentAvatarUrl || null;

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
          {imageUrl ? (
            <img
              src={imageUrl}
              alt="Avatar"
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-200 to-gray-300">
              <CameraIcon className={`${iconSizes[size]} text-gray-500`} />
            </div>
          )}
          
          {/* Overlay al hacer hover */}
          <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all duration-200 flex items-center justify-center">
            <CameraIcon className={`${iconSizes[size]} text-white opacity-0 group-hover:opacity-100 transition-opacity duration-200`} />
          </div>
        </button>

        {/* Botón para eliminar si hay imagen */}
        {imageUrl && (
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
