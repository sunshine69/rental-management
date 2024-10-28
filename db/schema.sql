CREATE TABLE IF NOT EXISTS "tenant" (
    "id" integer NOT NULL PRIMARY KEY,
    "first_name" varchar(254) NOT NULL,
    "last_name" varchar(254) NOT NULL,
    "address" varchar(254) NOT NULL,
    "contact_number" varchar(50) NOT NULL,
    "email" varchar(254),
    "join_date" text,
    "note" text,
    UNIQUE("email")
);

CREATE TABLE IF NOT EXISTS "property_manager" (
    "id" integer NOT NULL PRIMARY KEY,
    "first_name" varchar(254) NOT NULL,
    "last_name" varchar(254) NOT NULL,
    "address" varchar(254) NOT NULL,
    "contact_number" varchar(50) NOT NULL,
    "email" varchar(254),
    "join_date" text,
    "note" text,
    UNIQUE("email")
);

CREATE TABLE IF NOT EXISTS "property" (
    "id" integer NOT NULL PRIMARY KEY,
    "name" varchar(128) NOT NULL,
    "address" varchar(254) NOT NULL,
    "note" text,
    UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS "contract" (
    "id" integer NOT NULL PRIMARY KEY,
    "property_id" int NOT NULL REFERENCES "property" ("id"),
    "property_manager_id" int NOT NULL REFERENCES "property_manager" ("id"),
    "tenant_id" int NOT NULL REFERENCES "tenant" ("id"),
    "start_date" text NOT NULL,
    "end_date" text NOT NULL,
    "signed_date" text NOT NULL,
    "term" text DEFAULT "fixed",
    "rent" integer NOT NULL,
    "rent_period" text DEFAULT "week",
    "rent_paid_on" text,
    "water_charged" integer,
    "document_file_path" text,
    "url" text,
    "note" text,
    UNIQUE("property_id", "signed_date")
);

CREATE TABLE IF NOT EXISTS "account" (
    "id" integer NOT NULL PRIMARY KEY,
    "type" varchar(128),
    "balance" int,
    "contract_id" int NOT NULL REFERENCES "contract" ("id"),
    UNIQUE("contract_id")
);

CREATE TABLE IF NOT EXISTS "payment" (
    "id" integer NOT NULL PRIMARY KEY,
    "account_id" int NOT NULL REFERENCES "account" ("id"),
    "amount" int,
    "pay_date" text,
    "contract_id" int NOT NULL REFERENCES "contract" ("id"),
    "reference" varchar(256),
    UNIQUE("account_id", "pay_date")
);

CREATE TABLE IF NOT EXISTS "invoice" (
    "id" integer NOT NULL PRIMARY KEY,
    "date" text,
    "description" varchar(256),
    "amount" int,
    "number" varchar(128),
    "issuer" varchar(256),
    "payer" varchar(256),
    "property_id" int NOT NULL REFERENCES "property" ("id"),
    "due_date" text,
    UNIQUE("number", "issuer")
);

CREATE TABLE IF NOT EXISTS "maintenance_request" (
    "id" integer NOT NULL PRIMARY KEY,
    "request_date" text,
    "type" varchar(128),
    "status" varchar(256),
    "cost" int,
    "invoice_id" int  int NOT NULL REFERENCES "invoice" ("id"),
    "contract_id" int NOT NULL REFERENCES "contract" ("id"),
    UNIQUE("contract_id", "request_date")
);

