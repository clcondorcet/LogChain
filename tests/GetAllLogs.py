import requests
import json
import time

def make_http_request():
    url = "http://localhost:3333/querry"
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

def write_to_file(data, file_path):
    try:
        sorted_data = sorted(data, key=lambda x: x.get("timestamp", 0))
        with open(file_path, 'w') as file:
             for item in sorted_data:
                file.write(str(item["message"]) + '\n')
        print(f"Données écrites avec succès dans le fichier : {file_path}")
    except Exception as e:
        print(f"Une erreur est survenue lors de l'écriture dans le fichier : {e}")

def main():
    # Spécifiez le chemin du fichier où vous souhaitez enregistrer les données
    file_path = './tests/output.log'

    now = time.time()

    # Faire la requête HTTP
    response_data = make_http_request()

    if response_data:
        # Écrire la réponse dans un fichier JSON
        write_to_file(response_data, file_path)
    
    timed = time.time() - now
    print("InsertTime:", timed)

if __name__ == "__main__":
    main()
