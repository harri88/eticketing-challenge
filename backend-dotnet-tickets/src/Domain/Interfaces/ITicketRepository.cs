using Domain.Entities;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace Domain.Interfaces
{
    // DIP (Dependency Inversion): High-level modules define the contract.
    public interface ITicketRepository
    {
        Task<IEnumerable<Ticket>> GetAllAsync();
    }
}