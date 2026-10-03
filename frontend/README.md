# Travel Buddy customer site

This Next.js app runs on port **3001** and serves customers. Visitors can
browse tours and agencies, register or log in, book tours, and manage their own
bookings. Staff and platform admin pages run in the separate
[agency portal](../agency-portal/README.md).

## Run locally

Start the Go API first, then from this folder:

```bash
cp .env.example .env.local
npm ci
npm run dev
```

Open http://localhost:3001. Set `NEXT_PUBLIC_API_URL` in `.env.local` to
the API's actual URL. On a machine with `HTTP_PORT=3002`, use
`http://localhost:3002`. Set `NEXT_PUBLIC_AGENCY_PORTAL_URL` to the portal
URL (locally `http://localhost:3003`). Rebuild after changing these public
environment variables for a production deployment.

| Script | Does |
|---|---|
| `npm run dev` | Development server on port 3001 |
| `npm run build` / `npm start` | Production build and server on port 3001 |
| `npm run lint` | ESLint |
| `npm run typecheck` | TypeScript check |

Routes: `/`, `/search`, `/tours/[id]`, `/agencies/[id]`, `/login`,
`/register`, and `/bookings`. This app accepts customer logins only.
The API checks every protected request.
