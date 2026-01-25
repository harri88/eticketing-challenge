### 🚀 Local Development

1. **Prerequisites:** Ensure you have Docker and Docker Compose installed.
2. **Start Services:** Run the following command from the root directory:
   ```bash
   docker-compose up --build


   

graph TD
    subgraph Frontend
        A[React SPA]
        B[React Backoffice]
    end

    subgraph "Ticket Service (.NET Core)"
        C[Checkout API]
        D[(PostgreSQL/MySQL)]
    end

    subgraph "Payment Service (Go)"
        E[Payment Strategy Engine]
        F[CC/QR Handlers]
        G[(MySQL)]
    end

    subgraph "Ledger Service (Python)"
        H[FastAPI Clean Architecture]
        I[(PostgreSQL)]
    end

    %% Flow
    A -->|1. Reserve| C
    C -->|2. Order Created| D
    A -->|3. Pay| E
    E -->|4. Simulation| F
    F -->|5. Record| G
    E -->|6. Notify| C
    E -->|7. Record Entry| H
    H -->|8. Balanced Entry| I
    B -->|Audit| D
    B -->|Audit| G
    B -->|Audit| I