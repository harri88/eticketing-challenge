using Application.Services;
using Microsoft.AspNetCore.Mvc;
using System.Threading.Tasks;

namespace Api.Controllers
{
    [ApiController]
    [Route("api/v1/tickets")]
    public class TicketsController : ControllerBase
    {
        private readonly ITicketService _ticketService;

        public TicketsController(ITicketService ticketService)
        {
            _ticketService = ticketService;
        }

        [HttpGet]
        public async Task<IActionResult> GetTickets()
        {
            var result = await _ticketService.GetAvailableTicketsAsync();
            return Ok(result);
        }
    }
}