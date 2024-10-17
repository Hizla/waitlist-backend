
# Waitlist Backend

A lightweight, high-performance waitlist backend written in Go, using [Fiber v3](https://github.com/gofiber/fiber) as the web framework and [LevelDB](https://github.com/syndtr/goleveldb) for data storage. The backend supports email registration, rate limiting, and optional [hCaptcha](https://www.hcaptcha.com) verification for added security.

## Features

- **Email Registration:** Allows users to register their email for the waitlist, with each registration stored in LevelDB.
- **In-Memory Waitlist Count:** Tracks the total number of registered emails with a quick-access in-memory counter.
- **Rate Limiting:** Prevents abuse by limiting the number of registrations allowed from a single IP address within a specified time frame.
- **Optional hCaptcha Verification:** Adds an optional layer of security using hCaptcha to prevent bots from registering.
- **High Performance:** Built with Fiber v3, offering fast routing and efficient handling of concurrent requests.
- **Lightweight Storage:** Uses LevelDB to store email and timestamp data, ensuring a low overhead for database operations.

## Getting Started

### Prerequisites

- **Go 1.20+**
- **Optional:** hCaptcha API Key (if enabling hCaptcha verification)

### Installation

1. **Clone the repository:**
    ```bash
    git clone https://git.randomchars.net/hizla/waitlist/backend.git
    cd backend
    ```

2. **Build and run:**
    ```bash
    go build
    ./backend
    ```

3. **Configure the environment: (Optional)**
   - Create a `.env` file in the root directory and add the following configuration options:
     ```env
     DB=db
     VERBOSE=1
     HCAPTCHA_SECRET_KEY=unset # set your hCaptcha secret key
     HCAPTCHA_SITE_KEY=unset # set your hCaptcha site key
     LISTEN=unset
     LISTEN_ADDR=127.0.0.1:3000
     ALLOWED_URL=https://hizla.io
     ```


## Configuration Documentation

| Name                | Usage                                            | Default          |
|---------------------|--------------------------------------------------|------------------|
| DB                  | Filesystem path to database. Can be relative.    | db               |
| LISTEN              | Filesystem path to UNIX socket to listen on.     | unset            |
| LISTEN_ADDR         | Address to listen on if LISTEN is unset.         | 127.0.0.1:3000   |
| ALLOWED_URL         | Allowed URL enforced by CORS.                    | https://hizla.io |
| HCAPTCHA_SITE_KEY   | hCaptcha site key. Disables hCaptcha if unset.   | unset            |
| HCAPTCHA_SECRET_KEY | hCaptcha secret key. Disables hCaptcha if unset. | unset            |
| VERBOSE             | Verbose logging. Set 0 to disable.               | 1                |


## API Documentation

### Routes

#### 1. `POST /api/register`

Registers a new email to the waitlist.

- **Headers:**  
  `Content-Type: application/json`

- **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "hcaptcha_response": "optional-hcaptcha-response"
    }
    ```

- **Response:**
    - `200 OK`: Registration successful.
    - `400 Bad Request`: Invalid input data or hCaptcha validation failed.
    - `429 Too Many Requests`: Rate limit exceeded.
    - `500 Internal Server Error`: Server error.

- **Rate Limiting:**  
  Each IP address is limited to `5` registrations per week.

#### 2. `GET /api/count`

Returns the total number of registered emails in the waitlist.

- **Response:**
    ```
    1234
    ```
 
## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new feature branch (\`git checkout -b feature/my-feature\`).
3. Commit your changes (\`git commit -m "feat: add new feature"\`).
4. Push the branch (\`git push origin feature/my-feature\`).
5. Open a pull request.

## License

This project is licensed under the AGPLv3 License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or support, feel free to open an issue.
