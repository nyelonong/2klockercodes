import json
import dateutil.tz
import dateutil.parser
import os
import requests
import sys
from datetime import datetime as dt

tg_chat_id = os.environ["TELEGRAM_ID"]
tg_token = os.environ["TELEGRAM_TOKEN"]
host = os.environ["LOCKERCODES_HOST"]
path = os.environ["LOCKERCODES_PATH"]

if len(sys.argv) < 2:
    print("Need 2k version as an argument")
    exit()

version = sys.argv[1]
url = host + "/" + version + "/" + path
data = requests.get(url).json()

jakarta = dateutil.tz.gettz("Asia/Jakarta")
mountain = dateutil.tz.gettz("US/Mountain")

edges = data["result"]["data"]["allLockerCodes"]["edges"]

has_data = False
msg = f"NBA 2K{version}\n"
for edge in edges:
    code = edge["node"]["lockerCode"]
    title = edge["node"]["title"]
    create = dateutil.parser.parse(edge["node"]["dateAdded"])
    now = dt.now(tz=mountain)
    
    dur = now - create
    if dur.days > 1:
        continue
    
    expire_at = None
    if edge["node"]["expiration"] is not None:
        expire = dateutil.parser.parse(edge["node"]["expiration"])
        
        if now > expire:
            continue
        
        expire_at = expire.astimezone(tz=jakarta)
        
    msg += f"Title: {title}\nCode: {code}\nCreated At: {create.astimezone(tz=jakarta)}\nExpire At: {expire.astimezone(tz=jakarta)}\n\n"
    has_data = True

if not has_data:
    msg += 'No code today :('

params = {'chat_id': tg_chat_id, 'text': msg}
print(params)
res = requests.post(f"https://api.telegram.org/bot{tg_token}/sendMessage", data=params).json()
print(res["ok"])    
if not res["ok"]:
    print(params)
    print(res["description"])
