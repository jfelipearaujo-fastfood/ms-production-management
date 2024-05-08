#!/bin/sh

localstack_url=http://localhost:4566
queue_name=OrderProductionQueue

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

queue_url=$(aws sqs get-queue-url --endpoint-url "$localstack_url" --output text --queue-name "$queue_name")

if [ $? -eq 0 ]; then
    echo "Queue URL: $queue_url"
    echo "Sending a message..."

    message='{
        "order_id": "c3fdab1b-3c06-4db2-9edc-4760a2429462",
        "items": [
            {
                "id": "cfdab175-1f86-4fb0-9bcb-15f2c58df30c",
                "name": "Hamburger",
                "quantity": 1
            }
        ]
    }'

    # Publish the message to the queue
    aws sqs send-message \
        --endpoint-url "$localstack_url" \
        --queue-url "$queue_url" \
        --output text \
        --message-body "$message" > /dev/null

    # Check if the message publishing was successful
    if [ $? -eq 0 ]; then
        echo "Message published successfully."
    else
        echo "Failed to publish message."
    fi
else
    echo "Failed to retrieve the queue URL."
fi