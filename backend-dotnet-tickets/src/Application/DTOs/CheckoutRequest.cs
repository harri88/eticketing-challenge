using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Application.DTOs
{
    public class CheckoutRequest
    {
        [JsonPropertyName("customer_email")]
        public string CustomerEmail { get; set; }

        [JsonPropertyName("cart_items")]
        public List<CartItemDto> CartItems { get; set; }
    }

    public class CartItemDto
    {
        [JsonPropertyName("ticket_id")]
        public string TicketId { get; set; }

        [JsonPropertyName("quantity")]
        public int Quantity { get; set; }
    }
}