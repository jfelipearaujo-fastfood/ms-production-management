DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_items;

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