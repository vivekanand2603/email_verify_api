Sure! Here is a README file in table format for your Go application:

# Email Verification Service

## Overview

This service provides APIs for managing email lists and leads, performing CRUD operations, and various counts related to email verification statuses. It also supports CSV file uploads and downloads for leads.

## Routes

### Lists

| Method | Endpoint                  | Description                                    |
|--------|---------------------------|------------------------------------------------|
| GET    | /lists                    | Retrieve all lists.                            |
| POST   | /lists                    | Create a new list.                             |
| GET    | /lists/:id                | Retrieve a list by ID.                         |
| DELETE | /lists/:id                | Delete a list by ID.                           |

### Leads

| Method | Endpoint                      | Description                                            |
|--------|-------------------------------|--------------------------------------------------------|
| GET    | /leads                        | Retrieve all leads.                                    |
| POST   | /leads                        | Create a new lead.                                     |
| GET    | /leads/:id                    | Retrieve a lead by ID.                                 |
| DELETE | /leads/:id                    | Delete a lead by ID.                                   |
| GET    | /lists/:id/leads              | Retrieve all leads in a list.                          |
| GET    | /lists/:id/leads/count        | Count leads in a list.                                 |
| GET    | /lists/:id/leads/count/email_verified | Count email verified leads in a list.            |
| GET    | /lists/:id/leads/count/valid_emails   | Count valid email leads in a list.               |
| GET    | /lists/:id/leads/count/invalid_emails | Count invalid email leads in a list.             |
| GET    | /lists/:id/leads/count/unknown_emails | Count unknown email leads in a list.             |
| GET    | /count_all                   | Count all emails with a specific status.               |
| POST   | /lists/:id/queue              | Add a list to the queue.                               |
| GET    | /lists/:id/queue              | Check if a list is in the queue.                       |
| DELETE | /lists/:id/queue              | Remove a list from the queue.                          |
| GET    | /processQueue                 | Process the queue.                                     |
| POST   | /lists/:id/leads/csv          | Upload leads from a CSV file to a list.                |
| GET    | /lists/:id/leads/csv          | Download leads of a list as a CSV file.                |

## Installation

1. Clone the repository.
2. Install dependencies using `go mod tidy`.
3. Run the application using `go run main.go`.

## Usage

The application runs on port 30001. Use an API client like Postman to interact with the endpoints.

## Middleware

- **CORS**: Allows all origins and supports various HTTP methods.
- **Content-Type**: Sets the response content type to JSON.

## Database

This application uses MongoDB. Update the connection string in the `mongo.Connect` function call.

```go
client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://<username>:<password>@cluster0.isymvpw.mongodb.net/email_verify2"))
```

## Handlers

Handlers for each route are defined within the `main` function, interacting with MongoDB to perform CRUD operations and more.

## Contributing

Feel free to fork this repository and contribute by submitting a pull request.

## License

This project is licensed under the MIT License.

---

This table provides a clear and organized overview of your application's endpoints and functionalities, making it easier for developers to understand and use your API.