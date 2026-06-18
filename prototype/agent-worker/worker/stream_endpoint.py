from flask import Response, request as flask_request
from worker.llm_streaming import stream_invoke
from worker.agents import AGENTS
import json


def register_stream_routes(app):
    @app.route('/stream', methods=['POST'])
    def stream_chat():
        data = flask_request.get_json()
        message = data.get('message', '')
        agent_name = data.get('agent', 'dev')

        agent = AGENTS.get(agent_name)
        system_prompt = getattr(agent, 'system_prompt', 'You are a helpful assistant.') if agent else 'You are a helpful assistant.'

        def generate():
            for chunk in stream_invoke(system_prompt, message):
                yield f'data: {json.dumps({"text": chunk})}\n\n'
            yield 'data: {"done": true}\n\n'

        return Response(generate(), mimetype='text/event-stream',
                       headers={'Cache-Control': 'no-cache', 'X-Accel-Buffering': 'no'})
