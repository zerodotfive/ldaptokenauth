Example auth server config
```
SECRET=<jwt sign secret>
LDAP_SERVER=<ldap address:port>
BIND_DN=<bind dn>
BIND_PW=<bind password>
BASE_DN=<base dn>
GROUP_FILTER='(memberOf=cn=deployers,cn=groups,cn=accounts,dc=sf)'
TTL=168h
```

Nginx proxy example config:
```
server {
    listen 0.0.0.0:443 ssl;
    server_name <cluster address>;
    ssl_certificate     ...;
    ssl_certificate_key ...;
    ssl_prefer_server_ciphers on;
    ssl_session_cache    shared:SSL:64m;
    ssl_session_timeout  1h;

    location / {
        include /etc/nginx/internal_access.conf;
        auth_request /validate;

        proxy_pass https://<k8s api real address>;
        proxy_redirect off;
        proxy_set_header Host $http_host;
        proxy_set_header X-Remote-User $upstream_http_username;
        proxy_set_header Connection "Keep-Alive";
        proxy_http_version 1.1;

        proxy_ssl_name <k8s api real address>;
        proxy_ssl_server_name on;
        proxy_ssl_session_reuse on;
        proxy_ssl_certificate <k8s api cert>;
        proxy_ssl_certificate_key <k8s api key>;
    }

    location /validate {
        internal;
        proxy_pass              http://127.0.0.1:8090;
        proxy_set_header        Host $http_host;
        proxy_set_header        X-Forwarded-For $remote_addr;
    }

    location /auth {
        proxy_pass              http://127.0.0.1:8090;
        proxy_set_header        Host $http_host;
        proxy_set_header        X-Forwarded-For $remote_addr;
    }
}
```

How to login:
```
kubectl config set-cluster <cluster name> --server https://<cluster address>/
kubectl config set-credentials <cluster name> --token $(ldaptokenauth https://<cluster address>/auth)
kubectl config set-context <cluster name> --cluster <cluster name> --user=<cluster name>
kubectl config use-context <cluster name>
```
