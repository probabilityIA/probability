import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { StatusBadge, PriorityBadge, TypeBadge } from './TicketBadges';
import {
    STATUS_META,
    PRIORITY_META,
    TYPE_META,
    TICKET_STATUSES,
    TICKET_PRIORITIES,
    TICKET_TYPES,
} from '../../domain/types';

describe('StatusBadge', () => {
    it('muestra la etiqueta del estado recibido', () => {
        render(<StatusBadge status="open" />);

        expect(screen.getByText('Abierto')).toBeInTheDocument();
    });

    it('cambia la etiqueta al cambiar el estado', () => {
        const { rerender } = render(<StatusBadge status="open" />);
        expect(screen.getByText('Abierto')).toBeInTheDocument();

        rerender(<StatusBadge status="blocked" />);

        expect(screen.queryByText('Abierto')).not.toBeInTheDocument();
        expect(screen.getByText('Bloqueado')).toBeInTheDocument();
    });

    it('aplica los colores del estado', () => {
        render(<StatusBadge status="resolved" />);
        const meta = STATUS_META.resolved;

        const badge = screen.getByText(meta.label);
        expect(badge).toHaveClass(...meta.bg.split(' '));
        expect(badge).toHaveClass(...meta.color.split(' '));
        expect(badge).toHaveClass(meta.ring);
    });

    it('renderiza todos los estados declarados', () => {
        TICKET_STATUSES.forEach((status) => {
            const { unmount } = render(<StatusBadge status={status} />);
            expect(screen.getByText(STATUS_META[status].label)).toBeInTheDocument();
            unmount();
        });
    });
});

describe('PriorityBadge', () => {
    it('muestra la etiqueta de la prioridad recibida', () => {
        render(<PriorityBadge priority="critical" />);

        expect(screen.getByText('Critica')).toBeInTheDocument();
    });

    it('distingue prioridad baja de alta', () => {
        const { rerender } = render(<PriorityBadge priority="low" />);
        expect(screen.getByText('Baja')).toBeInTheDocument();

        rerender(<PriorityBadge priority="high" />);

        expect(screen.queryByText('Baja')).not.toBeInTheDocument();
        expect(screen.getByText('Alta')).toBeInTheDocument();
    });

    it('aplica los colores de la prioridad', () => {
        render(<PriorityBadge priority="high" />);
        const meta = PRIORITY_META.high;

        const badge = screen.getByText(meta.label);
        expect(badge).toHaveClass(...meta.bg.split(' '));
        expect(badge).toHaveClass(...meta.color.split(' '));
    });

    it('renderiza todas las prioridades declaradas', () => {
        TICKET_PRIORITIES.forEach((priority) => {
            const { unmount } = render(<PriorityBadge priority={priority} />);
            expect(screen.getByText(PRIORITY_META[priority].label)).toBeInTheDocument();
            unmount();
        });
    });
});

describe('TypeBadge', () => {
    it('muestra icono y etiqueta del tipo recibido', () => {
        render(<TypeBadge type="bug" />);

        expect(screen.getByText('BUG')).toBeInTheDocument();
        expect(screen.getByText('Bug')).toBeInTheDocument();
    });

    it('cambia icono y etiqueta al cambiar el tipo', () => {
        const { rerender } = render(<TypeBadge type="bug" />);

        rerender(<TypeBadge type="feature" />);

        expect(screen.queryByText('BUG')).not.toBeInTheDocument();
        expect(screen.getByText('NEW')).toBeInTheDocument();
        expect(screen.getByText('Nueva funcionalidad')).toBeInTheDocument();
    });

    it('renderiza el tipo con tilde escapada', () => {
        render(<TypeBadge type="integration" />);

        expect(screen.getByText('Integraci\u00f3n')).toBeInTheDocument();
        expect(screen.getByText('INT')).toBeInTheDocument();
    });

    it('renderiza todos los tipos declarados', () => {
        TICKET_TYPES.forEach((type) => {
            const { unmount } = render(<TypeBadge type={type} />);
            expect(screen.getByText(TYPE_META[type].label)).toBeInTheDocument();
            expect(screen.getByText(TYPE_META[type].icon)).toBeInTheDocument();
            unmount();
        });
    });
});
