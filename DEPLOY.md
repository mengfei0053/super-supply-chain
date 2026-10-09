# Super Supply Chain on 101

- Path: `/opt/super-supply-chain`
- Secrets: `.env.compose` (mode 600); original Mac env also at `configs/.env` (mode 600)
- URL: https://ssc.mengfei.tech/super-supply-chain/ (docs: https://ssc.mengfei.tech/docs/)
- Port: **172.17.0.1:8088** -> container 8081 (not on public NIC; NPM uses docker0). 8317 is CLI Proxy — do not reuse
- Stack: `ssc-mysql` + `ssc-app`, network `ssc-net`, restart unless-stopped
- DB name: `super_supply_chain` (imported dump, 18 tables)
- Dockerfile: public `node:22` + `golang:1.23.6-alpine` (Aliyun original saved as `Dockerfile.aliyun.bak`)

## Update

```sh
# from a machine that can reach this VPS
rsync -az --exclude frontend/node_modules --exclude .git ./ ubuntu@101.36.111.17:/opt/super-supply-chain/
ssh ubuntu@101.36.111.17 'cd /opt/super-supply-chain && sudo docker compose --env-file .env.compose build ssc-app && sudo docker compose --env-file .env.compose up -d ssc-app'
```

## Notes

- Public `http://101.36.111.17:8088` is intentionally closed; use the domain via Nginx Proxy Manager (`172.17.0.1:8088`).

- WebDAV/NAS (`UPLOAD_SERVER` points at home LAN 192.168.50.63) is unreachable from 101; uploads will fail until VPN/tunnel or new URL.
- Production logs print UPLOAD_* via `configs.LoadConfigFile()` — avoid sharing raw `docker logs`.
- Root `/` redirects on GET (301) to `/super-supply-chain`; HEAD may 404 (Gin route is GET-only).

## Uploads (production on 101)

`UPLOAD_SERVER=file:///data/ssc-uploads` stores Excel uploads on Docker volume `ssc-uploads`.
Home LAN WebDAV (`192.168.50.63:5000`) is unreachable from this VPS.
Set `UPLOAD_SERVER` to a reachable WebDAV URL to restore NAS mode; `file://` / `local://` keep local disk mode.

