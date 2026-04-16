'use client';

import React, { createContext, useContext, useState, ReactNode } from 'react';

interface NavbarContextType {
    actionButtons: React.ReactNode | null;
    setActionButtons: (buttons: React.ReactNode | null) => void;
}

const NavbarContext = createContext<NavbarContextType | undefined>(undefined);

export const NavbarProvider = ({ children }: { children: ReactNode }) => {
    const [actionButtons, setActionButtons] = useState<React.ReactNode | null>(null);

    return (
        <NavbarContext.Provider value={{ actionButtons, setActionButtons }}>
            {children}
        </NavbarContext.Provider>
    );
};

export const useNavbarActions = () => {
    const context = useContext(NavbarContext);
    if (!context) {
        throw new Error('useNavbarActions must be used within NavbarProvider');
    }
    return context;
};
