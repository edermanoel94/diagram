# Gen-sequence-diagram

API for generate sequence diagram

## Requirements

Golang 1.17+

## Endpoint

| Method                 | Router                                       | Request Body
| -----------------------| ---------------------------------------------| --------------------------------------------------------------------|
| GET                    | `https://diagram.edermanoel.net.br/health`   |                                x                                    |
| POST                   | `https://diagram.edermanoel.net.br/download` | `{"format": "png,pdf,svg", "message": "<DSL>", "style": "default"}` |


Example for message DSL:

```
title Untitled

Alice->Bob: Authentication Request
note right of Bob: Bob thinks about it
Bob->Alice: Authentication Response
```
