# Diagram

API for generate sequence diagram

## Requirements

Golang 1.17+

## Endpoint

| Method                 | Router                                       | Request Body
| -----------------------| ---------------------------------------------| -----------------------------------------------------------------------|
| GET                    | `https://diagram.edermanoel.net.br/health`   |                                x                                       |
| POST                   | `https://diagram.edermanoel.net.br/download` | `{"format": "png,pdf,svg", "message": "[]string", "style": "default"}` |


Example for message in Websequencediagrams.com:

```
title Untitled

Alice->Bob: Authentication Request
note right of Bob: Bob thinks about it
Bob->Alice: Authentication Response
```

How to pass the message in ```Request Body```:

```
{
"message": ["title Untitled", "Alice->Bob: Authentication Request", "note right of Bob: Bob thinks about it", "Bob->Alice: Authentication Response"],
...
}
```
