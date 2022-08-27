//go:build !noasm && amd64
// AUTO-GENERATED BY GOAT -- DO NOT EDIT

TEXT ·_mm256_mul_const_add_to(SB), $0-32
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ c+16(FP), DX
	MOVQ n+24(FP), CX
	BYTE $0x55               // pushq	%rbp
	WORD $0x8948; BYTE $0xe5 // movq	%rsp, %rbp
	WORD $0x5641             // pushq	%r14
	BYTE $0x53               // pushq	%rbx
	LONG $0xf8e48348         // andq	$-8, %rsp
	LONG $0x07418d48         // leaq	7(%rcx), %rax
	WORD $0x8548; BYTE $0xc9 // testq	%rcx, %rcx
	LONG $0xc1490f48         // cmovnsq	%rcx, %rax
	WORD $0x8949; BYTE $0xc0 // movq	%rax, %r8
	LONG $0x03f8c149         // sarq	$3, %r8
	LONG $0xf8e08348         // andq	$-8, %rax
	WORD $0x2948; BYTE $0xc1 // subq	%rax, %rcx
	WORD $0x8545; BYTE $0xc0 // testl	%r8d, %r8d
	JLE  LBB0_6
	LONG $0x01f88341         // cmpl	$1, %r8d
	JE   LBB0_4
	WORD $0x8944; BYTE $0xc0 // movl	%r8d, %eax
	WORD $0xe083; BYTE $0xfe // andl	$-2, %eax
	WORD $0xd8f7             // negl	%eax

LBB0_3:
	LONG $0x0710fcc5               // vmovups	(%rdi), %ymm0
	LONG $0x187de2c4; BYTE $0x0e   // vbroadcastss	(%rsi), %ymm1
	LONG $0xa87de2c4; BYTE $0x0a   // vfmadd213ps	(%rdx), %ymm0, %ymm1
	LONG $0x0a11fcc5               // vmovups	%ymm1, (%rdx)
	LONG $0x4710fcc5; BYTE $0x20   // vmovups	32(%rdi), %ymm0
	LONG $0x187de2c4; BYTE $0x0e   // vbroadcastss	(%rsi), %ymm1
	LONG $0xa87de2c4; WORD $0x204a // vfmadd213ps	32(%rdx), %ymm0, %ymm1
	LONG $0x4a11fcc5; BYTE $0x20   // vmovups	%ymm1, 32(%rdx)
	LONG $0x40c78348               // addq	$64, %rdi
	LONG $0x40c28348               // addq	$64, %rdx
	WORD $0xc083; BYTE $0x02       // addl	$2, %eax
	JNE  LBB0_3

LBB0_4:
	LONG $0x01c0f641             // testb	$1, %r8b
	JE   LBB0_6
	LONG $0x0710fcc5             // vmovups	(%rdi), %ymm0
	LONG $0x187de2c4; BYTE $0x0e // vbroadcastss	(%rsi), %ymm1
	LONG $0xa87de2c4; BYTE $0x0a // vfmadd213ps	(%rdx), %ymm0, %ymm1
	LONG $0x0a11fcc5             // vmovups	%ymm1, (%rdx)
	LONG $0x20c28348             // addq	$32, %rdx
	LONG $0x20c78348             // addq	$32, %rdi

LBB0_6:
	WORD $0xc985             // testl	%ecx, %ecx
	JLE  LBB0_18
	WORD $0x8941; BYTE $0xc8 // movl	%ecx, %r8d
	LONG $0x20f88349         // cmpq	$32, %r8
	JAE  LBB0_9
	WORD $0x3145; BYTE $0xd2 // xorl	%r10d, %r10d
	JMP  LBB0_14

LBB0_9:
	LONG $0x82048d4a             // leaq	(%rdx,%r8,4), %rax
	LONG $0x870c8d4e             // leaq	(%rdi,%r8,4), %r9
	LONG $0x04568d4c             // leaq	4(%rsi), %r10
	WORD $0x394c; BYTE $0xca     // cmpq	%r9, %rdx
	LONG $0xc6920f41             // setb	%r14b
	WORD $0x3948; BYTE $0xc7     // cmpq	%rax, %rdi
	WORD $0x920f; BYTE $0xc3     // setb	%bl
	WORD $0x394c; BYTE $0xd2     // cmpq	%r10, %rdx
	LONG $0xc3920f41             // setb	%r11b
	WORD $0x3948; BYTE $0xf0     // cmpq	%rsi, %rax
	LONG $0xc1970f41             // seta	%r9b
	WORD $0x3145; BYTE $0xd2     // xorl	%r10d, %r10d
	WORD $0x8441; BYTE $0xde     // testb	%bl, %r14b
	JNE  LBB0_14
	WORD $0x2045; BYTE $0xcb     // andb	%r9b, %r11b
	JNE  LBB0_14
	WORD $0x8941; BYTE $0xc9     // movl	%ecx, %r9d
	LONG $0x1fe18341             // andl	$31, %r9d
	WORD $0x894d; BYTE $0xc2     // movq	%r8, %r10
	WORD $0x294d; BYTE $0xca     // subq	%r9, %r10
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	WORD $0xc031                 // xorl	%eax, %eax

