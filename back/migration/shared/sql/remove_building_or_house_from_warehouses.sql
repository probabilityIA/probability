-- Eliminar columna building_or_house de la tabla warehouses
-- Esta columna fue agregada por error y no se utiliza en el código
ALTER TABLE warehouses DROP COLUMN IF EXISTS building_or_house;
