mov 	r1 #260
mov 	r2 r1
mov     r3 r1
add     r0 r1 #2
sub     r0 r0 #6

mov    r0 #1
mov    r1 #1
cmp    r0 r1
moveq  r0 #5

mov    r2 #1
mov    r3 #2
cmp    r2 r3
moveq  r4 #1

mov    r6 #1
subs   r6 r6 #1 ; if a == 0
moveq  r6 #2    ;    a = 2
