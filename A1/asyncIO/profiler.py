import asyncio
from time import perf_counter
import json
import matplotlib.pyplot as plt
import numpy as np

import aiohttp


async def fetch(s, cnt):
    async with s.get(f'http://localhost:5000/home') as r:
        if r.status != 200:
            r.raise_for_status()
        return await r.text()


async def fetch_all(s, cnt):
    tasks = []
    for c in cnt:
        task = asyncio.create_task(fetch(s, c))
        tasks.append(task)
    res = await asyncio.gather(*tasks)
    return res


async def main(n):
    cnt = range(0, 10000)
    async with aiohttp.ClientSession() as session:
        htmls = await fetch_all(session, cnt)
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
        
        plt.xlabel('Servers')  
        plt.ylabel('Load') 

        print("Average load per server: {0:5.2f}\n".format(np.mean(y_data)))
        print("Standard Deviation: {0:5.2f}\n".format(np.std(y_data)))
        
        plt.bar(x_data, y_data)
        plt.savefig(f"output_{n}.png")

if __name__ == '__main__':
    start = perf_counter()
    asyncio.run(main(3))
    stop = perf_counter()
    print("time taken:{0:5.2f} seconds\n".format(stop - start))