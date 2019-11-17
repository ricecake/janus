CREATE DATABASE janus_app;
\c janus_app
BEGIN;
CREATE EXTENSION citext;
CREATE EXTENSION ltree;

CREATE TABLE identity (
    code text  NOT NULL PRIMARY KEY,
    active boolean not null default true,
    email citext unique not null,
    preferred_name text not null,
    given_name text,
    family_name text
);
CREATE TABLE auth_password (
    identity text not null unique references identity(code),
    totp_active boolean not null default false,
    hash text NOT NULL,
    totp text,
    created_at timestamp with time zone DEFAULT now() not null
);


CREATE TABLE context (
    code text  NOT NULL PRIMARY KEY,
    name text  NOT NULL
);

CREATE TABLE action (
    id serial NOT NULL PRIMARY KEY,
    context text NOT NULL REFERENCES context(code),
    name ltree not null,
    unique(context, name)
);
-- need to add a constraint so that we always have two actions for each context: root, and system, and all actions must be children of an existing action
--  maybe a trigger on context create that adds them, and a constraint trigger on action create?
CREATE TABLE role (
    id serial NOT NULL PRIMARY KEY,
    context text NOT NULL REFERENCES context(code),
    name TEXT NOT NULL,
    unique(context, name)
);
CREATE TABLE clique (
    id serial NOT NULL PRIMARY KEY,
    context text NOT NULL REFERENCES context(code),
    name TEXT NOT NULL,
    unique(context, name)
);
CREATE TABLE client (
    context text NOT NULL REFERENCES context(code),
    display_name text NOT NULL,
    client_id text NOT NULL UNIQUE,
    secret text NOT NULL,
    base_uri text
);

CREATE TABLE role_to_action (
    role integer NOT NULL REFERENCES role(id),
    action integer not null references action(id),
    unique(role, action)
    -- Might be able to add a context column, and do a multi column foreign key, and not need to deal with int primary keys?
);

CREATE TABLE ratelimit_prototype (
    action integer unique NOT NULL references action(id),
    minimum integer NOT NULL,
    maximum integer NOT NULL,
    rate integer NOT NULL,
    unit interval NOT NULL,
    CONSTRAINT rate_limiter_template_check CHECK ((maximum > minimum)),
    CONSTRAINT rate_limiter_template_rate_check CHECK ((rate >= 0))
);
CREATE TABLE ratelimit_instance (
    action integer unique NOT NULL references action(id),
    value text NOT NULL,
    durable boolean DEFAULT false NOT NULL,
    minimum integer NOT NULL,
    maximum integer NOT NULL,
    available real NOT NULL,
    rate integer NOT NULL,
    unit interval NOT NULL,
    last_checked timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ratelimiter_instance_check CHECK ((maximum > minimum)),
    CONSTRAINT ratelimiter_instance_rate_check CHECK ((rate >= 0)),
    unique(action, value)
);


CREATE TABLE identity_clique_role (
    identity text not null references identity(code),
    clique integer not null references clique(id),
    role integer not null references role(id),
    unique(identity, clique, role)
);


CREATE TABLE session_token (
    code text not null primary key,
    identity text NOT NULL references identity(code),
    user_agent text,
    ip_address text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
);

CREATE TABLE access_context (
    code text not null primary key,
    session text  REFERENCES session_token(code),
    client text NOT NULL REFERENCES client(client_id),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE revocation (
    entity_code text not null PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
);

CREATE TABLE stash_data (
    uuid text not null primary key,
    data jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
);

ALTER TABLE identity OWNER to postgres;
ALTER TABLE auth_password OWNER to postgres;
ALTER TABLE context OWNER to postgres;
ALTER TABLE action OWNER to postgres;
ALTER TABLE role OWNER to postgres;
ALTER TABLE clique OWNER to postgres;
ALTER TABLE client OWNER to postgres;
ALTER TABLE role_to_action OWNER to postgres;
ALTER TABLE ratelimit_prototype OWNER to postgres;
ALTER TABLE ratelimit_instance OWNER to postgres;
ALTER TABLE identity_clique_role OWNER to postgres;
ALTER TABLE session_token OWNER to postgres;
ALTER TABLE access_context OWNER to postgres;
ALTER TABLE revocation OWNER to postgres;
ALTER TABLE stash_data OWNER to postgres;

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE identity TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE auth_password TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE context TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE action TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE role TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE clique TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE client TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE role_to_action TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE ratelimit_prototype TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE ratelimit_instance TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE identity_clique_role TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE session_token TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE access_context TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE revocation TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE stash_data TO janus;

COMMIT;
