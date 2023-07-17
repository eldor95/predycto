// Copyright 2022 gorse Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include <immintrin.h>
#include <stdint.h>
#include <math.h>

void _mm512_mul_const_add_to(float *a, float *b, float *c, int64_t n)
{
    int epoch = n / 16;
    int remain = n % 16;
    for (int i = 0; i < epoch; i++)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_set1_ps(*b);
        __m512 v3 = _mm512_loadu_ps(c);
        __m512 v = _mm512_fmadd_ps(v1, v2, v3);
        _mm512_storeu_ps(c, v);
        a += 16;
        c += 16;
    }
    if (remain >= 8)
    {
        __m256 v1 = _mm256_loadu_ps(a);
        __m256 v2 = _mm256_broadcast_ss(b);
        __m256 v3 = _mm256_loadu_ps(c);
        __m256 v = _mm256_add_ps(_mm256_mul_ps(v1, v2), v3);
        _mm256_storeu_ps(c, v);
        a += 8;
        c += 8;
        remain -= 8;
    }
    for (int i = 0; i < remain; i++)
    {
        c[i] += a[i] * b[0];
    }
}

void _mm512_mul_const_to(float *a, float *b, float *c, int64_t n)
{
    int epoch = n / 16;
    int remain = n % 16;
    for (int i = 0; i < epoch; i++)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_set1_ps(*b);
        __m512 v = _mm512_mul_ps(v1, v2);
        _mm512_storeu_ps(c, v);
        a += 16;
        c += 16;
    }
    if (remain >= 8)
    {
        __m256 v1 = _mm256_loadu_ps(a);
        __m256 v2 = _mm256_broadcast_ss(b);
        __m256 v = _mm256_mul_ps(v1, v2);
        _mm256_storeu_ps(c, v);
        a += 8;
        c += 8;
        remain -= 8;
    }
    for (int i = 0; i < remain; i++)
    {
        c[i] = a[i] * b[0];
    }
}

void _mm512_mul_const(float *a, float *b, int64_t n)
{
    int epoch = n / 16;
    int remain = n % 16;
    for (int i = 0; i < epoch; i++)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_set1_ps(*b);
        __m512 v = _mm512_mul_ps(v1, v2);
        _mm512_storeu_ps(a, v);
        a += 16;
    }
    if (remain >= 8)
    {
        __m256 v1 = _mm256_loadu_ps(a);
        __m256 v2 = _mm256_broadcast_ss(b);
        __m256 v = _mm256_mul_ps(v1, v2);
        _mm256_storeu_ps(a, v);
        a += 8;
        remain -= 8;
    }
    for (int i = 0; i < remain; i++)
    {
        a[i] *= b[0];
    }
}

void _mm512_mul_to(float *a, float *b, float *c, int64_t n)
{
    int epoch = n / 16;
    int remain = n % 16;
    for (int i = 0; i < epoch; i++)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_loadu_ps(b);
        __m512 v = _mm512_mul_ps(v1, v2);
        _mm512_storeu_ps(c, v);
        a += 16;
        b += 16;
        c += 16;
    }
    if (remain >= 8)
    {
        __m256 v1 = _mm256_loadu_ps(a);
        __m256 v2 = _mm256_loadu_ps(b);
        __m256 v = _mm256_mul_ps(v1, v2);
        _mm256_storeu_ps(c, v);
        a += 8;
        b += 8;
        c += 8;
        remain -= 8;
    }
    for (int i = 0; i < remain; i++)
    {
        c[i] = a[i] * b[i];
    }
}

