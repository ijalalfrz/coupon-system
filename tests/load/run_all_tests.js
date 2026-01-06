
import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

export const env = {
    host: 'http://host.docker.internal:8080',
}

// Custom metrics
const couponSuccessRace = new Counter('coupon_success_race');
const couponSuccessSameUser = new Counter('coupon_success_same_user');

export const options = {
    scenarios: {
        race_condition: {
            executor: 'per-vu-iterations',
            exec: 'raceConditionTest',
            vus: 50,
            iterations: 1,
            maxDuration: '10s',
            startTime: '0s',
        },
        same_user_race: {
            executor: 'per-vu-iterations',
            exec: 'sameUserRaceTest',
            vus: 10,
            iterations: 1,
            maxDuration: '10s',
            startTime: '15s', // Start after the first test finishes
        },
    },
    thresholds: {
        'coupon_success_race': ['count==5'],
        'coupon_success_same_user': ['count==1'],
    },
};

export function setup() {
    // Setup for Race Condition Test
    const promoNameRace = `RACE_${Date.now()}`;
    const resRace = http.post(
        `${env.host}/api/coupons`,
        JSON.stringify({
            name: promoNameRace,
            amount: 5,
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(resRace, {
        'race: coupon created': (r) => r.status === 201 || r.status === 200,
    });

    // Setup for Same User Test
    const promoNameSameUser = `SAME_USER_${Date.now()}`;
    const resSameUser = http.post(
        `${env.host}/api/coupons`,
        JSON.stringify({
            name: promoNameSameUser,
            amount: 10,
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(resSameUser, {
        'same_user: coupon created': (r) => r.status === 201 || r.status === 200,
    });

    return { promoNameRace, promoNameSameUser };
}

export function raceConditionTest(data) {
    const res = http.post(
        `${env.host}/api/coupons/claim`,
        JSON.stringify({
            coupon_name: data.promoNameRace,
            user_id: `${__VU}`,
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    if (res.status === 201) {
        couponSuccessRace.add(1);
    }

    check(res, {
        'race: status valid': (r) =>
            r.status === 201 || r.status === 409 || r.status === 400,
    });
}

export function sameUserRaceTest(data) {
    const res = http.post(
        `${env.host}/api/coupons/claim`,
        JSON.stringify({
            coupon_name: data.promoNameSameUser,
            user_id: "same-user-id",
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    if (res.status === 201) {
        couponSuccessSameUser.add(1);
    }

    check(res, {
        'same_user: status valid': (r) =>
            r.status === 201 || r.status === 409 || r.status === 400,
    });
}
