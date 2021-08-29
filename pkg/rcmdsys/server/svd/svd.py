import numpy as np
import time
import pickle
import copy
'''
    SVD: 令 bi 为用户 i 评分和平均的偏差, bu 为物品 u 得分和平均的偏差
    p(Ui, Mj) = u + bi + bu + (Ui ^ t) * Mj

    Loss: 1/2 Iij * (Vij - p(Ui, Mj)) ^ 2 + ku / 2 * Sum ||Ui|| ^ 2 + km / 2 * Sum ||Mj|| ^ 2

    static:
        delta Ui: Iij * ((Vij - p(Ui, Mj)) * Mj) - ku * Ui
        delta Mj: Iij * ((Vij - p(Ui, Mj)) * Ui) - km * Mj
    
    dynamic:
        delta Ui: (Vij - p(Ui, Mj)) * Mj - ku * Ui
        delta Mj: (Vij - p(Ui, Mj)) * Ui - km * Mj
'''
class svd:
    def __init__(self, epoch = 100, eta = 0.2, ku=0.001, km=0.001, \
                f=2, saveModel=False, savePath=None):
        super(svd, self).__init__()
        self.epoch = epoch
        self.eta = eta
        self.ku = ku
        self.km = km
        self.f = f

        self.U = None
        self.M = None
        self.uMap = {}; self.userNums = int(0)
        self.iMap = {}; self.itemNums = int(0)
        self.maxItem = 568

        self.uMapSaved = {}
        self.itemsSaved = {}
        self.history = {}
        self.hisSaved = {}
        self.uSaved = None
        self.mSaved = None
        self.buSaved = None
        self.biSaved = None
        self.meanVSaved = 0

    def trainAll(self, recordsData):
        records = copy.deepcopy(recordsData)

        self.meanV = 0
        self.uMap = {}; self.userNums = int(0)
        self.iMap = {}; self.itemNums = int(0)
        self.items = []
        self.history = {}
        for rec in records:
            self.meanV += rec["score"]
            if rec["userId"] not in self.uMap.keys():
                self.uMap[rec["userId"]] = int(self.userNums)
                rec["userId"] = int(self.userNums)
                self.userNums += int(1)
            else:
                rec["userId"] = self.uMap[rec["userId"]]

            if rec["probId"] not in self.iMap.keys():
                self.items.append(rec["probId"])
                self.iMap[rec["probId"]] = int(self.itemNums)
                rec["probId"] = int(self.itemNums)
                self.itemNums += int(1)
            else:
                rec["probId"] = self.iMap[rec["probId"]]
            
            self.history[tuple([rec["userId"], rec["probId"]])] = 1

        if len(records) == 0:
            return

        rateNums = len(records)
        self.meanV /= rateNums
        initv = np.sqrt((self.meanV - 1) / self.f)
        self.U = initv + np.random.uniform(-0.01,0.01,(self.userNums+1,self.f))
        self.M = initv + np.random.uniform(-0.01,0.01,(self.itemNums+1,self.f))
        self.bu = np.zeros(self.userNums + 1)
        self.bi = np.zeros(self.itemNums + 1)
        self.y = np.zeros((self.itemNums+1, self.f)) + 0.1
        
        start = time.time()
        for i in range(self.epoch):
            sumRmse = 0.0
            for sample in records:
                uid = sample["userId"]
                iid = sample["probId"]
                vij = float(sample['score'])
                # p(U_i,M_j) = mu + b_i + b_u + U_i^TM_j

                p = self.meanV + self.bu[uid] + self.bi[iid] + \
                    np.sum(self.U[uid] * self.M[iid])
                error = vij - p
                sumRmse += error**2
                # 计算Ui,Mj的梯度
                deltaU = error * self.M[iid] - self.ku * self.U[uid]
                deltaM = error * self.U[uid] - self.km * self.M[iid]
                # 更新参数
                self.U[uid] += self.eta *  deltaU
                self.M[iid] += self.eta *  deltaM

                self.bu[uid] += self.eta * (error - self.ku * self.bu[uid])
                self.bi[iid] += self.eta * (error - self.km * self.bi[iid])

            trainRmse = np.sqrt(sumRmse/rateNums)
            if i % 20 == 0:
                print("Epoch %d cost time %.4f, train RMSE: %.4f" % \
                    (i, time.time()-start, trainRmse))
        
        self.uMapSaved = copy.deepcopy(self.uMap)
        self.itemsSaved = copy.deepcopy(self.items)
        self.uSaved = copy.deepcopy(self.U)
        self.mSaved = copy.deepcopy(self.M)
        self.buSaved = copy.deepcopy(self.bu)
        self.biSaved = copy.deepcopy(self.bi)
        self.meanVSaved = copy.deepcopy(self.meanV)
        self.hisSaved = copy.deepcopy(self.history)


    def queryUser(self, userId, maxLength):
        if userId not in self.uMapSaved.keys():
            return []

        uId = self.uMapSaved[userId]
        scores = []

        for iId in range(len(self.itemsSaved)):
            if tuple([uId, iId]) in self.hisSaved:
                scores.append(-10)
            else:
                p = self.meanV + self.buSaved[uId] + self.biSaved[iId] \
                        + np.sum(self.uSaved[uId] * self.mSaved[iId])
                scores.append(p)

        scores = np.asarray(scores)
        rk = scores.argsort()[: maxLength]
        result = []
        for i in rk:
            if scores[i] > -1:
                result.append(self.itemsSaved[i])

        return result

