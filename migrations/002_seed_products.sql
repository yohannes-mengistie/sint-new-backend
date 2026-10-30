INSERT INTO stores (name, address)
SELECT 'Shoa Supermarket', 'Addis Ababa'
WHERE NOT EXISTS (
    SELECT 1 FROM stores WHERE name = 'Shoa Supermarket'
);

INSERT INTO stores (name, address)
SELECT 'Safeway', 'Addis Ababa'
WHERE NOT EXISTS (
    SELECT 1 FROM stores WHERE name = 'Safeway'
);

INSERT INTO products (name, variant, size, barcode)
VALUES
    ('Fresh Milk', 'Full Cream', '1L', 'sint-milk-1l'),
    ('Bananas', '', '1kg', 'sint-bananas-1kg'),
    ('Bread', '', '500g', 'sint-bread-500g')
ON CONFLICT (barcode) DO NOTHING;

INSERT INTO price_updates (
    product_id,
    store_id,
    price,
    status,
    observed_at
)
SELECT product.id, store.id, seed.price, 'approved', NOW()
FROM (
    VALUES
        ('Fresh Milk', 'sint-milk-1l', 'Shoa Supermarket', 50.00::NUMERIC),
        ('Bananas', 'sint-bananas-1kg', 'Safeway', 100.00::NUMERIC),
        ('Bread', 'sint-bread-500g', 'Shoa Supermarket', 12.50::NUMERIC)
) AS seed(name, barcode, store_name, price)
JOIN products product ON product.barcode = seed.barcode
JOIN stores store ON store.name = seed.store_name
WHERE NOT EXISTS (
    SELECT 1
    FROM price_updates existing
    WHERE existing.product_id = product.id
      AND existing.store_id = store.id
);
