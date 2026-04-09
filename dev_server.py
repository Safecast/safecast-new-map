import http.server
import socketserver
import json
import re

PORT = 8002
DIRECTORY = "public_html"

class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)

    def translate_path(self, path):
        # The Go server mounts public_html to /static for assets
        if path.startswith('/static/'):
            path = path.replace('/static/', '/', 1)
        return super().translate_path(path)

    def do_GET(self):
        if self.path == '/map.html' or self.path == '/':
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            
            with open("public_html/map.html", "r", encoding="utf-8") as f:
                content = f.read()
                
            try:
                with open("public_html/translations.json", "r", encoding="utf-8") as f:
                    translations = f.read()
            except:
                translations = "{}"
                
            # Simulate Go Templating for frontend UI tests
            content = content.replace("{{ .TranslationsJSON }}", translations)
            content = content.replace("{{ .Lang }}", "en")
            content = re.sub(r'\{\{\s*translate\s*"([^"]+)"\s*\}\}', r'\1', content)
            
            # Additional tags that block JS initialization
            content = content.replace('{{printf "%.6f" .DefaultLat}}', "44.08832")
            content = content.replace('{{printf "%.6f" .DefaultLon}}', "42.97577")
            content = content.replace("{{.DefaultZoom}}", "11")
            content = content.replace('{{ printf "%q" .DefaultLayer }}', '"OpenStreetMap"')
            content = content.replace('{{if .AutoLocateDefault}}true{{else}}false{{end}}', 'false')
            content = content.replace('{{if .SupportEmail}}{{printf "%q" .SupportEmail}}{{else}}""{{end}}', '""')
            content = content.replace('{{if .RealtimeAvailable}}true{{else}}false{{end}}', 'false')
            content = content.replace('{{ .Version }}', 'dev')
            content = content.replace('{{ .MarkersJSON }}', '[{"lat":35.6895,"lon":139.6917,"doseRate":0.05,"trackID":1,"time":"2023-01-01T00:00:00Z","dataType":1},{"lat":35.7100,"lon":139.8107,"doseRate":0.12,"trackID":2,"time":"2023-01-01T12:00:00Z","dataType":2},{"lat":35.6580,"lon":139.7016,"doseRate":0.25,"trackID":3,"time":"2023-01-02T00:00:00Z","dataType":3},{"lat":35.6266,"lon":139.7562,"doseRate":1.5,"trackID":4,"time":"2023-01-03T00:00:00Z","dataType":4},{"lat":35.6881,"lon":139.7710,"doseRate":0.08,"trackID":5,"time":"2023-01-04T00:00:00Z","dataType":1}]')

            self.wfile.write(content.encode("utf-8"))
        else:
            super().do_GET()

with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print("Serving at port", PORT)
    httpd.serve_forever()
