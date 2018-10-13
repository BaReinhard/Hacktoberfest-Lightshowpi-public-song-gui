import ndb
import webapp2
import json
import os


class DstoreObject(ndb.Model):
    """Models an individual Guestbook entry with content and date."""
    Value = ndb.StringProperty()

    @classmethod
    def query_state(cls, ancestor_key):
        return cls.query(ancestor=ancestor_key)


def GetState():
    # Instantiates a client

    # The kind for the new entity
    kind = 'State'
    # The name/ID for the new entity
    name = 'state'
    # The Cloud Datastore key for the new entity
    task_key = ndb.Key(kind, name)
    check = DstoreObject.query_state(task_key)
    print(check)
    # Prepares the new entity
    # task = datastore.Entity(key=task_key)
    # # task = datastore.Entity(key=task_key)
    # res = task.get()

    # Saves the entity
    # datastore_client.put(task)
    return check
    # print('Saved {}: {}'.format(task.key.name, task['description']))


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
