from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed
import requests
import json
import time

url_invoke = "http://localhost:3333/invoke"

def make_http_request():
    url = "http://localhost:3333/query"
    payload = {
        "function": "GetAllAssets",
        "args": []
    }

    # Faire la requête HTTP
    response = requests.post(url, json=payload)

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Erreur lors de la requête HTTP : {response.status_code}")
        return None

def delete_data(data):
        data = "{\"function\":\"DeleteAsset\",\"args\":[\"" + str(data["ID"]) + "\"]}"
        response = requests.post(url_invoke, data=data)

def main():
    # Spécifiez le chemin du fichier où vous souhaitez enregistrer les données
    now = time.time()

    # Faire la requête HTTP
    response_data = make_http_request()

    if response_data:
        with ThreadPoolExecutor(max_workers=100) as executor:
            futures = [executor.submit(delete_data, line) for line in response_data]

            # Wait for all tasks to complete
            for future in as_completed(futures):
                try:
                    future.result()
                except Exception as e:
                    print(f"An error occurred: {e}")
    
    timed = time.time() - now
    print("InsertTime:", timed)

if __name__ == "__main__":
    main()
