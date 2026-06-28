# 🚗 SpotSync – Smart Parking & EV Charging Reservation API

**Tech Stack:** Go 1.22 · Echo · GORM · PostgreSQL (NeonDB) · JWT · bcrypt

---

## 📁 Folder Structure

```
spotsync/
├── cmd/
│   └── main.go                        # Entry point – loads env, connects DB, starts server
├── internal/
│   ├── config/
│   │   ├── config.go                  # Loads env into typed Config struct
│   │   └── db.go                      # Opens GORM connection with connection pooling
│   ├── auth/
│   │   └── jwt.go                     # JWTService: GenerateToken + ParseToken
│   ├── middlewares/
│   │   └── auth.go                    # JWT middleware + AdminOnly guard + context helpers
│   ├── httpresponse/
│   │   └── response.go                # Standard JSON envelope (success/error)
│   ├── server/
│   │   └── http.go                    # Creates Echo, DI wiring, auto-migrate, route registration
│   └── domain/
│       ├── reservation_model/
│       │   └── model.go               # Shared Reservation GORM struct (avoids circular imports)
│       ├── user/
│       │   ├── entity.go              # User model + HashPassword + CheckPassword
│       │   ├── repository.go          # DB queries: Create, FindByEmail, FindByID
│       │   ├── service.go             # Register + Login business logic
│       │   ├── handler.go             # HTTP handlers: Register, Login
│       │   ├── register.go            # Route registration
│       │   └── dto/dto.go             # RegisterRequest, LoginRequest, UserResponse, LoginResponse
│       ├── parking_zone/
│       │   ├── entity.go              # ParkingZone model
│       │   ├── repository.go          # Create, FindAll, FindByID, CountActiveReservations
│       │   ├── service.go             # CreateZone, GetAllZones, GetZoneByID
│       │   ├── handler.go             # HTTP handlers: CreateZone, GetAllZones, GetZoneByID
│       │   ├── register.go            # Route registration
│       │   └── dto/dto.go             # CreateZoneRequest, ZoneResponse
│       └── reservation/
│           ├── entity.go              # Reservation with User + Zone relations for Preload
│           ├── repository.go          # CreateWithLock (TX + FOR UPDATE), FindByUserID, etc.
│           ├── service.go             # CreateReservation, GetMyReservations, Cancel, GetAll
│           ├── handler.go             # HTTP handlers
│           ├── register.go            # Route registration
│           └── dto/dto.go             # Request + response DTOs
├── .env.example
├── .air.toml
├── .gitignore
├── go.mod
└── spotsync.postman_collection.json
```

---

## 🚀 Local Setup

```bash
# 1. Clone and enter project
cd spotsync

# 2. Copy and fill in env
cp .env.example .env
# Set DATABASE_URL and JWT_SECRET

# 3. Install dependencies
go mod tidy

# 4. Run
go run ./cmd/main.go

# Optional: hot-reload with Air
go install github.com/air-verse/air@latest
air
```

---

## 🌐 API Endpoints

### Auth (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register user (role: driver or admin) |
| POST | `/api/v1/auth/login` | Login, returns JWT token |

### Parking Zones
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/api/v1/zones` | Public | List all zones with available_spots |
| GET | `/api/v1/zones/:id` | Public | Get single zone |
| POST | `/api/v1/zones` | Admin | Create parking zone |

### Reservations
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/api/v1/reservations` | Driver/Admin | Book a spot (concurrency-safe) |
| GET | `/api/v1/reservations/my-reservations` | Driver/Admin | Own reservations |
| DELETE | `/api/v1/reservations/:id` | Driver/Admin | Cancel reservation |
| GET | `/api/v1/reservations` | Admin only | All reservations |

### Response Envelope
```json
{ "success": true,  "message": "...", "data": { ... } }
{ "success": false, "message": "...", "errors": "..." }
```

---

## ⚡ Concurrency Solution

When two drivers race to claim the last EV spot, SpotSync uses a **DB transaction + `SELECT ... FOR UPDATE` row lock** to guarantee only one succeeds:

```go
db.Transaction(func(tx *gorm.DB) error {
    // Lock the zone row — concurrent TXs queue here
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
    tx.Model(&Reservation{}).Where("zone_id=? AND status=?", zoneID, "active").Count(&count)
    if count >= zone.TotalCapacity { return ErrZoneFull }
    return tx.Create(reservation).Error
})
```

The losing request gets a `409 Conflict` — no double-booking possible.
