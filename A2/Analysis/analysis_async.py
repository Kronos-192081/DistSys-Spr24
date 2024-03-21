import asyncio
from time import perf_counter
import json
import requests as rq
import sys
import aiohttp

async def req(s, d, op):
    if op == 'w':
        write_payload = {"data" : [d]}
        write_json = write_payload
        async with s.post(f'http://localhost:5000/write', json=write_json) as r:
            if r.status != 200:
                r.raise_for_status()
            # print(await r.text())
            return await r.text()
    elif op == 'r':
        read_payload = {
            "Stud_id" : {"low" : d["Stud_id"], "high" : d["Stud_id"]}
        }
        read_json = read_payload
        async with s.post("http://localhost:5000/read", json=read_json) as r:
            if r.status != 200:
                r.raise_for_status()
            # print(await r.text())
            return await r.text()

async def req_all(s, data, op):
    tasks = []
    for d in data:
        task = asyncio.create_task(req(s, d, op))
        tasks.append(task)
    res = await asyncio.gather(*tasks)
    return res
    
async def rw_check(data):

    start = perf_counter()
    async with aiohttp.ClientSession() as session:
        await req_all(session, data, 'w')
    stop = perf_counter()
    print("\nWrite Time: {0:5.2f} seconds\nWrite Speed : {1:5.2f} writes per second\n".format(stop - start, len(data) / (stop - start)))
    
    start = perf_counter()
    async with aiohttp.ClientSession() as session:
        await req_all(session, data, 'r')
    stop = perf_counter()
    print("\nRead Time: {0:5.2f} seconds\nRead Speed : {1:5.2f} reads per second\n".format(stop - start, len(data) / (stop - start)))
    
def main():
    if len(sys.argv) < 3:
        print("Usage: python3 analysis_async.py [1/2/3]")
    test_no = int(sys.argv[1])
    match test_no:
        case 1:
            f = open("data0.json", "r")
            data = json.load(f)
            f.close()
            print("Test 1: Sending 10000 write requests followed by 10000 read requests...\n")
            print("This test is performed with 3 shard replicas, 4 shards and 6 servers.\n")
            init_payload = {
                "N":6,
                "schema":{"columns":["Stud_id","Stud_name","Stud_marks"], "dtypes":["Number","String","String"]},
                "shards":[
                    {"Stud_id_low":0, "Shard_id": "sh1", "Shard_size":4096},
                    {"Stud_id_low":4096, "Shard_id": "sh2", "Shard_size":4096},
                    {"Stud_id_low":8192, "Shard_id": "sh3", "Shard_size":4096},
                    {"Stud_id_low":12288, "Shard_id": "sh4", "Shard_size":4096}],
                "servers":{
                    "Server0":["sh1","sh2"],
                    "Server1":["sh3","sh4"],
                    "Server3":["sh1","sh3"],
                    "Server4":["sh4","sh2"],
                    "Server5":["sh1","sh4"],
                    "Server6":["sh3","sh2"]}
            }
            res = rq.post("http://localhost:5000/init", json=init_payload)
            if res.status_code != 200:
                res.raise_for_status()

            print("Init completed successfully")

            start = perf_counter()
            asyncio.run(rw_check(data))
            stop = perf_counter()
            print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))
        case 2:
            f = open("data0.json", "r")
            data = json.load(f)
            f.close()
            print("Test 2: Sending 10000 write requests followed by 10000 read requests...\n")
            print("This test is performed with 7 shard replicas, 4 shards and 7 servers.\n")
            init_payload = {
                "N":7,
                "schema":{"columns":["Stud_id","Stud_name","Stud_marks"], "dtypes":["Number","String","String"]},
                "shards":[
                    {"Stud_id_low":0, "Shard_id": "sh1", "Shard_size":4096},
                    {"Stud_id_low":4096, "Shard_id": "sh2", "Shard_size":4096},
                    {"Stud_id_low":8192, "Shard_id": "sh3", "Shard_size":4096},
                    {"Stud_id_low":12288, "Shard_id": "sh4", "Shard_size":4096}],
                "servers":{
                    "Server0":["sh1","sh2","sh3","sh4"],
                    "Server1":["sh1","sh2","sh3","sh4"],
                    "Server2":["sh1","sh2","sh3","sh4"],
                    "Server3":["sh1","sh2","sh3","sh4"],
                    "Server4":["sh1","sh2","sh3","sh4"],
                    "Server5":["sh1","sh2","sh3","sh4"],
                    "Server6":["sh1","sh2","sh3","sh4"]}
            }
            res = rq.post("http://localhost:5000/init", json=init_payload)
            if res.status_code != 200:
                res.raise_for_status()

            print("Init completed successfully")

            start = perf_counter()
            asyncio.run(rw_check(data))
            stop = perf_counter()
            print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))
        case 3:
            f = open("data1.json", "r")
            data = json.load(f)
            f.close()
            print("Test 3: Sending 10000 write requests followed by 10000 read requests...\n")
            print("This test is performed with 8 shard replicas, 6 shards and 10 servers.\n")
            init_payload = {
                "N":10,
                "schema":{"columns":["Stud_id","Stud_name","Stud_marks"], "dtypes":["Number","String","String"]},
                "shards":[
                    {"Stud_id_low":0, "Shard_id": "sh1", "Shard_size":4096},
                    {"Stud_id_low":4096, "Shard_id": "sh2", "Shard_size":4096},
                    {"Stud_id_low":8192, "Shard_id": "sh3", "Shard_size":4096},
                    {"Stud_id_low":12288, "Shard_id": "sh4", "Shard_size":4096},
                    {"Stud_id_low":16384, "Shard_id": "sh5", "Shard_size":4096},
                    {"Stud_id_low":20480, "Shard_id": "sh6", "Shard_size":4096}],
                "servers":{
                    "Server0":["sh1","sh2","sh3","sh4","sh6"],
                    "Server1":["sh1","sh2","sh3","sh4","sh6"],
                    "Server2":["sh1","sh2","sh3","sh5","sh6"],
                    "Server3":["sh1","sh2","sh3","sh5","sh6"],
                    "Server4":["sh1","sh2","sh4","sh5","sh6"],
                    "Server5":["sh1","sh2","sh4","sh5","sh6"],
                    "Server6":["sh1","sh3","sh4","sh5","sh6"],
                    "Server7":["sh1","sh3","sh4","sh5","sh6"],
                    "Server8":["sh2","sh3","sh4","sh5"],
                    "Server9":["sh2","sh3","sh4","sh5"]}
            }
            res = rq.post("http://localhost:5000/init", json=init_payload)
            if res.status_code != 200:
                res.raise_for_status()

            print("Init completed successfully")

            start = perf_counter()
            asyncio.run(rw_check(data))
            stop = perf_counter()
            print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))
        case _:
            print("Invalid test no. Exiting...")

if __name__ == '__main__':
    main()