LBB0_12:
	LONG $0x0c59fcc5; BYTE $0x87   // vmulps	(%rdi,%rax,4), %ymm0, %ymm1
	LONG $0x5459fcc5; WORD $0x2087 // vmulps	32(%rdi,%rax,4), %ymm0, %ymm2
	LONG $0x5c59fcc5; WORD $0x4087 // vmulps	64(%rdi,%rax,4), %ymm0, %ymm3
	LONG $0x6459fcc5; WORD $0x6087 // vmulps	96(%rdi,%rax,4), %ymm0, %ymm4
	LONG $0x0c58f4c5; BYTE $0x82   // vaddps	(%rdx,%rax,4), %ymm1, %ymm1
	LONG $0x5458ecc5; WORD $0x2082 // vaddps	32(%rdx,%rax,4), %ymm2, %ymm2
	LONG $0x5c58e4c5; WORD $0x4082 // vaddps	64(%rdx,%rax,4), %ymm3, %ymm3
	LONG $0x6458dcc5; WORD $0x6082 // vaddps	96(%rdx,%rax,4), %ymm4, %ymm4
	LONG $0x0c11fcc5; BYTE $0x82   // vmovups	%ymm1, (%rdx,%rax,4)
	LONG $0x5411fcc5; WORD $0x2082 // vmovups	%ymm2, 32(%rdx,%rax,4)
	LONG $0x5c11fcc5; WORD $0x4082 // vmovups	%ymm3, 64(%rdx,%rax,4)
	LONG $0x6411fcc5; WORD $0x6082 // vmovups	%ymm4, 96(%rdx,%rax,4)
	LONG $0x20c08348               // addq	$32, %rax
	WORD $0x3949; BYTE $0xc2       // cmpq	%rax, %r10
	JNE  LBB0_12
	WORD $0x854d; BYTE $0xc9       // testq	%r9, %r9
	JE   LBB0_18

LBB0_14:
	WORD $0x2944; BYTE $0xd1       // subl	%r10d, %ecx
	LONG $0x01428d49               // leaq	1(%r10), %rax
	WORD $0xc1f6; BYTE $0x01       // testb	$1, %cl
	JE   LBB0_16
	LONG $0x107aa1c4; WORD $0x9704 // vmovss	(%rdi,%r10,4), %xmm0
	LONG $0x0659fac5               // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x587aa1c4; WORD $0x9204 // vaddss	(%rdx,%r10,4), %xmm0, %xmm0
	LONG $0x117aa1c4; WORD $0x9204 // vmovss	%xmm0, (%rdx,%r10,4)
	WORD $0x8949; BYTE $0xc2       // movq	%rax, %r10

LBB0_16:
	WORD $0x3949; BYTE $0xc0 // cmpq	%rax, %r8
	JE   LBB0_18

LBB0_17:
	LONG $0x107aa1c4; WORD $0x9704             // vmovss	(%rdi,%r10,4), %xmm0
	LONG $0x0659fac5                           // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x587aa1c4; WORD $0x9204             // vaddss	(%rdx,%r10,4), %xmm0, %xmm0
	LONG $0x117aa1c4; WORD $0x9204             // vmovss	%xmm0, (%rdx,%r10,4)
	LONG $0x107aa1c4; WORD $0x9744; BYTE $0x04 // vmovss	4(%rdi,%r10,4), %xmm0
	LONG $0x0659fac5                           // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x587aa1c4; WORD $0x9244; BYTE $0x04 // vaddss	4(%rdx,%r10,4), %xmm0, %xmm0
	LONG $0x117aa1c4; WORD $0x9244; BYTE $0x04 // vmovss	%xmm0, 4(%rdx,%r10,4)
	LONG $0x02c28349                           // addq	$2, %r10
	WORD $0x394d; BYTE $0xd0                   // cmpq	%r10, %r8
	JNE  LBB0_17

LBB0_18:
	LONG $0xf0658d48         // leaq	-16(%rbp), %rsp
	BYTE $0x5b               // popq	%rbx
	WORD $0x5e41             // popq	%r14
	BYTE $0x5d               // popq	%rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // retq

TEXT ·_mm256_mul_const_to(SB), $0-32
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ c+16(FP), DX
	MOVQ n+24(FP), CX
	BYTE $0x55               // pushq	%rbp
	WORD $0x8948; BYTE $0xe5 // movq	%rsp, %rbp
	BYTE $0x53               // pushq	%rbx
	LONG $0xf8e48348         // andq	$-8, %rsp
	LONG $0x07418d48         // leaq	7(%rcx), %rax
	WORD $0x8548; BYTE $0xc9 // testq	%rcx, %rcx
	LONG $0xc1490f48         // cmovnsq	%rcx, %rax
	WORD $0x8949; BYTE $0xc0 // movq	%rax, %r8
	LONG $0x03f8c149         // sarq	$3, %r8
	LONG $0xf8e08348         // andq	$-8, %rax
	WORD $0x2948; BYTE $0xc1 // subq	%rax, %rcx
	WORD $0x8545; BYTE $0xc0 // testl	%r8d, %r8d
	JLE  LBB1_6
	LONG $0x01f88341         // cmpl	$1, %r8d
	JE   LBB1_4
	WORD $0x8944; BYTE $0xc0 // movl	%r8d, %eax
	WORD $0xe083; BYTE $0xfe // andl	$-2, %eax
	WORD $0xd8f7             // negl	%eax