void _mm512_sqrt(float *a, int64_t n)
{
    int epoch = n / 16;
    int remain = n % 16;
    for (int i = 0; i < epoch; i++)
    {
        __m512 v = _mm512_loadu_ps(a);
        v = _mm512_sqrt_ps(v);
        _mm512_storeu_ps(a, v);
        a += 16;
    }
    if (remain >= 8)
    {
        __m256 v = _mm256_loadu_ps(a);
        v = _mm256_sqrt_ps(v);
        _mm256_storeu_ps(a, v);
        a += 8;
        remain -= 8;
    }
    for (int i = 0; i < remain; i++)
    {
        a[i] = __builtin_sqrtf(a[i]);
    }
}

void _mm512_dot(float *a, float *b, int64_t n, float *ret)
{
    int epoch = n / 16;
    int remain = n % 16;
    __m512 s;
    if (epoch > 0)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_loadu_ps(b);
        s = _mm512_mul_ps(v1, v2);
        a += 16;
        b += 16;
    }
    for (int i = 1; i < epoch; i++)
    {
        __m512 v1 = _mm512_loadu_ps(a);
        __m512 v2 = _mm512_loadu_ps(b);
        s = _mm512_fmadd_ps(v1, v2, s);
        a += 16;
        b += 16;
    }
    __m256 sf_e_d_c_b_a_9_8 = _mm512_extractf32x8_ps(s, 1);
    __m256 s7_6_5_4_3_2_1_0 = _mm512_castps512_ps256(s);
    __m256 s7f_6e_5d_4c_3b_2a_19_08 = _mm256_add_ps(sf_e_d_c_b_a_9_8, s7_6_5_4_3_2_1_0);
    __m128 s7f_6e_5d_4c = _mm256_extractf128_ps(s7f_6e_5d_4c_3b_2a_19_08, 1);
    __m128 s3b_2a_19_08 = _mm256_castps256_ps128(s7f_6e_5d_4c_3b_2a_19_08);
    __m128 s37bf_26ae_159d_048c = _mm_add_ps(s7f_6e_5d_4c, s3b_2a_19_08);
    __m128 sxx_159d_048c = s37bf_26ae_159d_048c;
    __m128 sxx_37bf_26ae = _mm_movehl_ps(sxx_159d_048c, s37bf_26ae_159d_048c);
    const __m128 sxx_13579bdf_02468ace = _mm_add_ps(sxx_159d_048c, sxx_37bf_26ae);
    const __m128 sxxx_02468ace = sxx_13579bdf_02468ace;
    const __m128 sxxx_13579bdf = _mm_shuffle_ps(sxx_13579bdf_02468ace, sxx_13579bdf_02468ace, 0x1);
    __m128 sxxx_0123456789abcdef = _mm_add_ss(sxxx_02468ace, sxxx_13579bdf);
    *ret = _mm_cvtss_f32(sxxx_0123456789abcdef);

    if (remain >= 8)
    {
        __m256 s;
        __m256 v1 = _mm256_loadu_ps(a);
        __m256 v2 = _mm256_loadu_ps(b);
        s = _mm256_mul_ps(v1, v2);
        a += 8;
        b += 8;
        __m128 s7_6_5_4 = _mm256_extractf128_ps(s, 1);
        __m128 s3_2_1_0 = _mm256_castps256_ps128(s);
        __m128 s37_26_15_04 = _mm_add_ps(s7_6_5_4, s3_2_1_0);
        __m128 sxx_15_04 = s37_26_15_04;
        __m128 sxx_37_26 = _mm_movehl_ps(s37_26_15_04, s37_26_15_04);
        const __m128 sxx_1357_0246 = _mm_add_ps(sxx_15_04, sxx_37_26);
        const __m128 sxxx_0246 = sxx_1357_0246;
        const __m128 sxxx_1357 = _mm_shuffle_ps(sxx_1357_0246, sxx_1357_0246, 0x1);
        __m128 sxxx_01234567 = _mm_add_ss(sxxx_0246, sxxx_1357);
        *ret += _mm_cvtss_f32(sxxx_01234567);
        remain -= 8;
    }

    for (int i = 0; i < remain; i++)
    {
        *ret += a[i] * b[i];
    }
}
