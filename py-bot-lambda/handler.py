import json
import os
import sys

from botocore.vendored import requests
import logging
from datetime import datetime
import boto3

today = datetime.now()
## add loging to CW
logger = logging.getLogger()
logger.setLevel(logging.INFO)

DB_TABLE = os.environ['DB_TABLE']
client = boto3.client('dynamodb')

TOKEN = os.environ['TELEGRAM_TOKEN']
BASE_URL = "https://api.telegram.org/bot{}".format(TOKEN)

def create_record(time_stamp):
    record_date= time_stamp.strftime("%Y-%b-%d")
    resp = client.put_item(
        TableName=DB_TABLE,
        Item={
            "date": record_date
        }
    )

def hello(event, context):
    logger.info('START...')
    logger.info(event)

    # data = json.loads(event)
    message = str(event["message"]["text"])
    chat_id = event["message"]["chat"]["id"]
    response = ""
    # first_name = event["message"]["chat"]["first_name"]

    # response = "Please /start, {}".format(first_name)

    if "/add" in message:
        create_record(today)
        response = "new new expences was added successful >>> "
    elif "/report" in message:
        response = "the expences is >>> till {}".format(today)
    elif "/help" in message:
        response = "List of commands: >>> "
    else:
        response = "please use /help to get commands info"

    data = {"text": response.encode("utf8"), "chat_id": chat_id}
    url = BASE_URL + "/sendMessage"
    requests.post(url, data)
    logger.info(' GOOD, event: {}'.format(event))


    return {"statusCode": 200,
    "type": "GOOD"
        
    }