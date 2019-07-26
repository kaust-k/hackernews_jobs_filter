# Hacker News Post filter

Command line arguments:
- POSTGRES_URI: Postgres database URL
- MIGRATION_VERSION: Database migration version (should be 1 for now)

$ POSTGRES_URL="insert_url_here" MIGRATION_VERSION=1 go run main.go


## Use case:

Filter and show remote jobs from Hacker News posts [Ask HN: Who is hiring](https://www.google.com?q=hn+who+is+hiring)

Main post looks something like this: https://news.ycombinator.com/item?id=20325925

- Main post and all of its children (only 1 level of depth) are stored locally in database.
- New posts (not present in db) are fetched on every run.

- Add new story ID in main.go#31
- Add new conditions to change filtered results in services/handler/http.go#64

- Open [localhost](http://localhost:9999) to see filtered results.