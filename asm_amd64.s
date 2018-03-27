#include "textflag.h"

// func sum(s uint, b []byte) uint
TEXT ·sum(SB),NOSPLIT,$0
    MOVQ s+0(FP), AX
    MOVQ b_base+8(FP), BX
    MOVQ b_len+16(FP), CX

    TESTQ $1, CX // & 0b 0000 0001
    JZ   loop2

    DECQ CX
    MOVBQZX (BX)(CX*1), SI
    ADDQ SI, AX
    ADCQ $0, AX

loop2:
    TESTQ $6, CX // & 0b 0000 0110
    JZ   loop8

    SUBQ $2, CX
    MOVWQZX (BX)(CX*1), SI
    ADDQ SI, AX
    ADCQ $0, AX
    JMP  loop2

loop8:
    TESTQ $56, CX // & 0b 0011 1000
    JZ   loop64

    SUBQ $8, CX
    MOVQ (BX)(CX*1), SI
    ADDQ SI, AX
    ADCQ $0, AX
    JMP  loop8

done:
    MOVQ AX, r+32(FP)
    RET

loop64:
    JCXZQ done
    SUBQ $64, CX
    LEAQ 0(BX)(CX*1), DX
    MOVQ  0(DX), R8
    MOVQ  8(DX), R9
    MOVQ 16(DX), R10
    MOVQ 24(DX), R11
    MOVQ 32(DX), R12
    MOVQ 40(DX), R13
    MOVQ 48(DX), R14
    MOVQ 56(DX), R15

    ADDQ R8, R9
    ADCQ $0, R9
    ADDQ R10, R11
    ADCQ $0, R11
    ADDQ R12, R13
    ADCQ $0, R13
    ADDQ R14, R15
    ADCQ $0, R15

    ADDQ R9, R11
    ADCQ $0, R11
    ADDQ R13, R15
    ADCQ $0, R15

    ADDQ R11, AX
    ADCQ $0, AX
    ADDQ R15, AX
    ADCQ $0, AX

    JMP loop64
