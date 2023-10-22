# Go Base Backend

The Go Base Backend serves as a solid foundation for your Go-based projects. You can use it as a starting point and customize it according to your needs.

## Features

This base backend offers the following key features:

- Code generation from SQL with [sqlc](https://docs.sqlc.dev/en/stable/index.html).
- Database migration capabilities using [goose](https://github.com/pressly/goose).
- Integration of PostgreSQL and Redis with Docker.
- Web server powered by [Echo](https://echo.labstack.com/).
- Session-based authentication.
- Role-based permissions for users.
- Support for Cron Jobs.
- Comprehensive logging system.

## Dependencies

Please make sure you have the following dependencies installed:

- [Git](https://github.com/git-guides/install-git)
- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/) (for running databases in a development environment)
- [Make](https://www.gnu.org/software/make/) (used for predefined commands in [Makefile](/Makefile))

## Installation and Setup

To get started with the Go Base Backend, follow these steps:

1. Install all the [dependencies](#dependencies).

2. Clone the project repository from GitHub and install the required packages:

   ```bash
   git clone https://github.com/JK-1117/go-base.git
   cd go-base
   go mod tidy
   ```

3. Configure the project to match your requirements. Refer to the [Congifurations Section](#configurations).

4. Start the database services with Docker:

   ```bash
   docker compose up
   ```

5. Build and run the application:

   ```bash
   make server
   ```

6. Access the application in your web browser at `http://localhost:8080`.

## Configurations

Here are some optional configurations that you may want to customize for your project:

- [Makefile](/Makefile)

  - Configure the Postgres URL, change the user, password, and DB name.
  - Modify the main folder and the binary name to match your app's name.

- [.env](/backend/.env)

  - Required variables include:
    - APPNAME - the name of your project (e.g., `APPNAME=base`).
    - PORT - the port where your backend will run (e.g., `PORT=8080`).
    - DOMAIN - the domain for your production backend (e.g., `DOMAIN=localhost`).
    - REDIS_URL - the URL for connecting to your Redis instance (e.g., `REDIS_URL=redis://localhost:6379/0`).
    - DB_URL - the URL for connecting to your PostgreSQL database (e.g., `DB_URL=postgres://ops:OnlyADevPasswOrD@localhost:5432/ops?sslmode=disable`).
    - COOKIE_HASHKEY - the hashKey used to authenticate the cookie value using HMAC (for [more information](https://github.com/gorilla/securecookie#examples), e.g., `COOKIE_HASHKEY=veryverylongsecretwith64length`).
    - COOKIE_BLOCKKEY - the blockKey used to encrypt the cookie value (for [more information](https://github.com/gorilla/securecookie#examples) e.g., `COOKIE_BLOCKKEY=secretwith32length`).

- [docker-compose.yaml](/docker-compose.yaml)

  - Configure environment variables and ports for the services.

- Module name

  - Replace "github.com/jk1117/go-base" with your module name.

- Session Cookie Name
  - You can change the session ID's cookie name by modifying the SESSIONCOOKIE constant in [session.go](/backend//internal/server/session.go)

## Contributing

Contributions to this repository are welcome. If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the [MIT License](https://opensource.org/license/mit/).

## Acknowledgements

This project was inspired by an [Amazing Free Course](https://www.youtube.com/watch?v=un6ZyFkqFKo&t=32565s) by [freeCodeCamp.org](https://www.youtube.com/@freecodecamp) and [bootdotdev](https://www.youtube.com/@bootdotdev). Check out their amazing content to learn more. Special thanks to the developers and all the open-source contributors whose libraries and frameworks have been used in this project.

## Contact

For any questions or inquiries, please contact the project maintainer at `chun11197@gmail.com`.
