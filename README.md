# FantomRocks API

## UFW Firewall Rule for API server

    We assume the API server uses default port #8084.
    
    Add new application setup for "FantomRocks".
    
    `vim /etc/ufw/applications.d/fantom-rocks`
    
    ```
    [FantomRocks]
    title=Fantom Rocks Demo API Service
    description=Fantom Rocks demo application API backend service
    ports=8084/tcp
    ```
    
    Enable applications on firewall and turn the firewall on.
    
    ```
    ufw app update fantom-rocks
    ufw allow FantomRocks
    ```    

## RDBMS Setup (structured database storage)

    Please make sure to use **secure password** for your installation.
    We recommend using a password manager with random password generator 
    to create one.

    **Do not use the same password in your production setup!**

    ```sql
    create user fantom with encrypted password '<choose-secure-password>';
    create database fantom
        WITH OWNER 'fantom'
        ENCODING 'UTF8'
        LC_COLLATE = 'en_US.UTF-8'
        LC_CTYPE = 'en_US.UTF-8';
    grant all privileges on database fantom to fantom;
    ```
    Unicode is supported when the database character set is UTF8.
    In that case, Unicode characters can be used in any field,
    there is no reduction of field length for non-special fields.
    No special settings of driver is necessary.

## Links to Tools, Modules and Tutorials
* [KeyCloak Identity Management](https://www.keycloak.org/)
* [Graph-Gophers/GraphQL-Go](https://github.com/graph-gophers/graphql-go)
* [Viper](https://github.com/spf13/viper) + [Viper basics tutorial on Root.cz](https://www.root.cz/clanky/zpracovani-konfiguracnich-souboru-v-go-s-vyuzitim-knihovny-viper/)
* [SQLx extension to Go's database/sql](https://github.com/jmoiron/sqlx)
* [Golang Logging library](https://github.com/op/go-logging)
* [Database driver to PostgreSQL (maintained & supported)](https://github.com/jackc/pgx)
* [Excellent article about PostgreSQL in Go](https://medium.com/avitotech/how-to-work-with-postgres-in-go-bad2dabd13e4)
* [Arbitrary-precision fixed-point decimal numbers](https://github.com/shopspring/decimal)
* [OAuth 2 Server and OpenID Connect Certified Provider in Go](https://github.com/ory/hydra)
* [OpenID Connect support for golang.org/x/oauth2](https://github.com/coreos/go-oidc)
* [Big Cache; Efficient cache for gigabytes of data written in Go](https://github.com/allegro/bigcache)
