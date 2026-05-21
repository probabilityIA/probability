WITH RECURSIVE src AS (
    SELECT o.id AS order_id,
           ST_SetSRID(ST_MakePoint(o.shipping_lng, o.shipping_lat), 4326) AS p
    FROM orders o
    WHERE o.deleted_at IS NULL
      AND o.destination_geozone_id IS NULL
      AND o.shipping_lat IS NOT NULL
      AND o.shipping_lng IS NOT NULL
),
match AS (
    SELECT DISTINCT ON (src.order_id)
           src.order_id,
           g.id AS gid
    FROM src
    JOIN geozones g
      ON g.deleted_at IS NULL
     AND g.is_active = TRUE
     AND g.business_id = 0
     AND ST_Contains(g.geometry, src.p)
    ORDER BY src.order_id, CASE g.type
        WHEN 'barrio' THEN 1
        WHEN 'neighborhood' THEN 2
        WHEN 'admin_district' THEN 3
        WHEN 'locality' THEN 4
        WHEN 'city' THEN 5
        WHEN 'state' THEN 6
        WHEN 'country' THEN 7
        ELSE 9 END
),
chain AS (
    SELECT m.order_id, g.id, g.parent_id, g.type, ARRAY[g.id]::bigint[] AS path
    FROM match m
    JOIN geozones g ON g.id = m.gid
    UNION ALL
    SELECT c.order_id, g.id, g.parent_id, g.type, c.path || g.id
    FROM chain c
    JOIN geozones g ON g.id = c.parent_id AND g.deleted_at IS NULL
),
levels AS (
    SELECT c.order_id,
           MAX(c.id) FILTER (WHERE c.type = 'country')        AS country_id,
           MAX(c.id) FILTER (WHERE c.type = 'state')          AS state_id,
           MAX(c.id) FILTER (WHERE c.type = 'city')           AS city_id,
           MAX(c.id) FILTER (WHERE c.type = 'admin_district') AS admin_district_id,
           MAX(c.id) FILTER (WHERE c.type = 'locality')       AS locality_id,
           MAX(c.id) FILTER (WHERE c.type = 'neighborhood')   AS neighborhood_id,
           MAX(c.id) FILTER (WHERE c.type = 'barrio')         AS barrio_id,
           (SELECT to_jsonb(path) FROM chain c2
            WHERE c2.order_id = c.order_id
            ORDER BY array_length(path, 1) DESC LIMIT 1) AS path_json
    FROM chain c
    GROUP BY c.order_id
)
UPDATE orders o
SET destination_point = ST_SetSRID(ST_MakePoint(o.shipping_lng, o.shipping_lat), 4326)::geography,
    destination_geozone_id = m.gid,
    destination_geozone_path = l.path_json,
    geozone_country_id = l.country_id,
    geozone_state_id = l.state_id,
    geozone_city_id = l.city_id,
    geozone_admin_district_id = l.admin_district_id,
    geozone_locality_id = l.locality_id,
    geozone_neighborhood_id = l.neighborhood_id,
    geozone_barrio_id = l.barrio_id
FROM match m
JOIN levels l ON l.order_id = m.order_id
WHERE o.id = m.order_id;
