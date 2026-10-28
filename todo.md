# RosterPulse TODO

## 🔐 Auth Core & Flow (Completed)

- [x] Password hashing (`bcrypt`) and validation
- [x] Token manager (`JWT` access tokens + SHA-256 refresh tokens)
- [x] Session repository database operations with `DBExecutor` / transaction support
- [x] Login handler and `HttpOnly` cookie setup
- [x] `authMiddleware` for token validation & context injection
- [x] `GET /api/v1/auth/me` (Profile endpoint)
- [x] `POST /api/v1/auth/refresh` (Atomic session rotation & new token issuance)
- [x] `POST /api/v1/auth/logout` (Session revocation & secure cookie cleanup)

## 🚀 Next Up

- [ ] Create `internal/games` module (`GET /api/v1/games` & `GET /api/v1/games/{gameid}`)
- [ ] Implement game service, repository, and guarded route handlers using `authMiddleware`
