import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, within } from '@testing-library/react';
import { UserSelect, UserAvatar, UnassignedAvatar } from './user-select';

const users = [
    { id: 1, name: 'ana lopez', email: 'ana@test.com', avatar_url: 'avatars/ana.png' },
    { id: 2, name: 'Bruno Diaz', email: 'bruno@test.com' },
];

const many = Array.from({ length: 9 }, (_, i) => ({
    id: i + 1,
    name: 'Persona ' + (i + 1),
    email: 'p' + (i + 1) + '@test.com',
}));

const trigger = () => screen.getByRole('button', { name: 'Asignado a' });
const open = () => fireEvent.click(trigger());

describe('UserAvatar', () => {
    it('resuelve una ruta relativa contra S3', () => {
        render(<UserAvatar name="Ana" avatarUrl="avatars/ana.png" />);

        expect(screen.getByRole('img')).toHaveAttribute('src', expect.stringContaining('/avatars/ana.png'));
    });

    it('deja intacta una url absoluta', () => {
        render(<UserAvatar name="Ana" avatarUrl="https://cdn.test/ana.png" />);

        expect(screen.getByRole('img')).toHaveAttribute('src', 'https://cdn.test/ana.png');
    });

    it('cae a las iniciales del nombre cuando no hay avatar', () => {
        render(<UserAvatar name="ana maria lopez" />);

        expect(screen.queryByRole('img')).toBeNull();
        expect(screen.getByTitle('ana maria lopez')).toHaveTextContent('AL');
    });

    it('no expone nombre ni titulo cuando es decorativo', () => {
        render(<UserAvatar name="Ana" avatarUrl="https://cdn.test/ana.png" decorative />);

        const img = document.querySelector('img') as HTMLImageElement;
        expect(img).toHaveAttribute('alt', '');
        expect(img).not.toHaveAttribute('title');
    });
});

describe('UnassignedAvatar', () => {
    it('marca el hueco con el titulo por defecto', () => {
        render(<UnassignedAvatar />);

        expect(screen.getByTitle('Sin asignar')).toBeInTheDocument();
    });
});

describe('UserSelect', () => {
    it('muestra el vacio y no abre la lista hasta hacer click', () => {
        render(<UserSelect value={null} options={users} onChange={vi.fn()} />);

        expect(trigger()).toHaveTextContent('Sin asignar');
        expect(screen.queryByRole('listbox')).toBeNull();
    });

    it('pinta el avatar de cada opcion al abrir la lista', () => {
        render(<UserSelect value={null} options={users} onChange={vi.fn()} showEmail />);
        open();

        const conAvatar = screen.getByRole('option', { name: /ana lopez/ });
        expect(conAvatar.querySelector('img')).toHaveAttribute('src', expect.stringContaining('/avatars/ana.png'));
        expect(conAvatar).toHaveTextContent('ana@test.com');

        const sinAvatar = screen.getByRole('option', { name: /Bruno Diaz/ });
        expect(sinAvatar.querySelector('img')).toBeNull();
        expect(within(sinAvatar).getByText('BD')).toBeInTheDocument();
    });

    it('muestra el avatar del usuario elegido en el disparador', () => {
        render(<UserSelect value={1} options={users} onChange={vi.fn()} />);

        expect(trigger()).toHaveTextContent('ana lopez');
        expect(trigger().querySelector('img')).toHaveAttribute('src', expect.stringContaining('/avatars/ana.png'));
    });

    it('usa el nombre y avatar de respaldo cuando el usuario no esta en la lista', () => {
        render(
            <UserSelect
                value={77}
                options={users}
                onChange={vi.fn()}
                fallbackName="Externo"
                fallbackAvatarUrl="https://cdn.test/ext.png"
            />
        );

        expect(trigger()).toHaveTextContent('Externo');
        expect(trigger().querySelector('img')).toHaveAttribute('src', 'https://cdn.test/ext.png');
    });

    it('avisa del usuario elegido y cierra la lista', () => {
        const onChange = vi.fn();
        render(<UserSelect value={null} options={users} onChange={onChange} />);
        open();

        fireEvent.click(screen.getByRole('option', { name: /Bruno Diaz/ }));

        expect(onChange).toHaveBeenCalledWith(2);
        expect(screen.queryByRole('listbox')).toBeNull();
    });

    it('avisa con null al elegir el vacio', () => {
        const onChange = vi.fn();
        render(<UserSelect value={1} options={users} onChange={onChange} />);
        open();

        fireEvent.click(screen.getByRole('option', { name: 'Sin asignar' }));

        expect(onChange).toHaveBeenCalledWith(null);
    });

    it('no avisa si se vuelve a elegir el valor actual', () => {
        const onChange = vi.fn();
        render(<UserSelect value={1} options={users} onChange={onChange} />);
        open();

        fireEvent.click(screen.getByRole('option', { name: /ana lopez/ }));

        expect(onChange).not.toHaveBeenCalled();
    });

    it('no ofrece buscador con pocas opciones', () => {
        render(<UserSelect value={null} options={users} onChange={vi.fn()} />);
        open();

        expect(screen.queryByPlaceholderText('Buscar persona...')).toBeNull();
    });

    it('filtra por nombre y por correo cuando hay muchas opciones', () => {
        render(<UserSelect value={null} options={many} onChange={vi.fn()} />);
        open();

        const search = screen.getByPlaceholderText('Buscar persona...');
        fireEvent.change(search, { target: { value: 'persona 3' } });
        expect(screen.getAllByRole('option')).toHaveLength(2);

        fireEvent.change(search, { target: { value: 'p7@test.com' } });
        expect(screen.getByRole('option', { name: /Persona 7/ })).toBeInTheDocument();

        fireEvent.change(search, { target: { value: 'nadie' } });
        expect(screen.getByText('Sin resultados')).toBeInTheDocument();
    });

    it('olvida la busqueda al cerrar y volver a abrir', () => {
        render(<UserSelect value={null} options={many} onChange={vi.fn()} />);
        open();
        fireEvent.change(screen.getByPlaceholderText('Buscar persona...'), { target: { value: 'persona 3' } });
        open();
        open();

        expect(screen.getByPlaceholderText('Buscar persona...')).toHaveValue('');
    });

    it('cierra al hacer click afuera y con Escape', () => {
        render(<UserSelect value={null} options={users} onChange={vi.fn()} />);

        open();
        fireEvent.mouseDown(document.body);
        expect(screen.queryByRole('listbox')).toBeNull();

        open();
        fireEvent.keyDown(document, { key: 'Escape' });
        expect(screen.queryByRole('listbox')).toBeNull();
    });

    it('no abre la lista cuando esta deshabilitado', () => {
        render(<UserSelect value={null} options={users} onChange={vi.fn()} disabled />);

        expect(trigger()).toBeDisabled();
        open();
        expect(screen.queryByRole('listbox')).toBeNull();
    });
});
