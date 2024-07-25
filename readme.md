## Wallet Core - Event Driven Architecture

## Responsability:
- It produces to an Apache Kafka topic, notificating when new transactions are made

- It exposes HTTP endpoints responsable for creating clients, accounts e transactions

## Seting up environment
- It uses a specific docker network to comunicate with other containers, so before "compose up", make sure to create the network correctly
  > docker network create wallet-network

- Once the network is created, run containers
  > docker compose up -d

- Make sure Apache Kafka topics were created
  - go to http://localhost:9021 for Control Center or http://localhost:8000 for Kafka Ui  
    > create both 'balances' and 'transactions' topics, with default configs

## Start application
- access go container and run main.go
  > docker exec -it wallet-core bash
  > go run cmd/walletcore/main.go
