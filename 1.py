version1 = False

inputfile = open('data/in1.txt', 'r')
lines = inputfile.readlines()

numbers = []
for line in lines:
    numbers.append(int(line))

for i in range(len(numbers)):
    for j in range(i+1, len(numbers)):
        if version1:
            if numbers[i] + numbers[j] == 2020:
                print(numbers[i]*numbers[j])
                break
        else:
            for k in range(j+1, len(numbers)):
                if numbers[i] + numbers[j] + numbers[k] == 2020:
                    print(numbers[i]*numbers[j]*numbers[k])
                    break
        