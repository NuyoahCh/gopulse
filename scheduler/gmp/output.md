go run main.go

[Scheduler] 提交 G1(cpu-bound-1) 到 P2 本地队列
[Scheduler] 提交 G2(cpu-bound-2) 到 P2 本地队列
[Scheduler] 提交 G3(syscall-1) 到 P2 本地队列
[Scheduler] 提交 G4(syscall-2) 到 P2 本地队列
[Scheduler] 提交 G5(mixed-1) 到 P2 本地队列
[Scheduler] 提交 G6(mixed-2) 到 P1 本地队列

M2 启动，绑定到 P2
M2 在 P2 上开始调度 G1(cpu-bound-1)
G1(cpu-bound-1) 在 M2/P2 上执行第 1/10 步
M1 启动，绑定到 P1
M1 在 P1 上开始调度 G6(mixed-2)
G6(mixed-2) 在 M1/P1 上执行第 1/11 步
G1(cpu-bound-1) 在 M2/P2 上执行第 2/10 步
G6(mixed-2) 在 M1/P1 上执行第 2/11 步
G1(cpu-bound-1) 在 M2/P2 上执行第 3/10 步
G6(mixed-2) 在 M1/P1 上执行第 3/11 步
M1：G6(mixed-2) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G6(mixed-2)
G6(mixed-2) 在 M1/P1 上执行第 4/11 步
M2：G1(cpu-bound-1) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G2(cpu-bound-2)
G2(cpu-bound-2) 在 M2/P2 上执行第 1/8 步
G6(mixed-2) 在 M1/P1 上执行第 5/11 步
G2(cpu-bound-2) 在 M2/P2 上执行第 2/8 步
G6(mixed-2) 在第 5 步执行阻塞系统调用，阻塞 500ms
M1：G6(mixed-2) 在 P1 上阻塞，切换到其他 G

[Steal] P1 从 P2 偷到了 G3(syscall-1)
M1 在 P1 上开始调度 G3(syscall-1)
G3(syscall-1) 在 M1/P1 上执行第 1/12 步
G2(cpu-bound-2) 在 M2/P2 上执行第 3/8 步
G3(syscall-1) 在 M1/P1 上执行第 2/12 步
M2：G2(cpu-bound-2) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G4(syscall-2)
G4(syscall-2) 在 M2/P2 上执行第 1/9 步
G4(syscall-2) 在 M2/P2 上执行第 2/9 步
G3(syscall-1) 在 M1/P1 上执行第 3/12 步
G4(syscall-2) 在 M2/P2 上执行第 3/9 步
M1：G3(syscall-1) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G3(syscall-1)
G3(syscall-1) 在 M1/P1 上执行第 4/12 步
G4(syscall-2) 在第 3 步执行阻塞系统调用，阻塞 600ms
M2：G4(syscall-2) 在 P2 上阻塞，切换到其他 G
M2 在 P2 上开始调度 G5(mixed-1)
G5(mixed-1) 在 M2/P2 上执行第 1/7 步
G5(mixed-1) 在 M2/P2 上执行第 2/7 步
G3(syscall-1) 在第 4 步执行阻塞系统调用，阻塞 600ms
M1：G3(syscall-1) 在 P1 上阻塞，切换到其他 G

[Steal] P1 从 P2 偷到了 G1(cpu-bound-1)
M1 在 P1 上开始调度 G1(cpu-bound-1)
G1(cpu-bound-1) 在 M1/P1 上执行第 4/10 步
G5(mixed-1) 在第 2 步执行阻塞系统调用，阻塞 400ms
M2：G5(mixed-1) 在 P2 上阻塞，切换到其他 G
M2 在 P2 上开始调度 G2(cpu-bound-2)
G2(cpu-bound-2) 在 M2/P2 上执行第 4/8 步
G1(cpu-bound-1) 在 M1/P1 上执行第 5/10 步
G2(cpu-bound-2) 在 M2/P2 上执行第 5/8 步
G1(cpu-bound-1) 在 M1/P1 上执行第 6/10 步
M1：G1(cpu-bound-1) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G1(cpu-bound-1)
G1(cpu-bound-1) 在 M1/P1 上执行第 7/10 步
G2(cpu-bound-2) 在 M2/P2 上执行第 6/8 步
G1(cpu-bound-1) 在 M1/P1 上执行第 8/10 步
M2：G2(cpu-bound-2) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G2(cpu-bound-2)
G2(cpu-bound-2) 在 M2/P2 上执行第 7/8 步

