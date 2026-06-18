import json
import boto3
import logging
from typing import Generator
from worker.config import AWS_REGION, MODEL_ID

log = logging.getLogger(__name__)
_client = None


def get_client():
    global _client
    if _client is None:
        _client = boto3.client('bedrock-runtime', region_name=AWS_REGION)
    return _client


def stream_invoke(system_prompt: str, user_message: str, max_tokens: int = 4096) -> Generator[str, None, None]:
    """Stream tokens from Bedrock Claude. Yields text chunks as they arrive."""
    client = get_client()
    body = json.dumps({
        'anthropic_version': 'bedrock-2023-05-31',
        'max_tokens': max_tokens,
        'system': system_prompt,
        'messages': [{'role': 'user', 'content': user_message}],
    })
    try:
        response = client.invoke_model_with_response_stream(
            modelId=MODEL_ID, body=body, contentType='application/json'
        )
        for event in response['body']:
            chunk = json.loads(event['chunk']['bytes'])
            if chunk['type'] == 'content_block_delta':
                text = chunk['delta'].get('text', '')
                if text:
                    yield text
    except Exception as e:
        log.error(f'Streaming error: {e}')
        yield f'[Error: {e}]'
