# 🌧️ Raindrop MCP Server

> ## ⚠️ ALPHA VERSION
> **This is an early alpha release! Features may be incomplete, APIs may change, and bugs are expected.**
> **Use at your own risk. Feedback and contributions welcome!**

---

A full-featured MCP (Model Context Protocol) server for [Raindrop.io](https://raindrop.io) bookmark management.

## Features

**20 Tools:**
- **Bookmarks**: create, get, update, delete, search
- **Collections**: create, get, update, delete, merge, list
- **Tags**: list, rename, delete, merge, suggest
- **Highlights**: get, create, delete
- **Filters**: get filters for collection
- **User**: get user info

**4 Resources:**
- `raindrop://collections` - All collections
- `raindrop://tags` - All tags
- `raindrop://user` - User info
- `raindrop://collection/{id}/bookmarks` - Bookmarks in collection

## Installation

### Option 1: Download from Releases

Download from [GitHub Releases](https://github.com/FyziGo/Raindrop-MCP/releases/latest):

| File | Description |
|------|-------------|
| `raindrop-mcp.exe` | Windows binary |
| `raindrop-mcp-docker.zip` | Docker image |
| `raindrop-mcp.mcpb` | MCPB installer |

### Option 2: MCPB Installer

```bash
# Install using mcpb CLI
mcpb install raindrop-mcp.mcpb
```

Or double-click `raindrop-mcp.mcpb` if you have MCPB Desktop installed.

### Option 3: Docker

```bash
# From release (download and extract docker.zip first)
docker load -i raindrop-mcp-docker.tar
docker run -i --rm -e RAINDROP_TOKEN=your_token raindrop-mcp
```

### Option 4: Build from Source

```bash
git clone https://github.com/FyziGo/Raindrop-MCP.git
cd raindrop-mcp
go build -o raindrop-mcp.exe .
```

## Authentication

1. Go to https://app.raindrop.io/settings/integrations
2. Scroll to "Test token" section
3. Create and copy your token

## Claude Desktop Config

Add to `%APPDATA%\Claude\claude_desktop_config.json` (Windows) or `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):

```json
{
  "mcpServers": {
    "raindrop": {
      "command": "/path/to/raindrop-mcp.exe",
      "env": {
        "RAINDROP_TOKEN": "your_token_here"
      }
    }
  }
}
```

### With Docker

```json
{
  "mcpServers": {
    "raindrop": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-e", "RAINDROP_TOKEN=your_token_here", "fyzigo/raindrop-mcp"]
    }
  }
}
```

<details>
<summary>Advanced: OAuth2 Authentication</summary>

For users who prefer OAuth2 (requires creating your own Raindrop app):

1. Create app at https://app.raindrop.io/settings/integrations
2. Get your Client ID and Client Secret
3. Configure:

```json
{
  "mcpServers": {
    "raindrop": {
      "command": "/path/to/raindrop-mcp.exe",
      "env": {
        "RAINDROP_CLIENT_ID": "your_client_id",
        "RAINDROP_CLIENT_SECRET": "your_client_secret"
      }
    }
  }
}
```

On first run, browser opens for authorization. Token saved to `~/.raindrop-mcp/token.json`.

</details>

## All Tools

| Category | Tool | Description |
|----------|------|-------------|
| **Bookmarks** | `create-bookmark` | Create bookmark with URL, title, tags |
| | `get-bookmark` | Get bookmark by ID |
| | `update-bookmark` | Update title, note, tags, collection |
| | `delete-bookmark` | Delete bookmark |
| | `search-bookmarks` | Search with query, filters |
| **Collections** | `list-collections` | List all collections |
| | `create-collection` | Create new collection |
| | `get-collection` | Get collection by ID |
| | `update-collection` | Update collection |
| | `delete-collection` | Delete collection |
| | `merge-collections` | Merge multiple into one |
| **Tags** | `list-tags` | List all tags |
| | `rename-tag` | Rename a tag |
| | `delete-tags` | Delete tags |
| | `merge-tags` | Merge tags into one |
| | `suggest-tags` | Get tag suggestions for URL |
| **Highlights** | `get-highlights` | Get highlights from bookmark |
| | `create-highlight` | Create highlight |
| | `delete-highlight` | Delete highlight |
| **Other** | `get-filters` | Get collection filters |
| | `get-user` | Get user info |

## Example Prompts

- "List my Raindrop collections"
- "Search bookmarks for python tutorial"
- "Create collection 'Reading List'"
- "Rename tag 'dev' to 'development'"
- "Get highlights from bookmark 12345"

## Project Structure

```
raindrop-mcp/
├── main.go
├── api/
│   ├── raindrop.go
│   ├── collections.go
│   └── extended.go
├── tools/
│   ├── tools.go
│   └── extended.go
├── resources/
│   └── resources.go
└── types/
    ├── types.go
    └── extended.go
```

## License

MIT
