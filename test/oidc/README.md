# OIDC Testing

## Dex

Check config in ./dex/config/dex.conf and do a `docker-compose up -d`.

Use this gotify config.
```ini
GOTIFY_OIDC_ENABLED=true
GOTIFY_OIDC_ISSUER=http://127.0.0.1:5556/dex
GOTIFY_OIDC_CLIENTID=gotify
GOTIFY_OIDC_CLIENTSECRET=secret
GOTIFY_OIDC_REDIRECTURL=http://127.0.0.1:8080/auth/oidc/callback
```

When testing external apps like gotify/android change every occurence of
127.0.0.1 in ./dex/config/dex.conf and in the gotify config above to an IP that's
routed in your local network like 192.168.178.2.

## Authelia

Authelia requires SSL to work, so you'll have to create a valid certificate. This has to be executed in the directory this README resides.

```
openssl req -x509 -newkey rsa:4096 -nodes -keyout ./authelia/config/key -out ./authelia/config/cert -days 365 -subj "/CN=127.0.0.1" -addext "subjectAltName=IP:127.0.0.1"
```

Check config in ./authelia/config/configuration.yml and do a `docker-compose up -d`.

Use this gotify config.
```ini
GOTIFY_OIDC_ENABLED=true
GOTIFY_OIDC_ISSUER=https://127.0.0.1:9091
GOTIFY_OIDC_CLIENTID=gotify
GOTIFY_OIDC_CLIENTSECRET=secret
GOTIFY_OIDC_REDIRECTURL=http://127.0.0.1:8080/auth/oidc/callback
GOTIFY_OIDC_SCOPES=openid,profile,email,groups
# GOTIFY_OIDC_GROUPS_CLAIM=groups
# GOTIFY_OIDC_GROUPS_USER=
# GOTIFY_OIDC_GROUPS_ADMIN=authelia-group
```

When testing external apps like gotify/android change every occurence of
127.0.0.1 in ./authelia/config/configuration.yml and in the gotify config above
to an IP that's routed in your local network like 192.168.178.2. Also recreate
the certificate with the adjusted IP.
