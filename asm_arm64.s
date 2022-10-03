#include "textflag.h"

// func sum(s uint64, b []byte) uint64
TEXT ·sum(SB),NOSPLIT,$0
    MOVD s+0(FP), R1
    MOVD b_base+8(FP), R2
    MOVD b_len+16(FP), R3

    TST $1, R3 // & 0b 0000 0001
    BEQ loop2

    SUB $1, R3
    MOVBU (R2)(R3*1), R4

    ADDS R4, R1
    ADC $0, R1

loop2:
    TST $6, R3 // & 0b 0000 0110
    BEQ  loop8

    SUB $2, R3
    MOVHU (R2)(R3*1), R4
    ADDS R4, R1
    ADC $0, R1
    JMP  loop2

loop8:
    TST $56, R3 // & 0b 0011 1000
    BEQ  loop64

    SUB $8, R3
    MOVD (R2)(R3*1), R4
    ADDS R4, R1
    ADC $0, R1
    JMP  loop8

done:
    MOVD R1, ret+32(FP)
    RET

loop64:
    CMP $0, R3
    BEQ done
    SUB $64, R3
    ADD R2, R3, R5
    MOVD  0(R5), R8
    MOVD  8(R5), R9
    MOVD 16(R5), R10
    MOVD 24(R5), R11
    MOVD 32(R5), R12
    MOVD 40(R5), R13
    MOVD 48(R5), R14
    MOVD 56(R5), R15

    ADDS R8, R9
    ADC $0, R9
    ADDS R10, R11
    ADC $0, R11
    ADDS R12, R13
    ADC $0, R13
    ADDS R14, R15
    ADC $0, R15

    ADDS R9, R11
    ADC $0, R11
    ADDS R13, R15
    ADC $0, R15

    ADDS R11, R1
    ADC $0, R1
    ADDS R15, R1
    ADC $0, R1

    JMP loop64
