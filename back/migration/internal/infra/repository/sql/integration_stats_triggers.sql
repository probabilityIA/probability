CREATE OR REPLACE FUNCTION integration_order_bucket(status text) RETURNS text AS $$
    SELECT CASE
        WHEN status IN ('delivered', 'completed') THEN 'delivered'
        WHEN status IN ('cancelled', 'rejected', 'failed') THEN 'cancelled'
        WHEN status IN ('returned', 'refunded', 'return_in_transit') THEN 'returned'
        ELSE 'in_progress'
    END;
$$ LANGUAGE sql IMMUTABLE;

CREATE OR REPLACE FUNCTION trg_integration_stats_orders() RETURNS trigger AS $$
DECLARE
    old_bucket text;
    new_bucket text;
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.deleted_at IS NULL THEN
            new_bucket := integration_order_bucket(NEW.status);
            INSERT INTO integration_stats (
                integration_id, business_id, orders_total,
                orders_in_progress, orders_delivered, orders_cancelled, orders_returned,
                products_count, last_order_at, updated_at
            ) VALUES (
                NEW.integration_id, COALESCE(NEW.business_id, 0), 1,
                (new_bucket = 'in_progress')::int, (new_bucket = 'delivered')::int,
                (new_bucket = 'cancelled')::int, (new_bucket = 'returned')::int,
                0, NEW.created_at, NOW()
            )
            ON CONFLICT (integration_id) DO UPDATE SET
                orders_total       = integration_stats.orders_total + 1,
                orders_in_progress = integration_stats.orders_in_progress + (new_bucket = 'in_progress')::int,
                orders_delivered   = integration_stats.orders_delivered + (new_bucket = 'delivered')::int,
                orders_cancelled   = integration_stats.orders_cancelled + (new_bucket = 'cancelled')::int,
                orders_returned    = integration_stats.orders_returned + (new_bucket = 'returned')::int,
                last_order_at      = GREATEST(COALESCE(integration_stats.last_order_at, NEW.created_at), NEW.created_at),
                updated_at         = NOW();
        END IF;
        RETURN NEW;
    END IF;

    IF TG_OP = 'UPDATE' THEN
        IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
            old_bucket := integration_order_bucket(OLD.status);
            UPDATE integration_stats SET
                orders_total       = orders_total - 1,
                orders_in_progress = orders_in_progress - (old_bucket = 'in_progress')::int,
                orders_delivered   = orders_delivered - (old_bucket = 'delivered')::int,
                orders_cancelled   = orders_cancelled - (old_bucket = 'cancelled')::int,
                orders_returned    = orders_returned - (old_bucket = 'returned')::int,
                updated_at         = NOW()
            WHERE integration_id = OLD.integration_id;
            RETURN NEW;
        END IF;

        IF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
            new_bucket := integration_order_bucket(NEW.status);
            INSERT INTO integration_stats (
                integration_id, business_id, orders_total,
                orders_in_progress, orders_delivered, orders_cancelled, orders_returned,
                products_count, last_order_at, updated_at
            ) VALUES (
                NEW.integration_id, COALESCE(NEW.business_id, 0), 1,
                (new_bucket = 'in_progress')::int, (new_bucket = 'delivered')::int,
                (new_bucket = 'cancelled')::int, (new_bucket = 'returned')::int,
                0, NEW.created_at, NOW()
            )
            ON CONFLICT (integration_id) DO UPDATE SET
                orders_total       = integration_stats.orders_total + 1,
                orders_in_progress = integration_stats.orders_in_progress + (new_bucket = 'in_progress')::int,
                orders_delivered   = integration_stats.orders_delivered + (new_bucket = 'delivered')::int,
                orders_cancelled   = integration_stats.orders_cancelled + (new_bucket = 'cancelled')::int,
                orders_returned    = integration_stats.orders_returned + (new_bucket = 'returned')::int,
                last_order_at      = GREATEST(COALESCE(integration_stats.last_order_at, NEW.created_at), NEW.created_at),
                updated_at         = NOW();
            RETURN NEW;
        END IF;

        IF NEW.deleted_at IS NULL AND OLD.status IS DISTINCT FROM NEW.status THEN
            old_bucket := integration_order_bucket(OLD.status);
            new_bucket := integration_order_bucket(NEW.status);
            IF old_bucket <> new_bucket THEN
                UPDATE integration_stats SET
                    orders_in_progress = orders_in_progress + (new_bucket = 'in_progress')::int - (old_bucket = 'in_progress')::int,
                    orders_delivered   = orders_delivered + (new_bucket = 'delivered')::int - (old_bucket = 'delivered')::int,
                    orders_cancelled   = orders_cancelled + (new_bucket = 'cancelled')::int - (old_bucket = 'cancelled')::int,
                    orders_returned    = orders_returned + (new_bucket = 'returned')::int - (old_bucket = 'returned')::int,
                    updated_at         = NOW()
                WHERE integration_id = NEW.integration_id;
            END IF;
        END IF;
        RETURN NEW;
    END IF;

    IF TG_OP = 'DELETE' THEN
        IF OLD.deleted_at IS NULL THEN
            old_bucket := integration_order_bucket(OLD.status);
            UPDATE integration_stats SET
                orders_total       = orders_total - 1,
                orders_in_progress = orders_in_progress - (old_bucket = 'in_progress')::int,
                orders_delivered   = orders_delivered - (old_bucket = 'delivered')::int,
                orders_cancelled   = orders_cancelled - (old_bucket = 'cancelled')::int,
                orders_returned    = orders_returned - (old_bucket = 'returned')::int,
                updated_at         = NOW()
            WHERE integration_id = OLD.integration_id;
        END IF;
        RETURN OLD;
    END IF;

    RETURN NULL;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS integration_stats_orders_trg ON orders;
