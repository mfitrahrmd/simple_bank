CREATE TABLE IF NOT EXISTS "accounts" (
                                          "id" serial PRIMARY KEY,
                                          "owner" varchar NOT NULL,
                                          "balance" float NOT NULL,
                                          "currency" varchar NOT NULL,
                                          "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE IF NOT EXISTS "entries" (
                                         "id" serial PRIMARY KEY,
                                         "account_id" integer NOT NULL,
                                         "amount" float NOT NULL,
                                         "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE IF NOT EXISTS "transfers" (
                                           "id" serial PRIMARY KEY,
                                           "from_account_id" integer NOT NULL,
                                           "to_account_id" integer NOT NULL,
                                           "amount" float NOT NULL,
                                           "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
