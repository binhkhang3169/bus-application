CREATE TABLE
    "accounts" (
        "id" bigserial PRIMARY KEY,
        "owner_name" bigint NOT NULL UNIQUE,
        "balance" bigint NOT NULL DEFAULT 0,
        "currency" varchar NOT NULL DEFAULT 'VND',
        "status" varchar NOT NULL DEFAULT 'active', -- 'active', 'closed'
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        "updated_at" timestamptz NOT NULL DEFAULT (now ())
    );

CREATE INDEX ON "accounts" ("owner_name");

CREATE INDEX ON "accounts" ("status");

CREATE TABLE
    "transaction_history" (
        "id" bigserial PRIMARY KEY,
        "account_id" bigint NOT NULL,
        "transaction_type" varchar NOT NULL, -- e.g., 'CREATE_ACCOUNT', 'DEPOSIT', 'PAYMENT', 'CLOSE_ACCOUNT'
        "amount" bigint, -- Can be NULL for non-monetary transactions
        "currency" varchar, -- Can be NULL or same as account's currency
        "description" text,
        "created_at" timestamptz NOT NULL DEFAULT (now ()),
        CONSTRAINT "fk_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE
    );

CREATE INDEX ON "transaction_history" ("account_id");

CREATE INDEX ON "transaction_history" ("transaction_type");

CREATE INDEX ON "transaction_history" ("created_at");