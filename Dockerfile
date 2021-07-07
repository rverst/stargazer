FROM ghcr.io/rverst/stargazer:v1.0.2

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
