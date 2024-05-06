package order_entity

type OrderState int

const (
	None       OrderState = iota
	Received              // When the order is received and is ready to be processed by the kitchen
	Processing            // When the order is being processed by the kitchen
	Completed             // When the order is completed and ready to be delivered
	Delivered             // When the order is delivered to the customer
	Cancelled             // When the order is cancelled
)

var (
	order_state_machine = map[OrderState][]OrderState{
		None:       {Received},
		Received:   {Processing, Cancelled},
		Processing: {Completed, Cancelled},
		Completed:  {Delivered},
	}
)

func NewOrderState(title string) OrderState {
	state, ok := map[string]OrderState{
		"Received":   Received,
		"Processing": Processing,
		"Completed":  Completed,
		"Delivered":  Delivered,
		"Cancelled":  Cancelled,
	}[title]
	if !ok {
		return None
	}

	return state
}

func (s OrderState) CanTransitionTo(to OrderState) bool {
	for _, allowed := range order_state_machine[s] {
		if to == allowed {
			return true
		}
	}
	return false
}

func (s OrderState) String() string {
	text, ok := map[OrderState]string{
		None:       "None",
		Received:   "Received",
		Processing: "Processing",
		Completed:  "Completed",
		Delivered:  "Delivered",
		Cancelled:  "Cancelled",
	}[s]
	if !ok {
		return "Unknown"
	}

	return text
}

func IsValidState(s OrderState) bool {
	return s >= Received && s <= Cancelled
}
