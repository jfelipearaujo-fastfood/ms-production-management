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
        "Type" : "Notification",
        "MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
        "TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
        "Message" : "{\"order_id\":\"c3fdab1b-3c06-4db2-9edc-4760a2429462\",\"items\":[{\"id\": \"cfdab175-1f86-4fb0-9bcb-15f2c58df30c\",\"name\": \"Hamburger\",\"quantity\": 1}]}",
        "Timestamp" : "2024-05-19T02:01:36.927Z",
        "SignatureVersion" : "1",
        "Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
        "SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
        "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
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