[Timer] G6(mixed-2) 阻塞结束，重新变为 runnable，重新入队
[Scheduler] 提交 G6(mixed-2) 到 P2 本地队列
G1(cpu-bound-1) 在 M1/P1 上执行第 9/10 步
G2(cpu-bound-2) 在 M2/P2 上执行第 8/8 步
M2：G2(cpu-bound-2) 在 P2 上全部完成
M2 在 P2 上开始调度 G6(mixed-2)
G6(mixed-2) 在 M2/P2 上执行第 6/11 步
M1：G1(cpu-bound-1) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G1(cpu-bound-1)
G1(cpu-bound-1) 在 M1/P1 上执行第 10/10 步
M1：G1(cpu-bound-1) 在 P1 上全部完成
G6(mixed-2) 在 M2/P2 上执行第 7/11 步
G6(mixed-2) 在 M2/P2 上执行第 8/11 步

[Timer] G5(mixed-1) 阻塞结束，重新变为 runnable，重新入队
[Scheduler] 提交 G5(mixed-1) 到 P2 本地队列

[Steal] P1 从 P2 偷到了 G5(mixed-1)
M1 在 P1 上开始调度 G5(mixed-1)
G5(mixed-1) 在 M1/P1 上执行第 3/7 步
M2：G6(mixed-2) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G6(mixed-2)
G6(mixed-2) 在 M2/P2 上执行第 9/11 步
G5(mixed-1) 在 M1/P1 上执行第 4/7 步
G5(mixed-1) 在 M1/P1 上执行第 5/7 步
G6(mixed-2) 在 M2/P2 上执行第 10/11 步

[Timer] G4(syscall-2) 阻塞结束，重新变为 runnable，重新入队
[Scheduler] 提交 G4(syscall-2) 到 P1 本地队列
G6(mixed-2) 在 M2/P2 上执行第 11/11 步
M1：G5(mixed-1) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G4(syscall-2)
G4(syscall-2) 在 M1/P1 上执行第 4/9 步

[Timer] G3(syscall-1) 阻塞结束，重新变为 runnable，重新入队
[Scheduler] 提交 G3(syscall-1) 到 P1 本地队列
M2：G6(mixed-2) 在 P2 上全部完成

[Steal] P2 从 P1 偷到了 G5(mixed-1)
M2 在 P2 上开始调度 G5(mixed-1)
G5(mixed-1) 在 M2/P2 上执行第 6/7 步
G4(syscall-2) 在 M1/P1 上执行第 5/9 步
G5(mixed-1) 在 M2/P2 上执行第 7/7 步
G4(syscall-2) 在 M1/P1 上执行第 6/9 步
M2：G5(mixed-1) 在 P2 上全部完成

[Steal] P2 从 P1 偷到了 G3(syscall-1)
M2 在 P2 上开始调度 G3(syscall-1)
G3(syscall-1) 在 M2/P2 上执行第 5/12 步
G3(syscall-1) 在 M2/P2 上执行第 6/12 步
M1：G4(syscall-2) 在 P1 上时间片用完，重新入队
M1 在 P1 上开始调度 G4(syscall-2)
G4(syscall-2) 在 M1/P1 上执行第 7/9 步
G4(syscall-2) 在 M1/P1 上执行第 8/9 步
G3(syscall-1) 在 M2/P2 上执行第 7/12 步
G4(syscall-2) 在 M1/P1 上执行第 9/9 步
M1：G4(syscall-2) 在 P1 上全部完成
M2：G3(syscall-1) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G3(syscall-1)
G3(syscall-1) 在 M2/P2 上执行第 8/12 步
G3(syscall-1) 在 M2/P2 上执行第 9/12 步
G3(syscall-1) 在 M2/P2 上执行第 10/12 步
M2：G3(syscall-1) 在 P2 上时间片用完，重新入队
M2 在 P2 上开始调度 G3(syscall-1)
G3(syscall-1) 在 M2/P2 上执行第 11/12 步
G3(syscall-1) 在 M2/P2 上执行第 12/12 步
M2：G3(syscall-1) 在 P2 上全部完成
M1 退出（绑定 P1）
M2 退出（绑定 P2）
所有 G 都执行完毕，调度器退出。