# URL Shortener 

This is the backend API for a small URL shortener application. The API allows users to submit a long URL and receive back a shorter code that can be used to redirect to the original link. The system saves all short URLs in a PostgreSQL database and also tracks when someone clicks a short link.

Right now the API supports the following features:

Creating a short URL (with an auto-generated code or with a custom code provided by the user)

Optional expiration date for each short link

Redirecting from /code to the original URL

Storing click information (IP, user agent, and a country placeholder)

Getting analytics about a specific short code (total clicks, number of unique IPs, clicks per country)

Generating a QR code for any short link

Basic rate limiting per IP (to avoid spamming the service)

Automatic table creation using GORM

The API is written in Go using the standard net/http package. PostgreSQL is used as the database, and GORM is used as the ORM so the tables are created automatically. Docker Compose is used only for running the database and Adminer.


Start the database using Docker Compose:
docker compose up -d



How to start the app:


Run the Go API: cd in api folder and
go run .