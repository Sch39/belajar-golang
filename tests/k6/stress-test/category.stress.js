import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '30s', target: 300 },
    { duration: '30s', target: 600 },
    { duration: '30s', target: 900 },
    { duration: '30s', target: 1500 },
    { duration: '30s', target: 1500 },
    { duration: '30s', target: 2000 },
    { duration: '30s', target: 2500 },
    { duration: '30s', target: 3000 },
    { duration: '30s', target: 3500 },
    { duration: '30s', target: 4000 },
    { duration: '30s', target: 5000 },
    { duration: '30s', target: 6000 },
    { duration: '30s', target: 7000 },
    { duration: '30s', target: 8000 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.3'],       // error boleh
    http_req_duration: ['p(95)<1500'],   // latency boleh naik
  },
};

export default function () {
  const payload = JSON.stringify({
    name: `Stress-${randomString(8)}`,
    description: `Stress test ${randomString(10)}`,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(`${BASE_URL}/api/categories`, payload, params);

  check(res, {
    'status is 201 or 429': (r) => r.status === 201 || r.status === 429,
  });
}