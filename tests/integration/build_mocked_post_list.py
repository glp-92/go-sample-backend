import requests
import time
from uuid import uuid4

BACKEND_URL = "http://localhost:8080"
USERNAME = "test"
PASSWORD = "changeme"
POST = {
    "title": None,
    "slug": None,
    "excerpt": "About a test post",
    "content": "Content of the post",
    "featuredImage": "https://static.chollometro.com/threads/raw/qet3j/1539472_1/re/202x202/qt/70/1539472_1.jpg",
    "categoryIds": [
        "d08a1e08-a757-4840-858f-d946e2ef7ca0",
        "3e8da162-4f51-46e7-8faa-1c3507ca9a32",
    ],
    "themeIds": ["4cb61578-3f7d-4faf-9609-ea8377e30aa5"],
}
N_POSTS = 20


def login():
    response = requests.post(
        f"{BACKEND_URL}/auth/login",
        json={"username": USERNAME, "password": PASSWORD},
    )
    return response.json()["access_token"]


def create_post(i: int, access_token: str):
    post = POST
    post["title"] = f"Test Post {i+1}"
    post["slug"] = str(uuid4())
    response = requests.post(
        f"{BACKEND_URL}/posts",
        json=post,
        headers={"Authorization": f"Bearer {access_token}"},
    )
    print(f"Post creating {i} : {response.status_code}")
    time.sleep(1)


access_token = login()
for i in range(N_POSTS):
    create_post(i=i, access_token=access_token)
