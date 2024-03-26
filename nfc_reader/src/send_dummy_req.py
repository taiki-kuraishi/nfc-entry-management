"""
デモンストレーション用のダミーリクエストを送信する
"""
import os
import time

import requests


def send_dummy_request() -> None:
    """
    ダミーリクエストを送信します。
    """
    url = os.environ["API_URL"]
    headers = {"Content-Type": "application/json"}
    body = {
        "student_number": 20122000,
        "name": "タナカ　タロウ",
        "timestamp": time.time(),
    }

    response = requests.post(url, json=body, headers=headers, timeout=5)
    print(response.status_code)
    print(response.text)

if __name__ == "__main__":
    send_dummy_request()
