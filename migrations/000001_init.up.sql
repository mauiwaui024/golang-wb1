CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    date_created TIMESTAMP,
    shardkey VARCHAR(10),
    sm_id INTEGER,
    oof_shard VARCHAR(10)
);


CREATE TABLE delivery_info (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR(255),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE payment_info (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(3),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt TIMESTAMP,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);
