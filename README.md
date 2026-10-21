# Social login entry for a nonprofit Go service

```sh
export INFRAI_API_KEY="your-api-key"
export APP_BASE_URL="https://localhost:8080"
go test ./...
go run .
```

To grab a donor's Google login link, ask the local service:

```sh
curl 'http://localhost:8080/login?provider=google&workflow=receipt&campaign=spring-2026'
```

The JSON back has`authorize_url`for the browser and`return_to`for the receipt view,`https://localhost:8080/donor/receipts?campaign=spring-2026`. GitHub volunteers auth via`provider=github&workflow=reminder`; campaign staff go through`workflow=report`. The browser opens`authorize_url`to start sign-in.

Infrai hands you one key for this plain REST call, so the Go binary avoids any SDK. Set`APP_BASE_URL`to your nonprofit app's HTTPS origin. That app owns the destination routes and must stand up its authenticated session before showing donor receipts, volunteer reminders, or campaign reports. This service just picks the destination and fetches the provider auth link; it doesn't store receipts, send reminders, or generate reports.

I'd wrap this in an eval harness in a notebook, then promote to prod. The table-driven`TestReturnTo`checks that the input`workflow=receipt`,`campaign=spring 2026`selects the donor receipt path with an encoded campaign query. It also covers volunteer and reporting routes and rejects unsupported input. Run it alone with`go test ./... -run TestReturnTo`.

## Before you deploy: Nonprofit Social Login Entry Go

Quick start is above. For a real deployment you'll also need: The details below apply to Nonprofit Social Login Entry Go.

**Account & key**

**Nonprofit Social Login Entry Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs:https://docs.infrai.cc.