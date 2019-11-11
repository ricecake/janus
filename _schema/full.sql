BEGIN;

CREATE TABLE identity (
    code text  NOT NULL PRIMARY KEY,
    email citext unique not null,
    preferred_name text not null,
    given_name text,
    family_name text
);
CREATE TABLE auth_password (
    identity text not null references identity(code),
    secret text NOT NULL,
    totp text,
    created timestamp with time zone DEFAULT now() not null
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
    client_id text NOT NULL,
    secret text NOT NULL,
    base_uri text
);

CREATE TABLE role_to_action (
    role integer NOT NULL REFERENCES role(id),
    action integer not null references action(id),
    unique(role, action)
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
    identity integer NOT NULL,
    user_agent text,
    ip_address text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
);

CREATE TABLE access_context (
    code text not null primary key,
    session text not null,
    client integer NOT NULL,
    scope text NOT NULL,
    redirect_uri text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
);
CREATE TABLE access_token (
    code text not null primary key,
    context text not null,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
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

ALTER TABLE context OWNER to postgres;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE context TO janus;


COMMIT;