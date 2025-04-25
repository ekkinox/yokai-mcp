# Yokai MCP Server Demo

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go version](https://img.shields.io/badge/Go-1.24-blue)](https://go.dev/)
[![Documentation](https://img.shields.io/badge/Doc-online-cyan)](https://ankorstore.github.io/yokai/)

> [MCP server](https://modelcontextprotocol.io/introduction) demo application, based on the [Yokai](https://github.com/ankorstore/yokai) Go framework.

<!-- TOC -->
* [Documentation](#documentation)
* [Overview](#overview)
  * [Layout](#layout)
  * [Makefile](#makefile)
* [Usage](#usage)
  * [Start the MCP server](#start-the-mcp-server)
  * [Configure your MCP client](#configure-your-mcp-client)
<!-- TOC -->

## Documentation

For more information about the [Yokai](https://github.com/ankorstore/yokai) framework, you can check its [documentation](https://ankorstore.github.io/yokai).

For MCP, you can check its [documentation](https://modelcontextprotocol.io/introduction).

## Overview

This demo application provides an application that offers books management, and:

- exposes via HTTP API the list of books
- exposes via MCP tools for LLMs to list, create and delete books


### Layout

This template is following the [recommended project layout](https://go.dev/doc/modules/layout#server-project):

- `cmd/`: entry points
- `configs/`: configuration files
- `db/`: database migration files
- `internal/`:
  - `domain/`: Books management domain (model, repo, service)
  - `http/`: HTTP API
  - `mcp/`: MCP API
  - `bootstrap.go`: bootstrap
  - `register.go`: dependencies registration
- `pkg`:
  - `mcp`: Work in progress Yokai MCP module

### Makefile

This template provides a [Makefile](Makefile):

```
make up     # start the docker compose stack
make down   # stop the docker compose stack
make logs   # stream the docker compose stack logs
make fresh  # refresh the docker compose stack
make test   # run tests
make lint   # run linter
```

## Usage

### Start the MCP server

After cloning the repository, simply run

```shell
make fresh && make logs
```

This will provide:
- [http://localhost:8080/books](http://localhost:8080/books): HTTP API endpoint, to list the books
- [http://localhost:3333/sse](http://localhost:3333/sse): MCP SSE server endpoint, for the MCP clients
- [http://localhost:8081](http://localhost:8081): Yokai dashboard
- [http://localhost:16686](http://localhost:16686): Jaeger

### Configure your MCP client

If you use MCP compatible applications like [Cursor](https://www.cursor.com/), or [Claude desktop](https://claude.ai/download), you can register this application as MCP server:

```json
{
  "mcpServers": {
    "yokai": {
      "url": "http://localhost:3333/sse"
    }
  }
}
```

Note, if you client does not support remote MCP servers, you can use a [local proxy](https://developers.cloudflare.com/agents/guides/test-remote-mcp-server/#connect-your-remote-mcp-server-to-claude-desktop-via-a-local-proxy):

```json
{
  "mcpServers": {
    "yokai": {
      "command": "npx",
      "args": ["mcp-remote", "http://localhost:3333/sse"]
    }
  }
}
```
