import pymongo
from pymongo import collection

'''
    score: {
        'user_id' string
        'prob_id' string
        'score' double
    }
'''

global dburl, dbName, visitTable, userTable, probTable, locked


def addUser():
    pass

def addProb():
    pass

def updScore(record):
    client = pymongo.MongoClient(dburl)
    mycol = client[dbName][userTable]
    
    result = mycol.find_one_and_update({
            'user_id': record.userId,
            'prob_id': record.probId
        }, {
            '$set': {'score': record.score}
        }
    )
    if result == None:
        mycol.insert_one(record)

def query():
    pass


dburl = "mongodb://xuyipei:123456@localhost:27017/"
dbName = "rcmdsys"
scoreTable = "score"
userTable = "user"
probTable = "prob"