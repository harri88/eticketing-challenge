using Domain.Entities;
using System.Threading.Tasks;

namespace Domain.Interfaces
{
    public interface ICheckoutRepository
    {
        Task<Ticket> GetTicketByIdAsync(string id);
        
        // Transactional method: Create order AND update quotas atomically
        Task<string> CreateOrderWithReservationAsync(Order order);
        Task<Order> GetOrderByIdAsync(string orderId);
        Task UpdateOrderAndTicketAsync(Order order);
    }
}