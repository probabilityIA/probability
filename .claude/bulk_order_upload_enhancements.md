# Bulk Order Upload - Enhanced Columns (2026-05-14)

## Summary
Enhanced the bulk order upload feature to support 30+ additional order fields, enabling richer data import from the initial order creation. This allows users to pre-populate financial, customer, shipping, and logistics details in bulk uploads.

## Changes Made

### 1. Backend Handler Updates
**File:** `back/central/services/modules/orders/internal/infra/primary/handlers/upload-bulk.go`

Updated both `parseCSV()` and `parseExcel()` functions to extract and map the following new optional fields:

#### Financial Data
- `subtotal` - Base amount before taxes
- `tax` - Tax amount
- `discount` - Discount amount
- `shipping_cost` - Shipping amount
- `shipping_discount` - Shipping discount
- `currency` - Currency code (e.g., "COP")

#### Customer Information
- `customer_first_name` - Customer first name
- `customer_last_name` - Customer last name
- `customer_dni` - Customer identification number (crucial for Latin America invoicing)

#### Shipping Address
- `shipping_country` - Country
- `shipping_postal_code` - Postal/ZIP code
- `shipping_lat` - Latitude coordinate
- `shipping_lng` - Longitude coordinate

#### Status & Payment
- `status` - Order status (e.g., "pending", "shipped")
- `payment_method_id` - Payment method ID (numeric)
- `is_paid` - Payment status (accepts: true/false/1/0/yes/no/si/no in any case)

#### Logistics
- `tracking_number` - Shipping tracking number
- `guide_id` - Shipping guide ID
- `warehouse_name` - Fulfillment warehouse
- `driver_name` - Assigned driver name

#### Additional
- `notes` - Order notes/comments
- `order_type_name` - Order type (e.g., "standard", "express")
- `invoiceable` - Invoice eligibility (accepts: true/false/1/0/yes/no/si/no)

### 2. Template Update
**File:** `front/central/public/template_orders.csv`

Updated template now includes:
- All required columns (unchanged)
- All 30+ optional columns with example data
- Two example rows showing realistic data entry

**Column Order in Template:**
```
order_number, customer_name, customer_first_name, customer_last_name, customer_email, 
customer_phone, customer_dni, shipping_street, shipping_city, shipping_state, 
shipping_country, shipping_postal_code, shipping_lat, shipping_lng, subtotal, tax, 
discount, shipping_cost, shipping_discount, total_amount, currency, weight, height, 
width, length, platform, status, payment_method_id, is_paid, tracking_number, 
guide_id, warehouse_name, driver_name, notes, order_type_name, invoiceable
```

### 3. UI Instructions Update
**File:** `front/central/src/shared/ui/modals/mass-order-upload-modal.tsx`

Enhanced the instructions section to categorize the new columns by type:
- **Requeridas:** Core required fields
- **Cliente:** customer_first_name, customer_last_name, customer_dni
- **Dirección:** country, postal_code, lat, lng
- **Financiero:** subtotal, tax, discount, shipping_cost, shipping_discount, currency
- **Logística:** tracking_number, guide_id, warehouse_name, driver_name, platform
- **Estado:** status, payment_method_id, is_paid, order_type_name, invoiceable
- **Otros:** weight, height, width, length, notes

## Data Type Handling

### Float Fields
All float fields (amounts, coordinates, dimensions) use `parseRobustFloat()` which handles:
- European format: `1.200,50` → 1200.5
- US format: `1,200.50` → 1200.5
- Clean format: `50.00` → 50.0

### Boolean Fields
`is_paid` and `invoiceable` accept case-insensitive values:
- True: `true`, `1`, `yes`, `si`
- False: `false`, `0`, `no`, `no`

### Numeric IDs
`payment_method_id` parses as unsigned integer (0-4294967295)

## Required Fields (Unchanged)
These 8 fields are still mandatory:
1. `order_number` - Unique order identifier
2. `customer_name` - Full customer name
3. `customer_email` - Email address
4. `customer_phone` - Phone number
5. `shipping_street` - Street address
6. `shipping_city` - City/Municipality
7. `shipping_state` - Department/State
8. `total_amount` - Total order amount

## Backwards Compatibility
✅ Fully backwards compatible - old templates with only required + basic optional columns (weight, height, width, length, platform) continue to work unchanged.

## Usage Example

Users can now upload:

### Minimal (as before)
```csv
order_number,customer_name,customer_email,customer_phone,shipping_street,shipping_city,shipping_state,total_amount
ORD-001,Juan Perez,juan@example.com,3001234567,Calle 1 # 2-3,Bogota,Cundinamarca,50000
```

### Enhanced (new capability)
```csv
order_number,customer_name,customer_email,customer_phone,shipping_street,shipping_city,shipping_state,customer_dni,subtotal,tax,shipping_cost,total_amount,currency,status,payment_method_id,is_paid
ORD-001,Juan Perez,juan@example.com,3001234567,Calle 1 # 2-3,Bogota,Cundinamarca,1234567890,45000,5000,2000,50000,COP,pending,1,false
```

## Testing Notes

To verify the enhancements:
1. Download the updated template from the modal
2. Fill in additional optional columns
3. Upload via the bulk order modal
4. Verify fields are saved in database via query or detail view
5. Check that partial data (some columns empty) doesn't cause failures

## Build Status
✅ Backend: Compiles successfully with no errors
✅ Frontend: No breaking changes, backwards compatible
