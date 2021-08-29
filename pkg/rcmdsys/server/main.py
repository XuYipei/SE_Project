import sys
sys.path.append(".")

import os
import json
import numpy as np
from concurrent import futures
import time

from google.protobuf import message
import grpc

from svd.svd import svd
import rcmdsys_pb2_grpc
import rcmdsys_pb2

records = []

class RcmdsysServicer(rcmdsys_pb2_grpc.RcmdsysServicer):
    # 实现 proto 文件中定义的 rpc 调用
    def Upd(self, request, context):
        global records
        records.append({
            "userId": request.userId, 
            "probId": request.probId,
            "score": request.score})
        if len(records) > 20000:
            records = records[:-10000]
        print(records)

        return rcmdsys_pb2.Status(status = 'sucess')
    
    def Query(self, request, context):
        print(request)

        result = svd.queryUser(request.id, request.maxLength)
        print(result)
        return rcmdsys_pb2.ProbAry(id = result)

if __name__ == '__main__':
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    rcmdsys_pb2_grpc.add_RcmdsysServicer_to_server(RcmdsysServicer(), server)
    server.add_insecure_port('localhost:5010')
    server.start()

    print("-----------recommend system starts-----------")

    svd = svd()
    lastTrain = time.time()

    recordsPath = './assets/recommendSystem/records.json'
    if os.path.exists(recordsPath):
        f = open(recordsPath, 'r')
        r = json.load(f)
        records = r['data']
        f.close()
        
        print("train begin:")
        print(records)
        svd.trainAll(records)
        print("train end.")

    try:
        while True:
            time.sleep(60)
            print("train begin:")
            print(records)
            svd.trainAll(records)
            print("train end.")

            # if os.path.exists(recordsPath):
            f = open(recordsPath, 'w')
            r = {'time': time.time(), 'data': records}
            json.dump(r, f)
            f.close()

    except KeyboardInterrupt:
        server.stop(0)
