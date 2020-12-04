import re
import sys

EYE_COLORS = {"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}

def HEIGHT_HANDLING(v):
    s = re.findall("([0-9]+)(in|cm)", v)
    if len(s) != 1:
        return False
    h, t = int(s[0][0]), s[0][1]
    return (t == 'in' and h >= 59 and h <= 76) or (t == 'cm' and h >= 150 and h <= 193)

EVL = { "byr": (0, lambda v: None != re.match("[0-9]{4}$", v) and int(v) >= 1920 and int(v) <= 2002),
        "iyr": (1, lambda v: None != re.match("[0-9]{4}$", v) and int(v) >= 2010 and int(v) <= 2020),
        "eyr": (2, lambda v: None != re.match("[0-9]{4}$", v) and int(v) >= 2020 and int(v) <= 2030),
        "hgt": (3, HEIGHT_HANDLING),
        "hcl": (4, lambda v: None != re.match("#[a-f0-9]{6}$", v)),
        "ecl": (5, lambda v: v in EYE_COLORS),
        "pid": (6, lambda v: None != re.match("[0-9]{9}$", v)),
        "cid": (7, lambda v: True)}

OPT = (7,)

lines = sys.stdin.read()
entries = lines.split("\n\n")

validated = [False] * len(EVL)
valid_passports = 0

for e in entries:
    for i in range(len(validated)):
        validated[i] = False
    invalid = False
    
    fields = re.findall("([a-z]+):([^\s]+)", e)
    
    for f in fields:
        label, value = f[0], f[1]
        id, validate = EVL.get(label, (-1, lambda v: False))

        if not validate(value) or validated[id] == True: #unexpected or wrong field or double entry
            invalid = True
            break
            
        validated[id] = True

    if not invalid:
        for o in OPT:
            validated[o] |= True
        if not False in validated:
            valid_passports += 1
        
print(valid_passports)