LBB1_3:
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x0759fcc5             // vmulps	(%rdi), %ymm0, %ymm0
	LONG $0x0211fcc5             // vmovups	%ymm0, (%rdx)
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x4759fcc5; BYTE $0x20 // vmulps	32(%rdi), %ymm0, %ymm0
	LONG $0x4211fcc5; BYTE $0x20 // vmovups	%ymm0, 32(%rdx)
	LONG $0x40c78348             // addq	$64, %rdi
	LONG $0x40c28348             // addq	$64, %rdx
	WORD $0xc083; BYTE $0x02     // addl	$2, %eax
	JNE  LBB1_3

LBB1_4:
	LONG $0x01c0f641             // testb	$1, %r8b
	JE   LBB1_6
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x0759fcc5             // vmulps	(%rdi), %ymm0, %ymm0
	LONG $0x0211fcc5             // vmovups	%ymm0, (%rdx)
	LONG $0x20c28348             // addq	$32, %rdx
	LONG $0x20c78348             // addq	$32, %rdi

LBB1_6:
	WORD $0xc985             // testl	%ecx, %ecx
	JLE  LBB1_18
	WORD $0x8941; BYTE $0xc8 // movl	%ecx, %r8d
	LONG $0x20f88349         // cmpq	$32, %r8
	JAE  LBB1_9
	WORD $0xc031             // xorl	%eax, %eax
	JMP  LBB1_14

LBB1_9:
	LONG $0x82048d4a             // leaq	(%rdx,%r8,4), %rax
	LONG $0x870c8d4e             // leaq	(%rdi,%r8,4), %r9
	LONG $0x04568d4c             // leaq	4(%rsi), %r10
	WORD $0x394c; BYTE $0xca     // cmpq	%r9, %rdx
	LONG $0xc3920f41             // setb	%r11b
	WORD $0x3948; BYTE $0xc7     // cmpq	%rax, %rdi
	WORD $0x920f; BYTE $0xc3     // setb	%bl
	WORD $0x394c; BYTE $0xd2     // cmpq	%r10, %rdx
	LONG $0xc1920f41             // setb	%r9b
	WORD $0x3948; BYTE $0xf0     // cmpq	%rsi, %rax
	LONG $0xc2970f41             // seta	%r10b
	WORD $0xc031                 // xorl	%eax, %eax
	WORD $0x8441; BYTE $0xdb     // testb	%bl, %r11b
	JNE  LBB1_14
	WORD $0x2045; BYTE $0xd1     // andb	%r10b, %r9b
	JNE  LBB1_14
	WORD $0x8941; BYTE $0xc9     // movl	%ecx, %r9d
	LONG $0x1fe18341             // andl	$31, %r9d
	WORD $0x894c; BYTE $0xc0     // movq	%r8, %rax
	WORD $0x294c; BYTE $0xc8     // subq	%r9, %rax
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	WORD $0x3145; BYTE $0xd2     // xorl	%r10d, %r10d

LBB1_12:
	LONG $0x597ca1c4; WORD $0x970c             // vmulps	(%rdi,%r10,4), %ymm0, %ymm1
	LONG $0x597ca1c4; WORD $0x9754; BYTE $0x20 // vmulps	32(%rdi,%r10,4), %ymm0, %ymm2
	LONG $0x597ca1c4; WORD $0x975c; BYTE $0x40 // vmulps	64(%rdi,%r10,4), %ymm0, %ymm3
	LONG $0x597ca1c4; WORD $0x9764; BYTE $0x60 // vmulps	96(%rdi,%r10,4), %ymm0, %ymm4
	LONG $0x117ca1c4; WORD $0x920c             // vmovups	%ymm1, (%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x9254; BYTE $0x20 // vmovups	%ymm2, 32(%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x925c; BYTE $0x40 // vmovups	%ymm3, 64(%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x9264; BYTE $0x60 // vmovups	%ymm4, 96(%rdx,%r10,4)
	LONG $0x20c28349                           // addq	$32, %r10
	WORD $0x394c; BYTE $0xd0                   // cmpq	%r10, %rax
	JNE  LBB1_12
	WORD $0x854d; BYTE $0xc9                   // testq	%r9, %r9
	JE   LBB1_18

LBB1_14:
	WORD $0xc129             // subl	%eax, %ecx
	WORD $0x8949; BYTE $0xc1 // movq	%rax, %r9
	WORD $0xf749; BYTE $0xd1 // notq	%r9
	WORD $0x014d; BYTE $0xc1 // addq	%r8, %r9
	LONG $0x03e18348         // andq	$3, %rcx
	JE   LBB1_16

LBB1_15:
	LONG $0x0410fac5; BYTE $0x87 // vmovss	(%rdi,%rax,4), %xmm0
	LONG $0x0659fac5             // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x82 // vmovss	%xmm0, (%rdx,%rax,4)
	LONG $0x01c08348             // addq	$1, %rax
	LONG $0xffc18348             // addq	$-1, %rcx
	JNE  LBB1_15

LBB1_16:
	LONG $0x03f98349 // cmpq	$3, %r9
	JB   LBB1_18

LBB1_17:
	LONG $0x0410fac5; BYTE $0x87   // vmovss	(%rdi,%rax,4), %xmm0
	LONG $0x0659fac5               // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x82   // vmovss	%xmm0, (%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0487 // vmovss	4(%rdi,%rax,4), %xmm0
	LONG $0x0659fac5               // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0482 // vmovss	%xmm0, 4(%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0887 // vmovss	8(%rdi,%rax,4), %xmm0
	LONG $0x0659fac5               // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0882 // vmovss	%xmm0, 8(%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0c87 // vmovss	12(%rdi,%rax,4), %xmm0
	LONG $0x0659fac5               // vmulss	(%rsi), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0c82 // vmovss	%xmm0, 12(%rdx,%rax,4)
	LONG $0x04c08348               // addq	$4, %rax
	WORD $0x3949; BYTE $0xc0       // cmpq	%rax, %r8
	JNE  LBB1_17

