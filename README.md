# webexec

**webexec** is a simple webhook server written in Go that securely **executes system commands via webhooks**, typically used with services like GitHub or GitLab.

---

## 📦 Installation

1. **Clone and build the project:**

```bash
git clone https://github.com/nsavinda/webexec.git
cd webexec
go build -o webexec
```

2. **(Optional) Move the binary to your PATH:**

```bash
sudo mv webexec /usr/local/bin/
```

---

## ⚙️ Configuration

Create a configuration file at `/etc/webexec/config.yaml`:

```yaml
key: your-hmac-secret-key
command: git pull origin main
dir: /path/to/your/repo
port: "8080"
```

* `key`: Secret key used to validate incoming webhook signatures (`X-Hub-Signature-256`).
* `command`: Command to execute upon successful signature verification.
* `dir`: Directory in which to execute the command.
* `port`: Port to run the HTTP server on.

> ⚠️ Make sure this file is only accessible by the service user: `chmod 600 /etc/webexec/config.yaml`

---

## 🚦 Running the Server

```bash
./webexec
```

Output:

```bash
Starting server on port 8080
```

---

## 📬 Sending a Webhook

### Minimal example (no data):

```bash
curl -X POST http://localhost:8080/webhook \
  -H "X-Hub-Signature-256: sha256=<generated-signature>" \
  -H "Content-Type: application/json"
```

### With JSON payload:

```bash
curl -X POST http://localhost:8080/webhook \
  -H "X-Hub-Signature-256: sha256=<generated-signature>" \
  -H "Content-Type: application/json" \
  -d '{"example":"data"}'
```

To generate the correct signature:

```bash
echo -n '{"example":"data"}' | openssl dgst -sha256 -hmac "your-hmac-secret-key"
```

Use the hex digest prefixed with `sha256=` in the header.

---

## 📁 Directory Structure

```
webexec/
├── main.go
├── go.mod
├── go.sum
└── README.md
```

---

## 📫 Author

[Nirmal Savinda](https://github.com/nsavinda)
