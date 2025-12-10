# Superpage Phase 1

## Running the Services

This project uses a [Procfile](https://devcenter.heroku.com/articles/procfile) to define and run all services together. A Procfile is a simple text file that declares the commands to start each process.

### Install a Process Manager

You'll need a Procfile runner. Choose one:

**Foreman (Ruby)**
```bash
gem install foreman
```

**Overmind (Go, recommended)**
```bash
# macOS
brew install overmind

# or with Go
go install github.com/DarthSim/overmind/v2@latest
```

**Hivemind (Go, lightweight)**
```bash
# macOS
brew install hivemind

# or with Go
go install github.com/DarthSim/hivemind@latest
```

### Start All Services

```bash
# Using foreman
foreman start

# Using overmind
overmind start

# Using hivemind
hivemind
```

### Access the UI

Once running, open your browser to:

```
http://localhost:3000
```

### Stop All Services

- **Foreman/Hivemind**: Press `Ctrl+C`
- **Overmind**: Run `overmind quit` or press `Ctrl+C` twice
