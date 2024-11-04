# GoTorrent - A Torrent Client and Tracker

GoTorrent is a lightweight torrent client and tracker written in Go. It supports basic torrenting functionalities, including peer management, file downloading, and web-based UI. This project is ideal for learning about the BitTorrent protocol, peer-to-peer networking, and distributed file sharing.

## Features

- **Torrent Tracker**: Tracks peers connected to a torrent and helps clients find each other.
- **Torrent Client**: Connects to peers, requests pieces, and assembles the complete file.
- **Web-Based UI**: Launches a local web server for an interactive user interface.
- **Torrent File Creation**: Allows users to create `.torrent` files for file sharing.

## Getting Started

### Prerequisites

- **Go**: [Download and install Go](https://golang.org/dl/) (Go 1.16 or higher recommended).
- **Git**: Version control system for cloning the repository.

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/GoTorrent.git
   cd GoTorrent
   ```

2. **Build the tracker and client**:
   ```bash
   go build -o tracker ./tracker
   go build -o client ./client
   ```

### Usage

#### Running the Tracker

Start the tracker to allow peers to connect and announce their presence.

```bash
./tracker
```

The tracker will run on `http://localhost:8080` by default. You can change the port in the `main.go` file if needed.

#### Running the Torrent Client

Start the client to connect to the tracker, download files, and manage connections.

```bash
./client -torrent ./path/to/file.torrent
```

#### Web UI

To launch the web-based UI, run the client, and a browser tab should open automatically at `http://localhost:3000`. From here, you can monitor active downloads, add new torrents, and view download progress.

### Creating a Torrent File

Generate a `.torrent` file for a specific file:

```bash
./client -create ./path/to/file
```

This will generate a `.torrent` file in the same directory, which you can then share with others.

### Example Workflow

1. **Run the tracker** on your server.
2. **Create a torrent file** for a file you want to share.
3. **Distribute the torrent file** to users.
4. **Users start the client** with the torrent file, connect to the tracker, and begin downloading from peers.

## Project Structure

- **`client/`**: Contains code for the torrent client, including peer communication and file downloading.
- **`tracker/`**: Contains code for the torrent tracker, managing peer connections and announcements.
- **`web-ui/`**: Code for the web-based UI to control and monitor downloads.
- **`testdata/`**: Sample `.torrent` files for testing purposes.

## Configuration

Modify settings in the `config.yaml` file for the tracker and client, such as:
- Tracker announce interval
- Max peers
- Listening ports

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Acknowledgments

- BitTorrent protocol documentation
- Go library contributors and open-source community

--- 

This README gives an overview of the project, explains how to set it up and use it, and describes the project structure. Let me know if you'd like to add any further details!
