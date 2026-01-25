using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Application.DTOs
{
    public class ConfirmPaymentReq
    {
        [JsonPropertyName("order_id")]
        public required string OrderId { get; set; }
    }
}