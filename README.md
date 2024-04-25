# Pokemon Card Trading API

This API is meant to provide a REST API to interact with a Pokemon card trading Discord bot.

## Usage

These are the endpoints for the API.

Do note that the following information is not maintained and may be out of date.

```bash
GET /cards       # Retrieves all of the Pokemon cards in the database;
                 # It can be filtered by rarity (ex. /cards?rarity=common)

GET /cards/{id}  # Retrieves a Pokemon card by ID

GET /sets        # Retrieves all of the Pokemon sets in the database

GET /sets/{name} # Retrieves a Pokemon set by name
```
