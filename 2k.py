import json
import dateutil.tz
import dateutil.parser
import os
import requests
from datetime import datetime as dt

data = requests.get("https://www.lockercodes.io/page-data/22/active-locker-codes/page-data.json").json()

tg_chat_id = os.environ["TELEGRAM_ID"]
tg_token = os.environ["TELEGRAM_TOKEN"]

jakarta = dateutil.tz.gettz("Asia/Jakarta")
mountain = dateutil.tz.gettz("US/Mountain")

edges = data["result"]["data"]["allLockerCodes"]["edges"]

msg = ""
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

if msg == '':
    msg = 'No code today :('

params = {'chat_id': tg_chat_id, 'text': msg}
res = requests.post(f"https://api.telegram.org/bot{tg_token}/sendMessage", data=params).json()
print(res["ok"])    
