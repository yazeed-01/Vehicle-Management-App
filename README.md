# Vehicle Management App

A secure, scalable vehicle management application built with a **Flutter frontend** and a **Go backend**. This app allows users to search for car owners using plate numbers, communicate with them via phone or email, and ensures security through mandatory login and contact info updates. It’s optimized for big data (tested with 1M+ rows) and tracks search history for security purposes.

## Features

- **Fast Search**: Quickly search for car owners by plate number, even with large datasets (1M+ rows), thanks to Redis caching and efficient database queries.
- **Secure Login**: Users must log in with their national number and a specific key, ensuring only authorized access.
- **Mandatory Contact Info**: After login, users without an email or phone number are prompted to add them, enhancing communication and security.
- **Government Integration**: Cars are linked to users via national numbers, simulating a government-connected database.
- **Communication**: Contact car owners via phone or email directly from the app.
- **Search History Tracking**: Silently saves search history (user ID, searched user, time, and count) for security auditing, without notifying the user.
- **Big Data Ready**: Tested with 1M randomly generated rows using Go scripts, ensuring scalability.

## Tech Stack

- **Frontend**: Flutter (Dart) - Cross-platform UI for Android, iOS, and Web.
- **Backend**: Go (Gin framework) - High-performance server with RESTful APIs.
- **Database**: PostgreSQL - Stores users, cars, and search history.
- **Caching**: Redis - Speeds up searches with in-memory caching.
- **Authentication**: JWT (JSON Web Tokens) - Secures API endpoints.
- **ORM**: GORM - Simplifies database interactions in Go.

