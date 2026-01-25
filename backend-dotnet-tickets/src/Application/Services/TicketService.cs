using Application.DTOs;
using Domain.Interfaces;
using System.Linq;
using System.Threading.Tasks;

namespace Application.Services
{
    public interface ITicketService
    {
        Task<TicketResponse> GetAvailableTicketsAsync();
    }

    public class TicketService : ITicketService
    {
        private readonly ITicketRepository _repository;

        // Constructor Injection (Dependency Injection)
        public TicketService(ITicketRepository repository)
        {
            _repository = repository;
        }

        public async Task<TicketResponse> GetAvailableTicketsAsync()
        {
            var tickets = await _repository.GetAllAsync();

            // Mapping Domain Entity -> DTO
            var dtos = tickets.Select(t => new TicketDto
            {
                Id = t.Id,
                Name = t.Name,
                Price = t.Price,
                Currency = t.Currency,
                Quota = t.Quota,
                IsAvailable = t.IsAvailable() // Using Domain Logic
            });

            return new TicketResponse { Data = dtos };
        }
    }
}