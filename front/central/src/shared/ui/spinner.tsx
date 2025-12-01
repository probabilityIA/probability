/**
 * Componente Spinner/Loader reutilizable
 * Usa los colores dinámicos del negocio
 */

'use client';

interface SpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  color?: 'primary' | 'secondary' | 'tertiary' | 'white';
  text?: string;
}

const sizeClasses = {
  sm: 'w-4 h-4 border-2',
  md: 'w-8 h-8 border-2',
  lg: 'w-12 h-12 border-3',
  xl: 'w-16 h-16 border-4',
};

export function Spinner({ size = 'md', color = 'primary', text }: SpinnerProps) {
  const colorStyle = color === 'white' 
    ? { borderColor: 'rgba(255, 255, 255, 0.3)', borderTopColor: 'white' }
    : { 
        borderColor: `rgba(0, 0, 0, 0.1)`, 
        borderTopColor: `var(--color-${color})` 
      };

  return (
    <div className="flex flex-col items-center justify-center gap-3">
      <div
        className={`${sizeClasses[size]} rounded-full animate-spin`}
        style={colorStyle}
      />
      {text && (
        <p className="text-sm text-gray-600 animate-pulse">{text}</p>
      )}
    </div>
  );
}

/**
 * Spinner de página completa (overlay)
 */
export function FullPageSpinner({ text = 'Cargando...' }: { text?: string }) {
  return (
    <div className="fixed inset-0 bg-white/80 backdrop-blur-sm z-50 flex items-center justify-center">
      <Spinner size="xl" color="primary" text={text} />
    </div>
  );
}

