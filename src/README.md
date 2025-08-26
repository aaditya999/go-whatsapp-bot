# Whatsmeow Go App

This project is a simple Go application that utilizes the Whatsmeow package to log in to a WhatsApp account and send messages to a group chat.

## Project Structure

```
whatsmeow-go-app
├── src
│   ├── main.go          # Entry point of the application
│   ├── config           # Configuration settings
│   │   └── config.go    # Configuration struct and loading function
│   ├── whatsapp         # WhatsApp client functionality
│   │   └── client.go     # WhatsApp client methods
│   └── utils            # Utility functions
│       └── logger.go     # Logging utilities
├── go.mod               # Module definition and dependencies
├── go.sum               # Dependency checksums
└── README.md            # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd whatsmeow-go-app
   ```

2. **Install dependencies:**
   Make sure you have Go installed, then run:
   ```
   go mod tidy
   ```

3. **Configuration:**
   Update the `src/config/config.json` file with your WhatsApp credentials and group chat identifiers.

4. **Run the application:**
   ```
   go run src/main.go
   ```

## Usage

After running the application, it will log in to the specified WhatsApp account and send a "Hello, World!" message to the configured group chat.
