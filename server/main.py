from google.appengine.ext import ndb
import webapp2
import json
import os
import base64


class State(ndb.Model):
    Value = ndb.BlobProperty()


def GetState():
    k = ndb.Key(State, 'state')
    entity = k.get()
    res = entity.Value
    return res


class MainPage(webapp2.RequestHandler):
    def get(self):
        # Set headers
        self.response.headers['Content-Type'] = 'application/json; charset=utf-8'
        self.response.write(json.dumps({'text': "Hello World"}))


class StateContent(webapp2.RequestHandler):
    def get(self):
        self.response.headers['Content-Type'] = 'application/json; charset=utf-8'
        self.response.write(json.dumps(GetState()))


app = webapp2.WSGIApplication([
    ('/', MainPage),
    ('/api/getState', StateContent),
], debug=True)