LBB1_18:
	LONG $0xf8658d48         // leaq	-8(%rbp), %rsp
	BYTE $0x5b               // popq	%rbx
	BYTE $0x5d               // popq	%rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // retq

TEXT ·_mm256_mul_const(SB), $0-32
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ n+16(FP), DX
	BYTE $0x55               // pushq	%rbp
	WORD $0x8948; BYTE $0xe5 // movq	%rsp, %rbp
	LONG $0xf8e48348         // andq	$-8, %rsp
	LONG $0x074a8d48         // leaq	7(%rdx), %rcx
	WORD $0x8548; BYTE $0xd2 // testq	%rdx, %rdx
	LONG $0xca490f48         // cmovnsq	%rdx, %rcx
	WORD $0x8948; BYTE $0xc8 // movq	%rcx, %rax
	LONG $0x03f8c148         // sarq	$3, %rax
	LONG $0xf8e18348         // andq	$-8, %rcx
	WORD $0x2948; BYTE $0xca // subq	%rcx, %rdx
	WORD $0xc085             // testl	%eax, %eax
	JLE  LBB2_6
	WORD $0xf883; BYTE $0x01 // cmpl	$1, %eax
	JE   LBB2_4
	WORD $0xc189             // movl	%eax, %ecx
	WORD $0xe183; BYTE $0xfe // andl	$-2, %ecx
	WORD $0xd9f7             // negl	%ecx

LBB2_3:
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x0759fcc5             // vmulps	(%rdi), %ymm0, %ymm0
	LONG $0x0711fcc5             // vmovups	%ymm0, (%rdi)
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x4759fcc5; BYTE $0x20 // vmulps	32(%rdi), %ymm0, %ymm0
	LONG $0x4711fcc5; BYTE $0x20 // vmovups	%ymm0, 32(%rdi)
	LONG $0x40c78348             // addq	$64, %rdi
	WORD $0xc183; BYTE $0x02     // addl	$2, %ecx
	JNE  LBB2_3

LBB2_4:
	WORD $0x01a8                 // testb	$1, %al
	JE   LBB2_6
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	LONG $0x0759fcc5             // vmulps	(%rdi), %ymm0, %ymm0
	LONG $0x0711fcc5             // vmovups	%ymm0, (%rdi)
	LONG $0x20c78348             // addq	$32, %rdi

LBB2_6:
	WORD $0xd285             // testl	%edx, %edx
	JLE  LBB2_19
	WORD $0x8941; BYTE $0xd1 // movl	%edx, %r9d
	LONG $0x20f98349         // cmpq	$32, %r9
	JB   LBB2_8
	LONG $0x04468d48         // leaq	4(%rsi), %rax
	WORD $0x3948; BYTE $0xc7 // cmpq	%rax, %rdi
	JAE  LBB2_12
	LONG $0x8f048d4a         // leaq	(%rdi,%r9,4), %rax
	WORD $0x3948; BYTE $0xf0 // cmpq	%rsi, %rax
	JBE  LBB2_12

LBB2_8:
	WORD $0xc931 // xorl	%ecx, %ecx

LBB2_15:
	WORD $0xca29             // subl	%ecx, %edx
	WORD $0x8949; BYTE $0xc8 // movq	%rcx, %r8
	WORD $0xf749; BYTE $0xd0 // notq	%r8
	WORD $0x014d; BYTE $0xc8 // addq	%r9, %r8
	LONG $0x03e28348         // andq	$3, %rdx
	JE   LBB2_17

LBB2_16:
	LONG $0x0610fac5             // vmovss	(%rsi), %xmm0
	LONG $0x0459fac5; BYTE $0x8f // vmulss	(%rdi,%rcx,4), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x8f // vmovss	%xmm0, (%rdi,%rcx,4)
	LONG $0x01c18348             // addq	$1, %rcx
	LONG $0xffc28348             // addq	$-1, %rdx
	JNE  LBB2_16

LBB2_17:
	LONG $0x03f88349 // cmpq	$3, %r8
	JB   LBB2_19

LBB2_18:
	LONG $0x0610fac5               // vmovss	(%rsi), %xmm0
	LONG $0x0459fac5; BYTE $0x8f   // vmulss	(%rdi,%rcx,4), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x8f   // vmovss	%xmm0, (%rdi,%rcx,4)
	LONG $0x0610fac5               // vmovss	(%rsi), %xmm0
	LONG $0x4459fac5; WORD $0x048f // vmulss	4(%rdi,%rcx,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x048f // vmovss	%xmm0, 4(%rdi,%rcx,4)
	LONG $0x0610fac5               // vmovss	(%rsi), %xmm0
	LONG $0x4459fac5; WORD $0x088f // vmulss	8(%rdi,%rcx,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x088f // vmovss	%xmm0, 8(%rdi,%rcx,4)
	LONG $0x0610fac5               // vmovss	(%rsi), %xmm0
	LONG $0x4459fac5; WORD $0x0c8f // vmulss	12(%rdi,%rcx,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0c8f // vmovss	%xmm0, 12(%rdi,%rcx,4)
	LONG $0x04c18348               // addq	$4, %rcx
	WORD $0x3949; BYTE $0xc9       // cmpq	%rcx, %r9
	JNE  LBB2_18
	JMP  LBB2_19

