# Multi-stage build for .NET 10 backend
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build

WORKDIR /app

# Copy project file and restore dependencies
COPY ["backend-dotnet-tickets.csproj", "./"]
RUN dotnet restore "backend-dotnet-tickets.csproj"

# Copy entire source code
COPY ["src/", "./src/"]

# Build the application
RUN dotnet build "backend-dotnet-tickets.csproj" -c Release -o /app/build

# Publish the application
RUN dotnet publish "backend-dotnet-tickets.csproj" -c Release -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime

WORKDIR /app

# Copy published application from build stage
COPY --from=build /app/publish .

# Expose ports
EXPOSE 5020

# Set environment for non-development
ENV ASPNETCORE_ENVIRONMENT=Production
ENV ASPNETCORE_HTTP_PORTS=5020
ENV ASPNETCORE_URLS=http://+:5020

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD dotnet --version || exit 1

# Start the application
ENTRYPOINT ["dotnet", "backend-dotnet-tickets.dll"]
