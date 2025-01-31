FROM gcr.io/distroless/static-debian12
COPY /bin/vault-plugin-secrets-naughty /bin/vault-plugin-secrets-naughty
ENTRYPOINT [ "/bin/vault-plugin-secrets-naughty" ]
