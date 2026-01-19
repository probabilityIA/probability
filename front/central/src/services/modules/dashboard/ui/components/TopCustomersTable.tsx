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
    VisibilityState,
    ColumnDef,
} from '@tanstack/react-table';
import {
    ArrowUpDown,
    ChevronDown,
    MoreHorizontal,
    Search,
    Columns,
    ChevronLeft,
    ChevronRight
} from 'lucide-react';
import { Button } from '@/shared/ui/shadcn'; // Usar nuestro Button existente

// Tipos
export type TopCustomer = {
    customer_name: string;
    customer_email: string;
    order_count: number;
    phone?: string; // Agregamos phone opcional
};

export const columns: ColumnDef<TopCustomer>[] = [
    {
        id: 'customer_name',
        accessorFn: (row) => `${row.customer_name} ${row.customer_email}`,
        header: ({ column }) => {
            return (
                <button
                    className="flex w-full items-center hover:text-gray-900 transition-colors"
                    onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
                >
                    Cliente
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </button>
            );
        },
        cell: ({ row }) => (
            <div>
                <div className="font-medium text-gray-900">{row.original.customer_name}</div>
                <div className="text-xs text-gray-500 lowercase">{row.original.customer_email}</div>
            </div>
        ),
    },
    {
        accessorKey: 'order_count',
        header: ({ column }) => (
            <div className="text-right w-[80px] ml-auto">
                <button
                    className="inline-flex items-center hover:text-gray-900 transition-colors"
                    onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
                >
                    Compras
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </button>
            </div>
        ),
        cell: ({ row }) => {
            return <div className="text-right font-medium text-gray-900 pr-4">{row.getValue('order_count')}</div>;
        },
    },
];

export function TopCustomersTable({ data }: { data: TopCustomer[] }) {
    const [sorting, setSorting] = React.useState<SortingState>([]);
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
    const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
    const [rowSelection, setRowSelection] = React.useState({});
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
        onColumnVisibilityChange: setColumnVisibility,
        onRowSelectionChange: setRowSelection,
        onGlobalFilterChange: setGlobalFilter,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
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
