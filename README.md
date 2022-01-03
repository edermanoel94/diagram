# Gen-sequence-diagram

API for generate sequence diagram

## Requirements

Golang 1.17+

## Endpoint

| Method                 | Router                           | Request Body
| -----------------------| ---------------------------------| --------------------------------------------------------------------|
| GET                    | `http://localhost:8080/health`   |                                x                                    |
| POST                   | `http://localhost:8080/download` | `{"format": "png,pdf,svg", "message": "<DSL>", "style": "default"}` |

