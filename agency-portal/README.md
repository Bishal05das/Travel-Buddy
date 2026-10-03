# Travel Buddy agency portal

This independent Next.js app runs on port **3003**. Agency staff manage
tours, guest bookings, agency bookings, and team permissions. Platform admins
use the same portal to onboard agencies. The customer site lives in
[`frontend/`](../frontend/README.md) and runs separately on port 3001.

## Run locally

Start the Go API first, then from this folder:

```bash
cp .env.example .env.local
npm ci
npm run dev
```

Open http://localhost:3003. Set `NEXT_PUBLIC_API_URL` to the API's actual
URL in `.env.local` (for example `http://localhost:3002` when
`HTTP_PORT=3002`). Set `NEXT_PUBLIC_CUSTOMER_URL` to the customer site's
URL. Rebuild after changing either public variable for production.

| Script | Does |
|---|---|
| `npm run dev` | Development server on port 3003 |
| `npm run build` / `npm start` | Production build and server on port 3003 |
| `npm run lint` | ESLint |
| `npm run typecheck` | TypeScript check |

Owners and staff log in at `/login`, then use `/dashboard` for tours, bookings and
team members. The header shows the signed-in person's name and role, with a link
to `/profile` for their email, phone, and agency. Viewing your profile does not
require permission to view the team list.

The guest booking form is at `/tours/[id]` and linked from
the tour list. Platform admins use `/admin` to create an agency and its
owner account. Owners have full agency access and can add or remove staff.
No account can remove an agency owner or change the owner's permissions.
These routes require the corresponding role; the API
enforces permissions for protected operations.

There are no seeded staff or admin login accounts. Follow the bootstrap steps
in the [repository README](../README.md#bootstrapping-a-new-environment)
to create the first platform admin, then use `/admin` to onboard an agency and
its owner. Owners manage staff at `/dashboard/members`.

Use **Settings** (`/dashboard/settings`) to upload or replace the agency image.
On **Tours → Edit**, use **Tour cover image → Replace image → Save image** to
replace a tour's cover without changing its other details. Both controls preview
the selected image and accept JPG, PNG and WebP files up to 10 MB. Owners have
access; staff need `agency:update` for agency images and `tour:update` for tour
images. Saved images appear on the customer site's agency and tour pages.
