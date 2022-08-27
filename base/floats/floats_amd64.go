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

package floats

import (
	"github.com/klauspost/cpuid/v2"
	"unsafe"
)

//go:generate go run ../../cmd/goat src/floats_avx.c -O3 -mavx -mfma
//go:generate go run ../../cmd/goat src/floats_avx512.c -O3 -mavx -mfma -mavx512f -mavx512dq

func init() {
	if cpuid.CPU.Supports(cpuid.AVX512F) && cpuid.CPU.Supports(cpuid.AVX512DQ) {
		impl = avx{}
	} else if cpuid.CPU.Supports(cpuid.AVX) && cpuid.CPU.Supports(cpuid.FMA3) {
		impl = avx{}
	}
}

type avx struct{}

func (avx) MulConstAddTo(a []float32, b float32, c []float32) {
	_mm256_mul_const_add_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx) MulConstTo(a []float32, b float32, c []float32) {
	_mm256_mul_const_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx) MulTo(a, b, c []float32) {
	_mm256_mul_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx) MulConst(a []float32, b float32) {
	_mm256_mul_const(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(uintptr(len(a))))
}

func (avx) Dot(a, b []float32) float32 {
	var ret float32
	_mm256_dot(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), unsafe.Pointer(uintptr(len(a))), unsafe.Pointer(&ret))
	return ret
}

type avx512 struct{}

func (avx512) MulConstAddTo(a []float32, b float32, c []float32) {
	_mm512_mul_const_add_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx512) MulConstTo(a []float32, b float32, c []float32) {
	_mm512_mul_const_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx512) MulTo(a, b, c []float32) {
	_mm512_mul_to(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), unsafe.Pointer(&c[0]), unsafe.Pointer(uintptr(len(a))))
}

func (avx512) MulConst(a []float32, b float32) {
	_mm512_mul_const(unsafe.Pointer(&a[0]), unsafe.Pointer(&b), unsafe.Pointer(uintptr(len(a))))
}

func (avx512) Dot(a, b []float32) float32 {
	var ret float32
	_mm512_dot(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), unsafe.Pointer(uintptr(len(a))), unsafe.Pointer(&ret))
	return ret
}
