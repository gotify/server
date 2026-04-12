# OIDC Testing

## Dex

Check config in ./dex/config/dex.conf and do a `docker-compose up -d`.

Use this gotify config.
```
oidc:
  enabled: true
  issuer: http://127.0.0.1:5556/dex
  clientid: gotify
  clientsecret: secret
  redirecturl: http://127.0.0.1:8080/auth/oidc/callback
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
```
oidc:
  enabled: true
  issuer: https://127.0.0.1:9091
  clientid: gotify
  clientsecret: secret
  redirecturl: http://127.0.0.1:8080/auth/oidc/callback
```

When testing external apps like gotify/android change every occurence of
127.0.0.1 in ./authelia/config/configuration.yml and in the gotify config above
to an IP that's routed in your local network like 192.168.178.2. Also recreate
the certificate with the adjusted IP.
