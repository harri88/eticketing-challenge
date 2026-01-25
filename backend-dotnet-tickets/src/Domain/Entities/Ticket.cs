namespace Domain.Entities
{
    // DDD: The Entity holds the logic, not just data.
    public class Ticket
    {
        // Private setters to ensure encapsulation (Ticket can only be modified via methods if we added write logic)
        public string Id { get; private set; }
        public string Name { get; private set; }
        public decimal Price { get; private set; }
        public string Currency { get; private set; }
        public int Quota { get; private set; }

        public Ticket() { } // Parameterless constructor for Dapper
        public Ticket(string id, string name, decimal price, string currency, int quota)
        {
            Id = id;
            Name = name;
            Price = price;
            Currency = currency;
            Quota = quota;
        }

        



        // DDD Logic: Availability is a derived state, not a database column.
        public bool IsAvailable() => Quota > 0;
    }
}