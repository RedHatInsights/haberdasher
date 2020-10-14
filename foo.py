import http.server
import json
import sys
import time

# This is useful for testing stdout vs. stderr
print('Python starting')

def heartbeat():
    i = 0
    while 1:
        time.sleep(2)
        if '--json' in sys.argv:
            sys.stderr.write(json.dumps(dict(i=i))+'\n')
        else:
            sys.stderr.write(f'{i}\n')
        sys.stderr.flush()
        i += 1

class FooHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(self.headers.items()).encode('utf8'))
    
def serve():
    server_address = ('', 8080)
    httpd = http.server.HTTPServer(server_address, FooHandler)
    httpd.serve_forever()

if '--serve' in sys.argv:
    serve()
else:
    heartbeat()
