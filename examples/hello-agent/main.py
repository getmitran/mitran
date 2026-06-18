import requests, time

class HelloAgent:
    BASE = "http://localhost:7780/api/v1/tasks"

    def run(self, description="Hello from example agent!"):
        resp = requests.post(self.BASE, json={"description": description})
        resp.raise_for_status()
        task_id = resp.json()["id"]
        print(f"Created task {task_id}, polling...")
        while True:
            status = requests.get(f"{self.BASE}/{task_id}").json()["status"]
            if status in ("completed", "failed"):
                print(f"Task {task_id} finished: {status}")
                return status
            time.sleep(2)

if __name__ == "__main__":
    HelloAgent().run()
