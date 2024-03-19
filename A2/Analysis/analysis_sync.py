import json
import requests as rq
from time import perf_counter

def req(d, op):
    if op == 'w':
        write_payload = {"data": [d]}
        write_json = write_payload
        # print(write_json)
        r = rq.post(f'http://localhost:5000/write', json=write_json)
        if r.status_code != 200:
            r.raise_for_status()
        return r.text
    elif op == 'r':
        read_payload = {
            "Stud_id": {"low": d["Stud_id"], "high": d["Stud_id"]}
        }
        read_json = read_payload
        r = rq.post("http://localhost:5000/read", json=read_json)
        if r.status_code != 200:
            r.raise_for_status()
        return r.text

def req_all(data, op):
    results = []
    for d in data:
        result = req(d, op)
        results.append(result)
    return results
    
def rw_check(data):
    start = perf_counter()
    req_all(data, 'w')
    stop = perf_counter()
    print("\nWrite Time: {0:5.2f} seconds\nWrite Speed : {1:5.2f} writes per second\n".format(stop - start, len(data) / (stop - start)))

    start = perf_counter()
    req_all(data, 'r')
    stop = perf_counter()
    print("\nRead Time: {0:5.2f} seconds\nRead Speed : {1:5.2f} reads per second\n".format(stop - start, len(data) / (stop - start)))

def main():
    f = open("data.json", "r")
    data = json.load(f)
    f.close()
    print("Test 1: Sending 10000 write requests followed by 10000 read requests...\n")
    print("This test is performed with 3 shard replicas.\n")
    init_payload = {
        "N": 3,
        "schema": {"columns" : ["Stud_id", "Stud_name", "Stud_marks"], "dtypes" : ["Number", "String", "Number"]},
        "shards" : [{"Stud_id_low" : 0, "Shard_id" : "sh1", "Shard_size" : 4096},
                    {"Stud_id_low" : 4096, "Shard_id" : "sh2", "Shard_size" : 4096},
                    {"Stud_id_low" : 8192, "Shard_id" : "sh3", "Shard_size" : 4096}],
        "servers" : {"Server0" : ["sh1", "sh2"],
                     "Server1" : ["sh2", "sh3"],
                     "Server2" : ["sh1", "sh3"]}
    }
    res = rq.post("http://localhost:5000/init", json=init_payload)
    if res.status_code != 200:
        res.raise_for_status()

    print("Init completed successfully")

    start = perf_counter()
    rw_check(data)
    stop = perf_counter()
    print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))

if __name__ == '__main__':
    main()
