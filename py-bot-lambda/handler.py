import json
import os
import sys

from botocore.vendored import requests
from boto3.dynamodb.conditions import Key, Attr
import logging
import datetime
import boto3

## add loging to CW
logger = logging.getLogger()
logger.setLevel(logging.INFO)

DB_TABLE = os.environ['DB_TABLE']
client = boto3.client('dynamodb')
type_list = ["Food","Gifts","Medical", "AppartmentRent","Transportation","Departmental", "WislianeTerasy", "WorkLunch" ,"Travel","DebtCar", "OneTime","Car", "Fun"]
TOKEN = os.environ['TELEGRAM_TOKEN']
BASE_URL = "https://api.telegram.org/bot{}".format(TOKEN)

def create_record(message, message_time):
    n_message = message.split(" ")
    if n_message[3] in type_list:
        resp = client.put_item(
            TableName=DB_TABLE,
            Item={
                'date': {'S':message_time},
                'amount':{'N':n_message[1]},
                'expencis': {'S':n_message[2]},
                'type':{'S':n_message[3]}
                }
            )
        return "the item added"
    else:
        return "the type should from the list: {}".format(type_list)

def get_records(message_time):
    response = client.scan(
        TableName=DB_TABLE
    )
    list_of_recodrs = ""
    for record in response["Items"]:
        list_of_recodrs += "{} :: {} :: {}\n".format(
            datetime.datetime.fromtimestamp(int(record["date"]["S"])).strftime('%Y-%m-%d'), 
            record["expencis"]["S"], 
            record["amount"]["N"]) 

    return list_of_recodrs

def hello(event, context):
    logger.info('START...')
    logger.info(event)

    message = str(event["message"]["text"])
    message_time = str(event["message"]["date"])
    chat_id = event["message"]["chat"]["id"]
    response = ""

    if "/add" in message:
        response = create_record(message, message_time)
    elif "/report" in message:
        response = str(get_records(message_time))
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