import requests
from datetime import datetime
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import as_completed
import time
import matplotlib.pyplot as plt
import numpy as np

url = "http://localhost:3333/invoke"

timesInsert = []

def getTimeStamp(log):
  timestamp_str = log.split()[0] + ' ' + log.split()[1]
  timestamp_format = "%Y-%m-%d %H:%M:%S,%f"
  timestamp = datetime.strptime(timestamp_str, timestamp_format)
  return (timestamp, (timestamp - datetime(1970, 1, 1)).total_seconds() - 3600)

def send_post_request(lines):
  now = time.time()

  jsoned = "["
  i = 0
  for line in lines:
      timestamp_str, timestamp = getTimeStamp(line)

      jsoned += "{\\\"hostname\\\":\\\"cleme\\\",\\\"message\\\":\\\"" + line + "\\\",\\\"timestamp\\\":\\\"" + str(round(int(timestamp))) + "\\\"}"

      if i != len(lines) - 1:
        jsoned += ","
      i += 1

  jsoned += "]"

  data = "{\"function\":\"AddAssets\",\"args\":[\"" + jsoned + "\"]}"

  # print(data)

  response = requests.post(url, data=data)

  # if response.status_code == 200:
  #     print(f"POST request successful for line: {line}")
  # else:
  #     print(f"Error in POST request for line: {line}")

  timeInsert = time.time() - now
  timesInsert.append(timeInsert)

def process_file(file_path):
  try:
    with open(file_path, 'r') as file:
        lines = [line.strip() for line in file]
    return lines
  except FileNotFoundError:
    print(f"Error: File not found - {file_path}")
    return []
  except Exception as e:
    print(f"An error occurred: {e}")
    return []

def main(file_path):
  lines = process_file(file_path)

  now = time.time()

  lines_batch = 100
  
  with ThreadPoolExecutor(max_workers=20) as executor:
    line_batches = [lines[i:i+lines_batch] for i in range(0, len(lines), lines_batch)]

    # Submit the tasks
    futures = [executor.submit(send_post_request, batch) for batch in line_batches]

    # Wait for all tasks to complete
    for future in as_completed(futures):
      try:
        future.result()
      except Exception as e:
        print(f"An error occurred: {e}")
          
  timetotal = time.time() - now

  print("inserted", len(lines), "in", timetotal, "seconds.")

  indices = np.arange(len(timesInsert))
  plt.plot(indices, timesInsert, marker='o', linestyle='-')
  plt.title('Time of Insertion vs Index')
  plt.xlabel('Index')
  plt.ylabel('Time of Insertion (seconds)')
  plt.grid(True)
  plt.show()

if __name__ == "__main__":
  # Replace 'your_file.txt' with the actual file path you want to read
  file_path = 'd:/dev/EngineeringProj/LogChain/tests/fail2ban.log'
  main(file_path)