class svd_ver2:
    def __init__(self, iter_num = 10000, lr = 0.001, debug = False, k = 10, max_rating = 5, min_rating = 0):
        self.iter_num = iter_num
        self.user_num = 0
        self.problem_num = 0
        self.lr = lr
        self.score_matrix = {}
        self.predict_score_matrix = {}
        self.user_map = {}
        self.problem_map = {}
        self.user_demap = {}
        self.problem_demap = {}
        self.debug = debug
        self.k = k
        self.p = {}
        self.q = {}
        self.max_rating = max_rating
        self.min_rating = min_rating
        self.result = {}
        pass
    
    def get_result(self, records):
        self._init_score_matrix(records);
        self._train()
        if self.debug:
            print(self.predict_score_matrix)
        self._count_avg()
        if self.debug:
            print(self.result)
        return self.result

    def _train(self):
        k = self.k
        if self.user_num < k:
            k = self.user_num
        if self.problem_num < k:
            k = self.problem_num
        self.p = np.zeros((self.user_num, k)) + 1.0 / self.user_num
        self.q = np.zeros((k, self.problem_num)) + 1.0 / self.problem_num
        for _ in range(self.iter_num):
            for i in range(self.user_num):
                for j in range(k):
                    for x in range(self.problem_num):
                        self.p[i][j] += self.lr * self._cal_diff(i, x, k) * self.q[j, x];
            for i in range(k):
                for j in range(self.problem_num):
                    for x in range(self.user_num):
                        self.q[i][j] += self.lr * self._cal_diff(x, j, k) * self.q[x, j];
        self.predict_score_matrix = self.p.dot(self.q)
        for i in range(self.user_num):
            for j in range(self.problem_num):
                if self.score_matrix[i][j] != -1:
                    self.predict_score_matrix[i][j] = self.score_matrix[i][j]
                if self.predict_score_matrix[i][j] > self.max_rating:
                    self.predict_score_matrix[i][j] = self.max_rating
                if self.predict_score_matrix[i][j] < self.min_rating:
                    self.predict_score_matrix[i][j] = self.max_rating
                self.predict_score_matrix[i][j] = round(self.predict_score_matrix[i][j], 4)
        pass

    def _cal_diff(self, i, j, k):
        if self.score_matrix[i][j] == -1:
            return 0.0
        ret = self.score_matrix[i][j]
        for z in range(k):
            ret -= self.p[i][z] * self.q[z][j];
        return ret;

    def _init_score_matrix(self, records):
        self.user_map = {}
        self.problem_map = {}
        self.user_num = 0
        self.problem_num = 0
        self.user_demap = {}
        self.problem_demap = {}
        for record in records:
            if record["userId"] not in self.user_map.keys():
                self.user_map[record["userId"]] = self.user_num
                self.user_demap[self.user_num] = record["userId"]
                self.user_num += 1

            if record["probId"] not in self.problem_map.keys():
                self.problem_map[record["probId"]] = self.problem_num
                self.problem_demap[self.problem_num] = record["probId"]
                self.problem_num += 1

        self.score_matrix = np.zeros((self.user_num, self.problem_num)) - 1
        self.predict_score_matrix = np.zeros((self.user_num, self.problem_num))
        for record in records:
            i = self.user_map[record['userId']]
            j = self.problem_map[record['probId']]
            self.score_matrix[i][j] = record['score']

        if self.debug:
            print(self.score_matrix)

    def _count_avg(self):
        self.result = np.zeros(self.problem_num)
        for i in range(self.problem_num):
            for j in range(self.user_num):
                self.result[i] += self.predict_score_matrix[j][i]
            self.result[i] /= self.user_num
        pass
        
if __name__ == '__main__':
    records = []
    scores = [[5,4,4.5,-1,3.9],
              [-1,4.5,-1,4.5,-1],
              [4.5,-1,4.4,4,4],
              [-1,4.8,-1,-1,4.5],
              [4,-1,4.5,5,-1]]
    scores = [
        [0,0,-1,0,-1],
        [5,5,5,0,-1],
        [2,3,0,-1,-1],
        [-1,-1,-1,-1,2],
        [0,2,-1,-1,-1]
    ]
    for i in range(5):
        for j in range(5):
            records.append({
            "userId": i, 
            "probId": j,
            "score": scores[i][j]})
    svd = svd_ver2(debug = True)
    
    records.append({
        "userId": 1, 
        "probId": 20,
        "score": 5
    })

    svd.get_result(records)