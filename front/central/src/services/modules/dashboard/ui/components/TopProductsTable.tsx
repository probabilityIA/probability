'use client';

import * as React from 'react';
import {
    flexRender,
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    useReactTable,
    SortingState,
    ColumnFiltersState,
    ColumnDef,
} from '@tanstack/react-table';
import {
    ArrowUpDown,
    Search,
    ChevronLeft,
    ChevronRight
} from 'lucide-react';

// Tipos adaptados para la tabla de productos (coincide con productsTableData en Dashboard)
export type TopProductRow = {
    name: string;
    sku: string;
    units: number;
    price: number | null;
    totalEarned: number;
};

export const columns: ColumnDef<TopProductRow>[] = [
    {
        accessorKey: 'name',
        header: ({ column }) => {
            return (
                <button
                    className="flex items-center hover:text-gray-900 transition-colors"
                    onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
                >
                    Producto
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </button>
            );
        },
        cell: ({ row }) => (
            <div>
                <div className="font-medium text-gray-900">{row.getValue('name')}</div>
                <div className="text-xs text-gray-500">{row.original.sku}</div>
            </div>
        ),
    },
    {
        accessorKey: 'units',
        header: ({ column }) => (
            <div className="text-right">
                <button
                    className="inline-flex items-center hover:text-gray-900 transition-colors"
                    onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
                >
                    Unidades
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </button>
            </div>
        ),
        cell: ({ row }) => <div className="text-right font-medium text-gray-700">{row.getValue<number>('units').toLocaleString()}</div>,
    },
    {
        accessorKey: 'price',
        header: ({ column }) => (
            <div className="text-right">
                <button
                    className="inline-flex items-center hover:text-gray-900 transition-colors"
                    onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
                >
                    Precio Unitario
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </button>
            </div>
        ),
        cell: ({ row }) => {
            const price = row.getValue<number | null>('price');
            const formatted = price != null
                ? new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(price)
                : '-';
            return <div className="text-right font-medium text-gray-900">{formatted}</div>;
        },
    },
];

export function TopProductsTable({ data }: { data: TopProductRow[] }) {
    const [sorting, setSorting] = React.useState<SortingState>([]);
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
    const [globalFilter, setGlobalFilter] = React.useState('');

    const table = useReactTable({
        data,
        columns,
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        onGlobalFilterChange: setGlobalFilter,
        state: {
            sorting,
            columnFilters,
            globalFilter,
        },
    });

    return (
        <div className="w-full space-y-4">


            <div className="rounded-xl border border-gray-100 bg-white shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full caption-bottom text-sm text-left">
                        <thead className="[&_tr]:border-b bg-gray-50/50">
                            {table.getHeaderGroups().map((headerGroup) => (
                                <tr key={headerGroup.id} className="border-b transition-colors data-[state=selected]:bg-gray-100/50">
                                    {headerGroup.headers.map((header) => {
                                        return (
                                            <th key={header.id} className="h-8 px-2 text-left align-middle font-medium text-gray-500 text-xs [&:has([role=checkbox])]:pr-0">
                                                {header.isPlaceholder
                                                    ? null
                                                    : flexRender(
                                                        header.column.columnDef.header,
                                                        header.getContext()
                                                    )}
                                            </th>
                                        );
                                    })}
                                </tr>
                            ))}
                        </thead>
                        <tbody className="[&_tr:last-child]:border-0">
                            {table.getRowModel().rows?.length ? (
                                table.getRowModel().rows.map((row) => (
                                    <tr
                                        key={row.id}
                                        data-state={row.getIsSelected() && 'selected'}
                                        className="border-b transition-colors hover:bg-gray-50/50 data-[state=selected]:bg-gray-100"
                                    >
                                        {row.getVisibleCells().map((cell) => (
                                            <td key={cell.id} className="px-2 py-1.5 align-middle [&:has([role=checkbox])]:pr-0 text-xs">
                                                {flexRender(
                                                    cell.column.columnDef.cell,
                                                    cell.getContext()
                                                )}
                                            </td>
                                        ))}
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td
                                        colSpan={columns.length}
                                        className="h-24 text-center text-gray-500"
                                    >
                                        No se encontraron resultados.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>


        </div>
    );
}
