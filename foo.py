import json
import sys
import time

i = 0
while 1:
    time.sleep(2)
    if '--json' in sys.argv:
        sys.stderr.write(json.dumps(dict(i=i))+'\n')
    else:
        sys.stderr.write(f'{i}\n')
    i += 1
