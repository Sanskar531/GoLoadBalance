from flask import Flask

app = Flask(__name__)

@app.route("/ping")
def ping():
    print("here")
    return "pong"

app.run(host='0.0.0.0')