LBB2_12:
	WORD $0x8941; BYTE $0xd0     // movl	%edx, %r8d
	LONG $0x1fe08341             // andl	$31, %r8d
	WORD $0x894c; BYTE $0xc9     // movq	%r9, %rcx
	WORD $0x294c; BYTE $0xc1     // subq	%r8, %rcx
	LONG $0x187de2c4; BYTE $0x06 // vbroadcastss	(%rsi), %ymm0
	WORD $0xc031                 // xorl	%eax, %eax

LBB2_13:
	LONG $0x0c59fcc5; BYTE $0x87   // vmulps	(%rdi,%rax,4), %ymm0, %ymm1
	LONG $0x5459fcc5; WORD $0x2087 // vmulps	32(%rdi,%rax,4), %ymm0, %ymm2
	LONG $0x5c59fcc5; WORD $0x4087 // vmulps	64(%rdi,%rax,4), %ymm0, %ymm3
	LONG $0x6459fcc5; WORD $0x6087 // vmulps	96(%rdi,%rax,4), %ymm0, %ymm4
	LONG $0x0c11fcc5; BYTE $0x87   // vmovups	%ymm1, (%rdi,%rax,4)
	LONG $0x5411fcc5; WORD $0x2087 // vmovups	%ymm2, 32(%rdi,%rax,4)
	LONG $0x5c11fcc5; WORD $0x4087 // vmovups	%ymm3, 64(%rdi,%rax,4)
	LONG $0x6411fcc5; WORD $0x6087 // vmovups	%ymm4, 96(%rdi,%rax,4)
	LONG $0x20c08348               // addq	$32, %rax
	WORD $0x3948; BYTE $0xc1       // cmpq	%rax, %rcx
	JNE  LBB2_13
	WORD $0x854d; BYTE $0xc0       // testq	%r8, %r8
	JNE  LBB2_15

LBB2_19:
	WORD $0x8948; BYTE $0xec // movq	%rbp, %rsp
	BYTE $0x5d               // popq	%rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // retq

TEXT ·_mm256_mul_to(SB), $0-32
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ c+16(FP), DX
	MOVQ n+24(FP), CX
	BYTE $0x55               // pushq	%rbp
	WORD $0x8948; BYTE $0xe5 // movq	%rsp, %rbp
	BYTE $0x53               // pushq	%rbx
	LONG $0xf8e48348         // andq	$-8, %rsp
	LONG $0x07418d4c         // leaq	7(%rcx), %r8
	WORD $0x8548; BYTE $0xc9 // testq	%rcx, %rcx
	LONG $0xc1490f4c         // cmovnsq	%rcx, %r8
	WORD $0x894c; BYTE $0xc0 // movq	%r8, %rax
	LONG $0x03f8c148         // sarq	$3, %rax
	LONG $0xf8e08349         // andq	$-8, %r8
	WORD $0x294c; BYTE $0xc1 // subq	%r8, %rcx
	WORD $0xc085             // testl	%eax, %eax
	JLE  LBB3_6
	LONG $0xff488d44         // leal	-1(%rax), %r9d
	WORD $0x8941; BYTE $0xc0 // movl	%eax, %r8d
	LONG $0x03e08341         // andl	$3, %r8d
	LONG $0x03f98341         // cmpl	$3, %r9d
	JB   LBB3_4
	WORD $0xe083; BYTE $0xfc // andl	$-4, %eax
	WORD $0xd8f7             // negl	%eax

LBB3_3:
	LONG $0x0710fcc5             // vmovups	(%rdi), %ymm0
	LONG $0x0659fcc5             // vmulps	(%rsi), %ymm0, %ymm0
	LONG $0x0211fcc5             // vmovups	%ymm0, (%rdx)
	LONG $0x4710fcc5; BYTE $0x20 // vmovups	32(%rdi), %ymm0
	LONG $0x4659fcc5; BYTE $0x20 // vmulps	32(%rsi), %ymm0, %ymm0
	LONG $0x4211fcc5; BYTE $0x20 // vmovups	%ymm0, 32(%rdx)
	LONG $0x4710fcc5; BYTE $0x40 // vmovups	64(%rdi), %ymm0
	LONG $0x4659fcc5; BYTE $0x40 // vmulps	64(%rsi), %ymm0, %ymm0
	LONG $0x4211fcc5; BYTE $0x40 // vmovups	%ymm0, 64(%rdx)
	LONG $0x4710fcc5; BYTE $0x60 // vmovups	96(%rdi), %ymm0
	LONG $0x4659fcc5; BYTE $0x60 // vmulps	96(%rsi), %ymm0, %ymm0
	LONG $0x4211fcc5; BYTE $0x60 // vmovups	%ymm0, 96(%rdx)
	LONG $0x80ef8348             // subq	$-128, %rdi
	LONG $0x80ee8348             // subq	$-128, %rsi
	LONG $0x80ea8348             // subq	$-128, %rdx
	WORD $0xc083; BYTE $0x04     // addl	$4, %eax
	JNE  LBB3_3

LBB3_4:
	WORD $0x8545; BYTE $0xc0 // testl	%r8d, %r8d
	JE   LBB3_6

