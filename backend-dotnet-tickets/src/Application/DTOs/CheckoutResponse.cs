using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Application.DTOs
{
    public class CheckoutResponse
    {
        [JsonPropertyName("order_id")]
        public string OrderId { get; set; }

        [JsonPropertyName("status")]
        public string Status { get; set; }

        [JsonPropertyName("summary")]
        public CheckoutSummary Summary { get; set; }

        [JsonPropertyName("redirect_url")]
        public string RedirectUrl { get; set; }
    }

    public class CheckoutSummary
    {
        [JsonPropertyName("subtotal")]
        public decimal Subtotal { get; set; }

        [JsonPropertyName("currency")]
        public string Currency { get; set; }

        [JsonPropertyName("item_count")]
        public int ItemCount { get; set; }
    }

    
}