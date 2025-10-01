from locust import HttpUser, task

class PrometheusUser(HttpUser):
    @task
    def prometheus(self):
        self.client.get("/graph")
        self.client.get("/metrics")

