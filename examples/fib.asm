; r0 = c
; r1 = next
; r2 = first
; r3 = second
; r4 = n

	mov	r4 #5 	; find number 5
	mov	r3 #1	; set r3 to 1
for_loop:
	lt 	r10 r0 r4
	jmpf 	r10 end
start_if:
	lteq 	r10 r0 #1
	jmpf 	r10 else

	mov 	r1 r0
	mov	r15 end_if
else:
	add 	r1 r2 r3
	mov 	r2 r3
	mov 	r3 r1
end_if:
	add 	r0 r0 #1
	mov 	r15 for_loop
end:
	mov 	r0 r1
