import re, sys

version1 = False

EVL = { "byr": (0, lambda v: 1920 <= int(v) <= 2002), "iyr": (1, lambda v: 2010 <= int(v) <= 2020),
        "eyr": (2, lambda v: 2020 <= int(v) <= 2030), "ecl": (3, lambda v: v in {"amb","blu","brn","gry","grn","hzl","oth"}),
        "hcl": (4, lambda v: re.match("#[a-f0-9]{6}$", v)), "pid": (5, lambda v: re.match("[0-9]{9}$", v)),
        "hgt": (6, lambda v: 59 <= int(v[:-2]) <= 76 if v[-2:] == 'in' else 150 <= int(v[:-2]) <= 193) }

valid_passports = 0
lines = sys.stdin.read()

for e in lines.split("\n\n"):
    valid = [False] * len(EVL)
    
    for f in e.split():
        factors = f.split(':')
        id, func = EVL.get(factors[0], (0, lambda v : False))
        try:
            valid[id] = True if version1 or func(factors[1]) else valid[id]
        except:
            continue
            
    valid_passports += 1 if not False in valid else 0
        
print(valid_passports)