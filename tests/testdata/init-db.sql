CREATE TABLE IF NOT EXISTS orders (
    order_id varchar(255) NOT NULL UNIQUE,
    state INT,
    state_updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (order_id)
);

CREATE TABLE IF NOT EXISTS order_items (
    id varchar(255) NOT NULL UNIQUE,
    order_id varchar(255),
    name varchar(255),
    quantity int,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id)
);

INSERT INTO orders(
	order_id, state, state_updated_at, created_at, updated_at)
	VALUES ('c3fdab1b-3c06-4db2-9edc-4760a2429462', 1, NOW(), NOW(), NOW());

INSERT INTO order_items(
	id, order_id, name, quantity)
	VALUES ('cfdab175-1f86-4fb0-9bcb-15f2c58df30c', 'c3fdab1b-3c06-4db2-9edc-4760a2429462', 'Hamburger', 1);