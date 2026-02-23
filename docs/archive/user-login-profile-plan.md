# Plan: Enable All Users to Login, See Tracks & API Keys

**Status:** ✅ IMPLEMENTED
**Date Archived:** 2026-02-22
**Implemented Features:**
- User registration and email verification
- Login/logout flows
- Password reset via email
- User profile pages with upload history
- API key display and management
- Authentication logging

See [README.md](../../README.md#user-authentication--api-keys) for current documentation.

---

## Context
There are ~8,400 CSV-migrated users who all have `requires_password_setup: true` and no passwords. They need a way to set passwords and, once logged in, see their tracks and API keys. Most of the auth infrastructure already exists — we mainly need a **user profile page** and a way for migrated users to set their passwords.

## What Already Exists
- Login/register/logout flows (working)
- Forgot password + reset token flow (working, emails now sending)
- Admin can trigger password reset emails per user
- API keys auto-generated for all users
- `/api/user/profile` endpoint returns user info including `api_key`
- Uploads table has `user_id` / `internal_user_id` to link tracks to users

## What's Missing
1. **No user profile page** — logged-in users only see their name + logout button
2. **No way for migrated users to self-serve password setup** — they'd need an admin to trigger reset emails one-by-one (8,400 users!)
3. **No user-facing track list** — tracks/uploads aren't queryable by the logged-in user

## Plan

### Step 1: Add "Claim Account" / Self-Service Password Setup
Allow migrated users to claim their account via the existing forgot-password flow. This already works — a user enters their email on the "Forgot password?" form, gets a reset link, and sets a password. The reset handler already clears `requires_password_setup`. **No code changes needed for this flow.**

However, we should make sure the forgot-password flow is discoverable:
- Add a prominent "Claim your account" or "First time? Set up password" link on the login modal that points to the forgot-password form

### Step 2: Create a User Profile/Dashboard Page (`/profile`)
A new page accessible to logged-in users showing:
- **Account info**: username, email, verified status
- **API key**: displayed with copy button
- **Their uploads/tracks**: list of tracks they've uploaded (queried via `internal_user_id` or `user_id` in uploads table)
- **Change password** button

Files to modify:
- `public_html/profile.html` (new file)
- `safecast-new-map.go` — add route for `/profile`

### Step 3: Add API Endpoint for User's Uploads
- `GET /api/user/uploads` — returns uploads belonging to the logged-in user
- Query uploads table where `internal_user_id` matches the user's ID
- File: `safecast-new-map.go` or new handler in `pkg/auth/`

### Step 4: Add Profile Link to User Menu
- Update the user menu in `map.html` (line 2239) to include a "Profile" link between the username and logout button

### Step 5: Add Change Password Endpoint
- `POST /api/user/change-password` — requires current password + new password
- File: `pkg/auth/handlers.go`

## Files to Create/Modify
| File | Action |
|------|--------|
| `public_html/profile.html` | **Create** — user dashboard page |
| `public_html/map.html` | Edit — add "Profile" link to user menu, add "First time?" link to login modal |
| `safecast-new-map.go` | Edit — add `/profile` route and `/api/user/uploads` endpoint |
| `pkg/auth/handlers.go` | Edit — add change-password handler |

## Verification
1. Open map, click Login, see "First time? Set up password" link
2. Click it, enter a migrated user's email, receive reset email
3. Set password via reset link, log in
4. See "Profile" link in user menu, click it
5. Profile page shows account info, API key with copy button, and user's tracks
6. Change password works from profile page
