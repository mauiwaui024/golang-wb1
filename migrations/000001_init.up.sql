CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created VARCHAR(255),
    oof_shard VARCHAR(10)
);

CREATE TABLE delivery (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR(255),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE payment (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(3),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    chrt_id INTEGER,
    track_number VARCHAR,
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(20),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
);
