#!/usr/bin/env python3
"""Expose only n8n webhook routes to a local HTTPS tunnel."""

from http.client import HTTPConnection
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import os


LISTEN_PORT = int(os.environ.get("PAWSEAR_WEBHOOK_PROXY_PORT", "5677"))
N8N_HOST = os.environ.get("PAWSEAR_N8N_HOST", "127.0.0.1")
N8N_PORT = int(os.environ.get("PAWSEAR_N8N_PORT", "5678"))
MAX_BODY_BYTES = 5 * 1024 * 1024
ALLOWED_PREFIXES = ("/webhook/", "/webhook-test/")
HOP_BY_HOP_HEADERS = {
    "connection",
    "keep-alive",
    "proxy-authenticate",
    "proxy-authorization",
    "te",
    "trailers",
    "transfer-encoding",
    "upgrade",
}


class WebhookProxyHandler(BaseHTTPRequestHandler):
    server_version = "PawsearWebhookProxy/1.0"

    def do_GET(self) -> None:
        self._proxy()

    def do_POST(self) -> None:
        self._proxy()

    def do_PUT(self) -> None:
        self._proxy()

    def do_PATCH(self) -> None:
        self._proxy()

    def do_DELETE(self) -> None:
        self._proxy()

    def do_HEAD(self) -> None:
        self._proxy()

    def _proxy(self) -> None:
        if not self.path.startswith(ALLOWED_PREFIXES):
            self.send_error(404)
            return

        content_length = int(self.headers.get("Content-Length", "0"))
        if content_length > MAX_BODY_BYTES:
            self.send_error(413)
            return
        body = self.rfile.read(content_length) if content_length else None
        headers = {
            name: value
            for name, value in self.headers.items()
            if name.lower() not in HOP_BY_HOP_HEADERS and name.lower() != "host"
        }
        headers["Host"] = f"{N8N_HOST}:{N8N_PORT}"

        connection = HTTPConnection(N8N_HOST, N8N_PORT, timeout=30)
        try:
            connection.request(self.command, self.path, body=body, headers=headers)
            response = connection.getresponse()
            response_body = response.read()
            self.send_response(response.status)
            for name, value in response.getheaders():
                if name.lower() not in HOP_BY_HOP_HEADERS and name.lower() != "content-length":
                    self.send_header(name, value)
            self.send_header("Content-Length", str(len(response_body)))
            self.end_headers()
            if self.command != "HEAD":
                self.wfile.write(response_body)
        except OSError:
            self.send_error(502)
        finally:
            connection.close()

    def log_message(self, format_string: str, *args: object) -> None:
        print(f"webhook-proxy: {format_string % args}")


if __name__ == "__main__":
    server = ThreadingHTTPServer(("127.0.0.1", LISTEN_PORT), WebhookProxyHandler)
    print(f"Webhook-only proxy listening on http://127.0.0.1:{LISTEN_PORT}")
    server.serve_forever()
