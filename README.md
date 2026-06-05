# BookPlay

System rezerwacji obiektów sportowych z modułem zarządzania drużynami – na przykładzie Grodziska Wielkopolskiego.

## 🏗️ Architektura

| Warstwa | Technologia | Port |
|---------|-------------|------|
| Frontend | Nuxt 3 (Vue 3 + TypeScript) | `3000` |
| Backend API | Go (Gin + GORM) | `8080` |
| Baza danych | PostgreSQL 16 | `5432` |
| Cache / Czat | Redis 7 | `6379` |

```
apps/
├── api/         # Backend Go (REST API)
│   ├── cmd/server/        # Punkt startowy
│   └── internal/
│       ├── config/        # Konfiguracja, seed, DB
│       ├── middleware/     # CORS, JWT, logowanie
│       ├── models/        # Modele GORM
│       └── modules/       # Moduły biznesowe
│           ├── auth/      # Rejestracja, logowanie, JWT, OTP
│           ├── facility/  # Obiekty sportowe, dostępność
│           ├── reservation/ # Rezerwacje, statusy
│           ├── team/      # Drużyny, członkowie, rekrutacja
│           ├── user/      # Profile użytkowników
│           ├── chat/      # Czat drużynowy (Redis)
│           └── mail/      # Wysyłka OTP, przypomnienia SMTP
└── web/          # Frontend Nuxt 3
    ├── pages/             # Strony (index, login, facilities, itp.)
    ├── stores/            # Pinia (auth, theme)
    ├── layouts/           # Layout domyślny
    └── middleware/        # Auth guard
```

## ✨ Główne funkcje

- **Rezerwacja obiektów** – przeglądanie dostępności w czasie rzeczywistym, wybór daty i godziny, kalkulacja ceny
- **Obsługa przedpłat** – możliwość wymagania przedpłaty z danymi konta bankowego i tytułem przelewu
- **Zarządzanie drużynami** – tworzenie, edycja, logo, dodawanie/usuwanie członków, role (kapitan/członek)
- **Rekrutacja do drużyn** – otwarte/zamknięte zgłoszenia, akceptowanie/odrzucanie aplikacji
- **Czat drużynowy** – komunikacja w czasie rzeczywistym (Redis, polling co 5s)
- **Autoryzacja JWT** – access + refresh token, OTP przez email, reset hasła
- **Przypomnienia email** – automatyczne powiadomienia 3 dni przed rezerwacją (SMTP)
- **Role użytkowników** – PLAYER, FACILITY_OWNER, ADMIN (różne dashboardy i uprawnienia)
- **Testowe dane** – automatyczny seeder z kontami testowymi i przykładowymi obiektami w Grodzisku Wlkp.

## 🚀 Uruchamianie

