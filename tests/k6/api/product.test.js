import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

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

// Setup: Create a category first since products depend on it
export function setup() {
  const payload = JSON.stringify({
    name: `Setup-Category-${randomString(5)}`,
    description: 'Category created for product tests',
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_URL}/api/categories`, payload, params);
  
  if (res.status === 201) {
    return { categoryId: res.json('data.id') };
  }
  console.error(`Setup failed: Unable to create category. Status: ${res.status}`);
  return { categoryId: null };
}

export default function (data) {
  const categoryId = data.categoryId;
  let productId;

  if (!categoryId) {
    console.error('Setup failed: Could not create category');
    return;
  }

  group('Product API Functional Test', function () {
    // 1. Create Product (POST)
    group('Create Product', function () {
      const payload = JSON.stringify({
        name: `Product-${randomString(5)}`,
        price: randomIntBetween(1000, 100000),
        stock: randomIntBetween(1, 100),
        category_id: categoryId,
      });

      const params = {
        headers: {
          'Content-Type': 'application/json',
        },
      };

      const res = http.post(`${BASE_URL}/api/products`, payload, params);

      check(res, {
        'is status 201': (r) => r.status === 201,
        'has id': (r) => r.json('data.id') !== undefined,
        'has correct category': (r) => r.json('data.category_id') === categoryId,
      });

      if (res.status === 201) {
        productId = res.json('data.id');
      } else {
        console.error(`Failed to create product. Status: ${res.status}, Body: ${res.body}`);
      }
    });

    sleep(SLEEP_DURATION);

    // 2. Get All Products (GET)
    group('Get All Products', function () {
      const res = http.get(`${BASE_URL}/api/products`);

      check(res, {
        'is status 200': (r) => r.status === 200,
        'is array': (r) => Array.isArray(r.json('data')),
      });
    });

    sleep(SLEEP_DURATION);

    if (productId) {
      // 3. Get Product By ID (GET)
      group('Get Product By ID', function () {
        const res = http.get(`${BASE_URL}/api/products/${productId}`);

        check(res, {
          'is status 200': (r) => r.status === 200,
          'id matches': (r) => r.json('data.id') === productId,
          // Verify detailed response includes category
          'has category details': (r) => r.json('data.category') !== undefined,
        });
      });

      sleep(SLEEP_DURATION);

      // 4. Update Product (PUT)
      group('Update Product', function () {
        const payload = JSON.stringify({
          name: `Updated-Product-${randomString(5)}`,
          price: randomIntBetween(5000, 50000),
          stock: randomIntBetween(10, 50),
          category_id: categoryId,
        });

        const params = {
          headers: {
            'Content-Type': 'application/json',
          },
        };

        const res = http.put(`${BASE_URL}/api/products/${productId}`, payload, params);

        check(res, {
          'is status 200': (r) => r.status === 200,
          'name updated': (r) => r.json('data.name').startsWith('Updated-Product'),
        });
      });

      sleep(SLEEP_DURATION);

      // 5. Delete Product (DELETE)
      group('Delete Product', function () {
        const res = http.del(`${BASE_URL}/api/products/${productId}`);

        check(res, {
          'is status 204': (r) => r.status === 204,
        });
      });

      sleep(SLEEP_DURATION);

      // 6. Verify Deletion (GET should fail)
      group('Verify Deletion', function () {
        const res = http.get(`${BASE_URL}/api/products/${productId}`);
        
        check(res, {
          'is status 404': (r) => r.status === 404,
        });
      });
    }
  });
}

// Teardown: Clean up the setup category (optional but good practice)
export function teardown(data) {
  if (data.categoryId) {
     http.del(`${BASE_URL}/api/categories/${data.categoryId}`);
  }
}
