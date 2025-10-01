### locust -f locustfile.py --headless --users 10 --spawn-rate 1 -H http://localhost:8000

import time
import random
import string
from locust import HttpUser, task, between

class QuickstartUser(HttpUser):
  wait_time = between(1, 5)

  @task(10)
  def home(self):
    self.client.get("/", name="get /home")

  @task(5)
  def io_task(self):
    self.client.get("/io_task", name="get /io_task")

  @task(5)
  def cpu_task(self):
    self.client.get("/cpu_task", name="get /cpu_task")

  @task(3)
  def random_sleep(self):
    self.client.get("/random_sleep", name="get random_sleep")

  @task(10)
  def random_status(self):
    self.client.get("/random_status", name="get /random_status")

  @task(3)
  def chain(self):
    self.client.get("/chain", name="get /chain")

  @task(1)
  def error_test(self):
    self.client.get("/error_test", name="get /error_test")

  @task(1)
  def error(self):
    self.client.get("/error", name="get /error")

  @task(1)
  def create_todo(self):
    self.client.post("/todos/", name="post /todos/", 
      json={
        "title": ''.join(random.choices(string.ascii_letters, k=20)),
        "description": ''.join(random.choices(string.ascii_letters, k=100)),
        "completed": False
      }
    )

  @task(5)
  def query_all_todos(self):
    self.client.get("/todos/", name="get /todos/")

  @task(10)
  def query_todo(self):
    self.client.get("/todos/" + str(random.randint(1, 10)), name="get /todos/(random_id)")

  @task(3)
  def update_todo(self):
    self.client.put("/todos/" + str(random.randint(1, 10)), name="put /todos/(random_id)",
      json={
        "completed": True
      }
    )
