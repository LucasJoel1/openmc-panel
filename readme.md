# OpenMC Panel

A comprehensive Minecraft server management panel with Discord bot integration and a modern web interface. This application allows you to manage your Minecraft server through both a Discord bot and a web-based control panel.

## Features

### 🎮 Server Management

- **Start/Stop/Restart** your Minecraft server remotely
- **Real-time server monitoring** with WebSocket connections
- **Server properties editor** for configuration management
- **Memory allocation controls** (min/max heap settings)
- **Server version tracking** (Minecraft, modloader, and modloader version)

### 🤖 Discord Integration

- **Chat bridge** between Discord and Minecraft
- **Console commands** through Discord channels
- **Command system** with prefix support (`!command`)
- **Server status notifications** and monitoring
- **Player management** through Discord

### 🌐 Web Interface

- **Modern React-based UI** built with TypeScript and Vite
- **Real-time server status** dashboard
- **Player management** interface
- **Server console** access through web browser
- **File manager** for server files
- **Dark/Light theme** support with Tailwind CSS
- **Responsive design** with mobile support

### 📊 Monitoring & Analytics

- **Server performance monitoring** with system resource tracking
- **Player activity tracking** (join times, ping, IP addresses)
- **Server ping monitoring** and uptime tracking
- **Crash report handling** and logging

## Tech Stack

### Backend (Go)

- **Go 1.24.1** - Main backend language
- **Discord API** integration via `bwmarrin/discordgo`
- **WebSocket support** for real-time communication
- **HTTP server** for web API and static file serving
- **System monitoring** with `shirou/gopsutil`
- **Image processing** with `disintegration/imaging`

### Frontend (React + TypeScript)

- **React 19** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **Radix UI** components for accessible UI elements
- **TanStack Query** for data fetching and caching
- **Lucide React** for icons
- **Sonner** for toast notifications

## Prerequisites

- **Go 1.24.1** or higher
- **Node.js** and **pnpm** for frontend development
- **Discord Bot Token** (for Discord integration)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/lucasjoel1/openmc-panel.git
cd openmc-panel
```

### 2. Backend Setup

```bash
# Install Go dependencies
go mod download

# Create environment file
cp .env.example .env
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory:

```env
DISCORD_TOKEN=your_discord_bot_token_here
CONSOLE_CHANNEL=discord_channel_id_for_console
CHAT_CHANNEL=discord_channel_id_for_chat
```

### 4. Configure Server Settings

Edit `server.json` to match your server setup:

```json
{
  "serverName": "Your Minecraft Server",
  "jarName": "server.jar",
  "path": "/path/to/your/minecraft/server/",
  "memory": ["2G", "4G"],
  "gui": false
}
```

### 5. Frontend Setup

```bash
cd frontend
pnpm install
pnpm run build
```

### 6. Place Minecraft Server

Place your Minecraft server files in the `server/` directory, including:

- Server JAR file
- `server.properties`
- World files
- Mods (if using modded server)

## Usage

### Start the Application

```bash
# From the root directory
go run src/main.go
```

The application will:

- Start the Discord bot
- Launch the web server on `http://localhost` (port 80)
- Begin monitoring server status

### Web Interface

Navigate to `http://localhost` to access the web panel with these sections:

- **Status** - Real-time server information and player count
- **Server Manager** - Start/stop/restart controls
- **Console** - Direct server console access
- **Player Management** - View and manage connected players
- **Edit Properties** - Modify server.properties
- **File Manager** - Browse and edit server files

### Discord Commands

Use these commands in your Discord server:

```
!start     - Start the Minecraft server
!stop      - Stop the Minecraft server
!restart   - Restart the Minecraft server
!status    - Get current server status
!players   - List online players
```

## Development

### Backend Development

```bash
# Run with auto-reload during development
go run src/main.go
```

### Frontend Development

```bash
cd frontend
pnpm run dev
```

### Building for Production

```bash
# Build frontend
cd frontend
pnpm run build

# Build backend
go build -o openmc-panel src/main.go
```

## Project Structure

```
openmc-panel/
├── src/                    # Go backend source
│   ├── main.go            # Application entry point
│   ├── globals/           # Global state management
│   ├── listeners/         # Discord and web event handlers
│   ├── server/           # Minecraft server management
│   └── web/              # Web server and API routes
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom React hooks
│   │   └── lib/          # Utility functions
│   └── public/           # Static assets
├── server/               # Minecraft server files
└── server.json          # Server configuration
```

## API Endpoints

- `GET /` - Serve web interface
- `POST /api/startServer` - Start Minecraft server
- `POST /api/stopServer` - Stop Minecraft server
- `POST /api/restartServer` - Restart Minecraft server
- `WS /api/ws/serverInfo` - WebSocket for real-time server info

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Troubleshooting

### Common Issues

**Server won't start:**

- Check that the JAR file path in `server.json` is correct
- Ensure sufficient memory is allocated
- Verify Java is installed and accessible

**Discord bot not responding:**

- Verify `DISCORD_TOKEN` is correct in `.env`
- Check bot permissions in Discord server
- Ensure bot has message content intent enabled

**Web interface not loading:**

- Check if port 80 is available
- Run frontend build process
- Verify API endpoints are accessible

## Support

For support, please open an issue on GitHub.
