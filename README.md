# Receipt Processor Service

Welcome to the Receipt Processor project! This service is designed to process receipts, calculate points based on predefined rules, and allow retrieval of these points through a simple API.


## Say Hi

Hi, developers at Fetch Rewards,

Thank you for taking the time to look at my project. Developing the backend was a really fun experience, and I learned a lot in the process. I hope you find it interesting and enjoy exploring its features as much as I enjoyed building it.

Have a great day!


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

What you need to install the software:

- Go (version 1.22 or later)
- Docker (optional for containerization)

### Installing

A step-by-step guide to get a development environment running:

1. **Clone the repository**

   ```sh
   git clone https://yourrepository/receipt_processor.git
   cd receipt_processor
   ```

2. **Run the server**

   To run the server locally without Docker:

   ```sh
   go run main.go
   ```

   Your server should now be running on `http://localhost:8080`.

3. **Using Docker (Optional)**

   To build and run the application using Docker:

   ```sh
   docker build -t receipt-processor .
   docker run -p 8080:8080 receipt-processor
   ```

   This will make the service accessible at `http://localhost:8080`.

## API Usage

The service exposes two main endpoints:

- **POST** `/receipts/process` - Submit a receipt for processing.
  
  Example request:
  ```sh
  curl -X POST http://localhost:8080/receipts/process \
       -H "Content-Type: application/json" \
       -d '{"retailer":"Example Store", "purchaseDate":"2023-01-01", "purchaseTime":"15:00", "items":[{"shortDescription":"Item 1", "price":"10.00"}, {"shortDescription":"Item 2", "price":"20.00"}], "total":"30.00"}'
  ```

- **GET** `/receipts/{id}/points` - Retrieve the points for a processed receipt.
  
  Example request:
  
  ```sh
  curl http://localhost:8080/receipts/{id}/points
  ```

Replace `{id}` with the actual ID returned from the process receipt endpoint.



## Contact
Email: tyang328@wisc.edu
LinkedIn: https://www.linkedin.com/in/tongluy/

