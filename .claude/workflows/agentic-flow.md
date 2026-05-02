# Agentic Workflow - LicenseManager

## Flow
```
User Request
    ↓
[1. BRAINSTORM] -> approaches + trade-offs
    ↓
[2. PLAN]       -> concrete tasks, files, order
    ↓
[3. COOK]       -> write the code
    ↓
[4. VERIFY]     -> security + correctness check
    ↓
[5. FIX]        -> fix issues found in verify
    ↓
Update session.md
```

## Stage Agents

| Stage | Lead | Support |
|-------|------|---------|
| BRAINSTORM | CTO | - |
| PLAN | CTO | backend-dev, frontend-dev, database |
| COOK | backend-dev + frontend-dev | database |
| VERIFY | security | - |
| FIX | original implementer | - |

## Common Features

### Add user self-service portal
- Backend: `handlers/user.go` - GET /user/licenses, DELETE /user/hwids/:id
- Middleware: user JWT (separate from admin JWT)
- Frontend: `pages/user/` - MyLicenses, Devices
- DB: `users` table already exists

### Add Stripe payments
- Backend: `services/stripe.go` - checkout session, webhook handler
- DB: `users.stripe_customer_id` already exists
- New migration: `subscriptions` table
- Frontend: redirect to Stripe checkout

### Add multi-admin support
- DB: `admins` table already exists with `role` column
- Backend: seed multiple admins, or add POST /admin/admins endpoint
- Frontend: Admin management page

## Gates
- BRAINSTORM -> PLAN: user confirms approach
- PLAN -> COOK: user reviews tasks
- COOK -> VERIFY: all tasks checked off
- VERIFY -> DONE: no critical/high issues (else -> FIX)