LBB3_5:
	LONG $0x0710fcc5 // vmovups	(%rdi), %ymm0
	LONG $0x0659fcc5 // vmulps	(%rsi), %ymm0, %ymm0
	LONG $0x0211fcc5 // vmovups	%ymm0, (%rdx)
	LONG $0x20c78348 // addq	$32, %rdi
	LONG $0x20c68348 // addq	$32, %rsi
	LONG $0x20c28348 // addq	$32, %rdx
	LONG $0xffc08341 // addl	$-1, %r8d
	JNE  LBB3_5

LBB3_6:
	WORD $0xc985             // testl	%ecx, %ecx
	JLE  LBB3_18
	WORD $0x8941; BYTE $0xc8 // movl	%ecx, %r8d
	LONG $0x20f88349         // cmpq	$32, %r8
	JAE  LBB3_9
	WORD $0xc031             // xorl	%eax, %eax
	JMP  LBB3_14

LBB3_9:
	LONG $0x82048d4a         // leaq	(%rdx,%r8,4), %rax
	LONG $0x870c8d4e         // leaq	(%rdi,%r8,4), %r9
	LONG $0x86148d4e         // leaq	(%rsi,%r8,4), %r10
	WORD $0x394c; BYTE $0xca // cmpq	%r9, %rdx
	LONG $0xc3920f41         // setb	%r11b
	WORD $0x3948; BYTE $0xc7 // cmpq	%rax, %rdi
	WORD $0x920f; BYTE $0xc3 // setb	%bl
	WORD $0x394c; BYTE $0xd2 // cmpq	%r10, %rdx
	LONG $0xc1920f41         // setb	%r9b
	WORD $0x3948; BYTE $0xc6 // cmpq	%rax, %rsi
	LONG $0xc2920f41         // setb	%r10b
	WORD $0xc031             // xorl	%eax, %eax
	WORD $0x8441; BYTE $0xdb // testb	%bl, %r11b
	JNE  LBB3_14
	WORD $0x2045; BYTE $0xd1 // andb	%r10b, %r9b
	JNE  LBB3_14
	WORD $0x8941; BYTE $0xc9 // movl	%ecx, %r9d
	LONG $0x1fe18341         // andl	$31, %r9d
	WORD $0x894c; BYTE $0xc0 // movq	%r8, %rax
	WORD $0x294c; BYTE $0xc8 // subq	%r9, %rax
	WORD $0x3145; BYTE $0xd2 // xorl	%r10d, %r10d

LBB3_12:
	LONG $0x107ca1c4; WORD $0x9704             // vmovups	(%rdi,%r10,4), %ymm0
	LONG $0x107ca1c4; WORD $0x974c; BYTE $0x20 // vmovups	32(%rdi,%r10,4), %ymm1
	LONG $0x107ca1c4; WORD $0x9754; BYTE $0x40 // vmovups	64(%rdi,%r10,4), %ymm2
	LONG $0x107ca1c4; WORD $0x975c; BYTE $0x60 // vmovups	96(%rdi,%r10,4), %ymm3
	LONG $0x597ca1c4; WORD $0x9604             // vmulps	(%rsi,%r10,4), %ymm0, %ymm0
	LONG $0x5974a1c4; WORD $0x964c; BYTE $0x20 // vmulps	32(%rsi,%r10,4), %ymm1, %ymm1
	LONG $0x596ca1c4; WORD $0x9654; BYTE $0x40 // vmulps	64(%rsi,%r10,4), %ymm2, %ymm2
	LONG $0x5964a1c4; WORD $0x965c; BYTE $0x60 // vmulps	96(%rsi,%r10,4), %ymm3, %ymm3
	LONG $0x117ca1c4; WORD $0x9204             // vmovups	%ymm0, (%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x924c; BYTE $0x20 // vmovups	%ymm1, 32(%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x9254; BYTE $0x40 // vmovups	%ymm2, 64(%rdx,%r10,4)
	LONG $0x117ca1c4; WORD $0x925c; BYTE $0x60 // vmovups	%ymm3, 96(%rdx,%r10,4)
	LONG $0x20c28349                           // addq	$32, %r10
	WORD $0x394c; BYTE $0xd0                   // cmpq	%r10, %rax
	JNE  LBB3_12
	WORD $0x854d; BYTE $0xc9                   // testq	%r9, %r9
	JE   LBB3_18

LBB3_14:
	WORD $0xc129             // subl	%eax, %ecx
	WORD $0x8949; BYTE $0xc1 // movq	%rax, %r9
	WORD $0xf749; BYTE $0xd1 // notq	%r9
	WORD $0x014d; BYTE $0xc1 // addq	%r8, %r9
	LONG $0x03e18348         // andq	$3, %rcx
	JE   LBB3_16

LBB3_15:
	LONG $0x0410fac5; BYTE $0x87 // vmovss	(%rdi,%rax,4), %xmm0
	LONG $0x0459fac5; BYTE $0x86 // vmulss	(%rsi,%rax,4), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x82 // vmovss	%xmm0, (%rdx,%rax,4)
	LONG $0x01c08348             // addq	$1, %rax
	LONG $0xffc18348             // addq	$-1, %rcx
	JNE  LBB3_15

LBB3_16:
	LONG $0x03f98349 // cmpq	$3, %r9
	JB   LBB3_18

