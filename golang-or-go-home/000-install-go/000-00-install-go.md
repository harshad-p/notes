# Lesson 0: Install Go

Before writing Go code, we need the Go toolchain.

## Step 1: Check if Go is already installed

Open terminal:

```bash
go version
```

If installed, you will see something like:

```
go version go1.24.5 darwin/arm64
```

If you see:

```
command not found: go
```

install it.

---

## Step 2: Install Go

### macOS (your likely environment)

Using Homebrew:

```bash
brew install go
```

Then verify:

```bash
go version
```

---

### Windows

Download the installer from:

https://go.dev/dl/

Run the `.msi` installer.

Verify in Command Prompt:

```cmd
go version
```

---

### Linux (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install golang-go
```

Verify:

```bash
go version
```

---

## Step 3: Check Go environment

Run:

```bash
go env
```

You will see settings like:

```
GOROOT=/usr/local/go
GOPATH=/Users/yourname/go
```

For now, don't worry about them. We will understand modules later.

---

## Your task before Lesson 1

Just do this:

1. Install Go
2. Run:

```bash
go version
```

3. Tell me the output.

Then we move to **Lesson 1: Creating your first Go project**.

No code yet. We set up the foundation first.