version1 = False

datafile = open('data/in3.txt', 'r')
lines = datafile.readlines()

trees_mult = 1

if version1:
    slopes = [(3,1)]
else:
    slopes = [(1,1), (3,1), (5, 1), (7, 1), (1,2)]

for slope in slopes:
    trees = 0
    curr_x = 0
    for l in range(0,len(lines),slope[1]):
        line = lines[l].strip()
        pos = curr_x % len(line)
        
        if line[pos] == '#':
            trees += 1
            
        curr_x += slope[0]
        
    trees_mult *= trees
    
print(trees_mult)