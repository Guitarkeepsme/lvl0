CREATE TABLE delivery (
    id TEXT PRIMARY KEY,
    "name" VARCHAR(30) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(30) NOT NULL,
    address VARCHAR(50) NOT NULL,
    region VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL
);

CREATE TABLE payment (
    id TEXT PRIMARY KEY,
    transaction VARCHAR(50) NOT NULL,
    request_id VARCHAR(30) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(30) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost NUMERIC(10, 2) NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee NUMERIC(10, 2) NOT NULL
);

CREATE TABLE "order" (
    order_uid TEXT PRIMARY KEY,
    track_number VARCHAR(30) UNIQUE NOT NULL,
    entry VARCHAR(30) NOT NULL,
    delivery_id TEXT REFERENCES delivery,
    payment_id TEXT REFERENCES payment,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(50) NOT NULL,
    customer_id VARCHAR(30) NOT NULL,
    delivery_service VARCHAR(30) NOT NULL,
    shardkey VARCHAR(30) NOT NULL,
    sm_id BIGINT NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(30) NOT NULL
);

CREATE TABLE item (
    chrt_id BIGINT PRIMARY KEY,
    order_uid TEXT REFERENCES "order",
    track_number VARCHAR(30) UNIQUE NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    rid VARCHAR(50) NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    sale INTEGER NOT NULL CHECK(sale >=0 AND sale <= 100),
    "size" VARCHAR(10) NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(50) NOT NULL,
    status INTEGER NOT NULL
);