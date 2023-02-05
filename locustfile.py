import string
import random
import queue
from typing import Any
from locust import FastHttpUser, between, task

class APIUser(FastHttpUser):
    """Send api user requests."""

    wait_time = between(1, 1)

    def __init__(self, *args: Any) -> None:
        """Init."""
        super().__init__(*args)
        self.que = queue.Queue(maxsize=100)

    @task
    def shorten_url(self) -> None:
        """Shorten a long url."""
        url = ''.join(
            random.choice(string.ascii_letters + string.digits) for _ in range(10)
        )
        resp = self.client.post("/shorten", {"url": url})
        key = resp.json()["key"]
        self.que.put(key)

    @task
    def original_url(self) -> None:
        """Get the long url."""
        if self.que.qsize() == 0:
            return
        key = self.que.get()
        resp = self.client.get(f"/{key}", name="/<key>", allow_redirects=False)
