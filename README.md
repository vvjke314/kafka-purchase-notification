# kafka-purchase-notification

### Homework in the discipline of Networks and Telecommunications
##### Task: Writing a payment confirmation email service using a message broker (Apache Kafka)


- [x] Zookeeper and kafka docker container
- [X] Producer to send notifications
- [X] Consumer to receive payment messages
- [X] Sending email via SMTP

To raise Zookeeper and Kafka servers. Kafka port is 29092

    docker compose up -d


Sending email's via SMTP package. Read more [here](https://pkg.go.dev/net/smtp)

Additionally, we deployed a service to simulate the purchase of goods. Service users and products have also been imitated.
