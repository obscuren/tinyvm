        jmp 	main
add:    ; add taket two arguments
	add 	r0 r0 r1
	ret

main:   ; main must be called with r0 and r1 set
	mov 	r0 3
	mov 	r1 2
	call 	add

	stop
