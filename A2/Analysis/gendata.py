import random
import json

def get_rand_name():
    name = ""
    for i in range(5):
        r = random.randint(0, 25)
        name += chr(ord('a') + r)
    return name

def gen_data(n, lower, upper):
    data = []
    id_map = {}
    for i in range(n):
        while True:
            r = random.randint(lower, upper)
            if r not in id_map:
                id_map[r] = True
                name = get_rand_name()
                marks = random.randint(0, 100)
                entry = {}
                entry['id'] = r
                entry['name'] = name
                entry['marks'] = marks
                data.append(entry)
                break
    return data


def main():
    n = int(input("Enter the number of data points: "))
    lower = int(input("Enter the lower bound: "))
    upper = int(input("Enter the upper bound: "))
    if upper - lower + 1 < n:
        print("The range is too small for the number of data points.")
        return
    data = gen_data(n, lower, upper)

    with open('data.json', 'w') as f:
        json.dump(data, f)


if __name__ == '__main__':
    main()