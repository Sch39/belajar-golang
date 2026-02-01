import http from 'k6/http';
import { check } from 'k6';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '30s', target: 300 },
    { duration: '30s', target: 600 },
    { duration: '30s', target: 1000 },
    { duration: '30s', target: 1500 },
    { duration: '30s', target: 2000 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.3'],      // error boleh
    http_req_duration: ['p(95)<2000'],  // latency boleh naik
  },
};

// SETUP: minimal & sekali saja
export function setup() {
  const payload = JSON.stringify({
    name: `Stress-Cat-${randomString(5)}`,
    description: 'Category for product stress test',
  });

  const res = http.post(`${BASE_URL}/api/categories`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (res.status !== 201) {
    throw new Error('Setup failed: cannot create category');
  }

  return { categoryId: res.json('data.id') };
}

export default function (data) {
  const payload = JSON.stringify({
    name: `Stress-Prod-${randomString(8)}`,
    price: randomIntBetween(1000, 100000),
    stock: randomIntBetween(1, 100),
    category_id: data.categoryId,
  });

  const res = http.post(`${BASE_URL}/api/products`, payload, {
    headers: { 'Content-Type': 'application/json' },
    tags: {
      feature: 'product',
      type: 'stress',
      expected_response: 'true',
    },
  });

  check(res, {
    'status 201 or 429': (r) =>
      r.status === 201 || r.status === 429,
  });
}

// TEARDOWN (optional, boleh di-skip di stress test)
export function teardown(data) {
  http.del(`${BASE_URL}/api/categories/${data.categoryId}`);
}