import sys
import time

i = 0
while 1:
    time.sleep(2)
    sys.stderr.write(f'{i}\n')
    i += 1
