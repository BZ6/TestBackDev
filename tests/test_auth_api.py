import requests
import pytest


BASE_URL = "http://host.docker.internal:8080"

def test_generate_tokens():
    user_id = "123e4567-e89b-12d3-a456-426614174000"
    response = requests.get(f"{BASE_URL}/auth/token", params={"user_id": user_id})

    assert response.status_code == 200

    data = response.json()
    assert "access_token" in data
    assert "refresh_token" in data
    assert data["access_token"] != None
    assert data["refresh_token"] != None

    return data

def test_refresh_tokens():
    tokens = test_generate_tokens()

    payload = {
        "access_token": tokens["access_token"],
        "refresh_token": tokens["refresh_token"],
    }
    response = requests.post(f"{BASE_URL}/auth/refresh", json=payload)

    assert response.status_code == 200

    data = response.json()
    assert "access_token" in data
    assert "refresh_token" in data
    assert data["access_token"] != None
    assert data["refresh_token"] != None
