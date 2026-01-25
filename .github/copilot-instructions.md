# Copilot Instructions for E-Ticketing Challenge

## Architecture Overview

This is a full-stack e-ticketing application with:
- **Backend**: .NET 10 with Clean Architecture (Domain → Application → Infrastructure → API)
- **Frontend**: React 19 with mock API integration (not yet wired to backend)
- **Database**: PostgreSQL with Dapper ORM for data access
- **Infrastructure**: Docker Compose for local development (PostgreSQL container)

**Key insight**: Clean Architecture pattern separates concerns. The Domain layer contains business logic (`IsAvailable()` method), Application layer orchestrates via services, Infrastructure handles data persistence with Dapper, and API exposes REST endpoints.

## Data Flow

1. React Frontend → (Currently mocked in `src/apiService`)
2. .NET API Layer (`TicketsController`) → injects `ITicketService`
3. Application Service (`TicketService`) → calls `ITicketRepository` (interface in Domain)
4. Infrastructure (`DapperTicketRepository`) → executes SQL via Dapper against PostgreSQL
5. Domain `Ticket` entity returned → mapped to DTO in Application layer
6. DTO serialized as JSON with snake_case properties (e.g., `is_available`)

## Critical Developer Workflows

### Local Development Setup
```bash
# Terminal 1: Start PostgreSQL container
docker-compose up -d

# Terminal 2: Run .NET backend
cd backend-dotnet-tickets
dotnet build
dotnet run

# Terminal 3: Run React frontend
cd frontend-react
npm start
```

### Backend Build & Run
- Build: `dotnet build` (from `backend-dotnet-tickets/`)
- Run: `dotnet run`
- Swagger UI available at `https://localhost:7175/swagger`
- Connection string in `appsettings.json` points to `localhost:5432`

### Frontend Setup
- Dependencies: `npm install` (React 19, react-scripts 5.0.1)
- Start dev: `npm start` (runs on port 3000)
- Build: `npm run build`
- Currently uses **mock API** in `apiService` object—must wire to actual backend endpoint

## Project-Specific Conventions

### Naming & Architecture Patterns
- **Namespaces**: Mirror folder structure (`Domain.Entities`, `Application.Services`, `Infrastructure.Data`)
- **Routing**: API routes use `api/v1/[controller]` convention (e.g., `GET /api/v1/tickets`)
- **DTOs**: JSON property names use snake_case (`IsAvailable` → `is_available`) via `[JsonPropertyName]` attribute
- **Response Wrapper**: API returns `{ "data": [...] }` object, not bare arrays (see `TicketResponse` class)

### Database Schema
- Single `tickets` table with columns: `id` (PK, VARCHAR), `name`, `price`, `currency`, `quota`
- Initial data seeded in `init.sql`: Gold/Premium/VIP tickets with AED pricing
- Dapper maps columns directly to `Ticket` constructor parameters in order

### Dependency Injection
- Registered in `Program.cs`:
  - `ITicketRepository` → `DapperTicketRepository` (Scoped)
  - `ITicketService` → `TicketService` (Scoped)
- Controllers receive `ITicketService` via constructor injection

## Key Files & Their Responsibilities

| File | Purpose |
|------|---------|
| [backend-dotnet-tickets/src/Domain/Entities/Ticket.cs](backend-dotnet-tickets/src/Domain/Entities/Ticket.cs) | Core business entity with `IsAvailable()` logic |
| [backend-dotnet-tickets/src/Domain/Interfaces/ITicketRepository.cs](backend-dotnet-tickets/src/Domain/Interfaces/ITicketRepository.cs) | DIP contract—defines what data access should support |
| [backend-dotnet-tickets/src/Application/Services/TicketService.cs](backend-dotnet-tickets/src/Application/Services/TicketService.cs) | Orchestrates repository + maps to DTO |
| [backend-dotnet-tickets/src/Api/Controllers/TicketsController.cs](backend-dotnet-tickets/src/Api/Controllers/TicketsController.cs) | HTTP endpoint—delegates to service |
| [backend-dotnet-tickets/src/Infrastructure/Data/DapperTicketRepository.cs](backend-dotnet-tickets/src/Infrastructure/Data/DapperTicketRepository.cs) | Dapper SQL queries + connection management |
| [frontend-react/src/App.js](frontend-react/src/App.js) | Main React app with `TicketSelection` & `ThankYouPage` components |

## Integration Points & External Dependencies

- **PostgreSQL**: Must be running (via docker-compose) for backend to function
- **Dapper 2.1.66**: Used for parameterized queries; maps query results to entities automatically
- **Npgsql 10.0.1**: PostgreSQL driver for ADO.NET
- **Swagger/OpenAPI**: Enabled in Development; remove from Production via `if (app.Environment.IsDevelopment())`

## Known Gaps / TODO

- Frontend not yet wired to actual backend API (still using mock `apiService`)
- No payment processing integration (checkout flow exists but no actual payments)
- No error handling for database connection failures beyond thrown exception
- No authentication/authorization layer
- No request validation or business rule constraints (e.g., quota enforcement on checkout)
