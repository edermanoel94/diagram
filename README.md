# Gen-sequence-diagram

API for generate sequence diagram

## Requirements

Golang 1.17+

## Endpoint

| Method                 | Router                   | Request Body
| -----------------------| -------------------------| --------------------------------------------------------------------|
| POST                   | `http://localhost:8080/` | `{"format": "png,pdf,svg", "message": "<DSL>", "style": "default"}` |

