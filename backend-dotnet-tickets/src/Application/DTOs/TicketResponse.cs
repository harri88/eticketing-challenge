using System.Text.Json.Serialization;

namespace Application.DTOs
{
    // Response Wrapper to match the { "data": [...] } requirement
    public class TicketResponse
    {
        [JsonPropertyName("data")]
        public IEnumerable<TicketDto> Data { get; set; }
    }

    public class TicketDto
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public decimal Price { get; set; }
        public string Currency { get; set; }
        public int Quota { get; set; }
        
        [JsonPropertyName("is_available")]
        public bool IsAvailable { get; set; }
    }
}