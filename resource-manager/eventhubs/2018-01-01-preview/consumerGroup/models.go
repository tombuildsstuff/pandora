package consumerGroup

type CreateConsumerGroupInput struct {
	Properties CreateConsumerGroupProperties `json:"properties"`
}

type CreateConsumerGroupProperties struct {
	UserMetadata *string `json:"userMetadata,omitempty"`
}

type GetConsumerGroup struct {
	Properties GetConsumerGroupProperties `json:"properties"`
}

type GetConsumerGroupProperties struct {
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	UserMetadata string `json:"userMetadata"`
}
