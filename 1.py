import re

page = open('data/in1.txt', 'r')

data = page.readlines()

numbers = []
for d in data:
    numbers.append(int(d))

for i in range(len(numbers)):
    for j in range(i+1, len(numbers)):
        for k in range(j+1, len(numbers)):
            if numbers[i] + numbers[j] + numbers[k] == 2020:
                print(numbers[i]*numbers[j]*numbers[k])
        