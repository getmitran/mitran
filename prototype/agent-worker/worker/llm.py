import json
import boto3
from worker.config import AWS_REGION, MODEL_ID


_client = None


def get_client():
    global _client
    if _client is None:
        _client = boto3.client("bedrock-runtime", region_name=AWS_REGION)
    return _client


def invoke(system_prompt: str, user_message: str, max_tokens: int = 4096) -> str:
    """Call AWS Bedrock Claude and return the text response."""
    client = get_client()
    body = json.dumps({
        "anthropic_version": "bedrock-2023-05-31",
        "max_tokens": max_tokens,
        "system": system_prompt,
        "messages": [{"role": "user", "content": user_message}],
    })
    try:
        response = client.invoke_model(modelId=MODEL_ID, body=body, contentType="application/json")
        result = json.loads(response["body"].read())
        return result["content"][0]["text"]
    except Exception as e:
        raise RuntimeError(f"Bedrock invocation failed: {e}") from e