### Wymagania

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go](https://go.dev/dl/) 1.22+
- [Node.js](https://nodejs.org/) 18+
- PowerShell (Windows) lub bash (Linux/Mac)

### Szybki start (Windows)

```powershell
# Klonuj repozytorium
git clone https://github.com/janmaciejewski/book-play.git
cd book-play

# Uruchom wszystko jednym poleceniem (baza, backend, frontend)
# Jeśli PowerShell blokuje uruchamianie skryptów, użyj:
# powershell -ExecutionPolicy Bypass -File ./start.ps1
.\start.ps1
```

Aplikacja będzie dostępna pod adresami:
- 🌐 **Frontend:** http://localhost:3000
- 🔌 **Backend API:** http://localhost:8080

### Ręczne uruchamianie krok po kroku

```powershell
# 1. Uruchom bazę danych i Redis
docker-compose up -d postgres redis

# 2. Zainstaluj zależności
npm install

# 3. Uruchom backend (hot reload przez Air)
go install github.com/cosmtrek/air@latest
air

# 4. W drugim terminalu – frontend
npm run dev:web
```

### Uruchamianie przez Docker (wszystkie usługi)

```powershell
docker-compose up -d
```

## 🔐 Konta testowe (seeder)

Przy pierwszym uruchomieniu baza jest automatycznie wypełniana danymi testowymi:

| Email | Hasło | Rola |
|-------|-------|------|
| admin@bookplay.com | admin123 | Administrator |
| owner@bookplay.com | owner123 | Właściciel obiektu |
| player@bookplay.com | player123 | Gracz |
| captain@bookplay.com | captain123 | Kapitan drużyny |

## ⚙️ Konfiguracja

Aplikacja ładuje ustawienia z pliku `.env` w głównym katalogu (ładowany przez `github.com/joho/godotenv` w `config/config.go`).

Zmienne używane przez backend:

| Zmienna | Domyślnie | Opis |
|---------|-----------|------|
| `APP_NAME` | `book-play` | Nazwa aplikacji |
| `APP_ENV` | `development` | Środowisko (`development` / `production`) |
| `API_PORT` | `8080` | Port serwera API |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `bookplay` | Użytkownik bazy |
| `DB_PASSWORD` | – | Hasło bazy |
| `DB_NAME` | `bookplay` | Nazwa bazy |
| `DB_SSL_MODE` | `disable` | Tryb SSL |
| `DATABASE_URL` | – | Pełny URL bazy (opcjonalny – nadpisuje powyższe) |
| `REDIS_HOST` | `localhost` | Host Redis |
| `REDIS_PORT` | `6379` | Port Redis |
| `REDIS_PASSWORD` | – | Hasło Redis |
| `REDIS_DB` | `0` | Indeks bazy Redis |
| `JWT_SECRET` | (wbudowany fallback) | Klucz do podpisywania tokenów JWT – **wymagane na produkcji** |
| `JWT_ACCESS_TOKEN_EXPIRY` | `15m` | Czas życia access tokena |
| `JWT_REFRESH_TOKEN_EXPIRY` | `168h` | Czas życia refresh tokena |
| `SMTP_HOST` | – | Host SMTP (do OTP i przypomnień) |
| `SMTP_PORT` | `587` | Port SMTP |
| `SMTP_USER` | – | Użytkownik SMTP |
| `SMTP_PASSWORD` | – | Hasło SMTP |
| `SMTP_FROM` | `Book-Play <noreply@bookplay.com>` | Adres nadawcy email |

Gdzie w kodzie są używane:
- `JWT_SECRET` → `config/config.go` ładuje do `JWTConfig.Secret` → `auth/service.go` podpisuje nim token (`SignedString`) i weryfikuje (`jwt.Parse`)
- `DATABASE_URL` → jeśli nie podany osobno, automatycznie składany z `DB_*` w `config/config.go`
- `SMTP_*` → `mail/service.go` – wysyłka OTP i przypomnień o rezerwacji

## 📡 API Endpointy

| Metoda | Endpoint | Dostęp | Opis |
|--------|----------|--------|------|
| POST | `/api/v1/auth/register` | Publiczny | Rejestracja |
| POST | `/api/v1/auth/login` | Publiczny | Logowanie |
| POST | `/api/v1/auth/refresh` | Publiczny | Odświeżenie tokenu |
| POST | `/api/v1/auth/logout` | Publiczny | Wylogowanie |
| POST | `/api/v1/auth/send-otp` | Publiczny | Wyślij kod OTP |
| POST | `/api/v1/auth/verify-otp` | Publiczny | Zweryfikuj OTP |
| POST | `/api/v1/auth/reset-password` | Publiczny | Reset hasła |
| GET | `/api/v1/auth/me` | JWT | Dane zalogowanego |
| GET | `/api/v1/facilities` | Publiczny | Lista obiektów |
| GET | `/api/v1/facilities/:id` | Publiczny | Szczegóły obiektu |
| GET | `/api/v1/facilities/:id/availability` | Publiczny | Dostępne godziny |
| GET | `/api/v1/facilities/mine` | JWT | Moje obiekty |
| POST | `/api/v1/facilities` | JWT | Dodaj obiekt |
| PUT | `/api/v1/facilities/:id` | JWT | Edytuj obiekt |
| PUT | `/api/v1/facilities/:id/slots` | JWT | Aktualizuj godziny |
| PUT | `/api/v1/facilities/:id/close` | JWT | Zamknij/Otwórz obiekt |
| GET | `/api/v1/reservations` | JWT | Moje rezerwacje |
| POST | `/api/v1/reservations` | JWT | Utwórz rezerwację |
| PUT | `/api/v1/reservations/:id/cancel` | JWT | Anuluj rezerwację |
| PUT | `/api/v1/reservations/:id/status` | JWT | Zmień status (właściciel) |
| GET | `/api/v1/facilities/my/reservations` | JWT | Rezerwacje na moje obiekty |
| GET | `/api/v1/teams` | JWT | Lista drużyn |
| GET | `/api/v1/teams/:id` | JWT | Szczegóły drużyny |
| POST | `/api/v1/teams` | JWT | Utwórz drużynę |
| PUT | `/api/v1/teams/:id` | JWT | Edytuj drużynę |
| DELETE | `/api/v1/teams/:id` | Admin | Usuń drużynę |
| POST | `/api/v1/teams/:id/members` | JWT | Dodaj członka |
| DELETE | `/api/v1/teams/:id/members/:memberId` | JWT | Usuń członka |
| PUT | `/api/v1/teams/:id/recruitment` | JWT | Otwórz/zamknij rekrutację |
| POST | `/api/v1/teams/:id/apply` | JWT | Aplikuj do drużyny |
| GET | `/api/v1/teams/:id/chat` | JWT | Pobierz wiadomości czatu |
| POST | `/api/v1/teams/:id/chat` | JWT | Wyślij wiadomość |
| GET | `/api/v1/users` | Admin | Lista użytkowników |
| GET | `/api/v1/users/:id` | JWT | Profil użytkownika |
| PUT | `/api/v1/users/:id` | JWT | Edytuj profil |
| PUT | `/api/v1/users/:id/role` | Admin | Zmień rolę |
| DELETE | `/api/v1/users/:id` | Admin | Usuń użytkownika |

## 🔄 Hot Reload

- **Backend:** [Air](https://github.com/cosmtrek/air) – automatycznie przeładowuje serwer Go przy zmianach w `apps/api/`
- **Frontend:** Nuxt Dev Server – HMR (Hot Module Replacement) przez Vite

## 🧱 Stack technologiczny

**Backend:** Go 1.22, Gin, GORM, PostgreSQL, Redis, JWT, bcrypt, SMTP
**Frontend:** Nuxt 3, Vue 3, TypeScript, Tailwind CSS, Pinia
**Infrastruktura:** Docker Compose, Air (hot reload)