// scripts/performance/load-test.js

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// --- Metrics ---
const GetArticleLatency = new Trend('get_article_latency');
const RegisterUserLatency = new Trend('register_user_latency');
const ErrorRate = new Rate('error_rate');
const SuccessRate = new Rate('success_rate');
const TotalRequests = new Counter('total_requests');

// --- Test Options ---
export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp-up to 20 virtual users
    { duration: '1m', target: 20 },   // Stay at 20 VUs
    { duration: '10s', target: 100 }, // Spike to 100 VUs
    { duration: '30s', target: 100 }, // Stay at 100 VUs
    { duration: '10s', target: 0 },    // Ramp-down
  ],
  thresholds: {
    'http_req_failed': ['rate<0.01'], // < 1% errors
    'http_req_duration': ['p(95)<500'], // 95% of requests must be below 500ms
    'get_article_latency{status:200}': ['p(95)<300'],
  },
};

export default function (data) {
  // --- Scenario 1: یک کاربر ناشناس مقالات را می‌خواند ---
  const res = http.get(`${__ENV.TARGET_URL}/api/v1/articles/some-test-slug`);
  const getArticleCheck = check(res, {
    'GET /articles/:slug status is 200': (r) => r.status === 200,
  });
  GetArticleLatency.add(res.timings.duration, { status: res.status });
  ErrorRate.add(!getArticleCheck);
  SuccessRate.add(getArticleCheck);
  TotalRequests.add(1);

  sleep(1);

  // --- Scenario 2: یک کاربر جدید ثبت‌نام می‌کند ---
  const email = `user-${__VU}-${__ITER}@test.com`;
  const registerPayload = JSON.stringify({
    email: email,
    password: 'Password123!',
    username: `user_${__VU}_${__ITER}`,
    fullName: 'Load Test User',
  });
  const params = { headers: { 'Content-Type': 'application/json' } };
  
  const registerRes = http.post(`${__ENV.TARGET_URL}/api/v1/auth/register`, registerPayload, params);
  const registerCheck = check(registerRes, {
    'POST /auth/register status is 201': (r) => r.status === 201,
  });
  RegisterUserLatency.add(registerRes.timings.duration);
  ErrorRate.add(!registerCheck);
  SuccessRate.add(registerCheck);
  TotalRequests.add(1);

  sleep(2);
}