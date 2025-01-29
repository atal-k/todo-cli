# todo-cli

A simple, lightweight command-line task manager written in Go.

## Features

- ✨ Simple and intuitive CLI interface
- 📝 Create and manage tasks
- ✅ Mark tasks as complete/incomplete
- 🔍 Filter tasks by status
- 💾 Persistent storage using SQLite
- 🚀 No external configuration needed

## Installation

### Prerequisites
- Go 1.16 or higher
- SQLite3

### From Source
```bash
# Clone the repository
git clone https://github.com/atal-k/todo-cli
cd todo-cli

# Build the application
go build -o todo cmd/main.go

# Optional: Move to PATH (Linux/macOS)
sudo mv todo /usr/local/bin/
```

### Using Go Install
```bash
go install github.com/atal-k/todo-cli@latest
```

## Usage

```bash
# Add a new task
todo add "Read notes"

# List all tasks
todo list

# List pending tasks
todo list --pending

# List completed tasks
todo list --completed

# Mark task as done
todo done 1

# Delete a task
todo delete 1

# Clear all tasks
todo clear
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `add`   | `a`   | Add a new task |
| `list`  | `ls`  | List tasks |
| `done`  | `d`   | Mark task as completed |
| `delete`| `rm`  | Delete a task |
| `clear` | -     | Remove all tasks |

## Project Structure
```
todo-cli/
├── cmd/
│   └── main.go         # Entry point
├── internal/
│   ├── db/            # Database operations
│   └── todo/          # Business logic
└── README.md
```

## Development

```bash
# Get dependencies
go mod tidy

# Run tests
go test ./...

# Build for development
go build -o todo cmd/main.go
```

## Storage

Tasks are stored in SQLite database at:
- Linux/macOS: `~/.todo.db`
- Windows: `%USERPROFILE%\.todo.db`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with [urfave/cli](https://github.com/urfave/cli)
- SQLite storage using [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)