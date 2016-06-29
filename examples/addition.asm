	mov 	r15 main
add:    ; add taket two arguments
	add 	r0 r0 r1
	ret

main:   ; main must be called with r0 and r1 set
	call 	add

	stop
