// load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// カスタムメトリクスの定義
const errorRate = new Rate('errors');

// テスト設定
export const options = {
    stages: [
        { duration: '10s', target: 200 },  // 徐々に200ユーザーまで増加
        { duration: '1m', target: 200 },   // 1分間200ユーザーを維持
        { duration: '10s', target: 0 },   // 徐々にユーザーを減少
    ],
    thresholds: {
        'http_req_duration': ['p(95)<500'], // 95%のリクエストが500ms以下
        'errors': ['rate<0.1'],             // エラー率10%以下
    },
};

const BASE_URL = 'http://host.docker.internal:3000';

// テストデータ
const testUser = {
    username: 'Test User',
    password_hash: 'hoge',
    email: 'test@example.com',
};

const testUser2 = {
    username: 'Test User2',
    password_hash: 'hoge',
    email: 'test@example.com',
};



export default function () {
    // 特定のユーザーの取得
    let userId = 1;
    // ユーザー一覧の取得
    const listResponse = http.get(`${BASE_URL}/users`);
    check(listResponse, {
        'ユーザー一覧取得成功': (r) => r.status === 200,
        'ユーザー一覧レスポンスタイム < 200ms': (r) => r.timings.duration < 200,
    });
    errorRate.add(listResponse.status !== 200);

    const getUserResponse = http.get(`${BASE_URL}/user/${userId}`);
    check(getUserResponse, {
        '特定ユーザー取得成功': (r) => r.status === 200,
        '特定ユーザーレスポンスタイム < 200ms': (r) => r.timings.duration < 200,
    });
    errorRate.add(getUserResponse.status !== 200);

    // 新規ユーザーの作成
    const createResponse = http.put(
        `${BASE_URL}/user`,
        JSON.stringify(testUser),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(createResponse, {
        'ユーザー作成成功': (r) => r.status === 200,
        'ユーザー作成レスポンスタイム < 300ms': (r) => r.timings.duration < 300,
    });
    errorRate.add(createResponse.status !== 200);
    
    userId = createResponse.body
    
    // ユーザーの更新
    const updateResponse = http.post(
        `${BASE_URL}/user/${userId}`,
        JSON.stringify(testUser2),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(updateResponse, {
        'ユーザー更新成功': (r) => r.status === 200,
        'ユーザー更新レスポンスタイム < 300ms': (r) => r.timings.duration < 300,
    });
    errorRate.add(updateResponse.status !== 200);

    // ユーザーの更新
    const deleteResponse = http.del(
        `${BASE_URL}/user/${userId}`,
        JSON.stringify(testUser2),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(deleteResponse, {
        'ユーザー削除成功': (r) => r.status === 200,
        'ユーザー削除レスポンスタイム < 300ms': (r) => r.timings.duration < 300,
    });
    errorRate.add(deleteResponse.status !== 200);
}