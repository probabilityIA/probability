'use client';

import { createContext, useContext, useState, ReactNode, useRef } from 'react';

interface SidebarContextType {
  primaryExpanded: boolean;
  setPrimaryExpanded: (expanded: boolean) => void;
  secondaryExpanded: boolean;
  setSecondaryExpanded: (expanded: boolean) => void;
  keepExpanded: () => void;
  releaseExpanded: () => void;
  requestExpand: () => void;
  requestCollapse: (hasSecondarySidebar?: boolean) => void;
  setHasSecondarySidebar: (has: boolean) => void;
  requestSecondaryExpand: () => void;
  requestSecondaryCollapse: () => void;
}

const SidebarContext = createContext<SidebarContextType | undefined>(undefined);

export function SidebarProvider({ children }: { children: ReactNode }) {
  const [primaryExpanded, setPrimaryExpanded] = useState(false);
  const [secondaryExpanded, setSecondaryExpanded] = useState(false);
  const collapseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const secondaryCollapseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const hoverRef = useRef(false); // Indica si algo está manteniendo el sidebar abierto
  const secondaryHoverRef = useRef(false); // Indica si algo está manteniendo el sidebar secundario abierto
  const hasSecondarySidebarRef = useRef(false); // Indica si hay sidebar secundario visible

  // Función para cancelar cualquier timeout pendiente
  const cancelCollapse = () => {
    if (collapseTimeoutRef.current) {
      clearTimeout(collapseTimeoutRef.current);
      collapseTimeoutRef.current = null;
    }
  };

  // Solicitar expansión (desde sidebar primario o secundario)
  const requestExpand = () => {
    cancelCollapse();
    hoverRef.current = true;
    setPrimaryExpanded(true);
  };

  // Solicitar colapso (con delay) - solo desde sidebar primario
  const requestCollapse = (hasSecondarySidebar?: boolean) => {
    cancelCollapse();
    
    const hasSecondary = hasSecondarySidebar ?? hasSecondarySidebarRef.current;
    
    // Si NO hay sidebar secundario, resetear hoverRef inmediatamente y cerrar rápido
    if (!hasSecondary) {
      hoverRef.current = false;
      collapseTimeoutRef.current = setTimeout(() => {
        setPrimaryExpanded(false);
      }, 200);
      return;
    }
    
    // Si hay sidebar secundario, verificar si el secundario también está siendo abandonado
    // Si el secundario no está en hover, cerrar ambos
    if (!secondaryHoverRef.current) {
      hoverRef.current = false;
      collapseTimeoutRef.current = setTimeout(() => {
        setPrimaryExpanded(false);
      }, 200);
      return;
    }
    
    // Si el secundario está en hover, dar un delay corto para permitir movimiento
    hoverRef.current = false;
    collapseTimeoutRef.current = setTimeout(() => {
      // Solo colapsar si nada está manteniendo el hover
      if (!hoverRef.current && !secondaryHoverRef.current) {
        setPrimaryExpanded(false);
      }
    }, 300);
  };

  // Establecer si hay sidebar secundario visible
  const setHasSecondarySidebar = (has: boolean) => {
    hasSecondarySidebarRef.current = has;
  };

  // Mantener expandido (desde sidebar secundario)
  const keepExpanded = () => {
    cancelCollapse();
    hoverRef.current = true;
    if (!primaryExpanded) {
      setPrimaryExpanded(true);
    }
  };

  // Liberar expansión (desde sidebar secundario)
  const releaseExpanded = () => {
    hoverRef.current = false;
    cancelCollapse();
    collapseTimeoutRef.current = setTimeout(() => {
      if (!hoverRef.current) {
        setPrimaryExpanded(false);
      }
    }, 500);
  };

  // Solicitar expansión del sidebar secundario (independiente del principal)
  const requestSecondaryExpand = () => {
    if (secondaryCollapseTimeoutRef.current) {
      clearTimeout(secondaryCollapseTimeoutRef.current);
      secondaryCollapseTimeoutRef.current = null;
    }
    secondaryHoverRef.current = true;
    setSecondaryExpanded(true);
  };

  // Solicitar colapso del sidebar secundario
  const requestSecondaryCollapse = () => {
    if (secondaryCollapseTimeoutRef.current) {
      clearTimeout(secondaryCollapseTimeoutRef.current);
      secondaryCollapseTimeoutRef.current = null;
    }
    secondaryHoverRef.current = false;
    
    // Si el principal tampoco está en hover, cerrar ambos
    if (!hoverRef.current) {
      // Cerrar el secundario
      secondaryCollapseTimeoutRef.current = setTimeout(() => {
        if (!secondaryHoverRef.current) {
          setSecondaryExpanded(false);
        }
      }, 200);
      
      // También cerrar el principal si no está en hover
      cancelCollapse();
      collapseTimeoutRef.current = setTimeout(() => {
        if (!hoverRef.current) {
          setPrimaryExpanded(false);
        }
      }, 200);
      return;
    }
    
    // Si el principal está en hover, solo cerrar el secundario
    secondaryCollapseTimeoutRef.current = setTimeout(() => {
      if (!secondaryHoverRef.current) {
        setSecondaryExpanded(false);
      }
    }, 200);
  };

  return (
    <SidebarContext.Provider value={{ 
      primaryExpanded, 
      setPrimaryExpanded,
      secondaryExpanded,
      setSecondaryExpanded,
      keepExpanded,
      releaseExpanded,
      requestExpand,
      requestCollapse,
      setHasSecondarySidebar,
      requestSecondaryExpand,
      requestSecondaryCollapse
    }}>
      {children}
    </SidebarContext.Provider>
  );
}

export function useSidebar() {
  const context = useContext(SidebarContext);
  if (context === undefined) {
    throw new Error('useSidebar must be used within a SidebarProvider');
  }
  return context;
}
