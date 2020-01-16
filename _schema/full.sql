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
    identity text not null unique references identity(code) ON DELETE CASCADE,
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
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    name ltree not null,
    unique(context, name)
);
-- need to add a constraint so that we always have two actions for each context: root, and system, and all actions must be children of an existing action
--  maybe a trigger on context create that adds them, and a constraint trigger on action create?
CREATE TABLE role (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    name TEXT NOT NULL,
    unique(context, name)
);
CREATE TABLE clique (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    name TEXT NOT NULL,
    unique(context, name)
);
CREATE TABLE client (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    display_name text NOT NULL,
    client_id text NOT NULL UNIQUE,
    secret text NOT NULL,
    base_uri text
);

CREATE TABLE role_to_action (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    role text NOT NULL,
    action ltree not null,
    unique(context, role, action),
    foreign key (context, role) references role(context, name) ON DELETE CASCADE,
    foreign key (context, action) references action(context, name) ON DELETE CASCADE
);

CREATE TABLE ratelimit_prototype (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    action ltree NOT NULL,
    minimum integer NOT NULL,
    maximum integer NOT NULL,
    rate integer NOT NULL,
    unit interval NOT NULL,
    foreign key (context, action) references action(context, name) ON DELETE CASCADE,
    CONSTRAINT rate_limiter_template_check CHECK ((maximum > minimum)),
    CONSTRAINT rate_limiter_template_rate_check CHECK ((rate >= 0)),
    unique(context, action)
);
CREATE TABLE ratelimit_instance (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    action ltree NOT NULL,
    value text NOT NULL,
    durable boolean DEFAULT false NOT NULL,
    minimum integer NOT NULL,
    maximum integer NOT NULL,
    available real NOT NULL,
    rate integer NOT NULL,
    unit interval NOT NULL,
    last_checked timestamp with time zone DEFAULT now() NOT NULL,
    foreign key (context, action) references action(context, name) ON DELETE CASCADE,
    CONSTRAINT ratelimiter_instance_check CHECK ((maximum > minimum)),
    CONSTRAINT ratelimiter_instance_rate_check CHECK ((rate >= 0)),
    unique(context, action, value)
);


CREATE TABLE identity_clique_role (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    identity text not null references identity(code) ON DELETE CASCADE,
    clique text not null,
    role text not null,
    unique(context, identity, clique, role),
    foreign key (context, clique) references clique(context, name) ON DELETE CASCADE,
    foreign key (context, role) references role(context, name) ON DELETE CASCADE
);

CREATE TABLE identity_role (
    context text NOT NULL REFERENCES context(code) ON DELETE CASCADE,
    identity text not null references identity(code) ON DELETE CASCADE,
    role text not null,
    unique(context, identity, role),
    foreign key (context, role) references role(context, name) ON DELETE CASCADE
);


CREATE TABLE session_token (
    code text not null primary key,
    identity text NOT NULL references identity(code) ON DELETE CASCADE,
    user_agent text,
    ip_address text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_in integer
);

CREATE TABLE access_context (
    code text not null primary key,
    session text  REFERENCES session_token(code) ON DELETE CASCADE,
    client text NOT NULL REFERENCES client(client_id) ON DELETE CASCADE,
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
ALTER TABLE identity_role OWNER TO postgres;
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
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE identity_role TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE session_token TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE access_context TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE revocation TO janus;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE stash_data TO janus;

create or replace view identity_access_summary as
    select i.code as identity, ir.context, null as clique, ir.role, rta.action
    from identity i
    join identity_role ir on ir.identity = i.code
    join role_to_action rta on rta.role = ir.role and ir.context = rta.context
union
    select i.code as identity, icr.context, icr.clique, icr.role, rta.action
    from identity i
    join identity_clique_role icr on icr.identity = i.code
    join role_to_action rta on rta.context = icr.context and rta.role = icr.role
;

ALTER view identity_access_summary OWNER TO postgres;
GRANT SELECT ON identity_access_summary TO janus;

COMMIT;
