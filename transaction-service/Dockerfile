services:
    transaction-service:
      build: ./transaction-service
      ports:
        - "8081:8081"
      environment:
        - MONGO_URI=mongodb+srv://admin:admin@cluster0.cnse2.mongodb.net/warehouse?retryWrites=true&w=majority
        - PORT=8081