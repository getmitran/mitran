import unittest
from unittest.mock import patch, MagicMock
from worker.grpc_client import MitranGRPCClient, TaskResult, TaskAssignment


class TestMitranGRPCClient(unittest.TestCase):
    def setUp(self):
        self.client = MitranGRPCClient(engine_host='localhost', grpc_port=7781, http_port=7780)

    @patch('requests.post')
    def test_register_success(self, mock_post):
        mock_post.return_value = MagicMock(status_code=200)
        assert self.client.register('dev-agent', 'development', 'http://localhost:8888/task') is True
        mock_post.assert_called_once()

    @patch('requests.post')
    def test_register_failure(self, mock_post):
        mock_post.return_value = MagicMock(status_code=500)
        assert self.client.register('dev-agent', 'development', 'http://localhost:8888/task') is False

    @patch('requests.post')
    def test_submit_result_success(self, mock_post):
        mock_post.return_value = MagicMock(status_code=200)
        result = TaskResult(
            task_id='task-001',
            summary='Created README.md',
            files=[{'path': 'README.md', 'action': 'create', 'content': '# Hello'}],
            status='complete',
        )
        assert self.client.submit_result(result) is True

    @patch('requests.post')
    def test_submit_result_network_error(self, mock_post):
        import requests
        mock_post.side_effect = requests.ConnectionError()
        result = TaskResult(task_id='task-001', summary='fail')
        assert self.client.submit_result(result) is False

    @patch('requests.get')
    def test_health_check_success(self, mock_get):
        mock_get.return_value = MagicMock(status_code=200)
        mock_get.return_value.json.return_value = {'status': 'ok', 'version': '0.1.0'}
        health = self.client.health_check()
        assert health['status'] == 'ok'

    @patch('requests.get')
    def test_health_check_unreachable(self, mock_get):
        import requests
        mock_get.side_effect = requests.ConnectionError()
        health = self.client.health_check()
        assert health['status'] == 'unreachable'

    def test_task_assignment_dataclass(self):
        task = TaskAssignment(
            task_id='t1', title='Setup', description='Init project',
            agent_type='development', project_id='p1', workspace_path='/tmp',
        )
        assert task.context == {}
        assert task.skill_paths == []

    def test_close_no_channel(self):
        self.client.close()  # Should not raise


if __name__ == '__main__':
    unittest.main()
