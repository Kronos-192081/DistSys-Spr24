import asyncio
from time import perf_counter
import json
import matplotlib.pyplot as plt
import numpy as np
import requests as re

import aiohttp

async def req(s):
    async with s.get(f'http://localhost:5000/home') as r:
        if r.status != 200:
            r.raise_for_status()
        return await r.text()

async def req_all(s, cnt):
    tasks = []
    for _ in cnt:
        task = asyncio.create_task(req(s))
        tasks.append(task)
    res = await asyncio.gather(*tasks)
    return res

async def run_async_req():
    cnt = range(0, 10000)
    async with aiohttp.ClientSession() as session:
        htmls = await req_all(session, cnt)
        serv_load = {}
        for ht in htmls:
            d = json.loads(ht)
            k = d["message"].split(" ")[-2]+ " " + d["message"].split(" ")[-1]
            if k not in serv_load.keys():
                serv_load[k] = 0
            serv_load[k]+=1
        sum = 0
        for vals in serv_load.values():
            sum += vals
        
        x_data = []
        y_data = []
        for key in serv_load.keys():
            x_data.append(key)
            y_data.append(serv_load[key])

        return x_data, y_data, np.mean(y_data), np.std(y_data)
    
async def main():
    headers = {'Content-Type': 'application/json'}
    print("Test 1: Sending 10000 requests for N = 3 ...")
    start = perf_counter()
    x_data, y_data, x_m, x_s = await run_async_req()
    stop = perf_counter()
    print("Time Taken:{0:5.2f} seconds\n".format(stop - start))
    plt.xlabel('Servers')  
    plt.ylabel('Load') 
    plt.title("Server vs Load plot")
    plt.bar(x_data, y_data)
    plt.savefig("A1.png")
    print("Bar chart saved in A1.png")
    print()
    plt.clf()

    print("Test 2: Running 10000 requests for N = 2 to 6 ...")
    data = {'n':2, "hostnames":[]}
    res = re.delete("http://localhost:5000/rm", headers = headers, json=data)
    if res.status_code != 200:
        res.raise_for_status()
    
    start = perf_counter()
    N_data = []
    means = []
    stds = []
    data = {'n':1, "hostnames":[]}
    for i in range(2, 7):
        res = re.post("http://localhost:5000/add", headers = headers, json=data)
        if res.status_code != 200:
            res.raise_for_status()
        x_d, y_d, x_m, x_s = await run_async_req()
        means.append(x_m)
        stds.append(x_s)
        print(f'Servers : {i} ==> Mean: {round(x_m, 3)}, Standard Deviation: {round(x_s, 3)}')
        N_data.append(i)
    stop = perf_counter()
    print("Time Taken: {0:5.2f} seconds\n".format(stop - start))
    plt.plot(N_data, means, label='Mean', color='blue', marker='o')
    # for x, y in zip(N_data, means):
    #     plt.text(x, y, f'({x}, {round(y, 2)})', ha='right', va='bottom')
    plt.plot(N_data, stds, label='Standard Deviation', color='red', marker='s')
    # for x, y in zip(N_data, stds):
    #     plt.text(x, y, f'({x}, {round(y, 2)})', ha='right', va='bottom')
    plt.xlabel('N')  
    plt.ylabel('Mean/Standard Deviation (Load)') 
    plt.title("N vs Mean/Standard Deviation")

    plt.legend()
    plt.savefig("A2.png")

    print("Line Chart saved in fig A2.png")
    print()

    data = {'n': 3, "hostnames": []}
    res = re.delete("http://localhost:5000/rm", headers = headers, json = data)
    if res.status_code != 200:
        res.raise_for_status()

if __name__ == '__main__':
    start = perf_counter()
    asyncio.run(main())
    stop = perf_counter()
    print("Analysis completed in {0:5.2f} seconds\n".format(stop - start))
