#!/bin/bash

# Load environment variables
source .env

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Seeding database...${NC}"

# Connect to database and insert seed data
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p ${DB_PORT:-5433} -U $DB_USER -d $DB_NAME <<EOF
-- Insert sample products
INSERT INTO products (name, description, price, stock_quantity, sku, category, image_url) VALUES
    ('Laptop Pro 15', 'High-performance laptop with 16GB RAM and 512GB SSD', 1299.99, 50, 'LP-15-001', 'Electronics', 'https://example.com/laptop-pro-15.jpg'),
    ('Wireless Mouse', 'Ergonomic wireless mouse with precision tracking', 29.99, 200, 'WM-001', 'Electronics', 'https://example.com/wireless-mouse.jpg'),
    ('USB-C Hub', '7-in-1 USB-C hub with HDMI, USB 3.0, and SD card reader', 49.99, 150, 'UCH-001', 'Electronics', 'https://example.com/usb-c-hub.jpg'),
    ('Mechanical Keyboard', 'RGB mechanical keyboard with blue switches', 89.99, 75, 'MK-RGB-001', 'Electronics', 'https://example.com/mech-keyboard.jpg'),
    ('4K Webcam', 'Ultra HD webcam with auto-focus and noise cancellation', 149.99, 60, 'WC-4K-001', 'Electronics', 'https://example.com/4k-webcam.jpg'),
    ('Desk Lamp', 'LED desk lamp with adjustable brightness and color temperature', 39.99, 100, 'DL-LED-001', 'Office', 'https://example.com/desk-lamp.jpg'),
    ('Standing Desk', 'Electric height-adjustable standing desk', 499.99, 30, 'SD-ELEC-001', 'Office', 'https://example.com/standing-desk.jpg'),
    ('Office Chair', 'Ergonomic office chair with lumbar support', 299.99, 40, 'OC-ERG-001', 'Office', 'https://example.com/office-chair.jpg'),
    ('Monitor Stand', 'Adjustable monitor stand with storage drawer', 59.99, 80, 'MS-ADJ-001', 'Office', 'https://example.com/monitor-stand.jpg'),
    ('Noise Cancelling Headphones', 'Premium wireless headphones with ANC', 249.99, 45, 'HP-ANC-001', 'Audio', 'https://example.com/anc-headphones.jpg')
ON CONFLICT (sku) DO NOTHING;

-- Show inserted products
SELECT COUNT(*) as product_count FROM products;
EOF

if [ $? -eq 0 ]; then
  echo -e "${GREEN}Database seeded successfully!${NC}"
else
  echo -e "${RED}Failed to seed database${NC}"
  exit 1
fi
