create schema if not exists privatekeeper;

create table if not exists privatekeeper.user
(
    id                      text,
    login                   text not null,
    password                bytea not null,
    crypt_key               bytea not null,
    created_at              timestamp not null,
    updated_at              timestamp not null,
    constraint pk_user primary key (id),
    constraint ux_user__login unique (login)
    );

create type privatekeeper.data_type as enum
    ('credit_card', 'text_data', 'credentials', 'binary_data');

create table if not exists privatekeeper.data
(
    id                      text,
    owner_id                text not null,
    type                    privatekeeper.data_type not null,
    data                    bytea not null,
    metadata                text,
    created_at              timestamp not null,
    updated_at              timestamp not null,
    constraint pk_credit_card primary key (id, type),
    constraint fk_owner_id foreign key (owner_id) references privatekeeper.user (id)
    ) partition by list (type);

create index if not exists privatekeeper_data_type_idx on privatekeeper.data (type);

create table if not exists privatekeeper.data_credit_card partition of privatekeeper.data
    for values in ('credit_card');

create table if not exists privatekeeper.data_text_data partition of privatekeeper.data
    for values in ('text_data');

create table if not exists privatekeeper.data_credentials partition of privatekeeper.data
    for values in ('credentials');

create table if not exists privatekeeper.data_binary_data partition of privatekeeper.data
    for values in ('binary_data');