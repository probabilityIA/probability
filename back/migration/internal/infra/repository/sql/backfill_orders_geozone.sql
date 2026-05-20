WITH RECURSIVE target AS (
    SELECT id,
           business_id,
           REGEXP_REPLACE(
             TRIM(unaccent(lower(COALESCE(shipping_city,  '')))),
             '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g'
           ) AS city_norm,
           REGEXP_REPLACE(
             TRIM(unaccent(lower(COALESCE(shipping_state, '')))),
             '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g'
           ) AS state_norm
    FROM orders
    WHERE geozone_state_id IS NULL
      AND COALESCE(shipping_city, '') <> ''
),
matched_city AS (
    SELECT DISTINCT ON (t.id) t.id AS order_id, g.id AS gid
    FROM target t
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.type = 'city'
     AND (g.business_id = 0 OR g.business_id = t.business_id)
     AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g') = t.city_norm
    ORDER BY t.id, g.business_id DESC, g.id
),
matched_state AS (
    SELECT DISTINCT ON (t.id) t.id AS order_id, g.id AS gid
    FROM target t
    LEFT JOIN matched_city mc ON mc.order_id = t.id
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.type = 'state'
     AND (g.business_id = 0 OR g.business_id = t.business_id)
     AND REGEXP_REPLACE(unaccent(lower(g.name)), '\s*[,\(]?\s*d\.?\s*c\.?\s*\)?\s*$', '', 'g') IN (t.state_norm, t.city_norm)
    WHERE mc.order_id IS NULL
    ORDER BY t.id, g.business_id DESC, g.id
),
picked AS (
    SELECT order_id, gid FROM matched_city
    UNION ALL
    SELECT order_id, gid FROM matched_state
),
chain AS (
    SELECT p.order_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM picked p
    JOIN geozones g ON g.id = p.gid
    UNION ALL
    SELECT c.order_id, g.id, g.parent_id, g.type, c.path || g.id
    FROM chain c
    JOIN geozones g ON g.id = c.parent_id AND g.deleted_at IS NULL
),
levels AS (
    SELECT c.order_id,
           MAX(c.id) FILTER (WHERE c.type = 'country') AS country_id,
           MAX(c.id) FILTER (WHERE c.type = 'state')   AS state_id,
           MAX(c.id) FILTER (WHERE c.type = 'city')    AS city_id,
           (SELECT to_jsonb(path) FROM chain c2
            WHERE c2.order_id = c.order_id
            ORDER BY array_length(path,1) DESC LIMIT 1) AS path_json
    FROM chain c
    GROUP BY c.order_id
)
UPDATE orders o
SET destination_geozone_id   = COALESCE(o.destination_geozone_id,   p.gid),
    destination_geozone_path = COALESCE(o.destination_geozone_path, l.path_json),
    geozone_country_id       = COALESCE(o.geozone_country_id,       l.country_id),
    geozone_state_id         = COALESCE(o.geozone_state_id,         l.state_id),
    geozone_city_id          = COALESCE(o.geozone_city_id,          l.city_id)
FROM picked p
JOIN levels l ON l.order_id = p.order_id
WHERE o.id = p.order_id;

WITH RECURSIVE target AS (
    SELECT o.id,
           o.business_id,
           o.geozone_city_id,
           TRIM(unaccent(lower(COALESCE(NULLIF(o.shipping_neighborhood, ''),
                                        SPLIT_PART(o.shipping_street, ' | ', 3))))) AS barrio_norm
    FROM orders o
    WHERE o.geozone_city_id IS NOT NULL
      AND o.geozone_barrio_id IS NULL
),
candidates AS (
    SELECT t.id AS order_id, g.id, g.parent_id, g.type
    FROM target t
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.type IN ('barrio','neighborhood')
     AND (g.business_id = 0 OR g.business_id = t.business_id)
     AND unaccent(lower(g.name)) = t.barrio_norm
    WHERE t.barrio_norm <> ''
),
anc AS (
    SELECT c.order_id, c.id AS leaf_id, c.id AS cur_id, c.parent_id, 0 AS depth
    FROM candidates c
    UNION ALL
    SELECT a.order_id, a.leaf_id, g.id, g.parent_id, a.depth + 1
    FROM anc a
    JOIN geozones g ON g.id = a.parent_id AND g.deleted_at IS NULL
    WHERE a.depth < 8
),
matched AS (
    SELECT DISTINCT ON (t.id) t.id AS order_id, a.leaf_id AS barrio_id
    FROM target t
    JOIN anc a ON a.order_id = t.id AND a.cur_id = t.geozone_city_id
    ORDER BY t.id, a.depth ASC
),
chain AS (
    SELECT m.order_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM matched m
    JOIN geozones g ON g.id = m.barrio_id
    UNION ALL
    SELECT c.order_id, g.id, g.parent_id, g.type, c.path || g.id
    FROM chain c
    JOIN geozones g ON g.id = c.parent_id AND g.deleted_at IS NULL
),
levels AS (
    SELECT c.order_id,
           MAX(c.id) FILTER (WHERE c.type = 'neighborhood')    AS neighborhood_id,
           MAX(c.id) FILTER (WHERE c.type = 'admin_district')  AS admin_district_id,
           MAX(c.id) FILTER (WHERE c.type = 'locality')        AS locality_id,
           MAX(c.id) FILTER (WHERE c.type = 'barrio')          AS barrio_id,
           (SELECT to_jsonb(path) FROM chain c2
            WHERE c2.order_id = c.order_id
            ORDER BY array_length(path,1) DESC LIMIT 1) AS path_json
    FROM chain c
    GROUP BY c.order_id
)
UPDATE orders o
SET destination_geozone_id    = COALESCE(l.barrio_id, o.destination_geozone_id),
    destination_geozone_path  = COALESCE(l.path_json, o.destination_geozone_path),
    geozone_barrio_id         = COALESCE(o.geozone_barrio_id,         l.barrio_id),
    geozone_neighborhood_id   = COALESCE(o.geozone_neighborhood_id,   l.neighborhood_id),
    geozone_admin_district_id = COALESCE(o.geozone_admin_district_id, l.admin_district_id),
    geozone_locality_id       = COALESCE(o.geozone_locality_id,       l.locality_id)
FROM levels l
WHERE o.id = l.order_id AND l.barrio_id IS NOT NULL;
