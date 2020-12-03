version1 = False

import re

inputfile = open('data/in2.txt', 'r')
lines = inputfile.readlines()

valid = 0

for line in lines:
    res = re.findall("([0-9]+)-([0-9]+) (.): (.*)", line.strip())
    occs_a, occs_b = int(res[0][0]), int(res[0][1])
    letter = res[0][2]
    pw = res[0][3]
    
    if version1:
        count = pw.count(letter)
        if count >= occs_a and count <= occs_b:
            valid += 1
    else:
        if (pw[occs_a-1] == letter) ^ (pw[occs_b-1] == letter):
            valid += 1
        
print(valid)