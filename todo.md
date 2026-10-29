# RosterPulse TODO

## 🚀 Next Up (Tomorrow's Sprint)

### Part 1: Player Profile APIs (`/api/player/profile`)
- [ ] Register guarded endpoints:
  - [ ] `POST /api/player/profile` (Create logged-in player's profile)
  - [ ] `GET /api/player/profile/me` (Get my profile)
  - [ ] `PATCH /api/player/profile/me` (Update my profile)
  - [ ] `GET /api/players/{playerID}` (View public player profile)

### Part 2: Player Game Profile APIs (`/api/player/game-profiles`)
- [ ] Register guarded endpoints:
  - [ ] `POST /api/player/game-profiles`
  - [ ] `GET /api/player/game-profiles`
  - [ ] `GET /api/player/game-profiles/{id}`
  - [ ] `PATCH /api/player/game-profiles/{id}`
  - [ ] `DELETE /api/player/game-profiles/{id}`

### Part 3: Postman Testing & Verification (0/11 Checklist)
- [ ] Verify profile creation, duplicate blocking, and updates for authenticated players
- [ ] Verify role restrictions (organization accounts blocked from player profile management)
- [ ] Test game profile validation against active/inactive games
- [ ] Confirm strict cross-player ownership security (Player A cannot access Player B's private records)