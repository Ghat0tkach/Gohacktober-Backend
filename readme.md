
# GoHacktober Backend

GoHacktober is a backend service designed for managing and tracking Hacktoberfest contributions. This service leverages the GitHub API to retrieve repositories and contributions tagged with Hacktoberfest.

```
    GGGG   OOO      H   H   AAAAA   CCCCC   K   K   TTTTT   OOO   BBBBB   EEEEE   RRRR  
   G      O   O     H   H   A   A   C       K  K      T    O   O  B    B  E       R   R
   G  GG  O   O     HHHHH   AAAAA   C       KKK       T    O   O  BBBBB   EEEE    RRRR 	
   G   G  O   O     H   H   A   A   C       K  K      T    O   O  B    B  E       R  R 
   GGGG    OOO      H   H   A   A   CCCCC   K   K     T     OOO   BBBBB   EEEEE   R   R

    BBBBB    AAAAA   CCCCC   K   K   EEEEE   N   N   DDDD  
    B    B   A   A   C       K  K    E       NN  N   D   D
    BBBBB    AAAAA   C       KKK     EEEE    N N N   D   D
    B    B   A   A   C       K  K    E       N  NN   D   D
    BBBBB    A   A   CCCCC   K   K   EEEEE   N   N   DDDD  
```


## Features

- Fetch Hacktoberfest-labeled repositories from both organizations and individual users.
- Track user contributions (issues, pull requests) to these repositories.


## Project Structure

```
Gohacktober-Backend/
├── cmd/                 # Command-line related files
│   └── server/          # Server entry point
├── config/              # Configuration files
├── internal/            # Internal libraries
│   ├── handlers/        # Handlers for GitHub API
│   └── github/          # GitHub API client
├── go.mod               # Module dependencies
└── .env                 # Environment variables (e.g., GitHub token)
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [GitHub OAuth Token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Ghat0tkach/Gohacktober-Backend.git
   cd Gohacktober-Backend
   ```

2. Set up environment variables in the `.env` file:
   ```
   GITHUB_TOKEN=your_personal_github_token
   GITHUB_ORG=Your organization or account for which you want contributions to be fetched from
   PORT=8080

   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

- **/api/hacktober-fest-contributions**: Fetch Hacktoberfest contributions for a user.
  
  Example:
  ```bash
  curl http://localhost:8080/api/hacktoberfest-contributions?username={Username}
  ```

## License

MIT License

---