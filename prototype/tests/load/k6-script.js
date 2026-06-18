import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:7780';

const healthErrors = new Rate('health_errors');
const chatErrors = new Rate('chat_errors');
const memoryErrors = new Rate('memory_errors');
const healthDuration = new Trend('health_duration');
const chatDuration = new Trend('chat_duration');
const memoryDuration = new Trend('memory_duration');

export const options = {
  scenarios: {
    health: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      maxVUs: 500,
      stages: [
        { duration: '30s', target: 1000 },
        { duration: '60s', target: 1000 },
        { duration: '10s', target: 0 },
      ],
      exec: 'healthCheck',
    },
    agent_chat: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 200,
      stages: [
        { duration: '30s', target: 50 },
        { duration: '60s', target: 50 },
        { duration: '10s', target: 0 },
      ],
      exec: 'agentChat',
    },
    memory_ops: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 300,
      stages: [
        { duration: '30s', target: 200 },
        { duration: '60s', target: 200 },
        { duration: '10s', target: 0 },
      ],
      exec: 'memoryOps',
    },
  },
  thresholds: {
    'health_duration': ['p(95)<500'],
    'chat_duration': ['p(95)<2000'],
    'memory_duration': ['p(95)<500'],
    'health_errors': ['rate<0.01'],
    'chat_errors': ['rate<0.05'],
    'memory_errors': ['rate<0.01'],
  },
};

export function healthCheck() {
  const res = http.get(`${BASE_URL}/health`);
  healthDuration.add(res.timings.duration);
  healthErrors.add(res.status !== 200);
  check(res, { 'health 200': (r) => r.status === 200 });
}

export function agentChat() {
  const payload = JSON.stringify({
    message: 'Hello, run tests',
    agent: 'dev',
    project_id: 'load-test',
  });
  const params = { headers: { 'Content-Type': 'application/json' } };
  const res = http.post(`${BASE_URL}/api/v1/chat`, payload, params);
  chatDuration.add(res.timings.duration);
  chatErrors.add(res.status !== 200);
  check(res, { 'chat 200': (r) => r.status === 200 });
}

export function memoryOps() {
  const ops = ['store', 'search'];
  const op = ops[Math.floor(Math.random() * ops.length)];

  let res;
  if (op === 'store') {
    const payload = JSON.stringify({
      key: `load-test-${Date.now()}`,
      value: 'test data for load testing',
      project_id: 'load-test',
    });
    const params = { headers: { 'Content-Type': 'application/json' } };
    res = http.post(`${BASE_URL}/api/v1/memory`, payload, params);
  } else {
    res = http.get(`${BASE_URL}/api/v1/memory?project_id=load-test&q=test`);
  }
  memoryDuration.add(res.timings.duration);
  memoryErrors.add(res.status !== 200);
  check(res, { 'memory 200': (r) => r.status === 200 });
}
