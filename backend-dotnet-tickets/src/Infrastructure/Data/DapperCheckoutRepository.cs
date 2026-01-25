using Dapper;
using Domain.Entities;
using Domain.Interfaces;
using Microsoft.Extensions.Configuration;
using Npgsql;
using System.Data;
using System.Threading.Tasks;

namespace Infrastructure.Data
{
    public class DapperCheckoutRepository : ICheckoutRepository
    {
        private readonly string _connectionString;

        public DapperCheckoutRepository(IConfiguration configuration)
        {
            _connectionString = configuration.GetConnectionString("DefaultConnection");
        }

        private IDbConnection CreateConnection() => new NpgsqlConnection(_connectionString);

        public async Task<Ticket> GetTicketByIdAsync(string id)
        {
            using var connection = CreateConnection();
            return await connection.QueryFirstOrDefaultAsync<Ticket>("SELECT * FROM tickets WHERE id = @Id", new { Id = id });
        }

        public async Task<string> CreateOrderWithReservationAsync(Order order)
        {
            using var connection = CreateConnection();
            connection.Open();

            // Begin Transaction
            using var transaction = connection.BeginTransaction();

            try
            {
                // 1. Insert Order
                var sqlOrder = @"INSERT INTO orders (id, customer_email, total_amount, currency, status, created_at) 
                                 VALUES (@Id, @CustomerEmail, @TotalAmount, @Currency, @Status, @CreatedAt)";
                await connection.ExecuteAsync(sqlOrder, order, transaction);

                // 2. Process Items (Insert Item & Update Ticket Quota)
                foreach (var item in order.Items)
                {
                    // A. Insert Order Item
                    var sqlItem = @"INSERT INTO order_items (order_id, ticket_id, quantity, price_at_purchase)
                                    VALUES (@OrderId, @TicketId, @Quantity, @PriceAtPurchase)";
                    await connection.ExecuteAsync(sqlItem, new 
                    { 
                        OrderId = order.Id, 
                        item.TicketId, 
                        item.Quantity, 
                        item.PriceAtPurchase 
                    }, transaction);

                    // B. Update Ticket (Decrease Quota, Increase Held)
                    // This is the specific requirement: "quota decreased... add... hold quota"
                    var sqlUpdateTicket = @"UPDATE tickets 
                                            SET quota = quota - @Qty, 
                                                held_quota = held_quota + @Qty 
                                            WHERE id = @TId";
                    await connection.ExecuteAsync(sqlUpdateTicket, new { Qty = item.Quantity, TId = item.TicketId }, transaction);
                }

                // Commit if all good
                transaction.Commit();
                return order.Id;
            }
            catch
            {
                // Rollback if anything fails (e.g., negative quota constraint)
                transaction.Rollback();
                throw;
            }
        }

        public async Task<Order> GetOrderByIdAsync(string orderId)
        {
            using var connection = CreateConnection();
            var sql = "SELECT * FROM orders WHERE id = @Id";
            return await connection.QueryFirstOrDefaultAsync<Order>(sql, new { Id = orderId });
        }

        public async Task UpdateOrderAndTicketAsync(Order order)
        {
            using var connection = CreateConnection();
            connection.Open();

            using var transaction = connection.BeginTransaction();

            try
            {
                // Update Order Status
                var sqlUpdateOrder = @"UPDATE orders 
                                       SET status = @Status, 
                                           paid_at = @PaidAt 
                                       WHERE id = @Id";
                await connection.ExecuteAsync(sqlUpdateOrder, new 
                { 
                    order.Status, 
                    order.PaidAt, 
                    order.Id 
                }, transaction);

                // Release Held Quota for Tickets in the Order
                var sqlUpdateTickets = @"UPDATE tickets 
                                          SET held_quota = held_quota - oi.quantity 
                                          FROM order_items oi WHERE tickets.id = oi.ticket_id 
                                        AND oi.order_id = @OrderId";
                await connection.ExecuteAsync(sqlUpdateTickets, new { OrderId = order.Id }, transaction);   


                transaction.Commit();
            }
            catch
            {
                transaction.Rollback();
                throw;
            }
        }
    }
}