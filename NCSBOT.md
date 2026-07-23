# NCSBot

NCSBot is the National Council of Sports branded customer-engagement deployment.

## Local production preview

```powershell
docker compose -f docker-compose.ncsbot.yaml build
docker compose -f docker-compose.ncsbot.yaml run --rm ncsbot-web bundle exec rails db:chatwoot_prepare
docker compose -f docker-compose.ncsbot.yaml up -d
```

Open `http://localhost:3100/app/login` for the administration interface.
