# RosterPulse TODO

## 🔐 Auth Core (Completed)
- [x] Password hashing (`bcrypt`) and validation
- [x] Token manager (`JWT` access tokens + SHA-256 refresh tokens)
- [x] Session repository database operations
- [x] Login handler and `HttpOnly` cookie setup

## 🚀 Next Up (Tomorrow)
- [ ] Implement `authMiddleware` for token validation & context injection
- [ ] Implement `GET /api/v1/auth/me` (Profile endpoint)
- [ ] Implement `POST /api/v1/auth/refresh` (Session rotation & new token issuance)
- [ ] Implement `POST /api/v1/auth/logout` (Session revocation & cookie cleanup)
- [ ] Create `internal/games` module (`GET /api/v1/games` & `GET /api/v1/games/{gameid}`)