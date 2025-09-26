# Getting Started with Vortex

This guide will walk you through the process of setting up your development environment and creating your first Vortex application.

## Prerequisites

Before you begin, ensure you have the following installed:

-   [Go](https://golang.org/doc/install) version 1.21 or later.
-   The Vortex CLI.

## Installation

First, install the Vortex command-line tool, which you'll use to create and manage your projects.

```bash
go install github.com/AureClai/vortex/cmd/vortex@latest
```

## Creating a New Application

To create a new Vortex application, use the `vortex init` command:

```bash
vortex init my-app
cd my-app
```

This will create a new directory called `my-app` with the basic structure of a Vortex project, including a simple "Hello, World!" application.

## Running the Development Server

To see your application in action, run the development server:

```bash
vortex dev
```

This command will compile your Go code to WebAssembly, start a local web server, and automatically rebuild your application when you make changes to the source code.

You can now open your browser and navigate to `http://localhost:8080` to see your running application.

## Building for Production

When you are ready to deploy your application, you can create an optimized production build:

```bash
vortex build
```

This will generate a `dist/` directory containing the compiled `app.wasm` file and the necessary HTML and JavaScript files to run your application. You can then serve these files from any static file server.