LBB3_17:
	LONG $0x0410fac5; BYTE $0x87   // vmovss	(%rdi,%rax,4), %xmm0
	LONG $0x0459fac5; BYTE $0x86   // vmulss	(%rsi,%rax,4), %xmm0, %xmm0
	LONG $0x0411fac5; BYTE $0x82   // vmovss	%xmm0, (%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0487 // vmovss	4(%rdi,%rax,4), %xmm0
	LONG $0x4459fac5; WORD $0x0486 // vmulss	4(%rsi,%rax,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0482 // vmovss	%xmm0, 4(%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0887 // vmovss	8(%rdi,%rax,4), %xmm0
	LONG $0x4459fac5; WORD $0x0886 // vmulss	8(%rsi,%rax,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0882 // vmovss	%xmm0, 8(%rdx,%rax,4)
	LONG $0x4410fac5; WORD $0x0c87 // vmovss	12(%rdi,%rax,4), %xmm0
	LONG $0x4459fac5; WORD $0x0c86 // vmulss	12(%rsi,%rax,4), %xmm0, %xmm0
	LONG $0x4411fac5; WORD $0x0c82 // vmovss	%xmm0, 12(%rdx,%rax,4)
	LONG $0x04c08348               // addq	$4, %rax
	WORD $0x3949; BYTE $0xc0       // cmpq	%rax, %r8
	JNE  LBB3_17

LBB3_18:
	LONG $0xf8658d48         // leaq	-8(%rbp), %rsp
	BYTE $0x5b               // popq	%rbx
	BYTE $0x5d               // popq	%rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // retq

TEXT ·_mm256_dot(SB), $0-32
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ n+16(FP), DX
	MOVQ ret+24(FP), CX
	BYTE $0x55                             // pushq	%rbp
	WORD $0x8948; BYTE $0xe5               // movq	%rsp, %rbp
	WORD $0x5641                           // pushq	%r14
	BYTE $0x53                             // pushq	%rbx
	LONG $0xf8e48348                       // andq	$-8, %rsp
	LONG $0x07428d48                       // leaq	7(%rdx), %rax
	WORD $0x8548; BYTE $0xd2               // testq	%rdx, %rdx
	LONG $0xc2490f48                       // cmovnsq	%rdx, %rax
	WORD $0x8949; BYTE $0xc1               // movq	%rax, %r9
	LONG $0x03f9c149                       // sarq	$3, %r9
	LONG $0xf8e08348                       // andq	$-8, %rax
	WORD $0x2948; BYTE $0xc2               // subq	%rax, %rdx
	WORD $0x8545; BYTE $0xc9               // testl	%r9d, %r9d
	JLE  LBB4_1
	LONG $0x0710fcc5                       // vmovups	(%rdi), %ymm0
	LONG $0x0659fcc5                       // vmulps	(%rsi), %ymm0, %ymm0
	LONG $0x20c78348                       // addq	$32, %rdi
	LONG $0x20c68348                       // addq	$32, %rsi
	LONG $0x01f98341                       // cmpl	$1, %r9d
	JE   LBB4_9
	QUAD $0x0007fffffff0b849; WORD $0x0000 // movabsq	$34359738352, %r8
	LONG $0xc8048d4b                       // leaq	(%r8,%r9,8), %rax
	LONG $0x08c88349                       // orq	$8, %r8
	WORD $0x2149; BYTE $0xc0               // andq	%rax, %r8
	LONG $0xff598d45                       // leal	-1(%r9), %r11d
	LONG $0xfe418d41                       // leal	-2(%r9), %eax
	WORD $0xf883; BYTE $0x03               // cmpl	$3, %eax
	JAE  LBB4_16
	WORD $0x8949; BYTE $0xfa               // movq	%rdi, %r10
	WORD $0x8948; BYTE $0xf0               // movq	%rsi, %rax
	JMP  LBB4_5

LBB4_1:
	JMP LBB4_9

LBB4_16:
	WORD $0x8944; BYTE $0xdb // movl	%r11d, %ebx
	WORD $0xe383; BYTE $0xfc // andl	$-4, %ebx
	WORD $0xdbf7             // negl	%ebx
	WORD $0x8949; BYTE $0xfa // movq	%rdi, %r10
	WORD $0x8948; BYTE $0xf0 // movq	%rsi, %rax

LBB4_17:
	LONG $0x107cc1c4; BYTE $0x0a   // vmovups	(%r10), %ymm1
	LONG $0x107cc1c4; WORD $0x2052 // vmovups	32(%r10), %ymm2
	LONG $0x107cc1c4; WORD $0x405a // vmovups	64(%r10), %ymm3
	LONG $0x107cc1c4; WORD $0x6062 // vmovups	96(%r10), %ymm4
	LONG $0x987de2c4; BYTE $0x08   // vfmadd132ps	(%rax), %ymm0, %ymm1
	LONG $0xb86de2c4; WORD $0x2048 // vfmadd231ps	32(%rax), %ymm2, %ymm1
	LONG $0xb865e2c4; WORD $0x4048 // vfmadd231ps	64(%rax), %ymm3, %ymm1
	LONG $0xc128fcc5               // vmovaps	%ymm1, %ymm0
	LONG $0xb85de2c4; WORD $0x6040 // vfmadd231ps	96(%rax), %ymm4, %ymm0
	LONG $0x80ea8349               // subq	$-128, %r10
	LONG $0x80e88348               // subq	$-128, %rax
	WORD $0xc383; BYTE $0x04       // addl	$4, %ebx
	JNE  LBB4_17

