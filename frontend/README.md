# Travel Buddy — web app

Next.js (App Router) frontend for the Travel Buddy API. Customers browse
tours and manage their bookings. Agency staff run their agency from a
dashboard, and platform admins onboard agencies.

## Running it

Start the API first (see the repository README), then:

```bash
cp .env.example .env.local   # points NEXT_PUBLIC_API_URL at the API
npm install
npm run dev                  # http://localhost:3001
```

The API listens on port 3000, so the web app uses 3001. `NEXT_PUBLIC_API_URL`
is read at build time; rebuild after changing it.

| Script | Does |
|---|---|
| `npm run dev` | Development server on port 3001 |
| `npm run build` / `npm start` | Production build and server |
| `npm run lint` | ESLint |
| `npm run typecheck` | TypeScript without emitting |

## What each role can do

| Role | Pages |
|---|---|
| Anyone | `/` home, `/search`, `/tours/[id]`, `/agencies/[id]`, `/login`, `/register` |
| Customer | Book from a tour page. `/bookings`: track and cancel bookings, see payment status and why a booking was cancelled |
| Agency staff | `/dashboard`: overview, tours (create, edit, open/close, cancel, delete), bookings (confirm payment, complete, cancel), team (add members, edit permissions). Guest bookings from a tour page |
| Platform admin | `/admin`: create an agency and its first staff account |

Staff see only what their permissions allow. A section they lack a permission
for names the permission to ask for.

## How it's organised

```
src/
  app/            routes; page files read params/searchParams and render a view
  components/
    ui/           buttons, fields, badges, loading / error / empty states
    tours/        tour card, detail, create/edit form
    bookings/     booking form, booking card, "my bookings"
    dashboard/    agency dashboard shell and its pages
    members/      member form, permission picker
  hooks/useAsync  loading / error / data for a request, with reload
  lib/
    api.ts        fetch wrapper: auth header, JSON or text errors, friendly messages
    endpoints.ts  one typed function per API endpoint
    session.ts    token storage (localStorage) shared by the API client and React
    auth.tsx      useSession / login / logout, auto-logout at token expiry
    format.ts     money, dates, and the same price formula the API checks
```

Data is fetched in the browser, not during server rendering. The API rate
limits per client IP, so server-side fetches would make every visitor share
the Next.js server's budget.

Every screen has explicit loading, error (with retry) and empty states. Forms
validate on the client before submitting, and turn API validation errors into
readable messages.
