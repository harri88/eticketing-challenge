using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Application.DTOs
{
    public class ConfirmPaymentResponse
    {
        public string OrderId { get; set; }
        public string Status { get; set; }
    }
}