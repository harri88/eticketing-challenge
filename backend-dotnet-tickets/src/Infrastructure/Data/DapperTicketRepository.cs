using Dapper;
using Domain.Entities;
using Domain.Interfaces;
using Microsoft.Extensions.Configuration;
using Npgsql;
using System.Collections.Generic;
using System.Data;
using System.Threading.Tasks;

namespace Infrastructure.Data
{
    public class DapperTicketRepository : ITicketRepository
    {
        private readonly string _connectionString;

        public DapperTicketRepository(IConfiguration configuration)
        {
            _connectionString = configuration.GetConnectionString("DefaultConnection");
            // ADD THIS DEBUG CHECK
            if (string.IsNullOrWhiteSpace(_connectionString))
            {
                // This will print to your console so you know exactly what's wrong
                throw new InvalidOperationException("CRITICAL ERROR: Connection string 'DefaultConnection' is null or empty. Check appsettings.json.");
            }
        }

        private IDbConnection CreateConnection() => new NpgsqlConnection(_connectionString);

        public async Task<IEnumerable<Ticket>> GetAllAsync()
        {
            const string sql = @"
                SELECT id, name, price, currency, quota 
                FROM tickets";

            using var connection = CreateConnection();
            
            // Dapper automatically maps columns to the Constructor of the Ticket entity
            return await connection.QueryAsync<Ticket>(sql);
        }
    }
}