LBB4_5:
	LONG $0x08708d4d // leaq	8(%r8), %r14
	LONG $0x03c3f641 // testb	$3, %r11b
	JE   LBB4_8
	LONG $0xffc18041 // addb	$-1, %r9b
	LONG $0xc9b60f45 // movzbl	%r9b, %r9d
	LONG $0x03e18341 // andl	$3, %r9d
	LONG $0x05e1c149 // shlq	$5, %r9
	WORD $0xdb31     // xorl	%ebx, %ebx

LBB4_7:
	LONG $0x107cc1c4; WORD $0x1a0c // vmovups	(%r10,%rbx), %ymm1
	LONG $0xb875e2c4; WORD $0x1804 // vfmadd231ps	(%rax,%rbx), %ymm1, %ymm0
	LONG $0x20c38348               // addq	$32, %rbx
	WORD $0x3941; BYTE $0xd9       // cmpl	%ebx, %r9d
	JNE  LBB4_7

LBB4_8:
	LONG $0x873c8d4a // leaq	(%rdi,%r8,4), %rdi
	LONG $0x20c78348 // addq	$32, %rdi
	LONG $0xb6348d4a // leaq	(%rsi,%r14,4), %rsi

LBB4_9:
	LONG $0x197de3c4; WORD $0x01c1 // vextractf128	$1, %ymm0, %xmm1
	LONG $0xc058f0c5               // vaddps	%xmm0, %xmm1, %xmm0
	LONG $0x0579e3c4; WORD $0x01c8 // vpermilpd	$1, %xmm0, %xmm1
	LONG $0xc158f8c5               // vaddps	%xmm1, %xmm0, %xmm0
	LONG $0xc816fac5               // vmovshdup	%xmm0, %xmm1
	LONG $0xc158fac5               // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5               // vmovss	%xmm0, (%rcx)
	WORD $0xd285                   // testl	%edx, %edx
	JLE  LBB4_15
	WORD $0x8941; BYTE $0xd0       // movl	%edx, %r8d
	LONG $0xff408d49               // leaq	-1(%r8), %rax
	WORD $0xe283; BYTE $0x03       // andl	$3, %edx
	LONG $0x03f88348               // cmpq	$3, %rax
	JAE  LBB4_18
	WORD $0xc031                   // xorl	%eax, %eax
	JMP  LBB4_12

LBB4_18:
	WORD $0x2949; BYTE $0xd0 // subq	%rdx, %r8
	WORD $0xc031             // xorl	%eax, %eax

LBB4_19:
	LONG $0x0c10fac5; BYTE $0x87   // vmovss	(%rdi,%rax,4), %xmm1
	LONG $0x0c59f2c5; BYTE $0x86   // vmulss	(%rsi,%rax,4), %xmm1, %xmm1
	LONG $0xc158fac5               // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5               // vmovss	%xmm0, (%rcx)
	LONG $0x4c10fac5; WORD $0x0487 // vmovss	4(%rdi,%rax,4), %xmm1
	LONG $0x4c59f2c5; WORD $0x0486 // vmulss	4(%rsi,%rax,4), %xmm1, %xmm1
	LONG $0xc158fac5               // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5               // vmovss	%xmm0, (%rcx)
	LONG $0x4c10fac5; WORD $0x0887 // vmovss	8(%rdi,%rax,4), %xmm1
	LONG $0x4c59f2c5; WORD $0x0886 // vmulss	8(%rsi,%rax,4), %xmm1, %xmm1
	LONG $0xc158fac5               // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5               // vmovss	%xmm0, (%rcx)
	LONG $0x4c10fac5; WORD $0x0c87 // vmovss	12(%rdi,%rax,4), %xmm1
	LONG $0x4c59f2c5; WORD $0x0c86 // vmulss	12(%rsi,%rax,4), %xmm1, %xmm1
	LONG $0xc158fac5               // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5               // vmovss	%xmm0, (%rcx)
	LONG $0x04c08348               // addq	$4, %rax
	WORD $0x3949; BYTE $0xc0       // cmpq	%rax, %r8
	JNE  LBB4_19

LBB4_12:
	WORD $0x8548; BYTE $0xd2 // testq	%rdx, %rdx
	JE   LBB4_15
	LONG $0x86348d48         // leaq	(%rsi,%rax,4), %rsi
	LONG $0x87048d48         // leaq	(%rdi,%rax,4), %rax
	WORD $0xff31             // xorl	%edi, %edi

LBB4_14:
	LONG $0x0c10fac5; BYTE $0xb8 // vmovss	(%rax,%rdi,4), %xmm1
	LONG $0x0c59f2c5; BYTE $0xbe // vmulss	(%rsi,%rdi,4), %xmm1, %xmm1
	LONG $0xc158fac5             // vaddss	%xmm1, %xmm0, %xmm0
	LONG $0x0111fac5             // vmovss	%xmm0, (%rcx)
	LONG $0x01c78348             // addq	$1, %rdi
	WORD $0x3948; BYTE $0xfa     // cmpq	%rdi, %rdx
	JNE  LBB4_14

LBB4_15:
	LONG $0xf0658d48         // leaq	-16(%rbp), %rsp
	BYTE $0x5b               // popq	%rbx
	WORD $0x5e41             // popq	%r14
	BYTE $0x5d               // popq	%rbp
	WORD $0xf8c5; BYTE $0x77 // vzeroupper
	BYTE $0xc3               // retq
