Here is the API documentation in table format:

### Lists API

| HTTP Method | Endpoint                 | Description                     | Request Body                            | Response Body                 |
|-------------|--------------------------|---------------------------------|-----------------------------------------|-------------------------------|
| GET         | /lists                   | Get all lists                   | None                                    | List of all lists             |
| POST        | /lists                   | Create a new list               | `{"name": "list_name"}`                 | Created list with ID          |
| GET         | /lists/:id               | Get a list by ID                | None                                    | List with specified ID        |
| DELETE      | /lists/:id               | Delete a list by ID             | None                                    | None                          |
| GET         | /lists/:id/leads         | Get all leads in a list by ID   | None                                    | List of leads in the list     |
| GET         | /lists/:id/leads/count   | Get lead count by list ID       | None                                    | Lead count                    |
| GET         | /lists/:id/leads/count/email_verified | Get verified lead count by list ID | None                                    | Verified lead count           |
| POST        | /lists/:id/queue         | Add a list to the queue         | None                                    | None                          |
| POST        | /lists/:id/leads/csv     | Upload leads from CSV to a list | CSV file                                | None                          |
| GET         | /lists/:id/leads/csv     | Download leads as CSV           | None                                    | CSV file                      |

### Leads API

| HTTP Method | Endpoint                 | Description                     | Request Body                            | Response Body                 |
|-------------|--------------------------|---------------------------------|-----------------------------------------|-------------------------------|
| GET         | /leads                   | Get all leads                   | None                                    | List of all leads             |
| POST        | /leads                   | Create a new lead               | `{"email": "email", "list_id": "list_id"}` | Created lead with ID          |
| GET         | /leads/:id               | Get a lead by ID                | None                                    | Lead with specified ID        |
| DELETE      | /leads/:id               | Delete a lead by ID             | None                                    | None                          |

### Queue Processing API

| HTTP Method | Endpoint                 | Description                     | Request Body                            | Response Body                 |
|-------------|--------------------------|---------------------------------|-----------------------------------------|-------------------------------|
| GET         | /processQueue            | Process the lead queue          | None                                    | None                          |

### Common Responses

| Status Code | Description                     |
|-------------|---------------------------------|
| 200         | OK                              |
| 201         | Created                         |
| 204         | No Content                      |
| 400         | Bad Request                     |
| 500         | Internal Server Error           |

This table provides an overview of the available API endpoints, their HTTP methods, descriptions, request and response bodies, and common response status codes.