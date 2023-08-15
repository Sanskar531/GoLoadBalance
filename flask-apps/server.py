from flask import Flask

app = Flask(__name__)

@app.route("/ping")
def ping():
    return "pong"

@app.route("/hello_world")
def hello_world():
    return "<h1>Hello!!! How are you?</h1>"
