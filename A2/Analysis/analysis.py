import asyncio
from time import perf_counter
import json
import requests as re

import aiohttp

async def req(s, d, op):
    if op == 'w':
        write_payload = {"data" : [d]}
        write_json = json.dumps(write_payload)
        async with s.post(f'http://localhost:5000/write', json=write_json) as r:
            if r.status != 200:
                r.raise_for_status()
            return await r.text()
    elif op == 'r':
        read_payload = {
            "Stud_id" : {"low" : d["id"], "high" : d["id"]}
        }
        read_json = json.dumps(read_payload)
        async with s.post("http://localhost:5000/read", json=read_json) as r:
            if r.status != 200:
                r.raise_for_status()
            return await r.text()

async def req_all(s, data, op):
    tasks = []
    for d in data:
        task = asyncio.create_task(req(s, d, op))
        tasks.append(task)
    res = await asyncio.gather(*tasks)
    return res
    
async def rw_check(init_payload):
    f = open("data.json", "r")
    data = json.load(f)
    f.close()

    init_json = json.dumps(init_payload)
    res = re.post("http://localhost:5000/init", json=init_json)
    if res.status_code != 200:
        res.raise_for_status()

    start = perf_counter()
    async with aiohttp.ClientSession() as session:
        await req_all(session, data, 'w')
    stop = perf_counter()
    print("\nWrite Time: {0:5.2f} seconds\nWrite Speed : {0:5.2f} writes per second\n".format(stop - start, len(data) / (stop - start)))
    
    start = perf_counter()
    async with aiohttp.ClientSession() as session:
        await req_all(session, data, 'r')
    stop = perf_counter()
    print("\nRead Time: {0:5.2f} seconds\nRead Speed : {0:5.2f} reads per second\n".format(stop - start, len(data) / (stop - start)))
    
async def main():
    headers = {'Content-Type': 'application/json'}
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
    await rw_check(init_payload)
'''
    print("Test 2: Sending 10000 write requests followed by 10000 read requests...")
    print("This test is performed with 7 shard replicas.\n")
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
    await rw_check(init_payload)
'''
    
    

if __name__ == '__main__':
    start = perf_counter()
    asyncio.run(main())
    stop = perf_counter()
    print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))