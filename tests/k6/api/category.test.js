import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const SLEEP_DURATION = 1;

export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    checks: ['rate==1.00'], // All checks must pass (100% success)
  },
};

export default function () {
  let categoryId;

  group('Category API Functional Test', function () {
    // 1. Create Category (POST)
    group('Create Category', function () {
      const payload = JSON.stringify({
        name: `Category-${randomString(5)}`,
        description: `Description for category ${randomString(5)}`,
      });

      const params = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

      const res = http.post(`${BASE_URL}/api/categories`, payload, params);

      check(res, {
        'is status 201': (r) => r.status === 201,
        'has id': (r) => r.json('data.id') !== undefined,
        'has correct name': (r) => r.json('data.name') !== undefined,
      });

      if (res.status === 201) {
        categoryId = res.json('data.id');
      } else {
        console.error(`Failed to create category. Status: ${res.status}, Body: ${res.body}`);
      }
    });

    sleep(SLEEP_DURATION);

    // 2. Get All Categories (GET)
    group('Get All Categories', function () {
      const res = http.get(`${BASE_URL}/api/categories`);

      check(res, {
        'is status 200': (r) => r.status === 200,
        'is array': (r) => Array.isArray(r.json('data')),
      });
    });

    sleep(SLEEP_DURATION);

    if (categoryId) {
      // 3. Get Category By ID (GET)
      group('Get Category By ID', function () {
        const res = http.get(`${BASE_URL}/api/categories/${categoryId}`);

        check(res, {
          'is status 200': (r) => r.status === 200,
          'id matches': (r) => r.json('data.id') === categoryId,
        });
      });

      sleep(SLEEP_DURATION);

      // 4. Update Category (PUT)
      group('Update Category', function () {
        const payload = JSON.stringify({
          name: `Updated-Category-${randomString(5)}`,
          description: `Updated description ${randomString(5)}`,
        });

        const params = {
          headers: {
            'Content-Type': 'application/json',
          },
        };

        const res = http.put(`${BASE_URL}/api/categories/${categoryId}`, payload, params);

        check(res, {
          'is status 200': (r) => r.status === 200,
          'name updated': (r) => r.json('data.name').startsWith('Updated-Category'),
        });
      });

      sleep(SLEEP_DURATION);

      // 5. Delete Category (DELETE)
      group('Delete Category', function () {
        const res = http.del(`${BASE_URL}/api/categories/${categoryId}`);

        check(res, {
          'is status 204': (r) => r.status === 204,
        });
      });
      
      sleep(SLEEP_DURATION);

      // 6. Verify Deletion (GET should fail)
      group('Verify Deletion', function () {
        const res = http.get(`${BASE_URL}/api/categories/${categoryId}`);
        
        check(res, {
          'is status 404': (r) => r.status === 404,
        });
      });
    }
  });
}