CREATE TRIGGER integration_stats_orders_trg
AFTER INSERT OR UPDATE OR DELETE ON orders
FOR EACH ROW EXECUTE FUNCTION trg_integration_stats_orders();

CREATE OR REPLACE FUNCTION trg_integration_stats_products() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.deleted_at IS NULL THEN
        INSERT INTO integration_stats (
            integration_id, business_id, orders_total,
            orders_in_progress, orders_delivered, orders_cancelled, orders_returned,
            products_count, last_order_at, updated_at
        ) VALUES (NEW.integration_id, NEW.business_id, 0, 0, 0, 0, 0, 1, NULL, NOW())
        ON CONFLICT (integration_id) DO UPDATE SET
            products_count = integration_stats.products_count + 1,
            updated_at     = NOW();
        RETURN NEW;
    END IF;

    IF TG_OP = 'UPDATE' THEN
        IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
            UPDATE integration_stats SET
                products_count = products_count - 1,
                updated_at     = NOW()
            WHERE integration_id = OLD.integration_id;
        ELSIF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
            INSERT INTO integration_stats (
                integration_id, business_id, orders_total,
                orders_in_progress, orders_delivered, orders_cancelled, orders_returned,
                products_count, last_order_at, updated_at
            ) VALUES (NEW.integration_id, NEW.business_id, 0, 0, 0, 0, 0, 1, NULL, NOW())
            ON CONFLICT (integration_id) DO UPDATE SET
                products_count = integration_stats.products_count + 1,
                updated_at     = NOW();
        END IF;
        RETURN NEW;
    END IF;

    IF TG_OP = 'DELETE' THEN
        IF OLD.deleted_at IS NULL THEN
            UPDATE integration_stats SET
                products_count = products_count - 1,
                updated_at     = NOW()
            WHERE integration_id = OLD.integration_id;
        END IF;
        RETURN OLD;
    END IF;

    RETURN NULL;
END $$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS integration_stats_products_trg ON product_business_integrations;
CREATE TRIGGER integration_stats_products_trg
AFTER INSERT OR UPDATE OR DELETE ON product_business_integrations
FOR EACH ROW EXECUTE FUNCTION trg_integration_stats_products();

TRUNCATE integration_stats;

INSERT INTO integration_stats (
    integration_id, business_id, orders_total,
    orders_in_progress, orders_delivered, orders_cancelled, orders_returned,
    products_count, last_order_at, updated_at
)
SELECT
    COALESCE(o.integration_id, p.integration_id),
    COALESCE(o.business_id, p.business_id, 0),
    COALESCE(o.orders_total, 0),
    COALESCE(o.orders_in_progress, 0),
    COALESCE(o.orders_delivered, 0),
    COALESCE(o.orders_cancelled, 0),
    COALESCE(o.orders_returned, 0),
    COALESCE(p.products_count, 0),
    o.last_order_at,
    NOW()
FROM (
    SELECT
        integration_id,
        MAX(business_id) AS business_id,
        COUNT(*) AS orders_total,
        SUM((integration_order_bucket(status) = 'in_progress')::int) AS orders_in_progress,
        SUM((integration_order_bucket(status) = 'delivered')::int) AS orders_delivered,
        SUM((integration_order_bucket(status) = 'cancelled')::int) AS orders_cancelled,
        SUM((integration_order_bucket(status) = 'returned')::int) AS orders_returned,
        MAX(created_at) AS last_order_at
    FROM orders
    WHERE deleted_at IS NULL
    GROUP BY integration_id
) o
FULL OUTER JOIN (
    SELECT
        integration_id,
        MAX(business_id) AS business_id,
        COUNT(*) AS products_count
    FROM product_business_integrations
    WHERE deleted_at IS NULL
    GROUP BY integration_id
) p ON p.integration_id = o.integration_id;
