import sys

bags = dict()

class BagInfo:
    def __init__(self):
        self.contained = set()
        self.containing = dict()
    
def registerBag(desc):
    if not desc in bags:
        bags[desc] = BagInfo()

def visitUp(desc, to_visit):
    for b in bags[desc].contained:
        if b in to_visit:
            to_visit.remove(b)
            visitUp(b, to_visit)
            
def visitDown(desc):
    visited = 1
    for b, N in bags[desc].containing.items():
        visited += N * visitDown(b)
    return visited

for line in sys.stdin.read().split("\n"):
    infos = line.split("contain ")
    desc = infos[0][:-6]
    registerBag(desc)
    
    if infos[1] == "no other bags.":
        continue
        
    for info in infos[1].split(","):
        contain_desc = " ".join(info.split()[1:3])
        registerBag(contain_desc)
        bags[contain_desc].contained.add(desc)
        bags[desc].containing[contain_desc] = int(info.split()[0])

to_visit = set(bags.keys())
visitUp('shiny gold', to_visit)
print("All containing shiny gold:", len(bags) - len(to_visit))
print("All contained by shiny gold:", visitDown('shiny gold') - 1)