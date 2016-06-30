; r0 = c
; r1 = next
; r2 = first
; r3 = second
; r4 = n
	mov 	r4 	#5		; n = 5
	mov 	r3 	#1		; first = 1

for:
	cmp 	r0 	r4		; if c < n
	movgt   r15 	else		;	next = c

	mov 	r1 	r0
	mov 	r15 	end_if
else					; else:
	add 	r1 	r2 	r3	;	next = first + second
	mov	r2 	r3		; 	first = second
	mov 	r3 	r1		; 	second = next
end_if:
	add	r0 	r0 	#1
	cmp 	r0 	r4		; for c < n; c++
	movgt 	r15 	end

	mov 	r15 	for
end:
	mov 	r0 	r1
